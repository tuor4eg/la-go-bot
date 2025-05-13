package handler

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type UserSettings struct {
	Language string `json:"language"`
}

type UserAccounts struct {
	TelegramID int64  `json:"telegramId"`
	Phone      string `json:"phone"`
}

type User struct {
	ID        string       `json:"_id"`
	Name      string       `json:"name"`
	Email     string       `json:"email"`
	Role      string       `json:"role"`
	Settings  UserSettings `json:"settings"`
	Accounts  UserAccounts `json:"accounts"`
	UpdatedAt string       `json:"updatedAt"`
}

type UserResponse struct {
	User User `json:"user"`
}

type Camera struct {
	ID       string  `json:"_id"`
	Address  string  `json:"address"`
	Title    string  `json:"title"`
	Distance float64 `json:"distance"`
}

func formatDate(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr
	}
	return t.Format("02.01.2006 15:04:05")
}

func FormatUserInfo(userJSON string) (string, error) {
	var response UserResponse
	if err := json.Unmarshal([]byte(userJSON), &response); err != nil {
		return "", fmt.Errorf("error parsing user info: %w", err)
	}

	user := response.User
	lang := user.Settings.Language

	return fmt.Sprintf(`%s

%s: %s
%s: %s
%s: %s
%s: %s
%s: %d
%s: %s
%s: %s`,
		GetTranslation(lang, "user_info"),
		GetTranslation(lang, "name"), user.Name,
		GetTranslation(lang, "email"), user.Email,
		GetTranslation(lang, "role"), user.Role,
		GetTranslation(lang, "language"), user.Settings.Language,
		GetTranslation(lang, "telegram_id"), user.Accounts.TelegramID,
		GetTranslation(lang, "phone"), user.Accounts.Phone,
		GetTranslation(lang, "last_update"), formatDate(user.UpdatedAt),
	), nil
}

func FormatCameras(camerasJSON string, lang string) (string, error) {
	var response struct {
		Cameras []Camera `json:"cameras"`
	}

	if err := json.Unmarshal([]byte(camerasJSON), &response); err != nil {
		return "", fmt.Errorf("error parsing cameras: %w", err)
	}

	if len(response.Cameras) == 0 {
		return GetTranslation(lang, "no_cameras"), nil
	}

	formattedCameras := fmt.Sprintf("%s\n", GetTranslation(lang, "cameras_list"))

	for _, camera := range response.Cameras {
		name := camera.Title
		if camera.Address != "" {
			name = camera.Address
		}
		formattedCameras += fmt.Sprintf("%s (%.2f %s)\n", name, camera.Distance, GetTranslation(lang, "meter"))
	}
	return formattedCameras, nil
}

func ParseCoordinates(message string) (float64, float64, error) {
	// Regular expression to find latitude and longitude in the message
	// Matches patterns like: 12.345, -67.890 or 12.345 -67.890 or lat:12.345 lng:-67.890
	re := regexp.MustCompile(`[-+]?([0-9]*\.[0-9]+|[0-9]+)[^0-9\-+.]+[-+]?([0-9]*\.[0-9]+|[0-9]+)`)
	matches := re.FindStringSubmatch(message)

	if len(matches) < 3 {
		return 0, 0, fmt.Errorf("coordinates not found in message")
	}

	// Parse the matched coordinates
	lat, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid latitude format: %w", err)
	}

	lng, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid longitude format: %w", err)
	}

	// Validate latitude and longitude ranges
	if lat < -90 || lat > 90 {
		return 0, 0, fmt.Errorf("latitude must be between -90 and 90 degrees")
	}

	if lng < -180 || lng > 180 {
		return 0, 0, fmt.Errorf("longitude must be between -180 and 180 degrees")
	}

	return lat, lng, nil
}
