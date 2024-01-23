/*
 The MIT License (MIT)

Copyright © 2024 Kasianov Nikolai Alekseevich (Unbewohnte)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const ConfName string = "conf.json"

const Version string = "0.1-server"

var (
	port    *uint = flag.Uint("port", 80, "Set server port")
	version *bool = flag.Bool("version", false, "Print version information")
)

func main() {
	log.SetOutput(os.Stdout)
	flag.Parse()

	if *version {
		fmt.Printf("HTCPCP-server %s\nKasianov Nikolai Alekseevich (Unbewohnte)\n", Version)
		return
	}

	// Work out the working directory
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to retrieve executable's path")
	}
	wDir := filepath.Dir(exePath)

	// Open commands file, create if does not exist
	conf, err := ConfFromFile(filepath.Join(wDir, ConfName))
	if err != nil {
		_, err = CreateConf(filepath.Join(wDir, ConfName), DefaultConf())
		if err != nil {
			log.Fatalf("Failed to create a new commands file: %s", err)
		}
		log.Printf("Created a new commands file")
		os.Exit(0)
	}

	pot := NewPot(conf)

	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "BREW" || r.Method == "POST" {
			// Brew some coffee
			if r.Header.Get("Content-Type") != "application/coffee-pot-command" {
				http.Error(w, "Coffee content type is not set", http.StatusBadRequest)
				return
			}

			if r.Header.Get("Accept-Additions") != "" {
				// Additions were specified!
				http.Error(w, "Additions are not supported", http.StatusNotAcceptable)
				return
			}

			err := pot.Brew()
			if err != nil {
				log.Printf("Failed to BREW: %s", err)
				http.Error(w, "Brewing error", http.StatusInternalServerError)
			}
		} else if r.Method == "GET" {
			// Return Pot information
			w.Header().Add("Additions-List", "milk")
			w.Header().Add("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(&pot)
			if err != nil {
				log.Printf("Failed to answer a GET: %s", err)
				http.Error(w, "JSON encoding failed", http.StatusInternalServerError)
			}
		} else if r.Method == "PROPFIND" {
			// Write what king of coffee we're making
			w.Header().Add("Content-Type", "text/plain")
			w.Write([]byte(pot.CoffeeType))
		} else if r.Method == "WHEN" {
			if r.Header.Get("Content-Type") != "application/coffee-pot-command" {
				http.Error(w, "Coffee content type is not set", http.StatusBadRequest)
				return
			}

			err := pot.StopPouring()
			if err != nil {
				log.Printf("Failed to stop pouring milk: %s\n", err)
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: handler,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Fatal server error: %s!", err)
	}
}
