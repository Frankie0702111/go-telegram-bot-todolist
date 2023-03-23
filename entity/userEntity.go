package entity

import (
	"go-telegram-bot-todolist/model"

	"gorm.io/gorm"
)

func FindByEmail(db *gorm.DB, email string) bool {
	var user model.User
	checkExist := db.Where("email = ?", email).Take(&user).Error
	if checkExist == nil {
		return true
	}

	return false
}

func FindByChatID(db *gorm.DB, telegram_id int64) bool {
	var user model.User
	checkExist := db.Where("telegram_id = ?", telegram_id).Take(&user).Error
	if checkExist == nil {
		return true
	}

	return false
}

func UpdateTelegramIDByEmail(db *gorm.DB, email string, telegram_id int64) bool {
	var user model.User
	res := db.Model(&user).Where("email = ?", email).Update("telegram_id", telegram_id).Error
	if res == nil {
		return true
	}

	return false
}
