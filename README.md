# go-telegram-bot

MyTelegramBot is a flexible Go package for creating Telegram bots with optional state management. It supports both complex stateful bots and simple stateless logging bots, making it easy to integrate Telegram bot functionalities into your applications.

## Features

- **Stateful Mode**: Manage user states and interactions seamlessly.
- **Stateless Mode**: Quickly send updates or logs to a Telegram channel.
- **Configurable**: Customize bot behavior using a flexible configuration structure.
- **Easy to Use**: Simplified functions for common bot operations.

## Installation

To install the package, use:

```sh
go get -u github.com/yourusername/mytelegrambot
```

## Usage

Initial Setup

Create a configuration structure to initialize the bot:
```
package main

import (
    "log"
    "mytelegrambot"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
    // Example for stateful bot
    statefulConfig := mytelegrambot.BotConfig{
        Token:    "YOUR_TELEGRAM_BOT_TOKEN",
        DebugMode: true,
        OwnerID:  123456789, // Replace with your owner ID
        UseState: true,
    }
    statefulBot, err := mytelegrambot.NewBot(statefulConfig)
    if err != nil {
        log.Panic(err)
    }

    statefulBot.RegisterState("start", handleStartState)
    statefulBot.RegisterState("greeting", handleGreetingState)
    statefulBot.RegisterState("farewell", handleFarewellState)

    go statefulBot.Start()

    // Example for stateless logging bot
    statelessConfig := mytelegrambot.BotConfig{
        Token:     "YOUR_TELEGRAM_BOT_TOKEN",
        DebugMode: true,
        OwnerID:   123456789, // Replace with your owner ID
        ChannelID: -1001234567890, // Your channel ID
        UseState:  false,
    }
    statelessBot, err := mytelegrambot.NewBot(statelessConfig)
    if err != nil {
        log.Panic(err)
    }

    err = statelessBot.SendMessage(0, "This is a test log message.") // chatId is ignored because ChannelID is set
    if err != nil {
        log.Println("Error sending log message:", err)
    }
}

func handleStartState(bot *mytelegrambot.Bot, message *tgbotapi.Message, userState *mytelegrambot.UserState) {
    msg := tgbotapi.NewMessage(message.Chat.ID, "Hello! How can I assist you? (type 'greet' or 'bye')")
    bot.API.Send(msg)
    bot.StateManager.SetState(message.Chat.ID, "greeting")
}

func handleGreetingState(bot *mytelegrambot.Bot, message *tgbotapi.Message, userState *mytelegrambot.UserState) {
    switch message.Text {
    case "greet":
        msg := tgbotapi.NewMessage(message.Chat.ID, "Nice to meet you! How are you?")
        bot.API.Send(msg)
    case "bye":
        msg := tgbotapi.NewMessage(message.Chat.ID, "Goodbye! Have a great day!")
        bot.API.Send(msg)
        bot.StateManager.SetState(message.Chat.ID, "farewell")
    default:
        msg := tgbotapi.NewMessage(message.Chat.ID, "Please type 'greet' or 'bye'")
        bot.API.Send(msg)
    }
}

func handleFarewellState(bot *mytelegrambot.Bot, message *tgbotapi.Message, userState *mytelegrambot.UserState) {
    msg := tgbotapi.NewMessage(message.Chat.ID, "You have already said goodbye. Type 'start' to begin again.")
    bot.API.Send(msg)
    if message.Text == "start" {
        bot.StateManager.SetState(message.Chat.ID, "start")
    }
}
```

## Configuration

The BotConfig struct allows you to customize the bot’s behavior:

```
type BotConfig struct {
    Token     string
    DebugMode bool
    OwnerID   int64
    ChannelID int64 // Optional: Only needed for logging bot
    UseState  bool  // Flag to determine if state management is needed
}
```

# Sending Messages

Use the SendMessage method to send messages. If ChannelID is set in the configuration, it will be used instead of the provided chatId.

```
func (b *Bot) SendMessage(chatId int64, text string) error {
    if b.ChannelID != 0 {
        chatId = b.ChannelID
    }
    msg := tgbotapi.NewMessage(chatId, text)
    _, err := b.API.Send(msg)
    return err
}
```

# State Management

Register state handlers if you are using stateful mode:

```func (b *Bot) RegisterState(state string, handler StateHandler) {
    if b.useState {
        b.StateManager.RegisterState(state, handler)
    } else {
        log.Println("State management is not enabled for this bot.")
    }
}
```

# Example State Handlers

Implement state handlers to manage bot interactions:

```func handleStartState(bot *mytelegrambot.Bot, message *tgbotapi.Message, userState *mytelegrambot.UserState) {
    msg := tgbotapi.NewMessage(message.Chat.ID, "Hello! How can I assist you? (type 'greet' or 'bye')")
    bot.API.Send(msg)
    bot.StateManager.SetState(message.Chat.ID, "greeting")
}

func handleGreetingState(bot *mytelegrambot.Bot, message *tgbotapi.Message, userState *mytelegrambot.UserState) {
    switch message.Text {
    case "greet":
        msg := tgbotapi.NewMessage(message.Chat.ID, "Nice to meet you! How are you?")
        bot.API.Send(msg)
    case "bye":
        msg := tgbotapi.NewMessage(message.Chat.ID, "Goodbye! Have a great day!")
        bot.API.Send(msg)
        bot.StateManager.SetState(message.Chat.ID, "farewell")
    default:
        msg := tgbotapi.NewMessage(message.Chat.ID, "Please type 'greet' or 'bye'")
        bot.API.Send(msg)
    }
}

func handleFarewellState(bot *mytelegrambot.Bot, message *tgbotapi.Message, userState *mytelegrambot.UserState) {
    msg := tgbotapi.NewMessage(message.Chat.ID, "You have already said goodbye. Type 'start' to begin again.")
    bot.API.Send(msg)
    if message.Text == "start" {
        bot.StateManager.SetState(message.Chat.ID, "start")
    }
}
```

# License

This project is licensed under the MIT License - see the LICENSE file for details.
