package repository

import (
	"database/sql"

	"forum/internal/models"
	"forum/internal/repository/post"
	"forum/internal/repository/reaction"
	"forum/internal/repository/session"
	"forum/internal/repository/user"
)

type Repository struct {
	UserRepository     models.UserRepository
	SessionRepository  models.SessionRepository
	PostRepository     models.PostRepository
	ReactionRepository models.ReactionRepository
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		UserRepository:     user.NewUserRepo(db),
		SessionRepository:  session.NewSessionRepo(db),
		PostRepository:     post.NewPostRepo(db),
		ReactionRepository: reaction.NewReactionRepo(db),
	}
}
