package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Result struct {
	ID                           int `gorm:"primaryKey"`
	ConnectionType               string
	QueryType                    string
	InternalDurationMicroSeconds int
	TotalDurationMicroSeconds    int
	CreatedAt                    time.Time `gorm:"autoCreateTime"`
}

const dbName = "results.sqlite"

var db *gorm.DB

func resultDbInit() error {
	var err error
	if _, err = os.Stat(dbName); err == nil {
		err = os.Remove(dbName)
		log.Println("Removing existing database")
		if err != nil {
			return err
		}
	}

	db, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		return err
	}
	db.AutoMigrate(&Result{})
	return nil
}

func logResult(connection string, query string, internalDuration int, totalDuration int) {
	res := Result{
		ConnectionType:               connection,
		QueryType:                    query,
		InternalDurationMicroSeconds: internalDuration,
		TotalDurationMicroSeconds:    totalDuration,
	}
	db.Create(&res)
}
