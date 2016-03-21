package main

import (
	"fmt"
	"os"
)

func main() {
	client, err := NewClient("0.0.0.0:1234")

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	defer client.Close()

	client.Set("a", "1")

	a, _ := client.Get("a")
	a, _ = client.Incr("a", 2)

	fmt.Printf("a = %s", a)
}
