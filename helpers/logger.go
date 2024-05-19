package helpers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Pratham-Mishra04/interactWS/initializers"
	"github.com/golang-jwt/jwt/v5"
)

func LogInfo(customString string, path string) {
	initializers.Logger.Info(customString, "Path", path)
	LogToAdminLogger(customString, "warn", nil, path)
}

func LogWarn(customString string, err error, path string) {
	initializers.Logger.Warnw(customString, "Path", path, "Error", err)
	LogToAdminLogger(customString, "warn", err, path)
}

func LogError(customString string, err error, path string) {
	initializers.Logger.Errorw(customString, "Path", path, "Error", err)
	LogToAdminLogger(customString, "error", err, path)
}

func LogFatal(customString string, err error, path string) {
	initializers.Logger.Infow(customString, "Path", path, "Error", err)
	LogToAdminLogger(customString, "fatal", err, path)
}

type LogEntrySchema struct {
	Level       string `json:"level"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Path        string `json:"path"`
	Timestamp   string `json:"timestamp"`
}

func createAdminJWT() (string, error) {
	token_claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "sockets",
		"crt": time.Now().Unix(),
		"exp": time.Now().Add(30 * time.Second).Unix(),
	})

	token, err := token_claim.SignedString([]byte(initializers.CONFIG.LOGGER_SECRET))
	if err != nil {
		return "", err
	}

	return token, nil
}

func LogToAdminLogger(customString string, level string, err error, path string) {
	logEntry := LogEntrySchema{
		Level:       level,
		Title:       customString,
		Description: err.Error(),
		Path:        path,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		initializers.Logger.Errorw("Error Posting to Admin Logger", "Message", err.Error(), "Path", "LogToAdminLogger", "Error", err.Error())
		return
	}

	jwt, err := createAdminJWT()
	if err != nil {
		return
	}

	request, err := http.NewRequest("POST", initializers.CONFIG.LOGGER_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		initializers.Logger.Errorw("Error Posting to Admin Logger", "Message", err.Error(), "Path", "LogToAdminLogger", "Error", err.Error())
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+jwt)
	request.Header.Set("api-token", initializers.CONFIG.LOGGER_TOKEN)

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		initializers.Logger.Errorw("Error Adding to Admin Logger", "Message", err.Error(), "Path", "LogToAdminLogger", "Error", err.Error())
		return
	}
	defer response.Body.Close()
}
