package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var wg sync.WaitGroup

func main() {
	// Handle Ctrl+C to gracefully exit the program
	handleCtrlC()

	// Start the ScanKeyboard goroutine
	wg.Add(1)
	go ScanKeyboard()

	// Wait for the ScanKeyboard goroutine to finish (when Ctrl+C is pressed)
	wg.Wait()

	fmt.Println("Program exited.")
}

func ScanKeyboard() {
	defer wg.Done()

	fmt.Println("ScanKeyboard started. Press Ctrl+C to exit.")

	for {
		// Replace the following line with the actual code to scan for keyboard input
		// For simplicity, we use fmt.Scanf to simulate keyboard input
		var key string
		fmt.Scanf("%s", &key)

		// Call the ProcessKey function for each keystroke
		ProcessKey(key)
	}
}

func ProcessKey(key string) {
	// Replace this function with your actual key processing logic
	fmt.Printf("%s", key)
}

func handleCtrlC() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nCtrl+C pressed. Exiting...")
		// Signal the ScanKeyboard goroutine to exit
		wg.Done()
		os.Exit(0)
	}()
}
