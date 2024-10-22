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

// LIKE IN POST
// insert like in DB
func (r *ReactionRepo) InsertLikePost(like *models.Like) error {
	stmt := `INSERT INTO post_likes(like_id, post_id, user_id, created_at)VALUES(?,?,?,datetime('now','localtime'))`
	if _, err := r.db.Exec(stmt, like.LikeID, like.PostID, like.UserID); err != nil {
		return models.ErrNotCreated
	}
	return nil
}

// check exist like in post
func (r *ReactionRepo) LikeExistInPost(PostID, UserID string) bool {
	stmt := `SELECT COUNT(*) FROM post_likes WHERE user_id = ? AND post_id = ?`
	var count int
	err := r.db.QueryRow(stmt, UserID, PostID).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

// delete like
func (r *ReactionRepo) DeleteLike(PostID, UserID string) error {
	stmt := `DELETE FROM post_likes WHERE post_id = ? AND user_id = ?`
	if _, err := r.db.Exec(stmt, PostID, UserID); err != nil {
		return models.ErrDeleteLikeInPost
	}
	return nil
}

// DISLIKE IN POST
// insert dislike in post
func (r *ReactionRepo) InsertDisLikePost(dislike *models.Dislike) error {
	stmt := `INSERT INTO post_dislikes(dislike_id, post_id, user_id, created_at)VALUES(?,?,?,datetime('now','localtime'))`
	if _, err := r.db.Exec(stmt, dislike.DislikeID, dislike.PostID, dislike.UserID); err != nil {
		return models.ErrNotCreated
	}
	return nil
}

// check exist dislike in post
func (r *ReactionRepo) DislikeExistInPost(PostID, UserID string) bool {
	stmt := `SELECT COUNT(*) FROM post_dislikes WHERE user_id = ? AND post_id = ?`
	var count int
	err := r.db.QueryRow(stmt, UserID, PostID).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func (r *ReactionRepo) DeleteDislike(PostID, UserID string) error {
	stmt := `DELETE FROM post_dislikes WHERE post_id = ? AND user_id = ?`
	if _, err := r.db.Exec(stmt, PostID, UserID); err != nil {
		return models.ErrDeletDisLikeInPost
	}
	return nil
}

// COMMENT IN POST
// insert new comment
func (r *ReactionRepo) InsertCommentInPost(comment *models.Comment) error {
	stmt := `INSERT INTO comments(comment_id, post_id, author, comment_text ,likes_count, dislikes_count, created_at, updated_at)
	VALUES(?,?,?,?,?,?,datetime('now','localtime'),datetime('now','localtime'))`
	if _, err := r.db.Exec(stmt, comment.CommentID, comment.PostID, comment.Author, comment.CommentText, comment.LikesCount, comment.DislikesCount); err != nil {
		return models.ErrNotFound
	}
	return nil
}

// get comment list by ID
func (r *ReactionRepo) GetCommentsByID(id string) ([]models.Comment, error) {
	stmt := `SELECT * FROM comments WHERE post_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var commentList []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.CommentID, &comment.PostID, &comment.Author, &comment.CommentText, &comment.LikesCount, &comment.DislikesCount, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
			return nil, err
		}
		commentList = append(commentList, comment)
	}

	return commentList, nil
}

// delete comment
func (r *ReactionRepo) DeleteComment(CommentID string) error {
	stmt := `DELETE from comments WHERE comment_id = ?`
	if _, err := r.db.Exec(stmt, CommentID); err != nil {
		return models.ErrDeleteComment
	}
	return nil
}

// delete comment by POST ID
func (r *ReactionRepo) DeleteCommentByPostID(PostID string) error {
	stmt := `DELETE from comments WHERE post_id = ?`
	if _, err := r.db.Exec(stmt, PostID); err != nil {
		return models.ErrDeleteComment
	}
	return nil
}

// LIKE IN COMMENT
// insert new  like in comment
func (r *ReactionRepo) InsertLikeInComment(reaction *models.CommentLike) error {
	stmt := `INSERT INTO comment_likes(like_id,comment_id,user_id,created_at)VALUES(?,?,?,datetime('now','localtime'))`
	if _, err := r.db.Exec(stmt, reaction.LikeID, reaction.CommentID, reaction.UserID); err != nil {
		return models.ErrNotCreated
	}
	return nil
}

// increment like in comment
func (r *ReactionRepo) IncrementLikeInComment(CommentID string) error {
	stmt := `UPDATE comments SET likes_count = likes_count + 1 WHERE comment_id = ?`
	if _, err := r.db.Exec(stmt, CommentID); err != nil {
		return err
	}
	return nil
}

// decrement like in comment
func (r *ReactionRepo) DecrementLikeCountInComment(CommentID string) error {
	stmt := `UPDATE comments SET likes_count = likes_count - 1 WHERE comment_id = ? AND likes_count > 0`
	if _, err := r.db.Exec(stmt, CommentID); err != nil {
		return err
	}
	return nil
}

// check exist like in comment from current user
func (r *ReactionRepo) ExistLikeInComment(UserID, CommentID string) bool {
	stmt := `SELECT COUNT(*) from comment_likes WHERE user_id = ? AND comment_id = ?`
	var count int
	err := r.db.QueryRow(stmt, UserID, CommentID).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

// delete like in comment
func (r *ReactionRepo) DeleteLikeInComment(CommentID, UserID string) error {
	stmt := `DELETE FROM comment_likes WHERE comment_id = ? AND user_id = ?`
	if _, err := r.db.Exec(stmt, CommentID, UserID); err != nil {
		return models.ErrDeleteLikeInComment
	}
	return nil
}

func (r *ReactionRepo) DeleteLikeInCommentByUserID(UserID string) error {
	stmt := `DELETE FROM comment_likes WHERE user_id = ?`
	if _, err := r.db.Exec(stmt, UserID); err != nil {
		fmt.Println("error like com", err)
		return models.ErrDeleteLikeInComment
	}
	return nil
}

// DISLIKE IN COMMENT
func (r *ReactionRepo) InsertDislikeInComment(reaction *models.CommentDislike) error {
	stmt := `INSERT INTO comment_dislikes (dislike_id,comment_id,user_id,created_at) VALUES (?,?,?,datetime('now','localtime'))`
	if _, err := r.db.Exec(stmt, reaction.DislikeID, reaction.CommentID, reaction.UserID); err != nil {
		fmt.Println(err)
		return models.ErrNotCreated
	}
	return nil
}

// increment dislike in comment
func (r *ReactionRepo) IncrementDislikeCountInComment(CommentID string) error {
	stmt := `UPDATE comments SET dislikes_count = dislikes_count + 1 WHERE comment_id = ?`
	if _, err := r.db.Exec(stmt, CommentID); err != nil {
		return err
	}
	return nil
}

// decrement dislike in comment
func (r *ReactionRepo) DecrementDislikeCountInComment(CommentID string) error {
	stmt := `UPDATE comments SET dislikes_count = dislikes_count - 1 WHERE comment_id = ? AND dislikes_count > 0`
	if _, err := r.db.Exec(stmt, CommentID); err != nil {
		return err
	}
	return nil
}

// check exist dislike in comment from current user
func (r *ReactionRepo) ExistDisLikeInComment(UserID, CommentID string) bool {
	stmt := `SELECT COUNT(*) from comment_dislikes WHERE user_id = ? AND comment_id = ?`
	var count int
	err := r.db.QueryRow(stmt, UserID, CommentID).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

// delete dislike in comment
func (r *ReactionRepo) DeleteDisLikeInComment(CommentID, UserID string) error {
	stmt := `DELETE FROM comment_dislikes WHERE comment_id = ? AND user_id = ?`
	if _, err := r.db.Exec(stmt, CommentID, UserID); err != nil {
		return models.ErrDeleteLikeInComment
	}
	return nil
}

func (r *ReactionRepo) DeleteDisLikeInCommentByUserID(UserID string) error {
	stmt := `DELETE FROM comment_dislikes WHERE user_id = ?`
	if _, err := r.db.Exec(stmt, UserID); err != nil {
		return models.ErrDeleteLikeInComment
	}
	return nil
}
