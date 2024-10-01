package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	botToken := flag.String("bot_token", "", "Telegram Bot Token")
	chatID := flag.String("chat_id", "", "Telegram Chat ID")
	healthURL := flag.String("health_url", "", "Health Check URL")
	message := flag.String("message", "⚠️ *Attention!* The server is not responding properly.", "Error message to send")
	expectedStatus := flag.String("expected_status", "OK", "Expected status response")
	timeout := flag.Int("timeout", 5, "HTTP client timeout in seconds")

	flag.Parse()

	client := &http.Client{
		Timeout: time.Duration(*timeout) * time.Second,
	}

	resp, err := client.Get(*healthURL)
	if err != nil {
		sendTelegramMessage(*botToken, *chatID, *message)
		log.Fatalf("Error fetching health URL: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		sendTelegramMessage(*botToken, *chatID, *message)
		log.Fatalf("Error reading response body: %v", err)
	}

	status := strings.TrimSpace(string(body))
	if status != *expectedStatus {
		sendTelegramMessage(*botToken, *chatID, *message)
		log.Fatalf("Unexpected status: got %s, expected %s", status, *expectedStatus)
	}

	fmt.Println("Server is healthy.")
}

func sendTelegramMessage(botToken, chatID, message string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	payload := fmt.Sprintf(`{"chat_id":"%s","text":"%s","parse_mode":"Markdown"}`, chatID, message)

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		log.Printf("Error creating Telegram request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending Telegram message: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Non-OK response from Telegram: %s", string(body))
	}
}
