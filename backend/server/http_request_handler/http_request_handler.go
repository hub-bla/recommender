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

const Reset = "\033[0m"
const Red = "\033[31m"
const Green = "\033[32m"
const Yellow = "\033[33m"
const Blue = "\033[34m"
const Magenta = "\033[35m"
const Cyan = "\033[36m"
const Gray = "\033[37m"
const White = "\033[97m"

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
		httpRH.logger.logMessage(message, http.StatusInternalServerError, "error")
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
