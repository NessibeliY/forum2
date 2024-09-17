package handler

import (
	"net/http"

	"forum/internal/models"
)

func (h *Handler) ErrorHandler(w http.ResponseWriter, statusCode int, errorText string) {
	data := &models.ErrorPage{
		StatusCode: statusCode,
		TextError:  errorText,
	}

	h.Render(w, "error.html", data)
}
