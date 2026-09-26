package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WSHandler struct {
	gameService *GameService
	upgrader    websocket.Upgrader
}

func NewWSHandler(gs *GameService) *WSHandler {
	return &WSHandler{
		gameService: gs,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (h *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("clientid")
	if clientID == "" {
		http.Error(w, "clientid é obrigatório", http.StatusBadRequest)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro no Upgrade WS:", err)
		return
	}
	defer conn.Close()

	session := h.gameService.GetOrCreateSession(clientID)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var req Request
		if err := json.Unmarshal(message, &req); err != nil {
			conn.WriteJSON(WSResponse{Status: "error", Message: "Formato de JSON inválido"})
			continue
		}

		switch req.Action {
		case "wallet":
			data := h.gameService.GetBalance(session)
			conn.WriteJSON(WSResponse{Status: "success", Action: "wallet", Data: data})

		case "play":
			var payload PlayPayload
			if err := json.Unmarshal(req.Payload, &payload); err != nil {
				conn.WriteJSON(WSResponse{Status: "error", Message: "Payload de aposta inválido"})
				continue
			}

			data, err := h.gameService.Play(session, payload)
			if err != nil {
				conn.WriteJSON(WSResponse{Status: "error", Message: err.Error()})
				continue
			}
			conn.WriteJSON(WSResponse{Status: "success", Action: "play", Data: data})

		case "end_play":
			data, err := h.gameService.EndPlay(session)
			if err != nil {
				conn.WriteJSON(WSResponse{Status: "error", Message: err.Error()})
				continue
			}
			conn.WriteJSON(WSResponse{Status: "success", Action: "end_play", Data: data})

		default:
			conn.WriteJSON(WSResponse{Status: "error", Message: "Ação desconhecida"})
		}
	}
}