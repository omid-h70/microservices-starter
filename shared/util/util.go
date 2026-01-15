package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetRandomAvatar returns a random avatar URL from the randomuser.me API
func GetRandomAvatar(index int) string {
	return fmt.Sprintf("https://randomuser.me/api/portraits/lego/%d.jpg", index)
}

func GenerateRandomPlate() string {
	return "abc123"
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Contetnt-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
