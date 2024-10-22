package models

import (
	"time"

	"forum/internal/validator"
)

type Post struct {
	PostID        string
	UserID        string
	Author        string
	Title         string
	Description   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LikesCount    int
	DislikesCount int
	Comments      int
	Tags          []string
	IsLike        bool // в базе не хранится
	IsDisLike     bool // в базе не хранится
}

type PostRepository interface {
	GetAllPosts(postPerPage, offset int) ([]Post, error)
	AddPost(post Post) error // insert new post
	DeletePostByPostID(postID string) error
	IncrementPostLikeCount(postID string) error
	DecrementPostLikeCount(postID string) error
	IncrementPostDislikeCount(postID string) error
	DecrementPostDislikeCount(postID string) error
	IncrementCommentCount(postID string) error
	DecrementCommentCount(postID string) error
	GetAllCategories() (*[]Categories, error)
	GetPostsByUsername(username string) ([]Post, error)
	GetPostsLikedByUser(userID string) ([]Post, error)
	GetPostsDislikedByUser(userID string) ([]Post, error)
	GetPostsByCategory(tag string) ([]Post, error)
	GetPostByPostID(postID string) (*Post, error)
	GetPostsCount() (int, error)
}

type PostService interface {
	CreatePost(createPostRequest *CreatePostRequest) error
	DeletePostByPostID(postID string) error
	IncrementPostLikeCount(postID string) error
	DecrementPostLikeCount(postID string) error
	IncrementPostDislikeCount(postID string) error
	DecrementPostDislikeCount(postID string) error
	IncrementCommentCount(postID string) error
	DecrementCommentCount(postID string) error
	GetAllPosts(postPerPage, offset int) ([]Post, error)
	GetAllCategories() (*[]Categories, error)
	GetPostsByUsername(username string) ([]Post, error)
	GetPostsLikedByUser(userID string) ([]Post, error)
	GetPostsDislikedByUser(userID string) ([]Post, error)
	GetPostsByCategory(tag string) ([]Post, error)
	GetPostByPostID(postID string) (*Post, error)
	GetPostsCount() (int, error)
	PopulatePostData(postID string, data *Login) error
}

type Posts struct {
	Posts []Post
}

type Categories struct {
	CategoryID   int
	CategoryName string
}

type ErrorComment struct {
	EmptyCommentText string
}

func ValidateTitle(v *validator.Validator, title string) {
	v.Check(title != "", "title", "Title is required")
}

func ValidateDescription(v *validator.Validator, description string) {
	v.Check(description != "", "description", "Description is required")
}

func ValidateTags(v *validator.Validator, tags []string) {
	v.Check(len(tags) > 0, "tags", "Tags is required")
}

func ValidateCreatePostRequest(v *validator.Validator, createPostRequest *CreatePostRequest) {
	ValidateTitle(v, createPostRequest.Title)
	ValidateDescription(v, createPostRequest.Description)
	ValidateTags(v, createPostRequest.Tags)
}

type CreatePostRequest struct {
	UserID      string
	Author      string
	Title       string
	Description string
	Tags        []string
}

type UserPostsResponse struct {
	Posts    Post
	Comments []Comment
	OwnerID  string
}
