package handler

import (
	"net/http"

	"github.com/harry713j/vibe_writer/internal/utils"
)

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, 200, map[string]string{
		"message": "Server is running well",
	})
}
