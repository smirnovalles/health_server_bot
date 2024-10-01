package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Config struct {
	BotToken       string
	ChatID         string
	HealthURL      string
	Message        string
	ExpectedStatus string
}

type HealthResponse struct {
	Status string `json:"status"`
}

func loadConfig() Config {
	botToken := flag.String("bot_token", "000", "Telegram Bot Token")
	chatID := flag.String("chat_id", "000", "Telegram Chat ID")
	healthURL := flag.String("health_url", "http://0.0.0.0:8080/health", "Health Check URL")
	message := flag.String("message", "⚠️ *Внимание!* Сервер не отвечает должным образом.", "Telegram Message")
	expectedStatus := flag.String("expected_status", "OK", "Expected status from Health Check")
	flag.Parse()
	return Config{
		BotToken:       *botToken,
		ChatID:         *chatID,
		HealthURL:      *healthURL,
		Message:        *message,
		ExpectedStatus: *expectedStatus,
	}
}

func sendTelegramMessage(config Config, message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.BotToken)
	payload := map[string]interface{}{
		"chat_id":    config.ChatID,
		"text":       message,
		"parse_mode": "Markdown",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	client := &http.Client{
		timeout: 10 * time.Second,
	}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body), url, config.BotToken)
	return nil
}

func checkServer(config Config) error {
	client := &http.Client{
		timeout: 10 * time.Second,
	}
	resp, err := client.Get(config.HealthURL)
	if err != nil {
		return fmt.Errorf("ошибка при обращении к серверу: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("получен некорректный HTTP статус: %s", resp.Status)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка при чтении ответа сервера: %v", err)
	}
	var healthResp HealthResponse
	err = json.Unmarshal(bodyBytes, &healthResp)
	if err != nil {
		return fmt.Errorf("ошибка при разборе JSON ответа: %v", err)
	}
	if healthResp.Status != config.ExpectedStatus {
		return fmt.Errorf("статус здоровья не %s: %s", config.ExpectedStatus, healthResp.Status)
	}
	return nil
}

func main() {
	config := loadConfig()
	err := checkServer(config)
	if err != nil {
		errMsg := fmt.Sprintf("%sn*Ошибка:* %s", config.Message, err.Error())
		sendErr := sendTelegramMessage(config, errMsg)
		if sendErr != nil {
			fmt.Fprintf(os.Stderr, "Не удалось отправить сообщение в Telegram: %vn", sendErr)
		}
		os.Exit(1)
	} else {
		fmt.Println("✅ Сервер работает нормально.")
	}
}
