package models

import "time"

type Like struct {
	LikeID    string
	PostID    string
	UserID    string
	CreatedAt string
}

type Dislike struct {
	DislikeID string
	PostID    string
	UserID    string
	CreatedAt string
}

type Comment struct {
	CommentID    string
	PostID       string
	Author       string
	CommentText  string
	LikeCount    int
	DislikeCount int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	IsLiked      bool
	DisLiked     bool
	OwnerId      string
}

type CommentLike struct {
	LikeID    string
	CommentID string
	UserID    string
	CreatedAt string
}

type CommentDislike struct {
	DislikeID string
	CommentID string
	UserID    string
	CreatedAt string
}

type ReactionRepository interface {
	// like in post
	InsertLikePost(like *Like) error              // insert new like
	LikeExistInPost(post_id, user_id string) bool // check like current post from user
	DeleteLike(post_id, user_id string) error     // delete like
	// dislike in post
	InsertDisLikePost(disLike *Dislike) error        // insert new dislike
	DisLikeExistInPost(post_id, user_id string) bool // check dislike current post from user
	DeleteDisLike(post_id, user_id string) error     // delete dislike
	// comments in post
	InsertCommentInPost(comment *Comment) error   // insert new comment
	GetCommentsByID(id string) ([]Comment, error) // get comment list by post_id
	DeleteComment(comment_id string) error        // delete comment
	DeleteCommentByPostID(post_id string) error   // delete comment by POST ID
	// like in comment
	InsertLikeInComment(reaction *CommentLike) error      // insert new like in comment
	IncrementLikeInComment(comment_id string) error       // increment like in comment
	DecrementLikeInComment(comment_id string) error       // decrement like in comment
	ExistLikeInComment(user_id, comment_id string) bool   // check exist like in comment from current user
	DeleteLikeInComment(comment_id, user_id string) error // delete dislike in comment
	DeleteLikeInCommentByUserID(user_id string) error     // delete Like In Comment By Post ID
	// dislike in comment
	InsertDisLikeInComment(reaction *CommentDislike) error // insert dislike in comment
	IncrementDisLikeInComment(comment_id string) error     // increment dislike in comment
	DecrementDisLikeInComment(comment_id string) error     // decrement dislike in comment
	ExistDisLikeInComment(user_id, comment_id string) bool // check exist dislike in comment from current user
	DeleteDisLikeInComment(comment_id, user_id string) error
	DeleteDisLikeInCommentByUserID(user_id string) error
}

type ReactionService interface {
	// delete reaction
	DeleteReaction(post_id, user_id string) error
	// like in post
	CreateLikeInPost(like *Like) error              // create like in post
	LikeExistInPost(post_id, user_id string) bool   // check like exist in post
	DeleteLikeInPost(post_id, user_id string) error // delete like in post
	// dislike in post
	CreateDislikeInPost(dislike *Dislike) error      // create dislike in post
	DisLikeExistInPost(post_id, user_id string) bool // check dislike exist in post
	DeleteDisLike(post_id, user_id string) error     // delete dislike in post
	// comment
	CreateCommentInPost(comment *Comment) error        // create comment in post
	GetCommentsByID(post_id string) ([]Comment, error) // get comment list by POST ID
	DeleteComment(comment_id string) error
	// like in comment
	CreateLikeInComment(reaction *CommentLike) error      // create  like in comment
	IncrementLikeInComment(comment_id string) error       // increment like in comment
	DecrementLikeInComment(comment_id string) error       // decrement like in comment
	ExistLikeInComment(user_id, comment_id string) bool   // check exist like in comment from current user
	DeleteLikeInComment(comment_id, user_id string) error // delete like in comment
	// dislike in comment
	CreateDisLikeInComment(reaction *CommentDislike) error
	IncrementDisLikeInComment(comment_id string) error     // increment dislike in comment
	DecrementDisLikeInComment(comment_id string) error     // decrement dislike in comment
	ExistDisLikeInComment(user_id, comment_id string) bool // check exist dislike in comment from current user
	DeleteDisLikeInComment(comment_id, user_id string) error
}
