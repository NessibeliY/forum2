package reaction

import (
	"forum/internal/models"
	"forum/internal/repository/reaction"

	"github.com/gofrs/uuid"
)

type ReactionService struct {
	ReactionRepo reaction.IReactionRepo
}

func NewReactionService(reaction reaction.IReactionRepo) *ReactionService {
	return &ReactionService{
		ReactionRepo: reaction,
	}
}

type IReactionService interface {
	// delete reaction
	DeleteReaction(post_id, user_id string) error
	// like in post
	CreateLikeInPost(like *models.Likes) error      // create like in post
	LikeExistInPost(post_id, user_id string) bool   // check like exist in post
	DeleteLikeInPost(post_id, user_id string) error // delete like in post
	// dislike in post
	CreateDislikeInPost(dislike *models.DisLike) error // create dislike in post
	DisLikeExistInPost(post_id, user_id string) bool   // check dislike exist in post
	DeleteDisLike(post_id, user_id string) error       // delete dislike in post
	// comment
	CreateCommentInPost(comment *models.Comment) error        // create comment in post
	GetCommentsByID(post_id string) ([]models.Comment, error) // get comment list by POST ID
	DeleteComment(comment_id string) error
	// like in comment
	CreateLikeInComment(reaction *models.CommentLike) error // create  like in comment
	IncrementLikeInComment(comment_id string) error         // increment like in comment
	DecrementLikeInComment(comment_id string) error         // decrement like in comment
	ExistLikeInComment(user_id, comment_id string) bool     // check exist like in comment from current user
	DeleteLikeInComment(comment_id, user_id string) error   // delete like in comment
	// dislike in comment
	CreateDisLikeInComment(reaction *models.CommentDisLike) error
	IncrementDisLikeInComment(comment_id string) error     // increment dislike in comment
	DecrementDisLikeInComment(comment_id string) error     // decrement dislike in comment
	ExistDisLikeInComment(user_id, comment_id string) bool // check exist dislike in comment from current user
	DeleteDisLikeInComment(comment_id, user_id string) error
}

func (r *ReactionService) DeleteReaction(post_id, user_id string) error {
	err := r.DeleteLikeInPost(post_id, user_id)
	if err != nil {
		return err
	}

	err = r.DeleteDisLike(post_id, user_id)
	if err != nil {
		return err
	}

	err = r.ReactionRepo.DeleteCommentByPostID(post_id)
	if err != nil {
		return err
	}

	err = r.ReactionRepo.DeleteLikeInCommentByUserID(user_id)
	if err != nil {
		return err
	}

	err = r.ReactionRepo.DeleteDisLikeInCommentByUserID(user_id)
	if err != nil {
		return err
	}

	return nil
}

// LIKE
// create new like
func (r *ReactionService) CreateLikeInPost(like *models.Likes) error {
	// generate unique ID
	id, err := uuid.NewV4()
	if err != nil {
		return models.ErrUUIDCreate
	}
	// new like
	newReaction := &models.Likes{
		LikeId: id.String(),
		PostId: like.PostId,
		UserId: like.UserId,
	}

	// insert new like in DB
	err = r.ReactionRepo.InsertLikePost(newReaction)
	if err != nil {
		return err
	}

	return nil
}

// check exist like in post by POST_ID and USER_ID
func (r *ReactionService) LikeExistInPost(post_id, user_id string) bool {
	return r.ReactionRepo.LikeExistInPost(post_id, user_id)
}

// delete like
func (r *ReactionService) DeleteLikeInPost(post_id, user_id string) error {
	err := r.ReactionRepo.DeleteLike(post_id, user_id)
	if err != nil {
		return err
	}
	return nil
}

// DISLIKE
// create new dislike
func (r *ReactionService) CreateDislikeInPost(dislike *models.DisLike) error {
	// generate unique ID
	id, err := uuid.NewV4()
	if err != nil {
		return models.ErrUUIDCreate
	}

	// new dislike
	newReaction := &models.DisLike{
		DisLikeId: id.String(),
		PostId:    dislike.PostId,
		UserId:    dislike.UserId,
	}

	// insert new dislike in DB
	err = r.ReactionRepo.InsertDisLikePost(newReaction)
	if err != nil {
		return err
	}

	return nil
}

// check exist dislike in post by POST_ID and USER_ID
func (r *ReactionService) DisLikeExistInPost(post_id, user_id string) bool {
	return r.ReactionRepo.DisLikeExistInPost(post_id, user_id)
}

