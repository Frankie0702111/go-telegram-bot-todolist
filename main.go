package main

import (
	"fmt"
	"go-telegram-bot-todolist/entity"
	gorm_utils "go-telegram-bot-todolist/utils/gorm"
	logger "go-telegram-bot-todolist/utils/log"
	"go-telegram-bot-todolist/utils/response"
	"os"
	"regexp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
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

	// Telegram bot
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_KEY"))
	if err != nil {
		logger.Panic(err)
	}
	// log.Printf("Authorized on account %s", bot.Self.UserName)

	bot.Debug = true

	// Cron job
	c := cron.New()
	i := 1
	c.AddFunc("*/5 7-23 * * *", func() {
		now := time.Now()
		if now.Hour() < 7 || now.Hour() > 23 {
			return
		}

		if now.Minute()%5 != 0 {
			fmt.Print("You broke the rules, run the cron job.")
			return
		}

		tasks, tasksErr := entity.TaskNotify(db)
		if tasksErr != nil {
			logger.Error("Failed to send message for task notification : " + tasksErr.Error())
		}

		for _, task := range tasks {
			if len(tasks) > 0 {
				go sendNotifyMessage(bot, task, true)
			}
		}

		fmt.Println("Run every 5 minutes : ", i)
		i++
	})
	c.Start()
	defer c.Stop()

	// Telegram bot send message
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			// log.Printf("[%s] %s", update.Message.From.LastName, update.Message.Text)
			go sendMessage(bot, update, db)
		}
	}
}

func sendMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *gorm.DB) {
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
		case "tasklist":
			checkChatID := entity.FindByChatID(db, update.Message.Chat.ID)
			if checkChatID {
				notify, _ := entity.TaskNotifyByChatID(db, update.Message.Chat.ID)
				for _, task := range notify {
					if len(notify) > 0 {
						go sendNotifyMessage(bot, task, false)
					}
				}
				// jsonBytes, _ := json.Marshal(notify)
				// msg.Text = string(jsonBytes)
				return
			} else {
				msg.Text = "Please input your email (e.g. /setaccess test@test.com)"
			}
		default:
			msg.Text = "I don't know that command."
		}
	}
	bot.Send(msg)
}

func sendNotifyMessage(bot *tgbotapi.BotAPI, message response.Notify, is_update bool) {
	msg := tgbotapi.NewMessage(message.TelegramID, fmt.Sprintf("Notify\nTitle: %s\nDateTime: %s\nComplete: %t\n", message.Title, message.Time, message.Complete))
	_, err := bot.Send(msg)
	if err == nil {
		if is_update {
			updateErr := entity.TaskNotifyUpdate(db, message.TaskID)
			if updateErr != nil {
				logger.Error("The update task is_notify has failed : " + updateErr.Error())
			}
		}
	}
}
