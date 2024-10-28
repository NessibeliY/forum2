package reaction

import (
	"fmt"

	"github.com/gofrs/uuid"

	"forum/internal/models"
)

type ReactionService struct { //nolint:revive
	ReactionRepo models.ReactionRepository
}

func NewReactionService(reactionRepo models.ReactionRepository) *ReactionService {
	return &ReactionService{
		ReactionRepo: reactionRepo,
	}
}

func (r *ReactionService) HandleCommentLike(userID, commentID string) error {
	existLike := r.IsCommentLikedByUser(userID, commentID)

	if !existLike {
		like := &models.CommentLike{
			CommentID: commentID,
			UserID:    userID,
		}

		err := r.CreateLikeInComment(like)
		if err != nil {
			return fmt.Errorf("create like in comment: %w", err)
		}

		err = r.IncrementLikeCountInComment(commentID)
		if err != nil {
			return fmt.Errorf("increment like count in comment: %w", err)
		}

		err = r.RemoveCommentDislike(commentID, userID)
		if err != nil {
			return fmt.Errorf("remove comment dislike: %w", err)
		}
	}

	return nil
}

func (r *ReactionService) GetCommentsWithReactions(postID, username, userID string) []models.Comment {
	comments, err := r.GetCommentsByPostID(postID)
	if err != nil {
		return []models.Comment{}
	}

	for i := range comments {
		comments[i].OwnerID = username
		comments[i].IsLiked = r.IsCommentLikedByUser(userID, comments[i].CommentID)
		comments[i].DisLiked = r.IsCommentDislikedByUser(userID, comments[i].CommentID)
	}

	return comments
}

func (r *ReactionService) HandlePostLike(userID, postID string) error {
	isLiked := r.ReactionRepo.IsPostLikedByUser(postID, userID)
	if isLiked {
		return nil
	}

	id, err := uuid.NewV4()
	if err != nil {
		return models.ErrUUIDCreate
	}

	newLike := models.Like{
		LikeID: id.String(),
		PostID: postID,
		UserID: userID,
	}

	err = r.ReactionRepo.AddPostLike(newLike)
	if err != nil {
		return fmt.Errorf("add post like: %w", err)
	}

	err = r.ReactionRepo.RemovePostDislike(postID, userID)
	if err != nil {
		return fmt.Errorf("remove post dislike: %w", err)
	}

	return nil
}

func (r *ReactionService) AddPostLike(newLike models.Like) error {
	return r.ReactionRepo.AddPostLike(newLike) //nolint:wrapcheck
}

func (r *ReactionService) RemovePostLike(postID, userID string) error {
	return r.ReactionRepo.RemovePostLike(postID, userID) //nolint:wrapcheck
}

func (r *ReactionService) HandlePostDislike(userID, postID string) error {
	isDisliked := r.ReactionRepo.IsPostDislikedByUser(postID, userID)
	if isDisliked {
		return nil
	}

	id, err := uuid.NewV4()
	if err != nil {
		return models.ErrUUIDCreate
	}

	newDislike := models.Dislike{
		DislikeID: id.String(),
		PostID:    postID,
		UserID:    userID,
	}

	err = r.ReactionRepo.AddPostDislike(newDislike)
	if err != nil {
		return fmt.Errorf("add post dislike: %w", err)
	}

	err = r.ReactionRepo.RemovePostLike(postID, userID)
	if err != nil {
		return fmt.Errorf("remove post like: %w", err)
	}

	return nil
}

func (r *ReactionService) AddPostDislike(newDislike models.Dislike) error {
	return r.ReactionRepo.AddPostDislike(newDislike) //nolint:wrapcheck
}

func (r *ReactionService) DeleteReaction(postID, userID string) error {
	err := r.RemovePostLike(postID, userID)
	if err != nil {
		return fmt.Errorf("delete like in post: %w", err)
	}

	err = r.RemovePostDislike(postID, userID)
	if err != nil {
		return fmt.Errorf("remove post dislike: %w", err)
	}

	err = r.ReactionRepo.RemoveCommentByPostID(postID)
	if err != nil {
		return fmt.Errorf("remove comment dislike: %w", err)
	}

	err = r.ReactionRepo.RemoveAllCommentLikesByUser(userID)
	if err != nil {
		return fmt.Errorf("remove all comment likes: %w", err)
	}

	err = r.ReactionRepo.RemoveAllCommentDislikesByUser(userID)
	if err != nil {
		return fmt.Errorf("remove all comment dislikes: %w", err)
	}

	return nil
}

