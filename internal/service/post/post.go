package post

import (
	"fmt"

	"github.com/gofrs/uuid"

	"forum/internal/models"
)

type PostService struct { //nolint:revive
	PostRepo models.PostRepository
}

func NewPostService(postRepo models.PostRepository) *PostService {
	return &PostService{
		PostRepo: postRepo,
	}
}

func (p *PostService) GetPostsByCategory(tag string) ([]models.Post, error) {
	return p.PostRepo.GetPostsByCategory(tag) //nolint:wrapcheck
}

func (p *PostService) CreatePost(createPostRequest *models.CreatePostRequest) error {
	PostID, err := uuid.NewV4()
	if err != nil {
		return models.ErrUUIDCreate
	}

	newPost := &models.Post{
		PostID:        PostID.String(),
		UserID:        createPostRequest.UserID,
		Author:        createPostRequest.Author,
		Title:         createPostRequest.Title,
		Description:   createPostRequest.Description,
		LikesCount:    0,
		DislikesCount: 0,
		Comments:      0,
		Tags:          createPostRequest.Tags,
	}

	err = p.PostRepo.AddPost(*newPost)
	if err != nil {
		return fmt.Errorf("add post: %w", err)
	}
	return nil
}

func (p *PostService) DeletePostByPostID(postID string) error {
	return p.PostRepo.DeletePostByPostID(postID) //nolint:wrapcheck
}

func (p *PostService) IncrementPostLikeCount(postID string) error {
	err := p.PostRepo.IncrementPostLikeCount(postID)
	if err != nil {
		return fmt.Errorf("increment post like count: %w", err)
	}
	err = p.DecrementPostDislikeCount(postID)
	if err != nil {
		return fmt.Errorf("decrement post dislike count: %w", err)
	}
	return nil
}

func (p *PostService) DecrementPostLikeCount(postID string) error {
	return p.PostRepo.DecrementPostLikeCount(postID) //nolint:wrapcheck
}

func (p *PostService) IncrementPostDislikeCount(postID string) error {
	err := p.PostRepo.IncrementPostDislikeCount(postID)
	if err != nil {
		return fmt.Errorf("increment post dislike count: %w", err)
	}
	err = p.DecrementPostLikeCount(postID)
	if err != nil {
		return fmt.Errorf("decrement post like count: %w", err)
	}

	return nil
}

func (p *PostService) DecrementPostDislikeCount(postID string) error {
	return p.PostRepo.DecrementPostDislikeCount(postID) //nolint:wrapcheck
}

func (p *PostService) IncrementCommentCount(postID string) error {
	return p.PostRepo.IncrementCommentCount(postID) //nolint:wrapcheck
}

func (p *PostService) DecrementCommentCount(postID string) error {
	return p.PostRepo.DecrementCommentCount(postID) //nolint:wrapcheck
}

func (p *PostService) GetAllPosts(postPerPage, offset int) ([]models.Post, error) {
	return p.PostRepo.GetAllPosts(postPerPage, offset) //nolint:wrapcheck
}

func (p *PostService) GetAllCategories() (*[]models.Categories, error) {
	return p.PostRepo.GetAllCategories() //nolint:wrapcheck
}

func (p *PostService) GetPostsByUsername(username string) ([]models.Post, error) {
	return p.PostRepo.GetPostsByUsername(username) //nolint:wrapcheck
}

func (p *PostService) GetPostsLikedByUser(userID string) ([]models.Post, error) {
	return p.PostRepo.GetPostsLikedByUser(userID) //nolint:wrapcheck
}

func (p *PostService) GetPostsDislikedByUser(userID string) ([]models.Post, error) {
	return p.PostRepo.GetPostsDislikedByUser(userID) //nolint:wrapcheck
}

func (p *PostService) GetPostByPostID(postID string) (models.Post, error) {
	return p.PostRepo.GetPostByPostID(postID) //nolint:wrapcheck
}

func (p *PostService) GetPostsCount() (int, error) {
	return p.PostRepo.GetPostsCount() //nolint:wrapcheck
}

func (p *PostService) PopulatePostData(postID string, data *models.Login) error {
	categories, err := p.GetAllCategories()
	if err != nil {
		return fmt.Errorf("get all categories: %w", err)
	}
	data.Categories = *categories

	post, err := p.GetPostByPostID(postID)
	if err != nil {
		return fmt.Errorf("get post by post id: %w", err)
	}
	data.Post = post

	return nil
}

func (p *PostService) CheckPostByID(postID string) (bool, error) {
	exist, err := p.PostRepo.CheckPostByID(postID)
	if err != nil {
		return exist, err //nolint:wrapcheck
	}

	return exist, nil
}
