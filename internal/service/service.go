package service

import (
	"forum/internal/models"
	"forum/internal/repository"
	"forum/internal/service/post"
	"forum/internal/service/reaction"
	"forum/internal/service/session"
	"forum/internal/service/user"
)

type Service struct {
	UserService     models.UserService
	SessionService  models.SessionService
	PostService     models.PostService
	ReactionService models.ReactionService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		UserService:     user.NewUserService(repo.UserRepository),
		SessionService:  session.NewSessionService(repo.SessionRepository),
		PostService:     post.NewPostService(repo.PostRepository),
		ReactionService: reaction.NewReactionService(repo.ReactionRepository),
	}
}
