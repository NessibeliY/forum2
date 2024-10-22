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
	OwnerId       string
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
	InsertLikePost(like *Like) error
	LikeExistInPost(postID, userID string) bool
	DeleteLike(postID, userID string) error
	// dislike in post
	InsertDisLikePost(disLike *Dislike) error      // insert new dislike
	DislikeExistInPost(postID, userID string) bool // check dislike current post from user
	DeleteDislike(postID, userID string) error     // delete dislike
	// comments in post
	InsertCommentInPost(comment *Comment) error   // insert new comment
	GetCommentsByID(id string) ([]Comment, error) // get comment list by postID
	DeleteComment(commentID string) error         // delete comment
	DeleteCommentByPostID(postID string) error    // delete comment by POST ID
	// like in comment
	InsertLikeInComment(reaction *CommentLike) error    // insert new like in comment
	IncrementLikeInComment(commentID string) error      // increment like in comment
	DecrementLikeCountInComment(commentID string) error // decrement like in comment
	ExistLikeInComment(userID, commentID string) bool   // check exist like in comment from current user
	DeleteLikeInComment(commentID, userID string) error // delete dislike in comment
	DeleteLikeInCommentByUserID(userID string) error    // delete Like In Comment By Post ID
	// dislike in comment
	InsertDislikeInComment(reaction *CommentDislike) error // insert dislike in comment
	IncrementDislikeCountInComment(commentID string) error // increment dislike in comment
	DecrementDislikeCountInComment(commentID string) error // decrement dislike in comment
	ExistDisLikeInComment(userID, commentID string) bool   // check exist dislike in comment from current user
	DeleteDisLikeInComment(commentID, userID string) error
	DeleteDisLikeInCommentByUserID(userID string) error
}

type ReactionService interface {
	// delete reaction
	DeleteReaction(postID, userID string) error
	// like in post
	CreateLikeInPost(like *Like) error            // create like in post
	LikeExistInPost(postID, userID string) bool   // check like exist in post
	DeleteLikeInPost(postID, userID string) error // delete like in post
	// dislike in post
	CreateDislikeInPost(dislike *Dislike) error    // create dislike in post
	DislikeExistInPost(postID, userID string) bool // check dislike exist in post
	DeleteDislike(postID, userID string) error     // delete dislike in post
	// comment
	CreateCommentInPost(comment *Comment) error       // create comment in post
	GetCommentsByID(postID string) ([]Comment, error) // get comment list by POST ID
	DeleteComment(commentID string) error
	// like in comment
	CreateLikeInComment(reaction *CommentLike) error    // create  like in comment
	IncrementLikeInComment(commentID string) error      // increment like in comment
	DecrementLikeCountInComment(commentID string) error // decrement like in comment
	ExistLikeInComment(userID, commentID string) bool   // check exist like in comment from current user
	DeleteLikeInComment(commentID, userID string) error // delete like in comment
	// dislike in comment
	CreateDislikeInComment(reaction *CommentDislike) error
	IncrementDislikeCountInComment(commentID string) error // increment dislike in comment
	DecrementDislikeCountInComment(commentID string) error // decrement dislike in comment
	ExistDisLikeInComment(userID, commentID string) bool   // check exist dislike in comment from current user
	DeleteDisLikeInComment(commentID, userID string) error
}
