package post

import (
	"database/sql"
	"fmt"
	"strings"

	"forum/internal/models"
)

type PostRepo struct { //nolint:revive
	db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{
		db: db,
	}
}

func (p *PostRepo) GetAllPosts(postPerPage, offset int) ([]models.Post, error) {
	query := `SELECT * FROM posts ORDER BY created_at DESC LIMIT ? OFFSET ?`
	return p.fetchAndReturn(query, postPerPage, offset)
}

func (p *PostRepo) fetchAndReturn(query string, params ...interface{}) ([]models.Post, error) {
	rows, err := p.db.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("query data: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var tagsStr string
		err = rows.Scan(
			&post.PostID,
			&post.UserID,
			&post.Author,
			&post.Title,
			&post.Description,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.LikesCount,
			&post.DislikesCount,
			&post.Comments,
			&tagsStr,
		)
		if err != nil {
			return nil, fmt.Errorf("scan data: %w", err)
		}
		post.Tags = strings.Split(tagsStr, ",")
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return posts, nil
}

func (p *PostRepo) DeletePostByPostID(postID string) error {
	stmt := `DELETE FROM posts WHERE post_id = ?`
	if _, err := p.db.Exec(stmt, postID); err != nil {
		fmt.Println("repo: ", err)
		return models.ErrDeletePost
	}

	return nil
}

func (p *PostRepo) GetAllCategories() (*[]models.Categories, error) {
	stmt := `SELECT * FROM categories`
	rows, err := p.db.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("query data: %w", err)
	}
	defer rows.Close()

	var categories []models.Categories
	for rows.Next() {
		var category models.Categories
		if err := rows.Scan(&category.CategoryID, &category.CategoryName); err != nil {
			return nil, fmt.Errorf("scan data: %w", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return &categories, nil
}

func (p *PostRepo) AddPost(post models.Post) error {
	stmt := `INSERT INTO posts(post_id,user_id,author,title,description,created_at,updated_at,likes_count,dislikes_count,comments_count,tags)
	VALUES(?,?,?,?,?,datetime('now','localtime'),datetime('now','localtime'),?,?,?,?)`
	tagStr := strings.Join(post.Tags, ",")
	if _, err := p.db.Exec(stmt, post.PostID, post.UserID, post.Author, post.Title, post.Description, post.LikesCount, post.DislikesCount, post.Comments, tagStr); err != nil {
		fmt.Println("ERR", err)
		return models.ErrPostNotCreated
	}

	return nil
}

func (p *PostRepo) IncrementPostLikeCount(postID string) error {
	stmt := `UPDATE posts SET likes_count = likes_count + 1 WHERE post_id = ?`
	if _, err := p.db.Exec(stmt, postID); err != nil {
		return models.ErrIncrementLikeInPost
	}
	return nil
}

func (p *PostRepo) DecrementPostLikeCount(postID string) error {
	stmt := `UPDATE posts SET likes_count = likes_count - 1 WHERE post_id = ? AND likes_count > 0`
	if _, err := p.db.Exec(stmt, postID); err != nil {
		return models.ErrDecrementLikeInPost
	}
	return nil
}

func (p *PostRepo) IncrementPostDislikeCount(postID string) error {
	stmt := `UPDATE posts SET dislikes_count = dislikes_count + 1 WHERE post_id = ?`
	if _, err := p.db.Exec(stmt, postID); err != nil {
		return models.ErrIncrementLikeInPost
	}
	return nil
}

func (p *PostRepo) DecrementPostDislikeCount(postID string) error {
	stmt := `UPDATE posts SET dislikes_count = dislikes_count - 1 WHERE post_id = ? AND dislikes_count > 0`
	if _, err := p.db.Exec(stmt, postID); err != nil {
		return models.ErrDecrementDisLikeInPost
	}
	return nil
}

func (p *PostRepo) IncrementCommentCount(postID string) error {
	stmt := `UPDATE posts SET comments_count = comments_count + 1 WHERE post_id = ?`
	if _, err := p.db.Exec(stmt, postID); err != nil {
		return models.ErrIncrementCommentInPost
	}
	return nil
}

func (p *PostRepo) DecrementCommentCount(postID string) error {
	stmt := `UPDATE posts SET comments_count = comments_count - 1 WHERE post_id = ? AND comments_count > 0`
	if _, err := p.db.Exec(stmt, postID); err != nil {
		return models.ErrDecrementCommentInpost
	}

	return nil
}

func (p *PostRepo) GetPostsByUsername(username string) ([]models.Post, error) {
	stmt := `SELECT * FROM posts WHERE author = ?`
	rows, err := p.db.Query(stmt, username)
	if err != nil {
		return nil, models.ErrUserNotFound
	}
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var tagsStr string
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Author, &post.Title, &post.Description, &post.CreatedAt, &post.UpdatedAt, &post.LikesCount, &post.DislikesCount, &post.Comments, &tagsStr); err != nil {
			return nil, models.ErrNotFound
		}
		post.Tags = strings.Split(tagsStr, ",")
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}
	return posts, nil
}

func (p *PostRepo) GetPostsLikedByUser(userID string) ([]models.Post, error) {
	query := `SELECT p.post_id, p.user_id, p.author, p.title, p.description, p.created_at, p.updated_at, p.likes_count, p.dislikes_count, p.comments_count, p.tags
			FROM posts AS p
			JOIN post_likes AS l ON p.post_id = l.post_id
			WHERE l.user_id = ?
			ORDER BY p.created_at DESC`

	return p.fetchAndReturn(query, userID)
}

func (p *PostRepo) GetPostsDislikedByUser(userID string) ([]models.Post, error) {
	query := `SELECT p.post_id, p.user_id, p.author, p.title, p.description, p.created_at, p.updated_at, p.likes_count, p.dislikes_count, p.comments_count, p.tags
			FROM posts AS p
			JOIN post_dislikes AS l ON p.post_id = l.post_id
			WHERE l.user_id = ?
			ORDER BY p.created_at DESC`

	return p.fetchAndReturn(query, userID)
}

func (p *PostRepo) GetPostsByCategory(tag string) ([]models.Post, error) {
	query := `
		SELECT post_id, user_id, author, title, description, created_at, updated_at, likes_count, dislikes_count, comments_count, tags
		FROM posts
		WHERE tags LIKE ?;
	`
	rows, err := p.db.Query(query, "%"+tag+"%")
	if err != nil {
		return nil, fmt.Errorf("query data: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var tagsStr string
		if err := rows.Scan(
			&post.PostID,
			&post.UserID,
			&post.Author,
			&post.Title,
			&post.Description,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.LikesCount,
			&post.DislikesCount,
			&post.Comments,
			&tagsStr,
		); err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}
		post.Tags = strings.Split(tagsStr, ",")
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows scan: %w", err)
	}

	return posts, nil
}

func (p *PostRepo) GetPostByPostID(postID string) (models.Post, error) {
	var post models.Post
	stmt := `SELECT * FROM posts WHERE post_id = ?`
	row := p.db.QueryRow(stmt, postID)
	var tags string
	err := row.Scan(&post.PostID, &post.UserID, &post.Author, &post.Title, &post.Description, &post.CreatedAt, &post.UpdatedAt, &post.LikesCount, &post.DislikesCount, &post.Comments, &tags)
	if err == sql.ErrNoRows {
		return models.Post{}, nil
	} else if err != nil {
		return models.Post{}, fmt.Errorf("row scan: %w", err)
	}

	post.Tags = strings.Split(tags, ",")

	return post, nil
}

func (p *PostRepo) GetPostsCount() (int, error) {
	var totalPosts int
	stmt := `SELECT COUNT(*) FROM posts`
	err := p.db.QueryRow(stmt).Scan(&totalPosts)
	if err != nil {
		return 0, fmt.Errorf("row scan: %w", err)
	}
	return totalPosts, nil
}
