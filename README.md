# Dice Game

A backend service for an Even/Odd dice betting game implemented in Go using WebSockets.

---

## Overview

The backend allows players to bet on whether the next rolled dice number (1–6) will be **EVEN** or **ODD**. 
- Winning doubles the bet amount.
- Losing forfeits the bet.
- The game maintains state per player and prevents duplicate active rounds or negative balance exploits.

---

## Features & Safeguards

- **Persistent In-Memory State**: Tracks player balance, current state (`IDLE` vs `IN_PROGRESS`), and pending payouts.
- **Overlapping Play Protection**: A player cannot initiate a new bet while a round is already in progress without calling `end_play`.
- **Balance Validation**: Bets cannot be zero, negative, or greater than the player's available balance.
- **Thread Safety**: Concurrent read/write protection using mutex locks on player sessions and global registries.

---

## Requirements

- [Go](https://go.dev/) (version 1.20 or newer)
- [Postman Desktop App](https://www.postman.com/downloads) (or any WebSocket client)

---

## Getting Started

### 1. Install Dependencies
Run in your terminal:
```bash
go mod tidy
```

### 2. Run the Server
```bash
go run .
```

---

## WebSocket API Specification

### Connection Endpoint
Connect via WebSocket and provide a `clientid` query parameter:
```plaintext
ws://127.0.0.1:8080/ws?clientid=<YOUR_CLIENT_ID>
```
*Example:* `ws://127.0.0.1:8080/ws?clientid=player1`

New players automatically receive a starting balance of 100.0.

---

### Actions

#### 1. Retrieve Wallet Balance (`wallet`)
Fetches the current balance and state of the player session.

* **Request:**
  ```json
  {
    "action": "wallet"
  }
  ```

---

#### 2. Place a Bet (`play`)
Places a bet on the next dice roll. Changes session state to `IN_PROGRESS`.

* **Request:**
  ```json
  {
    "action": "play",
    "payload": {
      "bet_amount": 20,
      "bet_type": "EVEN"
    }
  }
  ```
  *(Supported `bet_type` values: `"EVEN"`, `"ODD"`)*

---

#### 3. Settle and Close Round (`end_play`)
Credits any winnings to the player's balance and resets state to `IDLE`.

* **Request:**
  ```json
  {
    "action": "end_play"
  }
  ```
---

## Error Handling

If an operation fails, the backend returns a standardized JSON error message:
```json
{
  "status": "error",
  "message": "<Error description>"
}
```

Common error cases:
- Starting a play while a round is still open: `"tem de finalizar a jogada anterior antes de jogar novamente"`
- Bet amount exceeding balance: `"saldo insuficiente"`
- Ending a play when no round is active: `"não existe nenhuma jogada em curso para finalizar"`
- Invalid or non-JSON input: `"Formato de JSON inválido"`

---

## Testing with Postman

1. Open Postman Desktop and add a **WebSocket** Request.
2. Enter the URL: `ws://127.0.0.1:8080/ws?clientid=player1` and click **Connect**.
3. Under the message box, switch the format from **Text** to **JSON**.
4. Test the actions in sequence as explained above.
