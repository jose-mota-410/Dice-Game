package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	gameService := NewGameService()
	wsHandler := NewWSHandler(gameService)

	http.Handle("/ws", wsHandler)

	fmt.Println("Servidor iniciado na porta :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}