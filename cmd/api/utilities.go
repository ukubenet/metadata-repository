package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"net/http"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, wrap string) error {
	wrapper := make(map[string]interface{})

	wrapper[wrap] = data

	js, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) errorJSON(w http.ResponseWriter, err error) {
	type jsonError struct {
		Message string `json:"message"`
	}

	theError := jsonError{
		Message: err.Error(),
	}

	app.writeJSON(w, http.StatusBadRequest, theError, "error")
}

func (app *application) writeRedis(key string, value string) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "eYVX4HwVmsxKPCDjwMtyKVge5oLd9t42",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	err = client.Set(key, value, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}
