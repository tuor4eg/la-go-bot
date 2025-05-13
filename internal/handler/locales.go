package handler

// Группируем переводы по категориям для лучшей организации
var translations = map[string]map[string]string{
	"en": {
		// Common UI elements
		"choose_action": "Choose action:",

		// Commands
		"help_cmd":     "🆘 Help",
		"info_cmd":     "📱 My Information",
		"start_cmd":    "🚀 Start",
		"send_geo":     "📍 Send Location or coordinates (latitude, longitude)",
		"send_geo_cmd": "Get cameras list in 500m radius",
		"unknown_cmd":  "Unknown command",

		// User Info Section
		"user_info":   "📱 User Information",
		"name":        "👤 Name",
		"email":       "📧 Email",
		"role":        "🎭 Role",
		"language":    "🌐 Language",
		"telegram_id": "📱 Telegram ID",
		"phone":       "📞 Phone",
		"last_update": "🕒 Last Update",

		// Cameras Section
		"cameras_list": "📸 Cameras List in 500m radius",
		"address":      "📍 Address",
		"title":        "🏷️ Title",
		"distance":     "🔍 Distance",
		"meter":        "m",
		"no_cameras":   "No cameras found in 500m radius",
	},
	"ru": {
		// Common UI elements
		"choose_action": "Выберите действие:",

		// Commands
		"info_cmd":     "📱 Моя информация",
		"help_cmd":     "🆘 Помощь",
		"start_cmd":    "🚀 Начать",
		"send_geo":     "📍 Отправить местоположение или координаты (широта, долгота)",
		"send_geo_cmd": "Получить список камер в радиусе 500м",
		"unknown_cmd":  "Неизвестная команда",
		// User Info Section
		"user_info":   "📱 Информация о пользователе",
		"name":        "👤 Имя",
		"email":       "📧 Почта",
		"role":        "🎭 Роль",
		"language":    "🌐 Язык",
		"telegram_id": "📱 Telegram ID",
		"phone":       "📞 Телефон",
		"last_update": "🕒 Последнее обновление",

		// Cameras Section
		"cameras_list": "📸 Список камер в радиусе 500м",
		"address":      "📍 Адрес",
		"title":        "🏷️ Название",
		"distance":     "🔍 Расстояние",
		"meter":        "м",
		"no_cameras":   "Камеры не найдены в радиусе 500м",
	},
}

func GetTranslation(lang, key string) string {
	if lang == "" {
		lang = "ru"
	}

	if translations[lang] == nil {
		lang = "ru"
	}

	if val, ok := translations[lang][key]; ok {
		return val
	}

	return translations["ru"][key]
}
