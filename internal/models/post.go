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
	GetPostList(postPerPage, offset int) ([]Post, error) //  get all post
	Insert(post Post) error                              // insert new post
	Delete(post_id string) error                         // delete post
	IncrementLike(post_id string) error                  // increment like in post
	DecrementLike(post_id string) error                  // decrement like in post
	IncrementDisLike(post_id string) error               // increment like in post
	DecrementDisLike(post_id string) error               // decrement dislike in post
	IncrementComment(post_id string) error               // increment comment count in post
	DecrementComment(post_id string) error               // decrement comment count in post
	GetCategoryList() (*[]Categories, error)             // get all category/tags
	GetPostByName(username string) ([]Post, error)       // get list post specified user has create
	GetPostByLiked(user_id string) ([]Post, error)       // get list post specified user has liked
	GetPostByDislike(user_id string) ([]Post, error)     // get list post specified user has disliked
	GetPostByTags(tag string) ([]Post, error)            // get post list by tag
	GetPostByID(post_id string) (*Post, error)           // get post by id
	GetCountPost() (int, error)
}

type ResPostModel struct {
	Posts    Post
	Comments []Comment
	OwnerId  string
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

func ValidatePost(v *validator.Validator, post *Post) {
	ValidateTitle(v, post.Title)
	ValidateDescription(v, post.Description)
	ValidateTags(v, post.Tags)
}
