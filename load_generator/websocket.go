package main

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type WebsocketSend struct {
	Id     int      `json:"id"`
	Method string   `json:"method"`
	Params []string `json:"params"`
}

type WebsocketReceive struct {
	Id     int                      `json:"id"`
	Result []map[string]interface{} `json:"result"`
}

func runWebsocketBenchmark(duration time.Duration, workers int) error {
	log.Printf("Starting Websocket benchmark with %d workers for %d minutes \n", workers, int(duration.Minutes()))

	ctx, ctxCancel := context.WithTimeout(context.Background(), duration)
	wg := new(sync.WaitGroup)

	for i := 0; i < workers; i++ {
		go websocketWorker(wg, ctx)
	}

	for range ctx.Done() {
		log.Println("Websocket benchmark timeout. Stopping workers.")
		break
	}
	wg.Wait()
	ctxCancel()

	log.Println("Websocket benchmark finished")
	return nil
}

func prepareWebsocket() (*websocket.Conn, error) {

	ws, err := websocket.Dial(wsUrl, "", url)
	if err != nil {
		return nil, err
	}
	ws.MaxPayloadBytes = 1024 * 1024 * 1024
	if _, err := ws.Write([]byte(`{"id":1,"method":"use","params":["` + db_ns + `", "` + db_name + `"]}`)); err != nil {
		return nil, err
	}
	var msg WebsocketReceive
	if err = websocket.JSON.Receive(ws, &msg); err != nil {
		return nil, err
	}

	if msg.Id != 1 {
		return nil, errors.New("unexpected websocket response id")
	}

	return ws, nil
}

func websocketWorker(wg *sync.WaitGroup, ctx context.Context) error {
	ws, err := prepareWebsocket()
	if err != nil {
		return err
	}
	nextId := 2
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			ws.Close()
			return nil
		default:
			wg.Add(1)

			start := time.Now()
			id, dur, err := websocketCreate(ws, nextId)
			if err != nil {
				ws.Close()
				wg.Done()
				log.Println(err)
				return err
			}
			final := time.Since(start)
			logResult("Websocket", "create", dur, int(final.Microseconds()))
			nextId++

			start = time.Now()
			dur, err = websocketRead(id, ws, nextId)
			if err != nil {
				ws.Close()
				wg.Done()
				log.Println(err)
				return err
			}
			final = time.Since(start)
			logResult("Websocket", "read", dur, int(final.Microseconds()))
			nextId++

			start = time.Now()
			dur, err = websocketUpdate(id, ws, nextId)
			if err != nil {
				ws.Close()
				wg.Done()
				log.Println(err)
				return err
			}
			final = time.Since(start)
			logResult("Websocket", "update", dur, int(final.Microseconds()))
			nextId++

			start = time.Now()
			dur, err = websocketDelete(id, ws, nextId)
			if err != nil {
				ws.Close()
				wg.Done()
				log.Println(err)
				return err
			}
			final = time.Since(start)
			logResult("Websocket", "delete", dur, int(final.Microseconds()))
			nextId++

			start = time.Now()
			dur, err = websocketSelect(ws, nextId)
			if err != nil {
				ws.Close()
				wg.Done()
				log.Println(err)
				return err
			}
			final = time.Since(start)
			logResult("Websocket", "select", dur, int(final.Microseconds()))
			nextId++

			start = time.Now()
			dur, err = websocketSimpleQuery(ws, nextId)
			if err != nil {
				ws.Close()
				wg.Done()
				log.Println(err)
				return err
			}
			final = time.Since(start)
			logResult("Websocket", "query", dur, int(final.Microseconds()))
			nextId++

			start = time.Now()
			dur, err = websocketJoinRelation(ws, nextId)
			if err != nil {
				ws.Close()
				wg.Done()
				log.Println(err)
				return err
			}
			final = time.Since(start)
			logResult("Websocket", "join_relation", dur, int(final.Microseconds()))
			nextId++

			start = time.Now()
			dur, err = websocketJoinGraph(ws, nextId)
			if err != nil {
				ws.Close()
				wg.Done()
				log.Println(err)
				return err
			}
			final = time.Since(start)
			logResult("Websocket", "join_graph", dur, int(final.Microseconds()))
			nextId++

			wg.Done()
		}
	}
}

