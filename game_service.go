package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type GameService struct {
	players map[string]*PlayerSession
	mu      sync.Mutex
}

func NewGameService() *GameService {
	rand.Seed(time.Now().UnixNano())
	return &GameService{
		players: make(map[string]*PlayerSession),
	}
}

func (s *GameService) GetOrCreateSession(clientID string) *PlayerSession {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.players[clientID]
	if !exists {
		session = &PlayerSession{
			ClientID: clientID,
			Balance:  100.0, // Saldo inicial
			State:    StateIdle,
		}
		s.players[clientID] = session
	}
	return session
}

func (s *GameService) GetBalance(session *PlayerSession) map[string]interface{} {
	session.Lock()
	defer session.Unlock()

	return map[string]interface{}{
		"balance": session.Balance,
		"state":   session.State,
	}
}

func (s *GameService) Play(session *PlayerSession, payload PlayPayload) (map[string]interface{}, error) {
	session.Lock()
	defer session.Unlock()

	if session.State == StateInProgress {
		return nil, fmt.Errorf("tem de finalizar a jogada anterior antes de jogar novamente")
	}

	if payload.BetAmount <= 0 {
		return nil, fmt.Errorf("o valor da aposta deve ser maior que 0")
	}

	if payload.BetAmount > session.Balance {
		return nil, fmt.Errorf("saldo insuficiente")
	}

	// Debitar saldo e alterar estado
	session.Balance -= payload.BetAmount
	session.State = StateInProgress

	// Sorteio
	drawnNumber := rand.Intn(6) + 1
	isEven := drawnNumber%2 == 0

	win := (isEven && payload.BetType == Even) || (!isEven && payload.BetType == Odd)

	if win {
		session.PendingWin = payload.BetAmount * 2
	} else {
		session.PendingWin = 0
	}

	resultStr := "LOSE"
	if win {
		resultStr = "WIN"
	}

	return map[string]interface{}{
		"drawn_number":    drawnNumber,
		"result":          resultStr,
		"current_balance": session.Balance,
		"potential_win":   session.PendingWin,
	}, nil
}

func (s *GameService) EndPlay(session *PlayerSession) (map[string]interface{}, error) {
	session.Lock()
	defer session.Unlock()

	if session.State != StateInProgress {
		return nil, fmt.Errorf("não existe nenhuma jogada em curso para finalizar")
	}

	credited := session.PendingWin
	session.Balance += credited
	session.PendingWin = 0
	session.State = StateIdle

	return map[string]interface{}{
		"credited_amount": credited,
		"new_balance":     session.Balance,
	}, nil
}