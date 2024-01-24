/*
 The MIT License (MIT)

Copyright © 2024 Kasianov Nikolai Alekseevich (Unbewohnte)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var version *bool = flag.Bool("version", false, "Print version information")

const (
	CommandGet      = "GET"
	CommandWhen     = "WHEN"
	CommandBrew     = "BREW"
	CommandPropfind = "PROPFIND"
)

const Version string = "0.1-client"

func makeRequest(address string, method string) {
	request, err := http.NewRequest(method, address, nil)
	if err != nil {
		log.Fatalf("Failed to form a request %s", err)
	}
	request.Header.Add("Content-Type", "application/coffee-pot-command")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatalf("Failed to %s: %s", method, err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Failed to read body: %s", err)
	}

	log.Printf("Status: %d; Headers: %+s\nData: %s", response.StatusCode, response.Header, body)
}

func main() {
	log.SetOutput(os.Stdout)
	log.Default().SetFlags(0)
	flag.Usage = func() {
		fmt.Printf(`HTCPCP-client (-version) [ADDR] [COMMAND]
  -version
        Print version information and exit
  ADDR string
    	Address of an HTCPCP server with port (ie: http://111.11.111.1:80 or http://coffeeserver:80)
  COMMAND string
    	Command to send (ie: GET, WHEN, BREW, PROPFIND)
`,
		)
	}
	flag.Parse()

	if *version {
		fmt.Printf("HTCPCP-client %s\n(C) 2024 Kasianov Nikolai Alekseevich (Unbewohnte)\n", Version)
		return
	}

	if len(os.Args) < 3 {
		log.Fatalf("Not enough arguments! Run with -help to see a help message")
	}

	address := os.Args[1]
	command := strings.ToUpper(os.Args[2])

	makeRequest(address, command)
}
