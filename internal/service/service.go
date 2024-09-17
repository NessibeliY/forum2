package service

import (
	"forum/internal/repository"
	"forum/internal/service/post"
	"forum/internal/service/reaction"
	"forum/internal/service/session"
	"forum/internal/service/user"
	"forum/pkg/logger"
)

type Service struct {
	User     user.IUserService
	Session  session.ISessionService
	Post     post.IPostService
	Reaction reaction.IReactionService
}

func NewService(l *logger.Logger, r *repository.Repository) *Service {
	return &Service{
		User:     user.NewUserService(r.User),
		Session:  session.NewSessionService(r.Session),
		Post:     post.NewPostService(r.Post),
		Reaction: reaction.NewReactionService(r.Reaction),
	}
}
