package directus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coder/websocket"
)

type WSMessage struct {
	Type         string                 `json:"type"`
	Collection   string                 `json:"collection,omitempty"`
	Status       string                 `json:"status,omitempty"`
	RefreshToken string                 `json:"refresh_token,omitempty"`
	AccessToken  string                 `json:"access_token,omitempty"`
	Error        any                    `json:"error,omitempty"`
	Event        string                 `json:"event,omitempty"`
	Query        map[string]interface{} `json:"query,omitempty"`
	UID          string                 `json:"uid,omitempty"`
	Data         interface{}            `json:"data,omitempty"`
}

func Watch(changes chan any) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	token := os.Getenv("DIRECTUS_TOKEN")
	if token == "" {
		return errors.New("DIRECTUS_TOKEN environment variable is not set")
	}

	wsURL := "ws://localhost:8055/websocket"

	// Establish WebSocket connection
	conn, _, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{})
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "closing")

	// Authenticate by sending the access token
	authMsg := WSMessage{
		Type:        "auth",
		AccessToken: token,
	}
	log.Printf("Sending auth message")
	if err := conn.Write(ctx, websocket.MessageText, mustJSON(authMsg)); err != nil {
		return fmt.Errorf("failed to send auth message: %w", err)
	}

	_, data, err := conn.Read(ctx)
	if err != nil {
		return fmt.Errorf("failed to read auth response: %w", err)
	}

	var msg WSMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return fmt.Errorf("failed to parse auth response: %w", err)
	}
	log.Printf("Auth response: %v", msg)

	if msg.Status != "ok" {
		return fmt.Errorf("authentication failed: %v", msg.Error)
	}

	time.Sleep(100 * time.Millisecond)

	// Subscribe to the "messages" collection
	subscribeMsg := WSMessage{
		Type:       "subscribe",
		Collection: collectionName,
	}
	log.Printf("Sending subscribe message: %v", subscribeMsg)
	if err := conn.Write(ctx, websocket.MessageText, mustJSON(subscribeMsg)); err != nil {
		return fmt.Errorf("failed to send subscribe message: %w", err)
	}

	// listening for init response
	_, data, err = conn.Read(ctx)
	if err != nil {
		return fmt.Errorf("failed to read subscribe response: %w", err)
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		return fmt.Errorf("failed to parse subscribe response: %w", err)
	}
	log.Printf("Subscribe response: %v", msg)

	// Listen for incoming messages
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		var msg WSMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		switch msg.Type {
		case "subscription":
			fmt.Printf("Event: %s\nData: %v\n", msg.Event, msg.Data)
			// write to the channel
			changes <- msg.Data
		case "ping":
			pongMsg := WSMessage{
				Type: "pong",
			}
			if err := conn.Write(ctx, websocket.MessageText, mustJSON(pongMsg)); err != nil {
				log.Printf("Failed to send pong message: %v", err)
			}
			log.Println("Received ping, sent pong")
		default:
			log.Printf("Received message: %+v\n", msg)
		}
	}
}

func mustJSON(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("JSON marshal error: %v", err)
	}
	return data
}
