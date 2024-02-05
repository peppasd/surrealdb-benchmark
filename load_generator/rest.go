package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

func runRestBenchmark(duration time.Duration, workers int) error {

	log.Printf("Starting REST benchmark with %d workers for %d minutes \n", workers, int(duration.Minutes()))

	ctx, ctxCancel := context.WithTimeout(context.Background(), duration)
	wg := new(sync.WaitGroup)

	for i := 0; i < workers; i++ {
		go restWorker(wg, ctx)
	}

	for range ctx.Done() {
		log.Println("REST benchmark timeout. Stopping workers.")
		break
	}
	wg.Wait()
	ctxCancel()

	log.Println("REST benchmark finished")
	return nil
}

func restWorker(wg *sync.WaitGroup, ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			wg.Add(1)

			start := time.Now()
			id, dur, err := restCreate()
			if err != nil {
				wg.Done()
				log.Println(err)
				return err
			}
			final := time.Since(start)
			logResult("REST", "create", dur, int(final.Microseconds()))

			start = time.Now()
			dur, err = restRead(id)
			if err != nil {
				wg.Done()
				log.Println(err)
				return err
			}
			final = time.Since(start)
			logResult("REST", "read", dur, int(final.Microseconds()))

			start = time.Now()
			dur, err = restUpdate(id)
			if err != nil {
				wg.Done()
				log.Println(err)
				return err
			}
			final = time.Since(start)
			logResult("REST", "update", dur, int(final.Microseconds()))

			start = time.Now()
			dur, err = restDelete(id)
			if err != nil {
				wg.Done()
				log.Println(err)
				return err
			}
			final = time.Since(start)
			logResult("REST", "delete", dur, int(final.Microseconds()))

			start = time.Now()
			dur, err = restSelect()
			if err != nil {
				wg.Done()
				log.Println(err)
				return err
			}
			final = time.Since(start)
			logResult("REST", "select", dur, int(final.Microseconds()))

			start = time.Now()
			dur, err = restSimpleQuery()
			if err != nil {
				wg.Done()
				log.Println(err)
				return err
			}
			final = time.Since(start)
			logResult("REST", "query", dur, int(final.Microseconds()))

			start = time.Now()
			dur, err = restJoinRelation()
			if err != nil {
				wg.Done()
				log.Println(err)
				return err
			}
			final = time.Since(start)
			logResult("REST", "join_relation", dur, int(final.Microseconds()))

			start = time.Now()
			dur, err = restJoinGraph()
			if err != nil {
				wg.Done()
				log.Println(err)
				return err
			}
			final = time.Since(start)
			logResult("REST", "join_graph", dur, int(final.Microseconds()))

			wg.Done()
		}
	}
}

func doRequest(method string, path string, body io.Reader) ([]map[string]interface{}, error) {
	req, err := http.NewRequest(method, url+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("NS", db_ns)
	req.Header.Set("DB", db_name)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %v", resp.Status)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("empty response")
	}
	if result[0]["status"] != "OK" {
		return nil, errors.New("status not OK")
	}
	return result, nil
}

func restRead(id string) (int, error) {
	resp, err := doRequest("GET", "/key/customer/"+id, nil)
	if err != nil {
		return 0, err
	}
	internalDur, err := time.ParseDuration(resp[0]["time"].(string))
	if err != nil {
		return 0, err
	}
	return int(internalDur.Microseconds()), nil
}

func restDelete(id string) (int, error) {
	resp, err := doRequest("DELETE", "/key/customer/"+id, nil)
	if err != nil {
		return 0, err
	}
	internalDur, err := time.ParseDuration(resp[0]["time"].(string))
	if err != nil {
		return 0, err
	}
	return int(internalDur.Microseconds()), nil
}

func restUpdate(id string) (int, error) {
	var body = strings.NewReader(`{"email":"test2@test.com"}`)
	resp, err := doRequest("PATCH", "/key/customer/"+id, body)
	if err != nil {
		return 0, err
	}
	internalDur, err := time.ParseDuration(resp[0]["time"].(string))
	if err != nil {
		return 0, err
	}
	return int(internalDur.Microseconds()), nil
}

func restCreate() (string, int, error) {
	var body = strings.NewReader(`{"first_name":"Test","last_name":"Tester","email":"test@test.com","country":"Germany","last_login":"2024-02-03T21:31:22+0000"}`)

	resp, err := doRequest("POST", "/key/customer", body)
	if err != nil {
		return "", 0, err
	}
	internalDur, err := time.ParseDuration(resp[0]["time"].(string))
	if err != nil {
		return "", 0, err
	}
	fullId := resp[0]["result"].([]interface{})[0].(map[string]interface{})["id"].(string)
	id := strings.Split(fullId, ":")[1]
	return id, int(internalDur.Microseconds()), nil
}

func restSelect() (int, error) {
	var body = strings.NewReader(`SELECT * FROM order LIMIT 1000`)
	resp, err := doRequest("POST", "/sql", body)
	if err != nil {
		return 0, err
	}
	internalDur, err := time.ParseDuration(resp[0]["time"].(string))
	if err != nil {
		return 0, err
	}
	return int(internalDur.Microseconds()), nil
}

func restSimpleQuery() (int, error) {
	var body = strings.NewReader(`SELECT * FROM order WHERE processed IS FALSE`)
	resp, err := doRequest("POST", "/sql", body)
	if err != nil {
		return 0, err
	}
	internalDur, err := time.ParseDuration(resp[0]["time"].(string))
	if err != nil {
		return 0, err
	}
	return int(internalDur.Microseconds()), nil
}

func restJoinRelation() (int, error) {
	var body = strings.NewReader(`SELECT books.title FROM order WHERE processed IS TRUE`)
	resp, err := doRequest("POST", "/sql", body)
	if err != nil {
		return 0, err
	}
	internalDur, err := time.ParseDuration(resp[0]["time"].(string))
	if err != nil {
		return 0, err
	}
	return int(internalDur.Microseconds()), nil
}

func restJoinGraph() (int, error) {
	var body = strings.NewReader(`SELECT <-ordered<-customer.first_name FROM order WHERE processed IS TRUE`)
	resp, err := doRequest("POST", "/sql", body)
	if err != nil {
		return 0, err
	}
	internalDur, err := time.ParseDuration(resp[0]["time"].(string))
	if err != nil {
		return 0, err
	}
	return int(internalDur.Microseconds()), nil
}
