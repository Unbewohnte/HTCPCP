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
	"io"
	"os"
)

type Commands struct {
	BrewCommand        string
	StopPouringCommand string
}

type Conf struct {
	Commands       Commands `json:"commands"`
	CoffeeType     string   `json:"coffee-type"`
	BrewTimeSec    uint     `json:"brew-time-sec"`
	MaxPourTimeSec uint     `json:"max-pour-time-sec"`
}

func DefaultConf() Conf {
	return Conf{
		Commands: Commands{
			BrewCommand:        "./brew.sh",
			StopPouringCommand: "./stopPouring.sh",
		},
		CoffeeType:     "Latte",
		BrewTimeSec:    10,
		MaxPourTimeSec: 5,
	}
}

// Tries to retrieve configuration structure from given json file
func ConfFromFile(path string) (*Conf, error) {
	confFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer confFile.Close()

	confBytes, err := io.ReadAll(confFile)
	if err != nil {
		return nil, err
	}

	var conf *Conf
	err = json.Unmarshal(confBytes, &conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

// Create a new configuration file
func CreateConf(path string, conf Conf) (*Conf, error) {
	confFile, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer confFile.Close()

	confJsonBytes, err := json.MarshalIndent(&conf, "", " ")
	if err != nil {
		return nil, err
	}

	_, err = confFile.Write(confJsonBytes)
	if err != nil {
		return nil, nil
	}

	return &conf, nil
}
