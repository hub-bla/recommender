package httprequesthandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type ErrorResponse struct {
	Error   string
	Message string
}

type Logger struct {
}

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func (logger *Logger) logMessage(message string, messageCode int, logType string) {
	if logType == "success" {
		fmt.Println("INFO:	" + "[" + Green + strconv.Itoa(messageCode) + Reset + "] " + message)
		return
	}
	fmt.Println("INFO:	" + "[" + Red + strconv.Itoa(messageCode) + Reset + "] " + message)

}

type HTTPRequestHandler struct {
	logger Logger
}

// HTTP/1.1 400 Bad Request
// Content-Type: application/json

// {
//   "error": "Bad Request",
//   "message": "Missing required fields: 'field1', 'field2'"
// }

type ErrorMessage struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (httpRH *HTTPRequestHandler) SendResponse(w *http.ResponseWriter,
	requestData any,
	log string,
	statusCode int,
	logType string) error {

	(*w).WriteHeader(statusCode)
	(*w).Header().Set("Content-Type", "application/json")
	err := json.NewEncoder((*w)).Encode(requestData)

	if err != nil {
		message := "failed to encode and send response"
		httpRH.logger.logMessage("Error: Failed to encode and send response", http.StatusInternalServerError, "error")
		errorMessage := ErrorMessage{
			Error:   "Internal Server Error",
			Message: message,
		}
		(*w).WriteHeader(http.StatusInternalServerError)
		(*w).Header().Set("Content-Type", "application/json")
		json.NewEncoder((*w)).Encode(errorMessage)
		return fmt.Errorf(message)
	}

	httpRH.logger.logMessage(log, statusCode, logType)
	return nil

}
