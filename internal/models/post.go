package models

import (
	"time"

	"forum/internal/validator"
)

type Post struct {
	PostID      string
	UserID      string
	Author      string
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Likes       int
	Dislikes    int
	Comments    int
	Tags        []string
	IsLike      bool
	IsDisLike   bool
}

type PostRepository interface {
	GetAllPosts(postPerPage, offset int) ([]Post, error)
	AddPost(post Post) error // insert new post
	DeletePostByPostID(post_id string) error
	IncrementLikeCount(post_id string) error
	DecrementLikeCount(post_id string) error
	IncrementDislikeCount(post_id string) error
	DecrementDislikeCount(post_id string) error
	IncrementCommentCount(post_id string) error
	DecrementCommentCount(post_id string) error
	GetAllCategories() (*[]Categories, error)
	GetPostsByUsername(username string) ([]Post, error)
	GetPostsLikedByUser(user_id string) ([]Post, error)
	GetPostsDislikedByUser(user_id string) ([]Post, error)
	GetPostsByCategory(tag string) ([]Post, error)
	GetPostByPostID(post_id string) (*Post, error)
	GetPostsCount() (int, error)
}

type PostService interface {
	CreatePost(createPostRequest *CreatePostRequest) error
	DeletePostByPostID(post_id string) error
	IncrementLikeCount(post_id string) error
	DecrementLikeCount(post_id string) error
	IncrementDislikeCount(post_id string) error
	DecrementDislikeCount(post_id string) error
	IncrementCommentCount(post_id string) error
	DecrementCommentCount(post_id string) error
	GetAllPosts(postPerPage, offset int) ([]Post, error)
	GetAllCategories() (*[]Categories, error)
	GetPostsByUsername(username string) ([]Post, error)
	GetPostsLikedByUser(user_id string) ([]Post, error)
	GetPostsDislikedByUser(user_id string) ([]Post, error)
	GetPostsByCategory(tag string) ([]Post, error)
	GetPostByPostID(post_id string) (*Post, error)
	GetPostsCount() (int, error)
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
	OwnerId  string
}
