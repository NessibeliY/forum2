package reaction

import (
	"database/sql"
	"fmt"

	"forum/internal/models"
)

type ReactionRepo struct {
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
	stmt := `INSERT INTO likes(like_id,post_id,user_id,created_at)VALUES(?,?,?,datetime('now','localtime'))`
	if _, err := r.db.Exec(stmt, like.LikeId, like.PostId, like.UserId); err != nil {
		return models.ErrNotCreated
	}
	return nil
}

// check exist like in post
func (r *ReactionRepo) LikeExistInPost(post_id, user_id string) bool {
	stmt := `SELECT COUNT(*) FROM likes WHERE user_id = ? AND post_id = ?`
	var count int
	err := r.db.QueryRow(stmt, user_id, post_id).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

// delete like
func (r *ReactionRepo) DeleteLike(post_id, user_id string) error {
	stmt := `DELETE FROM likes WHERE post_id = ? AND user_id = ?`
	if _, err := r.db.Exec(stmt, post_id, user_id); err != nil {
		return models.ErrDeleteLikeInPost
	}
	return nil
}

// DISLIKE IN POST
// insert dislike in post
func (r *ReactionRepo) InsertDisLikePost(disLike *models.Dislike) error {
	stmt := `INSERT INTO dislikes(dislike_id,post_id,user_id,created_at)VALUES(?,?,?,datetime('now','localtime'))`
	if _, err := r.db.Exec(stmt, disLike.DisLikeId, disLike.PostId, disLike.UserId); err != nil {
		return models.ErrNotCreated
	}
	return nil
}

// check exist dislike in post
func (r *ReactionRepo) DisLikeExistInPost(post_id, user_id string) bool {
	stmt := `SELECT COUNT(*) FROM dislikes WHERE user_id = ? AND post_id = ?`
	var count int
	err := r.db.QueryRow(stmt, user_id, post_id).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

// delete dislike
func (r *ReactionRepo) DeleteDisLike(post_id, user_id string) error {
	stmt := `DELETE FROM dislikes WHERE post_id = ? AND user_id = ?`
	if _, err := r.db.Exec(stmt, post_id, user_id); err != nil {
		return models.ErrDeletDisLikeInPost
	}
	return nil
}

// COMMENT IN POST
// insert new comment
func (r *ReactionRepo) InsertCommentInPost(comment *models.Comment) error {
	stmt := `INSERT INTO comments(comment_id,post_id,author,comment_text,likes,dislikes,created_at,updated_at)
	VALUES(?,?,?,?,?,?,datetime('now','localtime'),datetime('now','localtime'))`
	if _, err := r.db.Exec(stmt, comment.CommentId, comment.PostId, comment.Author, comment.CommentText, comment.Likes, comment.DisLikes); err != nil {
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
		if err := rows.Scan(&comment.CommentId, &comment.PostId, &comment.Author, &comment.CommentText, &comment.Likes, &comment.DisLikes, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
			return nil, err
		}
		commentList = append(commentList, comment)
	}

	return commentList, nil
}

// delete comment
func (r *ReactionRepo) DeleteComment(comment_id string) error {
	stmt := `DELETE from comments WHERE comment_id = ?`
	if _, err := r.db.Exec(stmt, comment_id); err != nil {
		return models.ErrDeleteComment
	}
	return nil
}

// delete comment by POST ID
func (r *ReactionRepo) DeleteCommentByPostID(post_id string) error {
	stmt := `DELETE from comments WHERE post_id = ?`
	if _, err := r.db.Exec(stmt, post_id); err != nil {
		return models.ErrDeleteComment
	}
	return nil
}

// LIKE IN COMMENT
// insert new  like in comment
func (r *ReactionRepo) InsertLikeInComment(reaction *models.CommentLike) error {
	stmt := `INSERT INTO comment_like(like_id,comment_id,user_id,created_at)VALUES(?,?,?,datetime('now','localtime'))`
	if _, err := r.db.Exec(stmt, reaction.LikeId, reaction.CommentId, reaction.UserId); err != nil {
		return models.ErrNotCreated
	}
	return nil
}

// increment like in comment
func (r *ReactionRepo) IncrementLikeInComment(comment_id string) error {
	stmt := `UPDATE comments SET likes = likes + 1 WHERE comment_id = ?`
	if _, err := r.db.Exec(stmt, comment_id); err != nil {
		return err
	}
	return nil
}

// decrement like in comment
func (r *ReactionRepo) DecrementLikeInComment(comment_id string) error {
	stmt := `UPDATE comments SET likes = likes - 1 WHERE comment_id = ? AND likes > 0`
	if _, err := r.db.Exec(stmt, comment_id); err != nil {
		return err
	}
	return nil
}

// check exist like in comment from current user
func (r *ReactionRepo) ExistLikeInComment(user_id, comment_id string) bool {
	stmt := `SELECT COUNT(*) from comment_like WHERE user_id = ? AND comment_id = ?`
	var count int
	err := r.db.QueryRow(stmt, user_id, comment_id).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

// delete like in comment
func (r *ReactionRepo) DeleteLikeInComment(comment_id, user_id string) error {
	stmt := `DELETE FROM comment_like WHERE comment_id = ? AND user_id = ?`
	if _, err := r.db.Exec(stmt, comment_id, user_id); err != nil {
		return models.ErrDeleteLikeInComment
	}
	return nil
}

func (r *ReactionRepo) DeleteLikeInCommentByUserID(user_id string) error {
	stmt := `DELETE FROM comment_like WHERE user_id = ?`
	if _, err := r.db.Exec(stmt, user_id); err != nil {
		fmt.Println("error like com", err)
		return models.ErrDeleteLikeInComment
	}
	return nil
}

// DISLIKE IN COMMENT
func (r *ReactionRepo) InsertDisLikeInComment(reaction *models.CommentDislike) error {
	stmt := `INSERT INTO comment_dislike (dislike_id,comment_id,user_id,created_at) VALUES (?,?,?,datetime('now','localtime'))`
	if _, err := r.db.Exec(stmt, reaction.DisLikeId, reaction.CommentId, reaction.UserId); err != nil {
		fmt.Println(err)
		return models.ErrNotCreated
	}
	return nil
}

// increment dislike in comment
func (r *ReactionRepo) IncrementDisLikeInComment(comment_id string) error {
	stmt := `UPDATE comments SET dislikes = dislikes + 1 WHERE comment_id = ?`
	if _, err := r.db.Exec(stmt, comment_id); err != nil {
		return err
	}
	return nil
}

// decrement dislike in comment
func (r *ReactionRepo) DecrementDisLikeInComment(comment_id string) error {
	stmt := `UPDATE comments SET dislikes = dislikes - 1 WHERE comment_id = ? AND dislikes > 0`
	if _, err := r.db.Exec(stmt, comment_id); err != nil {
		return err
	}
	return nil
}

// check exist dislike in comment from current user
func (r *ReactionRepo) ExistDisLikeInComment(user_id, comment_id string) bool {
	stmt := `SELECT COUNT(*) from comment_dislike WHERE user_id = ? AND comment_id = ?`
	var count int
	err := r.db.QueryRow(stmt, user_id, comment_id).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

// delete dislike in comment
func (r *ReactionRepo) DeleteDisLikeInComment(comment_id, user_id string) error {
	stmt := `DELETE FROM comment_dislike WHERE comment_id = ? AND user_id = ?`
	if _, err := r.db.Exec(stmt, comment_id, user_id); err != nil {
		return models.ErrDeleteLikeInComment
	}
	return nil
}

func (r *ReactionRepo) DeleteDisLikeInCommentByUserID(user_id string) error {
	stmt := `DELETE FROM comment_dislike WHERE user_id = ?`
	if _, err := r.db.Exec(stmt, user_id); err != nil {
		return models.ErrDeleteLikeInComment
	}
	return nil
}
