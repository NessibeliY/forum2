package reaction

import (
	"github.com/gofrs/uuid"

	"forum/internal/models"
)

type ReactionService struct { //nolint:revive
	ReactionRepo models.ReactionRepository
}

func NewReactionService(reaction models.ReactionRepository) *ReactionService {
	return &ReactionService{
		ReactionRepo: reaction,
	}
}

func (r *ReactionService) DeleteReaction(postID, userID string) error {
	err := r.DeleteLikeInPost(postID, userID)
	if err != nil {
		return err
	}

	err = r.DeleteDislike(postID, userID)
	if err != nil {
		return err
	}

	err = r.ReactionRepo.DeleteCommentByPostID(postID)
	if err != nil {
		return err
	}

	err = r.ReactionRepo.DeleteLikeInCommentByUserID(userID)
	if err != nil {
		return err
	}

	err = r.ReactionRepo.DeleteDisLikeInCommentByUserID(userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReactionService) CreateLikeInPost(like *models.Like) error {
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
func (r *ReactionService) LikeExistInPost(postID, userID string) bool {
	return r.ReactionRepo.LikeExistInPost(postID, userID)
}

func (r *ReactionService) DeleteLikeInPost(postID, userID string) error {
	err := r.ReactionRepo.DeleteLike(postID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReactionService) CreateDislikeInPost(dislike *models.Dislike) error {
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
func (r *ReactionService) DislikeExistInPost(postID, userID string) bool {
	return r.ReactionRepo.DislikeExistInPost(postID, userID)
}

func (r *ReactionService) DeleteDislike(postID, userID string) error {
	err := r.ReactionRepo.DeleteDislike(postID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReactionService) CreateCommentInPost(comment *models.Comment) error {
	id, err := uuid.NewV4()
	if err != nil {
		return models.ErrUUIDCreate
	}

	newComment := &models.Comment{
		CommentID:     id.String(),
		PostID:        comment.PostID,
		Author:        comment.Author,
		CommentText:   comment.CommentText,
		LikesCount:    0,
		DislikesCount: 0,
	}
	// insert new comment in DB
	err = r.ReactionRepo.InsertCommentInPost(newComment)
	if err != nil {
		return err
	}
	return nil
}

// get comments by ID
func (r *ReactionService) GetCommentsByID(postID string) ([]models.Comment, error) {
	comments, err := r.ReactionRepo.GetCommentsByID(postID)
	if err != nil {
		return nil, err
	}

	return comments, err
}

// delete comment
func (r *ReactionService) DeleteComment(commentID string) error {
	err := r.ReactionRepo.DeleteComment(commentID)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReactionService) CreateLikeInComment(reaction *models.CommentLike) error {
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
func (r *ReactionService) IncrementLikeInComment(commentID string) error {
	err := r.ReactionRepo.IncrementLikeInComment(commentID)
	if err != nil {
		return err
	}

	err = r.DecrementDislikeCountInComment(commentID)
	if err != nil {
		return err
	}

	return nil
}

// decrement like in comment
func (r *ReactionService) DecrementLikeCountInComment(commentID string) error {
	err := r.ReactionRepo.DecrementLikeCountInComment(commentID)
	if err != nil {
		return err
	}
	return nil
}

// check exist like in comment from current user
func (r *ReactionService) ExistLikeInComment(userID, commentID string) bool {
	return r.ReactionRepo.ExistLikeInComment(userID, commentID)
}

func (r *ReactionService) DeleteLikeInComment(commentID, userID string) error {
	err := r.ReactionRepo.DeleteLikeInComment(commentID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReactionService) CreateDislikeInComment(reaction *models.CommentDislike) error {
	id, err := uuid.NewV4()
	if err != nil {
		return models.ErrUUIDCreate
	}

	dislike := &models.CommentDislike{
		DislikeID: id.String(),
		CommentID: reaction.CommentID,
		UserID:    reaction.UserID,
	}

	err = r.ReactionRepo.InsertDislikeInComment(dislike)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReactionService) IncrementDislikeCountInComment(commentID string) error {
	err := r.ReactionRepo.IncrementDislikeCountInComment(commentID)
	if err != nil {
		return err
	}

	err = r.DecrementLikeCountInComment(commentID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReactionService) DecrementDislikeCountInComment(commentID string) error {
	err := r.ReactionRepo.DecrementDislikeCountInComment(commentID)
	if err != nil {
		return err
	}
	return err
}

func (r *ReactionService) ExistDisLikeInComment(userID, commentID string) bool {
	return r.ReactionRepo.ExistDisLikeInComment(userID, commentID)
}

func (r *ReactionService) DeleteDisLikeInComment(commentID, userID string) error {
	err := r.ReactionRepo.DeleteDisLikeInComment(commentID, userID)
	if err != nil {
		return err
	}
	return nil
}