// delete dislike
func (r *ReactionService) DeleteDisLike(post_id, user_id string) error {
	err := r.ReactionRepo.DeleteDisLike(post_id, user_id)
	if err != nil {
		return err
	}
	return nil
}

// COMMENT
// create new comment
func (r *ReactionService) CreateCommentInPost(comment *models.Comment) error {
	// generate unique ID
	id, err := uuid.NewV4()
	if err != nil {
		return models.ErrUUIDCreate
	}

	// new comment
	newComment := &models.Comment{
		CommentId:   id.String(),
		PostId:      comment.PostId,
		Author:      comment.Author,
		CommentText: comment.CommentText,
		Likes:       0,
		DisLikes:    0,
	}
	// insert new comment in DB
	err = r.ReactionRepo.InsertCommentInPost(newComment)
	if err != nil {
		return err
	}
	return nil
}

// get comments by ID
func (r *ReactionService) GetCommentsByID(post_id string) ([]models.Comment, error) {
	comments, err := r.ReactionRepo.GetCommentsByID(post_id)
	if err != nil {
		return nil, err
	}

	return comments, err
}

// delete comment
func (r *ReactionService) DeleteComment(comment_id string) error {
	err := r.ReactionRepo.DeleteComment(comment_id)
	if err != nil {
		return err
	}

	return nil
}

// LIKE IN COMMENT
// create like in comment
func (r *ReactionService) CreateLikeInComment(reaction *models.CommentLike) error {
	// generate unique ID
	id, err := uuid.NewV4()
	if err != nil {
		return models.ErrUUIDCreate
	}
	commentLike := &models.CommentLike{
		LikeId:    id.String(),
		CommentId: reaction.CommentId,
		UserId:    reaction.UserId,
	}

	err = r.ReactionRepo.InsertLikeInComment(commentLike)
	if err != nil {
		return err
	}
	return nil
}

// increment like in comment
func (r *ReactionService) IncrementLikeInComment(comment_id string) error {
	err := r.ReactionRepo.IncrementLikeInComment(comment_id)
	if err != nil {
		return err
	}

	err = r.DecrementDisLikeInComment(comment_id)
	if err != nil {
		return err
	}

	return nil
}

// decrement like in comment
func (r *ReactionService) DecrementLikeInComment(comment_id string) error {
	err := r.ReactionRepo.DecrementLikeInComment(comment_id)
	if err != nil {
		return err
	}
	return nil
}

// check exist like in comment from current user
func (r *ReactionService) ExistLikeInComment(user_id, comment_id string) bool {
	return r.ReactionRepo.ExistLikeInComment(user_id, comment_id)
}

// delete like in comment
func (r *ReactionService) DeleteLikeInComment(comment_id, user_id string) error {
	err := r.ReactionRepo.DeleteLikeInComment(comment_id, user_id)
	if err != nil {
		return err
	}

	return nil
}

// DISLIKE IN COMMENT
func (r *ReactionService) CreateDisLikeInComment(reaction *models.CommentDisLike) error {
	// generate unique ID
	id, err := uuid.NewV4()
	if err != nil {
		return models.ErrUUIDCreate
	}

	dislike := &models.CommentDisLike{
		DisLikeId: id.String(),
		CommentId: reaction.CommentId,
		UserId:    reaction.UserId,
	}

	err = r.ReactionRepo.InsertDisLikeInComment(dislike)
	if err != nil {
		return err
	}
	return nil
}

// increment dislike in comment
func (r *ReactionService) IncrementDisLikeInComment(comment_id string) error {
	err := r.ReactionRepo.IncrementDisLikeInComment(comment_id)
	if err != nil {
		return err
	}

	err = r.DecrementLikeInComment(comment_id)
	if err != nil {
		return err
	}
	return nil
}

// decrement dislike in comment
func (r *ReactionService) DecrementDisLikeInComment(comment_id string) error {
	err := r.ReactionRepo.DecrementDisLikeInComment(comment_id)
	if err != nil {
		return err
	}
	return err
}

// check exist dislike in comment from current user
func (r *ReactionService) ExistDisLikeInComment(user_id, comment_id string) bool {
	return r.ReactionRepo.ExistDisLikeInComment(user_id, comment_id)
}

// delete dislike in comment
func (r *ReactionService) DeleteDisLikeInComment(comment_id, user_id string) error {
	err := r.ReactionRepo.DeleteDisLikeInComment(comment_id, user_id)
	if err != nil {
		return err
	}
	return nil
}
