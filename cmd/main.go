package main

import (
	"log"

	"github.com/ashkenazi1/go-telegram-bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	statefulConfig := telegram.BotConfig{
		Token:     "YOUR_BOT_TOKEN",
		OwnerID:   1234567890, // Your user ID
		ChannelID: 123456789,
		UseState:  true,
	}
	statefulBot, err := telegram.NewBot(statefulConfig)
	if err != nil {
		log.Panic(err)
	}

	statefulBot.RegisterState("start", handleStartState)
	statefulBot.RegisterState("greeting", handleGreetingState)
	statefulBot.RegisterState("farewell", handleFarewellState)

	go statefulBot.Start()

	statelessConfig := telegram.BotConfig{
		Token:     "YOUR_BOT_TOKEN",
		ChannelID: 123456789, // Your channel ID
		OwnerID:   0,
		UseState:  false,
	}
	statelessBot, err := telegram.NewBot(statelessConfig)
	if err != nil {
		log.Panic(err)
	}

	err = statelessBot.SendMessage(0, "This is a test log message.")
	if err != nil {
		log.Println("Error sending log message:", err)
	}
}

func handleStartState(bot *telegram.Bot, message *tgbotapi.Message, userState *telegram.UserState) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Hello! How can I assist you? (type 'greet' or 'bye')")
	bot.API.Send(msg)
	bot.StateManager.SetState(message.Chat.ID, "greeting")
}

func handleGreetingState(bot *telegram.Bot, message *tgbotapi.Message, userState *telegram.UserState) {
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

func handleFarewellState(bot *telegram.Bot, message *tgbotapi.Message, userState *telegram.UserState) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "You have already said goodbye. Type 'start' to begin again.")
	bot.API.Send(msg)
	if message.Text == "start" {
		bot.StateManager.SetState(message.Chat.ID, "start")
	}
}
