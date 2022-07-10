package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type logMessage struct {
	ChatID  int64  `json:"chatID"`
	Message string `json:"message"`
}

var staticTelegramLogChan *chan logMessage

func startTelegramLog(telLogChan chan logMessage, stop chan bool, wg *sync.WaitGroup, settings *Settings) {
	defer wg.Done()
	staticTelegramLogChan = &telLogChan

	if settings.TelegramBotAPI == "" {
		log.Printf("We don't used tg bot")
		return
	}

	bot, err := tgbotapi.NewBotAPI(settings.TelegramBotAPI)
	if err != nil {
		log.Printf("error in tgbotapi %s", err.Error())
		return
	}

	msg := tgbotapi.NewMessage(settings.TelegramLogFeedChatID, "started bot")
	bot.Send(msg)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 20

	updates, err := bot.GetUpdatesChan(u)
	defer func() {
		msg := tgbotapi.NewMessage(settings.TelegramLogFeedChatID, "stopped bot")
		bot.Send(msg)
	}()
	startDate := time.Now().Unix()
	for i := 0; ; i++ {
		select {
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			if startDate > int64(update.Message.Date) {
				log.Printf("got old message \"%s\" from %s\n", update.Message.Text, update.Message.From.UserName)
				continue
			}

		case curLogMessage := <-telLogChan:
			if curLogMessage.ChatID == 0 {
				curLogMessage.ChatID = settings.TelegramLogFeedChatID
			}
			msg := tgbotapi.NewMessage(curLogMessage.ChatID, curLogMessage.Message)
			bot.Send(msg)
		case <-stop:
			log.Printf("closed telegramLog")
			return
		}
	}
}

func gotCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

}

func myLogChat(chatID int64, format string, a ...interface{}) {

	go func() {
		if staticTelegramLogChan != nil {
			*staticTelegramLogChan <- logMessage{ChatID: chatID, Message: fmt.Sprintf(format, a...)}
		}
	}()

}
func myLog(format string, a ...interface{}) {

	go func() {
		if staticTelegramLogChan != nil {
			*staticTelegramLogChan <- logMessage{Message: fmt.Sprintf(format, a...)}
		}
	}()

}

func myLogAlert(id int64, format string, a ...interface{}) {
	if id == 0 {
		id = 1
	}
	myLog("==========ALERT==========\n==========ALERT==========\n==========ALERT==========")
	myLogChat(id, "==========ALERT==========\n==========ALERT==========\n==========ALERT==========")
	time.Sleep(100 * time.Millisecond)
	myLog(fmt.Sprintf(format, a...))
	myLogChat(id, fmt.Sprintf(format, a...))
	time.Sleep(100 * time.Millisecond)
	myLog("==========ALERT==========\n==========ALERT==========\n==========ALERT==========")
	myLogChat(id, "==========ALERT==========\n==========ALERT==========\n==========ALERT==========")
}

func myLogAttention(id int64, format string, a ...interface{}) {
	if id == 0 {
		id = 1
	}
	myLog("==========START_ATTENTION==========\n==========START_ATTENTION==========\n==========START_ATTENTION==========")
	myLogChat(id, "==========START_ATTENTION==========\n==========START_ATTENTION==========\n==========START_ATTENTION==========")
	time.Sleep(100 * time.Millisecond)
	myLog(fmt.Sprintf(format, a...))
	myLogChat(id, fmt.Sprintf(format, a...))
	time.Sleep(100 * time.Millisecond)
	myLog("==========END_ATTENTION==========\n==========END_ATTENTION==========\n==========END_ATTENTION==========")
	myLogChat(id, "==========END_ATTENTION==========\n==========END_ATTENTION==========\n==========END_ATTENTION==========")
}
