package main

import (
	"os"
	"time"
)

type CDBStats struct {
	CDBSize int64   `prometheus:"data_cdb_size" prometheus_help:"in bytes"`
	CDBAge  float64 `prometheus:"data_cdb_age" prometheus_help:"in seconds"`
}

func checkCDB(path string) (CDBStats, error) {
	info, err := os.Stat(path)
	if err != nil {
		return CDBStats{}, err
	}

	return CDBStats{
		CDBSize: info.Size(),
		CDBAge:  time.Now().Sub(info.ModTime()).Seconds(),
	}, nil
}
