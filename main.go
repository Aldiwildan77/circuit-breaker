package main

import (
	"errors"
	"fmt"
	"time"
)

func main() {
	cb := NewCircuitBreaker(
		WithMaxFailures(2),
		WithResetTimeout(2*time.Second),
	)

	fmt.Printf("Initial Circuit is open: %t; Initial Failures: %d\n", !cb.IsAllowed(), cb.GetFailures())

	for i := 0; i < 13; i++ {
		fmt.Printf("Request %d\n", i+1)
		err := cb.Call(func() error {
			// Simulate failure in 50% of the cases
			if i%2 == 0 {
				return errors.New("simulated failure")
			}
			fmt.Println("Request succeeded")
			return nil
		})

		if err != nil {
			fmt.Printf("Request failed: %v\n", err)
		}

		fmt.Printf("Circuit is open: %t; Failures: %d\n\n", !cb.IsAllowed(), cb.GetFailures())

		time.Sleep(500 * time.Millisecond)
	}
}
