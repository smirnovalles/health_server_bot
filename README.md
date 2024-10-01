# Server Monitoring with Telegram Notifications

This tool is written in Go and is designed to monitor the status of a server. If the server does not respond properly, a notification is sent to Telegram.

## Table of Contents

- [Requirements](#requirements)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Example](#example)
- [License](#license)

## Requirements

- [Go](https://golang.org/dl/) version 1.16 or higher

## Installation

1. *Clone the repository:*

   
bash
   git clone https://github.com/your-username/monitoring-telegram.git
   cd monitoring-telegram
   

2. *Build the application:*

   
bash
   go build -o monitor main.go
   

## Configuration

The application accepts the following parameters via the command line:

- `-bot_token` : Your Telegram Bot Token.
- `-chat_id` : Telegram Chat ID where notifications will be sent.
- `-health_url` : URL to check the server's health status.
- `-message` : Message text to be sent in case of an error.
- `-expected_status` : Expected status response from the health check.

*Note:* It is important to replace the default values for `bot_token` and `chat_id` with your own.

## Usage

After building the application, you can run it with the necessary parameters. Below is an example of how to use it:

bash
./monitor 
  -bot_token "YOUR_TELEGRAM_BOT_TOKEN" 
  -chat_id "YOUR_CHAT_ID" 
  -health_url "http://yourserver.com/health" 
  -message "⚠️ Attention! The server is not responding properly." 
  -expected_status "OK"

### Parameter Descriptions:

- `-bot_token`: The token of your Telegram bot, which can be obtained through [BotFather](https://t.me/BotFather).
- `-chat_id`: The ID of the chat or user where notifications will be sent. You can obtain the `chat_id` in various ways, such as sending a message to the bot and using the API to retrieve updates.
- `-health_url`: The URL endpoint to check the status of your server.
- `-message`: The message that will be sent in case of a failure.
- `-expected_status`: The expected status in the response from the `health_url` endpoint.

## Example

Suppose you have a server with a health check endpoint at `http://example.com/health`, and you want to receive notifications in a chat with ID `123456789`. Your Telegram bot has the token `123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11`.

Running the application would look like this:

bash
./monitor 
  -bot_token "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11" 
  -chat_id "123456789" 
  -health_url "http://example.com/health" 
  -message "⚠️ Attention! The server is not responding properly." 
  -expected_status "OK"

If the endpoint `http://example.com/health` does not return the status `OK` or is unavailable, a warning message will be sent to the specified Telegram chat.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.