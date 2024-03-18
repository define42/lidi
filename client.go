package main

import (
	"fmt"
//	"math/rand"
	"crypto/rand"
	"net"
	"time"
)

func main() {

	time.Sleep(5 * time.Second) // Sleeps for 2 seconds

randomBytes := make([]byte, 100)

	// Fill the slice with random data
	if _, err := rand.Read(randomBytes); err != nil {
		// Handle the error here
		fmt.Println("Error generating random bytes:", err)
		return
	}

	// Seed the random number generator
//	rand.Seed(time.Now().UnixNano())

	// Specify the server's address and port
	serverAddress := "172.16.0.2:5000"

	// Establish a connection to the server
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		return
	}
	defer conn.Close()

	for {
		// Generate random data (in this case, a random integer)
//		randomData := fmt.Sprintf("%d", rand.Intn(100))

		// Send the random data to the server
		_, err = conn.Write(randomBytes)
		if err != nil {
			fmt.Println("Error sending data to the server:", err)
			return
		}

		// Optional: Print the data sent to the server (you might want to remove this in real use to avoid flooding the output)
		// fmt.Printf("Sent random data to the server: %s\n", randomData)

		// Note: In a real application, consider including a way to break out of this loop,
		// such as by checking for a specific condition or listening for a termination signal.
	}
}
