package handlers

import "github.com/manifoldco/promptui"

func (h *Handlers) GeneralMenu() {
	for {
		prompt := promptui.Select{
			Label: "Главное меню",
			Items: []string{
				"Выбрать проект",
				"Создать проект",
				"Сводка по всем проектам",
				"Выход",
			},
		}
		_, cmd, _ := prompt.Run()

		switch cmd {
		case "Выбрать проект":
			h.SelectAndManageProject()
		case "Создать проект":
			h.CreateProject()
		case "Сводка по всем проектам":
			h.ShowSummary()
		case "Выход":
			h.ProjectService.SaveData(h.Projects)
			return
		}
	}
}
