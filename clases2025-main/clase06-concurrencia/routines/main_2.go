package main

import (
    "fmt"
    "time"
    "math/rand"
)

func main() {

	// simple go routine
	go processUsers()
}

func processUsers() {

	i := 0

	for true {

		processUser(i)

		time.Sleep(time.Duration(200 + rand.Intn(300)) * time.Millisecond)

		i++
	}
}

func processUser(userID int) {

    fmt.Println("I'm processing user", userID)
}