package entity

import (
	"go-telegram-bot-todolist/model"
	"go-telegram-bot-todolist/utils/response"
	"time"

	"gorm.io/gorm"
)

func TaskNotify(db *gorm.DB) (notifies []response.Notify, err error) {
	var tasks []model.Task
	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), 0, 0, 0, time.Now().Location())
	beforeTime := (now.Add(-8 * time.Hour)).Format("2006-01-02 15:04:05")
	resErr := db.Preload("User").Where(
		"specify_datetime BETWEEN ? AND ? AND is_complete = ? AND is_notify = ?",
		beforeTime,
		now.Format("2006-01-02 15:04:05"),
		false,
		0,
	).Or(
		"specify_datetime < ? AND is_complete = ? AND is_notify = ?",
		now.Format("2006-01-02 15:04:05"),
		false,
		0,
	).Find(&tasks).Error
	if resErr != nil {
		return nil, resErr
	}

	for _, task := range tasks {
		notify := response.Notify{
			TaskID:     task.ID,
			TelegramID: task.User.TelegramID,
			Title:      task.Title,
			Time:       task.SpecifyDatetime.Format("2006-01-02 15:04:05"),
			Complete:   task.IsComplete,
		}
		notifies = append(notifies, notify)
	}

	return notifies, nil
}

func TaskNotifyUpdate(db *gorm.DB, task_id int64) error {
	var task model.Task

	update := db.Model(&task).Where("id=?", task_id).Update("is_notify", 1).Error

	if update != nil {
		return update
	}

	return nil
}

func TaskNotifyByChatID(db *gorm.DB, telegram_id int64) (notifies []response.Notify, err error) {
	var user model.User
	now := (time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), 0, 0, time.Now().Location())).Format("2006-01-02 15:04:05")
	userErr := db.Where("telegram_id = ? AND status = ?", telegram_id, 1).First(&user).Error
	if userErr != nil {
		return nil, userErr
	}

	var tasks []model.Task
	resErr := db.Where(
		"user_id = ? AND is_complete = ?",
		user.ID,
		false,
	).Where(
		"specify_datetime < ? OR specify_datetime > ?",
		now,
		now,
	).Find(&tasks).Error
	if resErr != nil {
		return nil, resErr
	}

	for _, task := range tasks {
		notify := response.Notify{
			Title:      task.Title,
			TelegramID: user.TelegramID,
			Time:       task.SpecifyDatetime.Format("2006-01-02 15:04:05"),
			Complete:   task.IsComplete,
		}
		notifies = append(notifies, notify)
	}

	return notifies, nil
}
