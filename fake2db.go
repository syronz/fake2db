package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"
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
			Type       string
			DSN        string
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

	var config Config

	if _, err := toml.DecodeFile(*configFile, &config); err != nil {
		log.Fatal("failed in decoding the toml file for terms", err)
	}

	//loadAvg, _ := load.Avg()
	//fmt.Println(" 1 min ave:", loadAvg.Load1)
	//fmt.Println(" 5 min ave:", loadAvg.Load5)
	//fmt.Println("15 min ave:", loadAvg.Load15)

	qMarks, columnCount := fake.ValueClauseToQuestionMark(config.Query.ValuesClause)
	fmt.Println(">>>>> qmarks", qMarks, columnCount)

	//qArr := make([]interface{}, config.Rows*columnCount)
	fakeFactory, _ := fake.NewFactory(config.Query.ValuesClause)
	fmt.Println(fakeFactory())

	for z := 0; z < 5; z++ {

		//for i := 0; i < config.Rows; i++ {
		//	row := fakeFactory()
		//	for j, v := range row {
		//		qArr[i*columnCount+j] = v
		//	}
		//}

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

		start := time.Now()
		for i := 0; i < 10; i++ {
			query := strings.Builder{}
			query.WriteString("INSERT INTO students(name, gender, code, dob, address, created_at) VALUES")

			/* stms approach 10,000 around 1.387s
			query.WriteString(strings.Repeat(qMarks+",", config.Rows-1))
			query.WriteString(qMarks)
			stmt, err := db.Prepare(query.String())
			if err != nil {
				log.Fatalln("error in creating stmt", err)
			}

			_, err = stmt.Exec(qArr...)
			if err != nil {
				log.Fatalln("error in executing the query", err)
			}*/

			/* 20,000 in 1.97
			var insideValues strings.Builder
			for i := 0; i < config.Rows; i++ {
				query.WriteString("(")

				insideValues.Reset()
				for j := 0; j < columnCount; j++ {
					insideValues.WriteString(fmt.Sprintf("'%v',", qArr[i*columnCount+j]))
				}
				query.WriteString(insideValues.String()[0 : len(insideValues.String())-1])
				query.WriteString("),")
			}
			_, err = db.Exec(query.String()[0 : len(query.String())-1])
			if err != nil {
				log.Fatalln("error in executing the query", err)
			}*/

			allRandoms := make([]string, config.Rows-1)
			for i := 0; i < config.Rows-1; i++ {
				allRandoms[i] = fakeFactory()
			}

			query.WriteString(strings.Join(allRandoms, ","))

			_, err = db.Exec(query.String())
			if err != nil {
				log.Fatalln("error in executing the query", err)
			}

		}
		fmt.Println("duration: ", time.Since(start))
	}

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
