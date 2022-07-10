package pkg

import "time"

const (
	MinNumber = 1
	MaxNumber = 1_000_000_000

	DateLayout = "2006-01-02T15:04:05Z"
)

var (
	MinDate = time.Date(1950, 01, 01, 0, 0, 0, 0, time.UTC)
	MaxDate = time.Now()
)
