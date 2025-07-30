package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var ErrTransient = errors.New("a transient error occurred")
var ErrPermanent = errors.New("a permanent error occurred")

func main() {
	fmt.Println(Retry(UnreliableAPICall, 4, time.Duration(100)))

}

func UnreliableAPICall() (string, error) {
	// Simulate random failures
	r := rand.Intn(10)

	if r < 2 { // 20% chance of permanent error
		fmt.Println("API call failed with a permanent error.")
		return "", ErrPermanent
	} else if r < 8 { // 40% chance of transient error
		fmt.Println("API call failed with a transient error.")
		return "", ErrTransient
	}

	fmt.Println("API call succeeded!")
	return "Success!", nil
}

func Retry(fn func() (string, error), retries int, delay time.Duration) (string, error) {
	resp, err := fn()
	fmt.Println(retries)
	if err != nil && retries > 0 && errors.Is(err, ErrTransient) {
		fmt.Printf("Retrying after %v delay...\n", delay)
		time.Sleep(delay)
		return Retry(fn, retries-1, delay*2)
	}
	return resp, nil
}
