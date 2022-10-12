// Copyright (c) 2022 EPAM Systems, Inc.
// 
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/agilestacks/stack-disposer/disposer/api"
	"github.com/agilestacks/stack-disposer/disposer/config"
)

func init() {
	defaultPort, exists := os.LookupEnv("PORT")
	if !exists || len(defaultPort) <= 0 {
		defaultPort = "8080"
	}

	flag.StringVar(&config.Port, "port", defaultPort, "port where to listen")
	flag.StringVar(&config.GitUrl, "gitUrl", "https://github.com/agilestacks/google-stacks.git", "Git URL with stacks")
	flag.StringVar(&config.GitDir, "gitDir", "/tmp/stacks", "directory where clone stacks to")
	flag.BoolVar(&config.Verbose, "verbose", false, "verbose logging")
	flag.DurationVar(&config.Timeout, "timeout", 1*time.Hour, "request timeout")

	flag.Parse()

	if config.Verbose {
		log.Print("VERBOSE logging enabled")
	}
}

func handleRequests() {
	router := mux.NewRouter()
	router.HandleFunc("/{sandboxId}/{stackId}", api.UndeployStackHandler).Methods(http.MethodDelete)

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%s", config.Port),
		WriteTimeout: config.Timeout,
		ReadTimeout:  config.Timeout,
	}

	log.Println("Server listens on port", config.Port)
	log.Fatal(srv.ListenAndServe())
}

func main() {
	handleRequests()
}
