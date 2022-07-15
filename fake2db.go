package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"

	"github.com/syronz/fake2db/pkg/fake"
)

var configFile = flag.String("config", "config.toml", "path to config file, default is config.toml")

type Config struct {
	ProcessName  string
	Worker       int
	ChunkSize    int
	Rows         int
	OutputPrefix string
	Database     struct {
		Source struct {
			Type string
			DSN  string
		}
		Destination struct {
			Type string
			DSN  string
		}
	}
	Query struct {
		InsertClause string `toml:"insert_clause"`
		ValuesClause string `toml:"values_clause"`
	}
}

func main() {
	flag.Parse()

	var config Config

	if _, err := toml.DecodeFile(*configFile, &config); err != nil {
		log.Fatal("failed in decoding the toml file for terms", err)
	}

	//loadAvg, _ := load.Avg()
	//fmt.Println(" 1 min ave:", loadAvg.Load1)
	//fmt.Println(" 5 min ave:", loadAvg.Load5)
	//fmt.Println("15 min ave:", loadAvg.Load15)

	fakeFactory, _ := fake.NewFactory(config.Query.ValuesClause)
	for i := 0; i < 1000; i++ {
		fmt.Printf(">>>>>>> fakerFactory %v - %+v\n", i, fakeFactory())
	}

	db, err := sql.Open(config.Database.Source.Type, config.Database.Source.DSN)
	if err != nil {
		log.Fatalln("error in opening database", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln("error in closing database gracefully", err)
		}
	}(db)

	fmt.Println(">>>>>", config.Database.Source.DSN, config.Rows)

	start := time.Now()
	for i := 0; i < 10000; i++ {
		query := "INSERT INTO students(first_name, last_name, gender, age, dob, description, created_at)"
		query += "VALUES('diako', 'sharifi', 'male', 35, '1987-02-17', 'he is a programmer', '2022-07-10 00:14:45')"
		_, err = db.Exec(query)
		if err != nil {
			log.Fatalln("error in executing the query", err)
		}
	}
	fmt.Println("duration: ", time.Since(start))

	//start := time.Now()
	//query := "INSERT INTO students(first_name, last_name, gender, age, dob, description, created_at)"
	//query += "VALUES(?, ?, ?, ?, ?, ?, ?)"
	//stmt, err := db.Prepare(query)
	//if err != nil {
	//	log.Fatalln("error in creating stmt", err)
	//}
	//
	//for i := 0; i < 10000; i++ {
	//	_, err = stmt.Exec("diako", "sharifi", "male", 35, "1987-02-17", "he is a programmer", "2022-07-10 00:14:45")
	//	if err != nil {
	//		log.Fatalln("error in executing the query", err)
	//	}
	//
	//}
	//fmt.Println("duration: ", time.Since(start))
}
