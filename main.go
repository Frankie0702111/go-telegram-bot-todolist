package main

import (
	"go-telegram-bot-todolist/entity"
	gorm_utils "go-telegram-bot-todolist/utils/gorm"
	logger "go-telegram-bot-todolist/utils/log"
	"os"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var (
	db *gorm.DB = gorm_utils.InitMySQL()
)

func main() {
	defer gorm_utils.Close(db)

	// Load .env file
	errEnv := godotenv.Load()
	if errEnv != nil {
		logger.Panic("Failed to load env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_KEY"))
	if err != nil {
		logger.Panic(err)
	}
	// log.Printf("Authorized on account %s", bot.Self.UserName)

	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 600
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			// log.Printf("[%s] %s", update.Message.From.LastName, update.Message.Text)
			go sendMessage(bot, update, db)
		}
	}
}

func sendMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *gorm.DB) {
	// user := model.User{}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			msg.Text = "Hello " + update.Message.From.LastName
		case "setaccess":
			args := update.Message.CommandArguments()
			match, _ := regexp.MatchString("^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", args)
			if match {
				if args != "" {
					checkEmail := entity.FindByEmail(db, args)
					updateTelegramID := entity.UpdateTelegramIDByEmail(db, args, update.Message.Chat.ID)
					if checkEmail && updateTelegramID {
						msg.Text = "Succesfully to set account."
					} else {
						msg.Text = "Failed to set account."
					}
				} else {
					msg.Text = "Please input your email (e.g. /setaccess test@test.com)"
				}
			} else {
				msg.Text = "Incorrect Email Format"
			}
		default:
			msg.Text = "I don't know that command."
		}
	}
	bot.Send(msg)
}
