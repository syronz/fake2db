package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

var configFile = flag.String("config", "config.toml", "path to config file, default is config.toml")

type Config struct {
	ProcessName  string
	Worker       int
	ChunkSize    int
	Rows         int
	OutputPrefix string
	Database     struct {
		Type string
		DSN  string
	}
	Query struct {
		InsertClause string
		ValuesClause string
	}
}

func main() {
	flag.Parse()

	var config Config

	if _, err := toml.DecodeFile(*configFile, &config); err != nil {
		log.Fatal("failed in decoding the toml file for terms", err)
	}

	db, err := sql.Open(config.Database.Type, config.Database.DSN)
	if err != nil {
		log.Fatalln("error in opening database", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln("error in closing database gracefully", err)
		}
	}(db)

	fmt.Println(">>>>>", config.Database.DSN, config.Rows)

}
