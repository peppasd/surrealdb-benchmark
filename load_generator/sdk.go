package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/surrealdb/surrealdb.go"
)

type SdkCustomer struct {
	ID        string    `json:"id,omitempty"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Country   string    `json:"country"`
	LastLogin time.Time `json:"last_login"`
}

type SdkQueryResult struct {
	Result []map[string]interface{} `json:"result"`
}

func runSdkBenchmark(duration time.Duration, workers int) error {

	log.Printf("Starting SDK benchmark with %d workers for %d minutes \n", workers, int(duration.Minutes()))

	ctx, ctxCancel := context.WithTimeout(context.Background(), duration)
	wg := new(sync.WaitGroup)

	for i := 0; i < workers; i++ {
		go sdkWorker(wg, ctx)
	}

	for range ctx.Done() {
		log.Println("SDK benchmark timeout. Stopping workers.")
		break
	}
	wg.Wait()
	ctxCancel()

	log.Println("SDK benchmark finished")
	return nil
}

func prepareSdk() (*surrealdb.DB, error) {
	db, err := surrealdb.New(wsUrl, surrealdb.UseWriteCompression(true))
	if err != nil {
		return nil, err
	}

	if _, err = db.Use(db_ns, db_name); err != nil {
		return nil, err
	}

	return db, nil
}

func sdkWorker(wg *sync.WaitGroup, ctx context.Context) error {
	db, err := prepareSdk()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			db.Close()
			return nil
		default:
			wg.Add(1)

			start := time.Now()
			id, err := sdkCreate(db)
			if err != nil {
				db.Close()
				wg.Done()
				log.Println(err)
				return err
			}
			final := time.Since(start)
			logResult("SDK", "create", -1, int(final.Microseconds()))

			start = time.Now()
			err = sdkRead(id, db)
			if err != nil {
				db.Close()
				wg.Done()
				log.Println(err)
				return err
			}
			final = time.Since(start)
			logResult("SDK", "read", -1, int(final.Microseconds()))

			start = time.Now()
			err = sdkUpdate(id, db)
			if err != nil {
				db.Close()
				wg.Done()
				log.Println(err)
				return err
			}
			final = time.Since(start)
			logResult("SDK", "update", -1, int(final.Microseconds()))

			start = time.Now()
			err = sdkDelete(id, db)
			if err != nil {
				db.Close()
				wg.Done()
				log.Println(err)
				return err
			}
			final = time.Since(start)
			logResult("SDK", "delete", -1, int(final.Microseconds()))

			start = time.Now()
			dur, err := sdkSelect(db)
			if err != nil {
				db.Close()
				wg.Done()
				log.Println(err)
				return err
			}
			final = time.Since(start)
			logResult("SDK", "select", dur, int(final.Microseconds()))

			start = time.Now()
			dur, err = sdkSimpleQuery(db)
			if err != nil {
				db.Close()
				wg.Done()
				log.Println(err)
				return err
			}
			final = time.Since(start)
			logResult("SDK", "query", dur, int(final.Microseconds()))

			start = time.Now()
			dur, err = sdkJoinRelation(db)
			if err != nil {
				db.Close()
				wg.Done()
				log.Println(err)
				return err
			}
			final = time.Since(start)
			logResult("SDK", "join_relation", dur, int(final.Microseconds()))

			start = time.Now()
			dur, err = sdkJoinGraph(db)
			if err != nil {
				db.Close()
				wg.Done()
				log.Println(err)
				return err
			}
			final = time.Since(start)
			logResult("SDK", "join_graph", dur, int(final.Microseconds()))

			wg.Done()
		}
	}
}

func sdkRead(id string, db *surrealdb.DB) error {
	data, err := db.Select(id)
	if err != nil {
		return err
	}
	selectedCustomer := new(SdkCustomer)
	err = surrealdb.Unmarshal(data, &selectedCustomer)
	if err != nil {
		panic(err)
	}
	return nil
}

func sdkDelete(id string, db *surrealdb.DB) error {
	_, err := db.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func sdkUpdate(id string, db *surrealdb.DB) error {
	changes := map[string]string{"email": "test2@test.com"}
	if _, err := db.Update(id, changes); err != nil {
		return err
	}
	return nil
}

func sdkCreate(db *surrealdb.DB) (string, error) {
	testCustomer := SdkCustomer{
		FirstName: "Test",
		LastName:  "Tester",
		Email:     "test@test.com",
		Country:   "Germany",
		LastLogin: time.Now(),
	}

	data, err := db.Create("customer", &testCustomer)
	if err != nil {
		return "", err
	}
	createdCustomer := make([]SdkCustomer, 1)
	err = surrealdb.Unmarshal(data, &createdCustomer)
	if err != nil {
		return "", err
	}
	if len(createdCustomer) == 0 {
		return "", nil
	}

	return createdCustomer[0].ID, nil
}

func sdkSelect(db *surrealdb.DB) (int, error) {
	query := `SELECT * FROM order LIMIT 1000`
	data, err := db.Query(query, nil)
	if err != nil {
		return 0, err
	}
	res := data.([]interface{})[0].(map[string]interface{})
	internalDur, err := time.ParseDuration(res["time"].(string))
	if err != nil {
		return 0, err
	}

	return int(internalDur.Microseconds()), nil
}

func sdkSimpleQuery(db *surrealdb.DB) (int, error) {
	query := `SELECT * FROM order WHERE processed IS FALSE LIMIT 1000`
	data, err := db.Query(query, nil)
	if err != nil {
		return 0, err
	}
	res := data.([]interface{})[0].(map[string]interface{})
	internalDur, err := time.ParseDuration(res["time"].(string))
	if err != nil {
		return 0, err
	}

	return int(internalDur.Microseconds()), nil
}

func sdkJoinRelation(db *surrealdb.DB) (int, error) {
	query := `SELECT books.title FROM order WHERE processed IS TRUE LIMIT 1000`
	data, err := db.Query(query, nil)
	if err != nil {
		return 0, err
	}
	res := data.([]interface{})[0].(map[string]interface{})
	internalDur, err := time.ParseDuration(res["time"].(string))
	if err != nil {
		return 0, err
	}

	return int(internalDur.Microseconds()), nil
}

func sdkJoinGraph(db *surrealdb.DB) (int, error) {
	query := `SELECT <-ordered<-customer.first_name FROM order WHERE processed IS TRUE LIMIT 1000`
	data, err := db.Query(query, nil)
	if err != nil {
		return 0, err
	}
	res := data.([]interface{})[0].(map[string]interface{})
	internalDur, err := time.ParseDuration(res["time"].(string))
	if err != nil {
		return 0, err
	}

	return int(internalDur.Microseconds()), nil
}
