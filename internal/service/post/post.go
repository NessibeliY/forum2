package post

import (
	"fmt"

	"github.com/gofrs/uuid"

	"forum/internal/models"
)

type PostService struct {
	PostRepo models.PostRepository
}

func NewPostService(postRepo models.PostRepository) *PostService {
	return &PostService{
		PostRepo: postRepo,
	}
}

func (p *PostService) GetPostsByCategory(tag string) ([]models.Post, error) {
	posts, err := p.PostRepo.GetPostsByCategory(tag)
	if err != nil {
		return nil, err
	}

	return posts, err
}

func (p *PostService) CreatePost(createPostRequest *models.CreatePostRequest) error {
	postID, err := uuid.NewV4()
	if err != nil {
		return models.ErrUUIDCreate
	}

	newPost := &models.Post{
		PostID:      postID.String(),
		UserID:      createPostRequest.UserID,
		Author:      createPostRequest.Author,
		Title:       createPostRequest.Title,
		Description: createPostRequest.Description,
		Likes:       0,
		Dislikes:    0,
		Comments:    0,
		Tags:        createPostRequest.Tags,
	}

	err = p.PostRepo.AddPost(*newPost)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostService) DeletePostByPostID(post_id string) error {
	err := p.PostRepo.DeletePostByPostID(post_id)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostService) IncrementLikeCount(post_id string) error {
	err := p.PostRepo.IncrementLikeCount(post_id)
	if err != nil {
		return err
	}
	err = p.DecrementDislikeCount(post_id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostService) DecrementLikeCount(post_id string) error {
	err := p.PostRepo.DecrementLikeCount(post_id)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostService) IncrementDislikeCount(post_id string) error {
	err := p.PostRepo.IncrementDislikeCount(post_id)
	if err != nil {
		return err
	}
	err = p.DecrementLikeCount(post_id)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostService) DecrementDislikeCount(post_id string) error {
	err := p.PostRepo.DecrementDislikeCount(post_id)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostService) IncrementCommentCount(post_id string) error {
	err := p.PostRepo.IncrementCommentCount(post_id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostService) DecrementCommentCount(post_id string) error {
	err := p.PostRepo.DecrementCommentCount(post_id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostService) GetAllPosts(postPerPage, offset int) ([]models.Post, error) {
	posts, err := p.PostRepo.GetAllPosts(postPerPage, offset)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *PostService) GetAllCategories() (*[]models.Categories, error) {
	posts, err := p.PostRepo.GetAllCategories()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *PostService) GetPostsByUsername(username string) ([]models.Post, error) {
	post, err := p.PostRepo.GetPostsByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return post, nil
}

func (p *PostService) GetPostsLikedByUser(user_id string) ([]models.Post, error) {
	posts, err := p.PostRepo.GetPostsLikedByUser(user_id)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return posts, nil
}

func (p *PostService) GetPostsDislikedByUser(user_id string) ([]models.Post, error) {
	posts, err := p.PostRepo.GetPostsDislikedByUser(user_id)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return posts, nil
}

func (p *PostService) GetPostByPostID(post_id string) (*models.Post, error) {
	post, err := p.PostRepo.GetPostByPostID(post_id)
	if err != nil {
		return nil, err
	}

	return post, err
}

func (p *PostService) GetPostsCount() (int, error) {
	count, err := p.PostRepo.GetPostsCount()
	if err != nil {
		return 0, err
	}
	return count, nil
}
