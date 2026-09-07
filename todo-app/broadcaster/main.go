package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nats-io/nats.go"
)

type config struct {
	natsURL        string
	natsSubject    string
	tbotAPIToken   string
	personalChatID string
}

func loadConfig() (*config, error) {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	cfg := &config{
		natsURL:        natsURL,
		natsSubject:    os.Getenv("NATS_SUBJECT"),
		tbotAPIToken:   os.Getenv("TBOT_API_TOKEN"),
		personalChatID: os.Getenv("PERSONAL_CHAT_ID"),
	}

	if cfg.natsSubject == "" || cfg.tbotAPIToken == "" ||
		cfg.personalChatID == "" {
		return nil, fmt.Errorf(
			`missing required environment variables
             (NATS_SUBJECT, TBOT_API_TOKEN, PERSONAL_CHAT_ID)`,
		)
	}

	return cfg, nil
}

type BroadcastMsg struct {
	Id    int    `json:"id"`
	Msg   string `json:"msg"`
	Title string `json:"title"`
}

type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	nc, err := nats.Connect(cfg.natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Drain()

	nc.Subscribe(cfg.natsSubject, func(msg *nats.Msg) {
		var bMsg BroadcastMsg
		if err := json.Unmarshal(msg.Data, &bMsg); err != nil {
			log.Printf("Failed to unmarshal broadcast message: %v", err)
			return
		}

		sendToTelegram(cfg.tbotAPIToken, cfg.personalChatID, bMsg)
	})

	log.Printf("Listening for messages on NATS subject: %s", cfg.natsSubject)

	select {}
}

func sendToTelegram(token, chatID string, bMsg BroadcastMsg) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	// Format the message nicely using Markdown
	text := fmt.Sprintf("*%s* (ID: %d)\n\n*%s*", bMsg.Title, bMsg.Id, bMsg.Msg)

	payload := TelegramMessage{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "Markdown",
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal telegram payload: %v", err)
		return
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Printf("HTTP request to Telegram failed: %v", err)
		return
	}
	defer resp.Body.Close()
}
