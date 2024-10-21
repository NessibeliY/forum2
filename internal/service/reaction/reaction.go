package reaction

import (
	"github.com/gofrs/uuid"

	"forum/internal/models"
)

type ReactionService struct {
	ReactionRepo models.ReactionRepository
}

func NewReactionService(reaction models.ReactionRepository) *ReactionService {
	return &ReactionService{
		ReactionRepo: reaction,
	}
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
func (r *ReactionService) CreateLikeInPost(like *models.Like) error {
	// generate unique ID
	id, err := uuid.NewV4()
	if err != nil {
		return models.ErrUUIDCreate
	}
	// new like
	newReaction := &models.Like{
		LikeID: id.String(),
		PostID: like.PostID,
		UserID: like.UserID,
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
func (r *ReactionService) CreateDislikeInPost(dislike *models.Dislike) error {
	// generate unique ID
	id, err := uuid.NewV4()
	if err != nil {
		return models.ErrUUIDCreate
	}

	// new dislike
	newReaction := &models.Dislike{
		DislikeID: id.String(),
		PostID:    dislike.PostID,
		UserID:    dislike.UserID,
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
		CommentID:    id.String(),
		PostID:       comment.PostID,
		Author:       comment.Author,
		CommentText:  comment.CommentText,
		LikeCount:    0,
		DislikeCount: 0,
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
		LikeID:    id.String(),
		CommentID: reaction.CommentID,
		UserID:    reaction.UserID,
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
func (r *ReactionService) CreateDisLikeInComment(reaction *models.CommentDislike) error {
	// generate unique ID
	id, err := uuid.NewV4()
	if err != nil {
		return models.ErrUUIDCreate
	}

	dislike := &models.CommentDislike{
		DislikeID: id.String(),
		CommentID: reaction.CommentID,
		UserID:    reaction.UserID,
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
