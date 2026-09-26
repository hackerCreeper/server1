package main

import (
	"fmt"
	"net/http"
)

const serverAddress = ":8080"

func main() {
	terminal := NewTerminal()
	if err := terminal.LoadPasswordHash("password.txt"); err != nil {
		fmt.Printf("Initialization warning: %v\n", err)
	}

	serverHandler := NewServerHandler(terminal)

	http.HandleFunc("/", serverHandler.ServeHome)
	http.HandleFunc("/api/cmd", serverHandler.HandleCommand)

	fmt.Printf("Secure hacking terminal listening on http://localhost%s\n", serverAddress)
	fmt.Println("Default emergency credential hint: *********")

	if err := http.ListenAndServe(serverAddress, nil); err != nil {
		fmt.Printf("Critical server failure: %v\n", err)
	}
}
