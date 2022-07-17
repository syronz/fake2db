package main

import (
	"constraints"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"

	"github.com/syronz/fake2db/pkg/fake"
)

var configFile = flag.String("config", "config.toml", "path to config file, default is config.toml")

type Config struct {
	ProcessName  string
	Worker       int
	ChunkSize    int `toml:"chunk_size"`
	Rows         int
	OutputPrefix string
	Database     struct {
		Source struct {
			Type string
			DSN  string
			Conn *sql.DB
		}
		Destination struct {
			Type       string
			DSN        string
			Conn       *sql.DB
			PreQueries []string `toml:"pre_queries"`
		}
	}
	Query struct {
		InsertClause string `toml:"insert_clause"`
		ValuesClause string `toml:"values_clause"`
	}
}

func main() {
	flag.Parse()
	start := time.Now()

	var config Config

	if _, err := toml.DecodeFile(*configFile, &config); err != nil {
		log.Fatal("failed in decoding the toml file for terms", err)
	}

	if config.ChunkSize > config.Rows {
		log.Fatalln("chunk_size is bigger than rows. it is forbidden!")
	}

	var err error
	config.Database.Destination.Conn, err = sql.Open(config.Database.Destination.Type, config.Database.Destination.DSN)
	if err != nil {
		log.Fatalln("error in opening database", err)
	}

	err = config.Database.Destination.Conn.Ping()
	if err != nil {
		log.Fatalln("error in ping to database", err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln("error in closing database gracefully", err)
		}
	}(config.Database.Destination.Conn)

	var wgConsumer, wgPublisher sync.WaitGroup
	countCh := make(chan int, 1)
	fakeFactory, _ := fake.NewFactory(config.Query.ValuesClause)

	wgPublisher.Add(1)
	go func(wgConsumer, wgPublisher *sync.WaitGroup, rows, chunkSize int, countCh chan int) {
		for i := 0; i < rows; i += min(chunkSize, rows-i) {
			fmt.Println("... find count >>>> ", i, chunkSize, rows, min(chunkSize, rows-0))
			wgConsumer.Add(1)
			countCh <- min(chunkSize, rows-i)
		}
		wgPublisher.Done()
	}(&wgConsumer, &wgPublisher, config.Rows, config.ChunkSize, countCh)

	for w := 1; w <= config.Worker; w++ {
		go Worker(&wgConsumer, config, fakeFactory, countCh)
	}

	//loadAvg, _ := load.Avg()
	//fmt.Println(" 1 min ave:", loadAvg.Load1)
	//fmt.Println(" 5 min ave:", loadAvg.Load5)
	//fmt.Println("15 min ave:", loadAvg.Load15)

	wgPublisher.Wait()
	wgConsumer.Wait()
	fmt.Println("duration: ", time.Since(start))

}

func Worker(wg *sync.WaitGroup, config Config, factory fake.Factory, countCh chan int) {

	for count := range countCh {
		fmt.Println(">>>>>! 000 ", count)
		query := strings.Builder{}
		query.WriteString(config.Query.InsertClause)

		var insideValues strings.Builder
		for i := 0; i < count; i++ {
			query.WriteRune('(')

			insideValues.Reset()
			insideValues.WriteString(fmt.Sprintf("'%v',", strings.Join(factory(), "','")))
			query.WriteString(insideValues.String()[0 : len(insideValues.String())-1])
			query.WriteString("),")
		}

		_, err := config.Database.Destination.Conn.Exec(query.String()[0 : len(query.String())-1])

		if err != nil {
			log.Println("error in executing the query", err)
			log.Println("query: ", query.String())
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
