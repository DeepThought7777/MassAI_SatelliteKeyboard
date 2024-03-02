package main

import (
	"fmt"
	"github.com/DeepThought7777/MassAI_SatelliteKeyboard/codebase"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/eiannone/keyboard"
)

func main() {
	// create your identity or read it in
	identity, created, err := codebase.GetOrCreateGUID()
	if err != nil {
		fmt.Errorf("Cannot create identity: %s", err.Error())
	}

	baseURL := "http://127.0.0.1:8080/v1"

	if created {
		// try to register the satellite
		requestURL := codebase.BuildRegisterURL(baseURL, identity, "1", "1")
		fmt.Printf("REGISTERING: %s\n", requestURL)
		if codebase.SendRequest(requestURL) != codebase.InfoEntityRegistered {
			fmt.Printf("Registration did not succeed")
		}
	}

	// try to connect the satellite
	requestURL := codebase.BuildConnectURL(baseURL, identity)
	fmt.Printf("CONNECTING: %s\n", requestURL)
	if codebase.SendRequest(requestURL) != codebase.InfoEntityConnected {
		fmt.Printf("Connection did not succeed")
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	keyCh := make(chan rune)

	wg.Add(1)
	go scanKeys(signalCh)
	go processKeys(baseURL, identity, keyCh, &wg)

	err = keyboard.Open()
	if err != nil {
		return
	}

	defer func() {
		err := keyboard.Close()
		if err != nil {
			return
		}
		wg.Wait() // Wait for the goroutines to finish
	}()

	fmt.Println("Press 'ESC' to exit.")

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			break
		}
		if key == keyboard.KeyEsc {
			break
		}

		keyCh <- char
	}

	// try to disconnect the satellite
	requestURL = codebase.BuildDisconnectURL(baseURL, identity)
	fmt.Printf("DISCONNECTING: %s\n", requestURL)
	if codebase.SendRequest(requestURL) != codebase.InfoEntityConnected {
		fmt.Printf("Disconnection did not succeed")
	}

	wg.Wait()
}

func scanKeys(signalCh chan os.Signal) {
	<-signalCh
	err := keyboard.Close()
	if err != nil {
		log.Println("Error closing keyboard:", err)
	}
}

func processKeys(baseURL, identity string, keyCh <-chan rune, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case char := <-keyCh:
			requestURL := codebase.BuildSendInputsURL(baseURL, identity, "1", string(char))
			fmt.Printf("SENDING INPUTS: %s\n", requestURL)
			if codebase.SendRequest(requestURL) != codebase.InfoInputDataSent {
				fmt.Printf("Sending did not succeed")
			}
		}
	}
}
