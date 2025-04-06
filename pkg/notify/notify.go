package notify

import (
	"fmt"
	"os/exec"
)

// Send - отправка уведомления
func Send(title, message string) error {
	return exec.Command("notify-send", title, message).Run()
}

// SendBreakReminder - отправка напоминания о перерыве
func SendBreakReminder(project string) error {
	return Send("Оповещение", fmt.Sprintf("Вы работаете над проектом '%s' уже 25 минут! Время сделать перерыв.", project))
}
