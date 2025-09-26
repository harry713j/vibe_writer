package handler

import (
	"errors"
	"net/http"

	"github.com/harry713j/vibe_writer/internal/middleware"
	"github.com/harry713j/vibe_writer/internal/service"
	"github.com/harry713j/vibe_writer/internal/utils"
)

// upload to cloud
func HandleUploadToCloud(w http.ResponseWriter, r *http.Request) {

	_, ok := middleware.GetUserID(r)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	file, fileHeader, err := r.FormFile("img")

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Image file required")
		return
	}

	defer file.Close()

	// check image type
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Couldn't read file")
		return
	}

	contentType := http.DetectContentType(buffer)
	if contentType != "image/png" && contentType != "image/jpeg" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid image type")
		return
	}

	// check file size
	if fileHeader.Size > 5*1024*1024 {
		utils.RespondWithError(w, http.StatusBadRequest, "Large image found")
		return
	}

	imgUrl, err := service.Upload(file, fileHeader.Filename)

	if err != nil {

		if errors.Is(err, service.ErrImageNotAllowed) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	type uploadResponse struct {
		PhotoUrl string `json:"photo_url"`
	}

	utils.RespondWithJSON(w, http.StatusCreated, uploadResponse{
		PhotoUrl: imgUrl,
	})
}
