package notify

import (
	"fmt"
	"os/exec"
)

func Notify(project string) {
	exec.Command("notify-send", "Оповещение", fmt.Sprintf("Вы работаете над проектом '%s' уже 25 минут! Время сделать перерыв.", project)).Run()
}
