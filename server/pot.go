/*
 The MIT License (MIT)

Copyright © 2024 Kasianov Nikolai Alekseevich (Unbewohnte)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"errors"
	"log"
	"os/exec"
	"time"
)

const (
	PotStatusErr     = "Error"
	PotStatusReady   = "Ready"
	PotStatusBrewing = "Brewing"
	PotStatusPouring = "Pouring"
)

type Pot struct {
	State          string `json:"state"`
	CoffeeType     string `json:"coffeeType"`
	BrewTimeSec    uint   `json:"brew-time-sec"`
	MaxPourTimeSec uint   `json:"max-pour-time-sec"`
	commands       Commands
}

func NewPot(conf *Conf) *Pot {
	return &Pot{
		State:          PotStatusReady,
		CoffeeType:     conf.CoffeeType,
		BrewTimeSec:    conf.BrewTimeSec,
		MaxPourTimeSec: conf.MaxPourTimeSec,
		commands:       conf.Commands,
	}
}

func run(command string, args ...string) error {
	output, err := exec.Command(command, args...).Output()
	if err != nil {
		return err
	}

	log.Printf("%s", string(output))
	return nil
}

// Brew some coffee!
func (p *Pot) Brew(args ...string) error {
	if p.State != PotStatusReady {
		return errors.New("pot is not yet ready")
	}

	p.State = PotStatusBrewing
	go func() {
		// Start pouring after brewing is done
		time.Sleep(time.Second * time.Duration(p.BrewTimeSec))
		p.State = PotStatusPouring
		log.Print("Pouring!")

		go func() {
			// Stop pouring...
			time.Sleep(time.Second * time.Duration(p.MaxPourTimeSec))
			if p.State == PotStatusPouring {
				// ...if it was not stopped earlier
				p.State = PotStatusReady
				log.Print("Poured at maximum capacity!")
			}
		}()
	}()
	return run(p.commands.BrewCommand, args...)
}

// No more pouring!
func (p *Pot) StopPouring() error {
	if p.State != PotStatusPouring {
		return errors.New("coffee is not brewed yet")
	}

	p.State = PotStatusReady
	return run(p.commands.StopPouringCommand)
}
