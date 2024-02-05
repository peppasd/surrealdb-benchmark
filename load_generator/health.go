package main

import (
	"errors"
	"net/http"
)

func runHealthcheck() error {
	resp, err := http.Get(url + "/health")
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("db is unhealthy")
	}
	return nil
}
