package models

type Login struct {
	Id                 string
	IsAuth             bool
	UserName           string
	Categories         []Categories
	Posts              []ResPostModel
	CurrentPage        string
	CommentError       ErrorComment
	Error              ErrorMessage
	ShowCreatePostForm bool
	Post               Post
	Comment            []Comment
	Page               int
	TotalPages         int
}
