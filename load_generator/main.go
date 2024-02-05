package main

import (
	"flag"
	"log"
	"time"
)

const (
	db_ns   = "benchmark"
	db_name = "benchmark"
)

var (
	url   = "http://localhost:8000"
	wsUrl = "ws://localhost:8000/rpc"
)

func main() {
	minutes := flag.Int("minutes", 1, "How many minutes to run each benchmark phase")
	workers := flag.Int("threads", 1, "How many workers/threads to use for each benchmark phase")
	flagUrl := flag.String("url", "localhost:8000", "URL of the server to benchmark. Example: localhost:8000 DO NOT INCLUDE THE PROTOCOL")
	flag.Parse()
	benchmarkDuration := time.Minute * time.Duration(*minutes)
	benchmarkWorkers := *workers
	url = "http://" + *flagUrl
	wsUrl = "ws://" + *flagUrl + "/rpc"
	log.Printf("Starting benchmark with phase duration %v and %v threads per phase on %s", benchmarkDuration, benchmarkWorkers, url)

	err := runHealthcheck()
	if err != nil {
		log.Fatalf("Healthcheck failed: %v", err)
	}
	log.Println("Surreal healthcheck passed")

	// creates the database if it doesn't exist, deletes old data if they exist
	err = resultDbInit()
	if err != nil {
		log.Fatalf("Failed to initialize results database: %v", err)
	}
	log.Println("Results database initialized")

	err = runRestBenchmark(benchmarkDuration, benchmarkWorkers)
	if err != nil {
		log.Fatalf("REST benchmark failed: %v", err)
	}

	err = runWebsocketBenchmark(benchmarkDuration, benchmarkWorkers)
	if err != nil {
		log.Fatalf("Websocket benchmark failed: %v", err)
	}

	err = runSdkBenchmark(benchmarkDuration, benchmarkWorkers)
	if err != nil {
		log.Fatalf("SDK benchmark failed: %v", err)
	}

	log.Println("Benchmark finished")
}
