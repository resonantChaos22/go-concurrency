package main

import (
	"sync"
	"testing"
	"time"
)

func Test_dine(t *testing.T) {
	eatTime = 0 * time.Second
	sleepTime = 0 * time.Second
	thinkTime = 0 * time.Second
	// eatChan = make(chan string)

	for i := 0; i < 5; i++ {
		eatChan = make(chan string)
		go func() {

			dine()
		}()

		eatingOrder := []string{}
		var mu sync.RWMutex

		for phil := range eatChan {
			mu.Lock()
			eatingOrder = append(eatingOrder, phil)
			mu.Unlock()
		}

		if len(eatingOrder) != 5 {
			t.Error("length less than 5")
		}
	}
}

func Test_dineWithVaryingDelays(t *testing.T) {
	tests := []struct {
		name  string
		delay time.Duration
	}{
		{name: "zero delay", delay: time.Second * 0},
		{name: "0.25 delay", delay: time.Millisecond * 250},
		{name: "0.5 delay", delay: time.Millisecond * 500},
	}

	for _, e := range tests {
		eatTime = e.delay
		sleepTime = e.delay
		thinkTime = e.delay
		eatChan = make(chan string)
		go func() {
			dine()
		}()

		eatingOrder := []string{}
		var mu sync.RWMutex

		for phil := range eatChan {
			mu.Lock()
			eatingOrder = append(eatingOrder, phil)
			mu.Unlock()
		}

		if len(eatingOrder) != 5 {
			t.Error("length less than 5")
		}
	}
}
