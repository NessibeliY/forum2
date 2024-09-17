package repository

import (
	"database/sql"

	"forum/internal/repository/post"
	"forum/internal/repository/reaction"
	"forum/internal/repository/session"
	"forum/internal/repository/user"
	"forum/pkg/logger"
)

type Repository struct {
	User     *user.UserRepo
	Session  *session.SessionRepo
	Post     *post.PostRepo
	Reaction *reaction.ReactionRepo
}

func NewRepo(l *logger.Logger, db *sql.DB) *Repository {
	return &Repository{
		User:     user.NewUserRepo(db),
		Session:  session.NewSessionRepo(db),
		Post:     post.NewPostRepo(db),
		Reaction: reaction.NewReactionRepo(db),
	}
}
