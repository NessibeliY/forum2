package reaction

import (
	"database/sql"
	"fmt"

	"forum/internal/models"
)

type ReactionRepo struct { //nolint:revive
	db *sql.DB
}

func NewReactionRepo(db *sql.DB) *ReactionRepo {
	return &ReactionRepo{
		db: db,
	}
}

func (r *ReactionRepo) AddPostLike(like models.Like) error {
	stmt := `INSERT INTO post_likes(like_id, post_id, user_id, created_at)VALUES(?,?,?,datetime('now','localtime'))`
	if _, err := r.db.Exec(stmt, like.LikeID, like.PostID, like.UserID); err != nil {
		return models.ErrNotCreated
	}
	return nil
}

func (r *ReactionRepo) IsPostLikedByUser(postID, userID string) bool {
	stmt := `SELECT COUNT(*) FROM post_likes WHERE user_id = ? AND post_id = ?`
	var count int
	err := r.db.QueryRow(stmt, userID, postID).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

func (r *ReactionRepo) RemovePostLike(postID, userID string) error {
	stmt := `DELETE FROM post_likes WHERE post_id = ? AND user_id = ?`
	if _, err := r.db.Exec(stmt, postID, userID); err != nil {
		return models.ErrDeleteLikeInPost
	}
	return nil
}

func (r *ReactionRepo) AddPostDislike(dislike models.Dislike) error {
	stmt := `INSERT INTO post_dislikes(dislike_id, post_id, user_id, created_at)VALUES(?,?,?,datetime('now','localtime'))`
	if _, err := r.db.Exec(stmt, dislike.DislikeID, dislike.PostID, dislike.UserID); err != nil {
		return models.ErrNotCreated
	}
	return nil
}

func (r *ReactionRepo) IsPostDislikedByUser(postID, userID string) bool {
	stmt := `SELECT COUNT(*) FROM post_dislikes WHERE user_id = ? AND post_id = ?`
	var count int
	err := r.db.QueryRow(stmt, userID, postID).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func (r *ReactionRepo) RemovePostDislike(postID, userID string) error {
	stmt := `DELETE FROM post_dislikes WHERE post_id = ? AND user_id = ?`
	if _, err := r.db.Exec(stmt, postID, userID); err != nil {
		return models.ErrDeletDisLikeInPost
	}
	return nil
}

func (r *ReactionRepo) AddComment(comment *models.Comment) error {
	stmt := `INSERT INTO comments(comment_id, post_id, author, comment_text ,likes_count, dislikes_count, created_at, updated_at)
	VALUES(?,?,?,?,?,?,datetime('now','localtime'),datetime('now','localtime'))`
	if _, err := r.db.Exec(stmt, comment.CommentID, comment.PostID, comment.Author, comment.CommentText, comment.LikesCount, comment.DislikesCount); err != nil {
		return models.ErrNotFound
	}
	return nil
}

func (r *ReactionRepo) GetCommentsByPostID(id string) ([]models.Comment, error) {
	stmt := `SELECT * FROM comments WHERE post_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(stmt, id)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()
	var commentList []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.CommentID, &comment.PostID, &comment.Author, &comment.CommentText, &comment.LikesCount, &comment.DislikesCount, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		commentList = append(commentList, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return commentList, nil
}

func (r *ReactionRepo) RemoveCommentByCommentID(commentID string) error {
	stmt := `DELETE from comments WHERE comment_id = ?`
	if _, err := r.db.Exec(stmt, commentID); err != nil {
		return models.ErrDeleteComment
	}
	return nil
}

func (r *ReactionRepo) CommentExistsByCommentID(commentID string) bool {
	stmt := `SELECT comment_id FROM comments WHERE comment_id = ?`
	row := r.db.QueryRow(stmt, commentID)

	var id string
	return row.Scan(&id) == nil
}

func (r *ReactionRepo) RemoveCommentByPostID(postID string) error {
	stmt := `DELETE from comments WHERE post_id = ?`
	if _, err := r.db.Exec(stmt, postID); err != nil {
		return models.ErrDeleteComment
	}
	return nil
}

func (r *ReactionRepo) AddCommentLike(reaction *models.CommentLike) error {
	stmt := `INSERT INTO comment_likes(like_id,comment_id,user_id,created_at)VALUES(?,?,?,datetime('now','localtime'))`
	if _, err := r.db.Exec(stmt, reaction.LikeID, reaction.CommentID, reaction.UserID); err != nil {
		return models.ErrNotCreated
	}
	return nil
}

func (r *ReactionRepo) IncrementLikeCountInComment(commentID string) error {
	stmt := `UPDATE comments SET likes_count = likes_count + 1 WHERE comment_id = ?`
	_, err := r.db.Exec(stmt, commentID)
	return err //nolint:wrapcheck
}

func (r *ReactionRepo) DecrementLikeCountInComment(commentID string) error {
	stmt := `UPDATE comments SET likes_count = likes_count - 1 WHERE comment_id = ? AND likes_count > 0`
	_, err := r.db.Exec(stmt, commentID)
	return err //nolint:wrapcheck
}

func (r *ReactionRepo) IsCommentLikedByUser(userID, commentID string) bool {
	stmt := `SELECT COUNT(*) from comment_likes WHERE user_id = ? AND comment_id = ?`
	var count int
	err := r.db.QueryRow(stmt, userID, commentID).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

func (r *ReactionRepo) RemoveCommentLike(commentID, userID string) error {
	stmt := `DELETE FROM comment_likes WHERE comment_id = ? AND user_id = ?`
	if _, err := r.db.Exec(stmt, commentID, userID); err != nil {
		return models.ErrDeleteLikeInComment
	}
	return nil
}

func (r *ReactionRepo) RemoveAllCommentLikesByUser(userID string) error {
	stmt := `DELETE FROM comment_likes WHERE user_id = ?`
	if _, err := r.db.Exec(stmt, userID); err != nil {
		return models.ErrDeleteLikeInComment
	}
	return nil
}

func (r *ReactionRepo) AddCommentDislike(reaction *models.CommentDislike) error {
	stmt := `INSERT INTO comment_dislikes (dislike_id,comment_id,user_id,created_at) VALUES (?,?,?,datetime('now','localtime'))`
	if _, err := r.db.Exec(stmt, reaction.DislikeID, reaction.CommentID, reaction.UserID); err != nil {
		return models.ErrNotCreated
	}
	return nil
}

func (r *ReactionRepo) IncrementDislikeCountInComment(commentID string) error {
	stmt := `UPDATE comments SET dislikes_count = dislikes_count + 1 WHERE comment_id = ?`
	_, err := r.db.Exec(stmt, commentID)
	return err //nolint:wrapcheck
}

func (r *ReactionRepo) DecrementDislikeCountInComment(commentID string) error {
	stmt := `UPDATE comments SET dislikes_count = dislikes_count - 1 WHERE comment_id = ? AND dislikes_count > 0`
	_, err := r.db.Exec(stmt, commentID)
	return err //nolint:wrapcheck
}

func (r *ReactionRepo) IsCommentDislikedByUser(userID, commentID string) bool {
	stmt := `SELECT COUNT(*) from comment_dislikes WHERE user_id = ? AND comment_id = ?`
	var count int
	err := r.db.QueryRow(stmt, userID, commentID).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

func (r *ReactionRepo) RemoveCommentDislike(commentID, userID string) error {
	stmt := `DELETE FROM comment_dislikes WHERE comment_id = ? AND user_id = ?`
	if _, err := r.db.Exec(stmt, commentID, userID); err != nil {
		return models.ErrDeleteLikeInComment
	}
	return nil
}

func (r *ReactionRepo) RemoveAllCommentDislikesByUser(userID string) error {
	stmt := `DELETE FROM comment_dislikes WHERE user_id = ?`
	if _, err := r.db.Exec(stmt, userID); err != nil {
		return models.ErrDeleteLikeInComment
	}
	return nil
}
