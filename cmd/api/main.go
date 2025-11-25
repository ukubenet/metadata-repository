package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	config "github.com/ukubenet/metadata-repository/config"
)

const version = "1.0.0"

type appConfig struct {
	port int
	env  string
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type application struct {
	appConfig appConfig
	logger    *log.Logger
}

func main() {
	os.Setenv("OPENAI_API_KEY", "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")
	os.Setenv("OPENAI_MODEL", "gpt-3.5-turbo")
	config.LoadConfig("./config", "app")
	var cfg appConfig

	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "dev", "Application environment (dev|prod)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		appConfig: cfg,
		logger:    logger,
	}

	// Set app name
	//entitysearch.Indexes.LoadAllIndexes()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Starting server on port", cfg.port)

	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
