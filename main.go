package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

var configFile = flag.String("config", "config.toml", "path to config file, default is config.toml")

type Config struct {
	ProcessName  string `toml:"process_name"`
	Worker       int    `toml:"worker"`
	ChunkSize    int    `toml:"chunk_size"`
	Query        string `toml:"query"`
	Values       string `toml:"values"`
	OutputPrefix string `toml:"output_prefix"`
}

func main() {
	flag.Parse()
	fmt.Println("hello", *configFile)

	var config Config

	if _, err := toml.DecodeFile(*configFile, &config); err != nil {
		log.Fatal("failed in decoding the toml file for terms", err)
	}

	fmt.Println(">>>>>", config)

}
