package telegram

import (
	"log"
	"os"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotConfig struct {
	Token             string
	DebugMode         bool
	OwnerID           int64
	ChannelID         int64
	UseState          bool
	UnswerOnlyToOwner bool
}

type Bot struct {
	API          *tgbotapi.BotAPI
	StateManager *StateManager
	OwnerID      int64
	ChannelID    int64
	useState     bool
}

func NewBot(config BotConfig) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, err
	}
	bot.Debug = config.DebugMode

	var stateManager *StateManager
	if config.UseState {
		stateManager = NewStateManager()
	}

	var ownerID int64
	if config.UnswerOnlyToOwner && config.OwnerID != 0 {
		ownerID = config.OwnerID
	}

	return &Bot{
		API:          bot,
		StateManager: stateManager,
		OwnerID:      ownerID,
		ChannelID:    config.ChannelID,
		useState:     config.UseState,
	}, nil
}

func (b *Bot) Start() error {
	if b.useState {
		return b.startStateful()
	} else {
		log.Println("Bot started in stateless mode. Ready to send updates.")
	}
	return nil
}

func (b *Bot) startStateful() error {
	updates, err := b.API.GetUpdatesChan(tgbotapi.UpdateConfig{
		Offset:  0,
		Timeout: 60,
	})
	if err != nil {
		return err
	}
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if b.OwnerID != 0 && update.Message.Chat.ID != b.OwnerID {
			continue
		}

		userID := update.Message.Chat.ID
		userState := b.StateManager.GetState(userID)

		handler, exists := b.StateManager.GetHandler(userState.State)
		if exists {
			handler(b, update.Message, userState)
		} else {
			log.Printf("No handler found for state: %s", userState.State)
		}
	}
	return nil
}

func (b *Bot) RegisterState(state string, handler StateHandler) {
	b.StateManager.RegisterState(state, handler)
}

func (b *Bot) SendMessage(chatId int64, text string) error {
	if b.ChannelID != 0 {
		chatId = b.ChannelID
	}
	msg := tgbotapi.NewMessage(chatId, text)
	_, err := b.API.Send(msg)
	return err
}

func (b *Bot) SendRichMessage(chatId int64, text string, parseMode string) error {
	if b.ChannelID != 0 {
		chatId = b.ChannelID
	}
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = parseMode
	_, err := b.API.Send(msg)
	return err
}

func (b *Bot) SendPhoto(chatId int64, photoPath string, caption string) {
	if b.ChannelID != 0 {
		chatId = b.ChannelID
	}
	relativePath := filepath.Join("uploads", photoPath)
	photoFile, err := os.ReadFile(relativePath)
	if err != nil {
		log.Panic(err)
	}

	msg := tgbotapi.NewPhotoUpload(chatId, tgbotapi.FileBytes{
		Name:  photoPath,
		Bytes: photoFile,
	})
	if caption != "" {
		msg.Caption = caption
	}
	b.API.Send(msg)
}

func (b *Bot) SendDocument(chatId int64, documentPath string) {
	if b.ChannelID != 0 {
		chatId = b.ChannelID
	}
	relativePath := filepath.Join("uploads", documentPath)
	documentFile, err := os.ReadFile(relativePath)
	if err != nil {
		log.Panic(err)
	}
	msg := tgbotapi.NewDocumentUpload(chatId, tgbotapi.FileBytes{
		Name:  documentPath,
		Bytes: documentFile,
	})
	b.API.Send(msg)
}

func (b *Bot) UpdateKeyboard(buttons ...tgbotapi.KeyboardButton) tgbotapi.ReplyKeyboardMarkup {
	var buttonRows [][]tgbotapi.KeyboardButton
	for _, button := range buttons {
		buttonRows = append(buttonRows, tgbotapi.NewKeyboardButtonRow(button))
	}

	return tgbotapi.NewReplyKeyboard(buttonRows...)
}