func wsSendMessage(ws *websocket.Conn, id int, query string) ([]map[string]interface{}, error) {
	sMsg := WebsocketSend{
		Id:     id,
		Method: "query",
		Params: []string{query},
	}
	if err := websocket.JSON.Send(ws, sMsg); err != nil {
		return nil, err
	}

	// bytes := make([]byte, 1024)
	// _, err := ws.Read(bytes)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Printf("Received: %s\n", string(bytes))

	var msg WebsocketReceive
	if err := websocket.JSON.Receive(ws, &msg); err != nil {
		return nil, err
	}
	if msg.Id != id {
		return nil, errors.New("unexpected websocket response id")
	}
	if len(msg.Result) == 0 {
		return nil, errors.New("empty response")
	}
	if msg.Result[0]["status"] != "OK" {
		return nil, errors.New("status not OK")
	}
	return msg.Result, nil
}

func websocketRead(id string, ws *websocket.Conn, msgId int) (int, error) {
	query := `SELECT * FROM customer:` + id + `;`
	resp, err := wsSendMessage(ws, msgId, query)
	if err != nil {
		return 0, err
	}
	internalDur, err := time.ParseDuration(resp[0]["time"].(string))
	if err != nil {
		return 0, err
	}
	return int(internalDur.Microseconds()), nil
}

func websocketDelete(id string, ws *websocket.Conn, msgId int) (int, error) {
	query := `DELETE customer:` + id + `;`
	resp, err := wsSendMessage(ws, msgId, query)
	if err != nil {
		return 0, err
	}
	internalDur, err := time.ParseDuration(resp[0]["time"].(string))
	if err != nil {
		return 0, err
	}
	return int(internalDur.Microseconds()), nil
}

func websocketUpdate(id string, ws *websocket.Conn, msgId int) (int, error) {
	query := `UPDATE customer:` + id + ` SET email = 'test2@test.com';`
	resp, err := wsSendMessage(ws, msgId, query)
	if err != nil {
		return 0, err
	}
	internalDur, err := time.ParseDuration(resp[0]["time"].(string))
	if err != nil {
		return 0, err
	}
	return int(internalDur.Microseconds()), nil
}

func websocketCreate(ws *websocket.Conn, msgId int) (string, int, error) {
	query := `CREATE customer SET first_name='Test', last_name = 'Tester', email = 'test@test.com', country = 'Germany', last_login = "2024-02-03T21:31:22+0000";`
	resp, err := wsSendMessage(ws, msgId, query)
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

func websocketSelect(ws *websocket.Conn, msgId int) (int, error) {
	query := `SELECT * FROM order LIMIT 1000`
	resp, err := wsSendMessage(ws, msgId, query)
	if err != nil {
		return 0, err
	}
	internalDur, err := time.ParseDuration(resp[0]["time"].(string))
	if err != nil {
		return 0, err
	}
	return int(internalDur.Microseconds()), nil
}

func websocketSimpleQuery(ws *websocket.Conn, msgId int) (int, error) {
	query := `SELECT * FROM order WHERE processed IS FALSE LIMIT 1000`
	resp, err := wsSendMessage(ws, msgId, query)
	if err != nil {
		return 0, err
	}
	internalDur, err := time.ParseDuration(resp[0]["time"].(string))
	if err != nil {
		return 0, err
	}
	return int(internalDur.Microseconds()), nil
}

func websocketJoinRelation(ws *websocket.Conn, msgId int) (int, error) {
	query := `SELECT books.title FROM order WHERE processed IS TRUE LIMIT 1000`
	resp, err := wsSendMessage(ws, msgId, query)
	if err != nil {
		return 0, err
	}
	internalDur, err := time.ParseDuration(resp[0]["time"].(string))
	if err != nil {
		return 0, err
	}
	return int(internalDur.Microseconds()), nil
}

func websocketJoinGraph(ws *websocket.Conn, msgId int) (int, error) {
	query := `SELECT <-ordered<-customer.first_name FROM order WHERE processed IS TRUE LIMIT 1000`
	resp, err := wsSendMessage(ws, msgId, query)
	if err != nil {
		return 0, err
	}
	internalDur, err := time.ParseDuration(resp[0]["time"].(string))
	if err != nil {
		return 0, err
	}
	return int(internalDur.Microseconds()), nil
}
