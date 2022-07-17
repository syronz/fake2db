package main

import (
	"constraints"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"

	"github.com/syronz/fake2db/pkg/config"
	"github.com/syronz/fake2db/pkg/fake"
)

var configFile = flag.String("config", "config.toml", "path to config file, default is config.toml")

func main() {
	flag.Parse()
	start := time.Now()

	var cfg config.Config
	if _, err := toml.DecodeFile(*configFile, &cfg); err != nil {
		log.Fatal("failed in decoding the toml file for terms", err)
	}

	fake.InitiatePattern(cfg)

	if cfg.ChunkSize > cfg.Rows {
		log.Fatalln("chunk_size is bigger than rows. it is forbidden!")
	}

	var err error
	cfg.Database.Destination.Conn, err = sql.Open(cfg.Database.Destination.Type, cfg.Database.Destination.DSN)
	if err != nil {
		log.Fatalln("error in opening database", err)
	}

	err = cfg.Database.Destination.Conn.Ping()
	if err != nil {
		log.Fatalln("error in ping to database", err)
	}

	for _, v := range cfg.Database.Destination.PreQueries {
		_, err := cfg.Database.Destination.Conn.Exec(v)
		if err != nil {
			fmt.Println("prequery: ", v)
			log.Fatalln("error in executing pre-query", err)
		}
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln("error in closing database gracefully", err)
		}
	}(cfg.Database.Destination.Conn)

	fmt.Printf("Start inserting data:\nrows:\t%d\nworker:\t%d\nchunk:\t%d\n", cfg.Rows, cfg.Worker, cfg.ChunkSize)

	var wgConsumer, wgPublisher sync.WaitGroup
	countCh := make(chan int, 1)
	fakeFactory, _ := fake.NewFactory(cfg.Query.ValuesClause)

	wgPublisher.Add(1)
	go func(wgConsumer, wgPublisher *sync.WaitGroup, rows, chunkSize int, countCh chan int) {
		for i := 0; i < rows; i += min(chunkSize, rows-i) {
			wgConsumer.Add(1)
			countCh <- min(chunkSize, rows-i)
			fmt.Printf("\rtime: %ds inserted: %v - %v%%",
				int64(time.Since(start)/time.Second),
				i+min(chunkSize, rows-i),
				math.Round(float64(i)/float64(rows)*100),
			)
		}
		wgPublisher.Done()
	}(&wgConsumer, &wgPublisher, cfg.Rows, cfg.ChunkSize, countCh)

	for w := 1; w <= cfg.Worker; w++ {
		go Worker(&wgConsumer, cfg, fakeFactory, countCh)
	}

	//loadAvg, _ := load.Avg()
	//fmt.Println(" 1 min ave:", loadAvg.Load1)
	//fmt.Println(" 5 min ave:", loadAvg.Load5)
	//fmt.Println("15 min ave:", loadAvg.Load15)

	wgPublisher.Wait()
	wgConsumer.Wait()
	fmt.Println("\nduration: ", time.Since(start))

}

func Worker(wg *sync.WaitGroup, cfg config.Config, factory fake.Factory, countCh chan int) {

	for count := range countCh {
		query := strings.Builder{}
		query.WriteString(cfg.Query.InsertClause)

		var insideValues strings.Builder
		for i := 0; i < count; i++ {
			query.WriteRune('(')

			insideValues.Reset()
			insideValues.WriteString(fmt.Sprintf("'%v',", strings.Join(factory(), "','")))
			query.WriteString(insideValues.String()[0 : len(insideValues.String())-1])
			query.WriteString("),")
		}

		_, err := cfg.Database.Destination.Conn.Exec(query.String()[0 : len(query.String())-1])
		if err != nil {
			log.Println("error in executing the query", err)
			//log.Println("query: ", query.String()[0:100])
		}
		wg.Done()
	}

}

func min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}
