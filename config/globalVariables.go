package config

import (
	"database/sql"
)

var (
	DaemonURL                 string
	APIPort                   string
	RateLimitTime             float64
	RateLimitRequests         int64
	DatabaseName              string
	DatabaseHost              string
	DatabaseUsername          string
	DatabasePassword          string
	DB                        *sql.DB
	EmissionHistoryFileName   string
	OracleNodeHistoryFileName string
)
