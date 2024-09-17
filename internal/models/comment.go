package models

import "time"

type Comment struct {
	CommentId   string
	PostId      string
	Author      string
	CommentText string
	Likes       int
	DisLikes    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	IsLiked     bool
	DisLiked    bool
	OwnerId     string
}
