package handler

import (
	"context"
	"net/http"

	"forum/internal/models"
)

func (h *Handler) SessionAuthMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(string(models.SessionKey))
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		session, err := h.service.SessionService.GetSessionByToken(cookie.Value)
		if err != nil {
			h.logger.Errorf("get session by token: %v", err)
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), models.SessionKey, *session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
