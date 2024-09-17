package handler

import (
	"text/template"

	"forum/internal/service"
	"forum/pkg/logger"
)

type Handler struct {
	cache   map[string]*template.Template
	service *service.Service
	logger  *logger.Logger
}

func NewHandler(l *logger.Logger, s *service.Service, cache map[string]*template.Template) *Handler {
	return &Handler{
		cache:   cache,
		service: s,
		logger:  l,
	}
}
