package handler

import (
	"context"
	"net/http"

	"forum/internal/models"
)

func (h *Handler) CheckSession(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get cookie
		c, err := r.Cookie(string(models.SessionKey))
		if err != nil {
			h.logger.Error("check session in middleware", err)
			next.ServeHTTP(w, r)
			return
		}

		// get session
		session, err := h.service.SessionService.GetSessionByToken(c.Value)
		if err != nil {
			h.logger.Error("Session not found in middleware", err)
			next.ServeHTTP(w, r)
			return
		}

		// context
		ctx := context.WithValue(r.Context(), models.SessionKey, *session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
