package post

import (
	"database/sql"
	"fmt"
	"strings"

	"forum/internal/models"
)

type PostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{
		db: db,
	}
}

func (p *PostRepo) GetAllPosts(postPerPage, offset int) ([]models.Post, error) {
	stmt := `SELECT * FROM posts ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := p.db.Query(stmt, postPerPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var tagsStr string
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Author, &post.Title, &post.Description, &post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Dislikes, &post.Comments, &tagsStr); err != nil {
			return nil, err
		}
		post.Tags = strings.Split(tagsStr, ",")
		posts = append(posts, post)
	}

	return posts, nil
}

func (p *PostRepo) DeletePostByPostID(post_id string) error {
	stmt := `DELETE FROM posts WHERE post_id = ?`
	if _, err := p.db.Exec(stmt, post_id); err != nil {
		fmt.Println("repo: ", err)
		return models.ErrDeletePost
	}

	return nil
}

func (p *PostRepo) GetAllCategories() (*[]models.Categories, error) {
	stmt := `SELECT * FROM categories`
	rows, err := p.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Categories
	for rows.Next() {
		var category models.Categories
		if err := rows.Scan(&category.CategoryID, &category.CategoryName); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return &categories, nil
}

func (p *PostRepo) AddPost(post models.Post) error {
	stmt := `INSERT INTO posts(post_id,user_id,author,title,description,created_at,updated_at,likes,dislikes,comments,tags)
	VALUES(?,?,?,?,?,datetime('now','localtime'),datetime('now','localtime'),?,?,?,?)`
	tagStr := strings.Join(post.Tags, ",")
	if _, err := p.db.Exec(stmt, post.PostID, post.UserID, post.Author, post.Title, post.Description, post.Likes, post.Dislikes, post.Comments, tagStr); err != nil {
		fmt.Println("ERR", err)
		return models.ErrPostNotCreated
	}

	return nil
}

func (p *PostRepo) IncrementLikeCount(post_id string) error {
	stmt := `UPDATE posts SET likes = likes + 1 WHERE post_id = ?`
	if _, err := p.db.Exec(stmt, post_id); err != nil {
		return models.ErrIncrementLikeInPost
	}
	return nil
}

func (p *PostRepo) DecrementLikeCount(post_id string) error {
	stmt := `UPDATE posts SET likes = likes - 1 WHERE post_id = ? AND likes > 0`
	if _, err := p.db.Exec(stmt, post_id); err != nil {
		return models.ErrDecrementLikeInPost
	}
	return nil
}

func (p *PostRepo) IncrementDislikeCount(post_id string) error {
	stmt := `UPDATE posts SET dislikes = dislikes + 1 WHERE post_id = ?`
	if _, err := p.db.Exec(stmt, post_id); err != nil {
		return models.ErrIncrementLikeInPost
	}
	return nil
}

func (p *PostRepo) DecrementDislikeCount(post_id string) error {
	stmt := `UPDATE posts SET dislikes = dislikes - 1 WHERE post_id = ? AND dislikes > 0`
	if _, err := p.db.Exec(stmt, post_id); err != nil {
		return models.ErrDecrementDisLikeInPost
	}
	return nil
}

func (p *PostRepo) IncrementCommentCount(post_id string) error {
	stmt := `UPDATE posts SET comments = comments + 1 WHERE post_id = ?`
	if _, err := p.db.Exec(stmt, post_id); err != nil {
		return models.ErrIncrementCommentInPost
	}
	return nil
}

func (p *PostRepo) DecrementCommentCount(post_id string) error {
	stmt := `UPDATE posts SET comments = comments - 1 WHERE post_id = ? AND comments > 0`
	if _, err := p.db.Exec(stmt, post_id); err != nil {
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
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Author, &post.Title, &post.Description, &post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Dislikes, &post.Comments, &tagsStr); err != nil {
			return nil, models.ErrNotFound
		}
		post.Tags = strings.Split(tagsStr, ",")
		posts = append(posts, post)
	}
	return posts, nil
}

func (p *PostRepo) GetPostsLikedByUser(user_id string) ([]models.Post, error) {
	stmt := `SELECT p.post_id, p.user_id, p.author, p.title, p.description, p.created_at, p.updated_at, p.likes, p.dislikes, p.comments, p.tags
			FROM posts AS p
			JOIN likes AS l ON p.post_id = l.post_id
			WHERE l.user_id = ?
			ORDER BY p.created_at DESC`

	rows, err := p.db.Query(stmt, user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post

		var tagsStr string
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Author, &post.Title, &post.Description, &post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Dislikes, &post.Comments, &tagsStr); err != nil {
			return nil, err
		}
		post.Tags = strings.Split(tagsStr, ",")
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *PostRepo) GetPostsDislikedByUser(user_id string) ([]models.Post, error) {
	stmt := `SELECT p.post_id, p.user_id, p.author, p.title, p.description, p.created_at, p.updated_at, p.likes, p.dislikes, p.comments, p.tags
			FROM posts AS p
			JOIN dislikes AS l ON p.post_id = l.post_id
			WHERE l.user_id = ?
			ORDER BY p.created_at DESC`

	rows, err := p.db.Query(stmt, user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post

		var tagsStr string
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Author, &post.Title, &post.Description, &post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Dislikes, &post.Comments, &tagsStr); err != nil {
			return nil, err
		}
		post.Tags = strings.Split(tagsStr, ",")
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *PostRepo) GetPostsByCategory(tag string) ([]models.Post, error) {
	query := `
		SELECT post_id, user_id, author, title, description, created_at, updated_at, likes, dislikes, comments, tags
		FROM posts
		WHERE tags LIKE ?;
	`
	rows, err := p.db.Query(query, "%"+tag+"%")
	if err != nil {
		return nil, err
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
			&post.Likes,
			&post.Dislikes,
			&post.Comments,
			&tagsStr,
		); err != nil {
			return nil, err
		}
		post.Tags = strings.Split(tagsStr, ",")
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *PostRepo) GetPostByPostID(post_id string) (*models.Post, error) {
	var post models.Post
	stmt := `SELECT * FROM posts WHERE post_id = ?`
	row := p.db.QueryRow(stmt, post_id)
	var tags string
	err := row.Scan(&post.PostID, &post.UserID, &post.Author, &post.Title, &post.Description, &post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Dislikes, &post.Comments, &tags)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	post.Tags = strings.Split(tags, ",")

	return &post, err
}

func (p *PostRepo) GetPostsCount() (int, error) {
	var totalPosts int
	stmt := `SELECT COUNT(*) FROM posts`
	err := p.db.QueryRow(stmt).Scan(&totalPosts)
	if err != nil {
		return 0, err
	}
	return totalPosts, nil
}
