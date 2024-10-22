package models

type Login struct {
	Id                 string
	IsAuth             bool
	Username           string
	Categories         []Categories
	Posts              []UserPostsResponse
	CurrentPage        string
	CommentError       ErrorComment
	ErrorMessages      ErrorMessage
	ShowCreatePostForm bool
	Post               Post
	Comment            []Comment
	Page               int
	TotalPages         int
}
