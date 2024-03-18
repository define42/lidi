package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", ":7000")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server started on :8080")

	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		wg.Add(1)
		go handleConnection(conn, &wg)
	}

	wg.Wait() // In a real server, you might handle shutdowns more gracefully.
}

func handleConnection(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()

	fmt.Printf("Client connected: %s\n", conn.RemoteAddr())

	const interval = 5 * time.Second
	timer := time.NewTimer(interval)
	defer timer.Stop()

	startTime := time.Now()
	var totalBytes int64

	buf := make([]byte, 40960) // Adjust buffer size to your needs

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading from %s: %v\n", conn.RemoteAddr(), err)
			} else {
				fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr())
			}
			break
		}
		totalBytes += int64(n)

		select {
		case <-timer.C:
			elapsed := time.Since(startTime)
			speed := float64(totalBytes) / elapsed.Seconds() / 1024 // Speed in KiB/s
			fmt.Printf("Data received from %s: %d bytes, Speed: %.2f KiB/s\n", conn.RemoteAddr(), totalBytes, speed)
			// Reset timer and counters
			timer.Reset(interval)
			startTime = time.Now()
			totalBytes = 0
		default:
			// Continue reading without blocking
		}
	}
}
