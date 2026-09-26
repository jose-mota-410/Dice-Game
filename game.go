package main

import (
	"encoding/json"
	"sync"
)

type BetType string

const (
	Even BetType = "EVEN"
	Odd  BetType = "ODD"
)

type GameState string

const (
	StateIdle       GameState = "IDLE"
	StateInProgress GameState = "IN_PROGRESS"
)

type PlayerSession struct {
	sync.Mutex
	ClientID   string    `json:"client_id"`
	Balance    float64   `json:"balance"`
	State      GameState `json:"state"`
	PendingWin float64   `json:"pending_win"`
}

type Request struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type PlayPayload struct {
	BetAmount float64 `json:"bet_amount"`
	BetType   BetType `json:"bet_type"`
}

type WSResponse struct {
	Status  string      `json:"status"`
	Action  string      `json:"action,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}