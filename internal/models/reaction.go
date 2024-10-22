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
	CommentID     string
	PostID        string
	Author        string
	CommentText   string
	LikesCount    int
	DislikesCount int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	IsLiked       bool
	DisLiked      bool
	OwnerID       string
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
	AddPostLike(like Like) error
	IsPostLikedByUser(postID, userID string) bool
	RemovePostLike(postID, userID string) error
	// dislike in post
	AddPostDislike(disLike *Dislike) error
	IsPostDislikedByUser(postID, userID string) bool
	RemovePostDislike(postID, userID string) error
	// comments in post
	AddComment(comment *Comment) error
	GetCommentsByPostID(id string) ([]Comment, error)
	RemoveCommentByCommentID(commentID string) error
	RemoveCommentByPostID(postID string) error
	// like in comment
	AddCommentLike(reaction *CommentLike) error
	IncrementLikeCountInComment(commentID string) error
	DecrementLikeCountInComment(commentID string) error
	IsCommentLikedByUser(userID, commentID string) bool
	RemoveCommentLike(commentID, userID string) error
	RemoveAllCommentLikesByUser(userID string) error
	// dislike in comment
	AddCommentDislike(reaction *CommentDislike) error
	IncrementDislikeCountInComment(commentID string) error
	DecrementDislikeCountInComment(commentID string) error
	IsCommentDislikedByUser(userID, commentID string) bool
	RemoveCommentDislike(commentID, userID string) error
	RemoveAllCommentDislikesByUser(userID string) error
}

type ReactionService interface {
	// delete reaction
	DeleteReaction(postID, userID string) error
	// like in post
	HandlePostLike(userID, postID string) error
	IsPostLikedByUser(postID, userID string) bool
	DeleteLikeInPost(postID, userID string) error
	// dislike in post
	HandlePostDislike(userID, postID string) error
	IsPostDislikedByUser(postID, userID string) bool
	RemovePostDislike(postID, userID string) error
	// comment
	GetCommentsWithReactions(postID, username, userID string) []Comment
	CreateCommentInPost(comment *Comment) error
	GetCommentsByPostID(postID string) ([]Comment, error)
	RemoveCommentByCommentID(commentID string) error
	// like in comment
	HandleCommentLike(userID, commentID string) error
	CreateLikeInComment(reaction *CommentLike) error
	IncrementLikeCountInComment(commentID string) error
	DecrementLikeCountInComment(commentID string) error
	IsCommentLikedByUser(userID, commentID string) bool
	RemoveCommentLike(commentID, userID string) error
	// dislike in comment
	HandleCommentDislike(userID, commentID string) error
	CreateDislikeInComment(reaction *CommentDislike) error
	IncrementDislikeCountInComment(commentID string) error
	DecrementDislikeCountInComment(commentID string) error
	IsCommentDislikedByUser(userID, commentID string) bool
	RemoveCommentDislike(commentID, userID string) error
}