func (r *ReactionService) IsPostLikedByUser(postID, userID string) bool {
	return r.ReactionRepo.IsPostLikedByUser(postID, userID) //nolint:wrapcheck
}

func (r *ReactionService) IsPostDislikedByUser(postID, userID string) bool {
	return r.ReactionRepo.IsPostDislikedByUser(postID, userID) //nolint:wrapcheck
}

func (r *ReactionService) RemovePostDislike(postID, userID string) error {
	return r.ReactionRepo.RemovePostDislike(postID, userID) //nolint:wrapcheck
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

	err = r.ReactionRepo.AddComment(newComment)
	if err != nil {
		return fmt.Errorf("add comment: %w", err)
	}
	return nil
}

func (r *ReactionService) GetCommentsByPostID(postID string) ([]models.Comment, error) {
	return r.ReactionRepo.GetCommentsByPostID(postID) //nolint:wrapcheck
}

func (r *ReactionService) RemoveCommentByCommentID(commentID string) error {
	return r.ReactionRepo.RemoveCommentByCommentID(commentID) //nolint:wrapcheck
}

func (r *ReactionService) CommentExists(commentID string) bool {
	return r.ReactionRepo.CommentExistsByCommentID(commentID)
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

	err = r.ReactionRepo.AddCommentLike(commentLike)
	if err != nil {
		return fmt.Errorf("add comment like: %w", err)
	}
	return nil
}

func (r *ReactionService) IncrementLikeCountInComment(commentID string) error {
	err := r.ReactionRepo.IncrementLikeCountInComment(commentID)
	if err != nil {
		return fmt.Errorf("increment like count in comment: %w", err)
	}

	err = r.DecrementDislikeCountInComment(commentID)
	if err != nil {
		return fmt.Errorf("increment dislike count in comment: %w", err)
	}

	return nil
}

func (r *ReactionService) DecrementLikeCountInComment(commentID string) error {
	return r.ReactionRepo.DecrementLikeCountInComment(commentID) //nolint:wrapcheck
}

func (r *ReactionService) IsCommentLikedByUser(userID, commentID string) bool {
	return r.ReactionRepo.IsCommentLikedByUser(userID, commentID) //nolint:wrapcheck
}

func (r *ReactionService) RemoveCommentLike(commentID, userID string) error {
	return r.ReactionRepo.RemoveCommentLike(commentID, userID) //nolint:wrapcheck
}

func (r *ReactionService) HandleCommentDislike(userID, commentID string) error {
	existDisLike := r.IsCommentDislikedByUser(userID, commentID)

	if !existDisLike {
		disLike := &models.CommentDislike{
			CommentID: commentID,
			UserID:    userID,
		}

		err := r.CreateDislikeInComment(disLike)
		if err != nil {
			return fmt.Errorf("create dislike in comment: %w", err)
		}

		err = r.IncrementDislikeCountInComment(commentID)
		if err != nil {
			return fmt.Errorf("increment dislike in comment: %w", err)
		}

		err = r.RemoveCommentLike(commentID, userID)
		if err != nil {
			return fmt.Errorf("remove comment like: %w", err)
		}
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

	err = r.ReactionRepo.AddCommentDislike(dislike)
	if err != nil {
		return fmt.Errorf("add comment dislike: %w", err)
	}
	return nil
}

func (r *ReactionService) IncrementDislikeCountInComment(commentID string) error {
	err := r.ReactionRepo.IncrementDislikeCountInComment(commentID)
	if err != nil {
		return fmt.Errorf("increment dislike count in comment: %w", err)
	}

	err = r.DecrementLikeCountInComment(commentID)
	if err != nil {
		return fmt.Errorf("decrement dislike count in comment: %w", err)
	}
	return nil
}

func (r *ReactionService) DecrementDislikeCountInComment(commentID string) error {
	return r.ReactionRepo.DecrementDislikeCountInComment(commentID) //nolint:wrapcheck
}

func (r *ReactionService) IsCommentDislikedByUser(userID, commentID string) bool {
	return r.ReactionRepo.IsCommentDislikedByUser(userID, commentID) //nolint:wrapcheck
}

func (r *ReactionService) RemoveCommentDislike(commentID, userID string) error {
	return r.ReactionRepo.RemoveCommentDislike(commentID, userID) //nolint:wrapcheck
}
