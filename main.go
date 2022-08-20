package main

import (
	"context"
	"golang-developer-test-task/redclient"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	port := "3000"

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() {
		err = logger.Sync()
	}()

	ctx := context.TODO()
	conf := redclient.RedisConfig{}
	conf.Load()

	client := redclient.NewRedisClient(ctx, conf)
	defer func() {
		err = client.Close()
		if err != nil {
			panic(err)
		}
	}()

	dbLogic := DBProcessor{client: client, logger: logger}
	mux := http.NewServeMux()

	mux.HandleFunc("/api/load_file", dbLogic.HandleLoadFile)

	mux.HandleFunc("/api/load_from_url", dbLogic.HandleLoadFromURL)

	//https://nimblehq.co/blog/getting-started-with-redisearch
	mux.HandleFunc("/api/search", dbLogic.HandleSearch)

	//mux.HandleFunc("/api/metrics", func(w http.ResponseWriter, r *http.Request) {
	//
	//})

	mux.HandleFunc("/", dbLogic.HandleMainPage)

	//mux := DBProcessor{}
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		panic(err)
	}
}
