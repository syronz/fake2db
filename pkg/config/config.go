package config

import "database/sql"

type Config struct {
	Worker      int
	ChunkSize   int `toml:"chunk_size"`
	Rows        int
	RandomLevel int `toml:"random_level"`
	Database    struct {
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
