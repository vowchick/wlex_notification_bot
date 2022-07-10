package main

// Settings struct that contains all settings
type Settings struct {
	fileName string

	TelegramBotAPI        string `json:"telegramBotApi"` // апи тг бота
	TelegramLogFeedChatID int64  `json:"logFeedChatID"`  // id чата
}
