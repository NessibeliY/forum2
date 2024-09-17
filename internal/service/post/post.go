package post

import (
	"fmt"

	"forum/internal/models"
	"forum/internal/repository/post"

	"github.com/gofrs/uuid"
)

type PostService struct {
	PostRepo post.IPostRepo
}

func NewPostService(postRepo post.IPostRepo) *PostService {
	return &PostService{
		PostRepo: postRepo,
	}
}

type IPostService interface {
	CreatePost(post *models.Post) error                         // create new post
	Delete(post_id string) error                                // delete post by id
	IncrementLike(post_id string) error                         // increment like in post
	DecrementLike(post_id string) error                         // decrement like in post
	IncrementDisLike(post_id string) error                      // increment dislike in post
	DecrementDisLike(post_id string) error                      // decrement dislike in post
	IncrementComment(post_id string) error                      // increment comment count in post
	DecrementComment(post_id string) error                      // decrement comment count in post
	GetPostList(postPerPage, offset int) ([]models.Post, error) // get all post
	GetCategoryList() (*[]models.Categories, error)             // get all category
	GetPostByName(username string) ([]models.Post, error)       // get post specified has user create
	GetPostByLiked(user_id string) ([]models.Post, error)       // get list post specified user has liked
	GetPostByDisLike(user_id string) ([]models.Post, error)     // get list post specified user has disliked
	GetPostByTags(tag string) ([]models.Post, error)
	GetPostByID(post_id string) (*models.Post, error) // get post by id
	GetCountPost() (int, error)
}

func (p *PostService) GetPostByTags(tag string) ([]models.Post, error) {
	posts, err := p.PostRepo.GetPostByTags(tag)
	if err != nil {
		return nil, err
	}

	return posts, err
}

// create new post
func (p *PostService) CreatePost(post *models.Post) error {
	postID, err := uuid.NewV4()
	if err != nil {
		return models.ErrUUIDCreate
	}

	newPost := &models.Post{
		PostID:      postID.String(),
		UserID:      post.UserID,
		Author:      post.Author,
		Title:       post.Title,
		Description: post.Description,
		Likes:       0,
		Dislikes:    0,
		Comments:    0,
		Tags:        post.Tags,
	}

	err = p.PostRepo.Insert(*newPost)
	if err != nil {
		return err
	}
	return nil
}

// delete post
func (p *PostService) Delete(post_id string) error {
	err := p.PostRepo.Delete(post_id)
	if err != nil {
		return err
	}

	return nil
}

// increment like in post
func (p *PostService) IncrementLike(post_id string) error {
	err := p.PostRepo.IncrementLike(post_id)
	if err != nil {
		return err
	}
	err = p.DecrementDisLike(post_id)
	if err != nil {
		return err
	}
	return nil
}

// decrement like  in post
func (p *PostService) DecrementLike(post_id string) error {
	err := p.PostRepo.DecrementLike(post_id)
	if err != nil {
		return err
	}

	return nil
}

// increment dislike  in post
func (p *PostService) IncrementDisLike(post_id string) error {
	err := p.PostRepo.IncrementDisLike(post_id)
	if err != nil {
		return err
	}
	err = p.DecrementLike(post_id)
	if err != nil {
		return err
	}

	return nil
}

// decrement dislike in post
func (p *PostService) DecrementDisLike(post_id string) error {
	err := p.PostRepo.DecrementDisLike(post_id)
	if err != nil {
		return err
	}

	return nil
}

// increment comment count in post
func (p *PostService) IncrementComment(post_id string) error {
	err := p.PostRepo.IncrementComment(post_id)
	if err != nil {
		return err
	}
	return nil
}

// decrement comment count in post
func (p *PostService) DecrementComment(post_id string) error {
	err := p.PostRepo.DecrementComment(post_id)
	if err != nil {
		return err
	}
	return nil
}

// get post list
func (p *PostService) GetPostList(postPerPage, offset int) ([]models.Post, error) {
	posts, err := p.PostRepo.GetPostList(postPerPage, offset)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

// get category list
func (p *PostService) GetCategoryList() (*[]models.Categories, error) {
	posts, err := p.PostRepo.GetCategoryList()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

// get post by name
func (p *PostService) GetPostByName(username string) ([]models.Post, error) {
	post, err := p.PostRepo.GetPostByName(username)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return post, nil
}

// get post liked user
func (p *PostService) GetPostByLiked(user_id string) ([]models.Post, error) {
	posts, err := p.PostRepo.GetPostByLiked(user_id)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return posts, nil
}

// get post disliked user
func (p *PostService) GetPostByDisLike(user_id string) ([]models.Post, error) {
	posts, err := p.PostRepo.GetPostByDislike(user_id)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return posts, nil
}

func (p *PostService) GetPostByID(post_id string) (*models.Post, error) {
	post, err := p.PostRepo.GetPostByID(post_id)
	if err != nil {
		return nil, err
	}

	return post, err
}

func (p *PostService) GetCountPost() (int, error) {
	count, err := p.PostRepo.GetCountPost()
	if err != nil {
		return 0, err
	}
	return count, nil
}
