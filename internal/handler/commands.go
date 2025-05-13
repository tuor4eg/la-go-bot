package handler

import (
	"fmt"
)

type Command struct {
	Name        string // название команды, например "/help"
	Description string // ключ для локализации описания
	IsCommand   bool   // true если команда, false если нет
}

var Commands = []Command{
	{
		Name:        "/help",
		Description: "help_cmd",
		IsCommand:   true,
	},
	{
		Name:        "/info",
		Description: "info_cmd",
		IsCommand:   true,
	},
	{
		Name:        "/start",
		Description: "start_cmd",
		IsCommand:   true,
	},
	{
		Name:        "send_geo",
		Description: "send_geo_cmd",
		IsCommand:   false,
	},
}

func FormatCommandList(lang string) string {
	var result string
	for _, cmd := range Commands {
		name := cmd.Name

		if !cmd.IsCommand {
			name = GetTranslation(lang, name)
		}

		result += fmt.Sprintf("%s - %s\n",
			name,
			GetTranslation(lang, cmd.Description),
		)
	}
	return result
}
