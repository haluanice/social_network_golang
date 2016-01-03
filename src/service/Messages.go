package service

import (
	"encoding/json"
	"model"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
)

var globalExecutionSuccessMessage atomic.Value
var globalExecutionErrorMessage atomic.Value

func SetHeaderParameter(w http.ResponseWriter) {
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Content-Type", "application/json")
}

func GetTokenHeader(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Asolole ")
	token := splitToken[len(splitToken)-1]
	return token == "jambu"
}

func StringtoInt(integer string) int {
	newInteger, _ := strconv.ParseInt(integer, 10, 0)
	return int(newInteger)
}

func OutputError(message string) string {
	globalExecutionErrorMessage.Store(model.ErrorMessage{message})
	dataErrorMessage := globalExecutionErrorMessage.Load().(model.ErrorMessage)
	output, _ := json.Marshal(dataErrorMessage)
	return string(output)
}

func OutputSuccess(message string, user model.User) string {
	globalExecutionSuccessMessage.Store(model.SuccessMessage{message, user})
	SuccessMessageMessage := globalExecutionSuccessMessage.Load().(model.SuccessMessage)
	output, _ := json.Marshal(SuccessMessageMessage)
	return string(output)
}
