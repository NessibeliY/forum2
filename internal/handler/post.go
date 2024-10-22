package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"forum/internal/models"
	"forum/internal/validator"
)

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(models.SessionKey).(models.Session)
	if !ok {
		h.ErrorHandler(w, http.StatusUnauthorized, "Unauthorized\nPlease sign in or sign up")
		return
	}

	user, err := h.service.UserService.GetUserByUserID(session.UserID)
	if err != nil {
		h.ErrorHandler(w, http.StatusUnauthorized, "Unauthorized\nPlease sign in or sign up")
		return
	}

	data := models.Login{IsAuth: true}

	categories, err := h.service.PostService.GetAllCategories()
	if err != nil {
		h.logger.Errorf("get all categories: %v", err)
		data.Categories = nil
	}

	if r.Method == http.MethodPost {
		err = r.ParseForm()
		if err != nil {
			h.logger.Errorf("parse form: %v", err)
			h.ErrorHandler(w, http.StatusBadRequest, "Bad request")
			return
		}

		title := r.Form.Get("title")
		description := r.Form.Get("description")
		tags := r.Form["tags"]

		v := validator.NewValidator()
		createPostRequest := &models.CreatePostRequest{
			UserID:      user.ID,
			Author:      user.Username,
			Title:       title,
			Description: description,
			Tags:        tags,
		}

		if models.ValidateCreatePostRequest(v, createPostRequest); !v.Valid() {
			if v.ErrorsMap["title"] != "" {
				data.ErrorMessages.Title = v.ErrorsMap["title"]
			}

			if v.ErrorsMap["description"] != "" {
				data.ErrorMessages.Description = v.ErrorsMap["description"]
			}

			if v.ErrorsMap["tags"] != "" {
				data.ErrorMessages.Tags = v.ErrorsMap["tags"]
			}
		}

		data.Post.Title = title
		data.Post.Description = description
		data.Post.Tags = tags

		if data.ErrorMessages.Title == "" && data.ErrorMessages.Tags == "" && data.ErrorMessages.Description == "" {
			err := h.service.PostService.CreatePost(createPostRequest)
			if err != nil {
				h.logger.Errorf("create post: %v", err)
				h.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}

	data.Username = user.Username
	data.Categories = *categories

	h.Render(w, "create_post.html", data)
}

func (h *Handler) Like(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ErrorHandler(w, http.StatusMethodNotAllowed, "Method not allowd")
		return
	}

	session, ok := r.Context().Value(models.SessionKey).(models.Session)
	if !ok {
		h.ErrorHandler(w, http.StatusUnauthorized, "Unauthorized\nPlease sign in or sign up")
		return
	}

	PostID := r.FormValue("PostID")

	err := h.service.ReactionService.HandlePostLike(session.UserID, PostID)
	if err != nil {
		h.logger.Errorf("handle post like: %v", err)
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err = h.service.PostService.IncrementPostLikeCount(PostID)
	if err != nil {
		h.logger.Errorf("increment post like count: %v", err)
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	currentPath := strings.TrimPrefix(r.Header.Get("Referer"), r.Header.Get("Origin"))

	http.Redirect(w, r, currentPath, http.StatusSeeOther)
}

func (h *Handler) DisLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	session, ok := r.Context().Value(models.SessionKey).(models.Session)
	if !ok {
		h.ErrorHandler(w, http.StatusUnauthorized, "Unauthorized\nPlease sign in or sign up")
		return
	}

	PostID := r.FormValue("PostID")

	err := h.service.ReactionService.HandlePostDislike(session.UserID, PostID)
	if err != nil {
		h.logger.Errorf("handle post dislike: %v", err)
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err = h.service.PostService.IncrementPostDislikeCount(PostID)
	if err != nil {
		h.logger.Errorf("increment post dislike count: %v", err)
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	currentPath := strings.TrimPrefix(r.Header.Get("Referer"), r.Header.Get("Origin"))

	http.Redirect(w, r, currentPath, http.StatusSeeOther)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(models.SessionKey).(models.Session)

	if r.Method == http.MethodGet {
		query, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			h.logger.Errorf("parse query: %v", err)
			h.ErrorHandler(w, http.StatusBadRequest, "Status Bad Request")
			return
		}

		postID := query.Get("post-id")
		if postID == "" {
			h.ErrorHandler(w, http.StatusBadRequest, "Missing post ID")
			return
		}
		data := models.Login{IsAuth: false} // default user not authenticated

		if ok {
			data.IsAuth = true
			user, err := h.service.UserService.GetUserByUserID(session.UserID)
			if err != nil {
				h.logger.Errorf("get user by user_id: %v", err)
				h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			data.Username = user.Username
			data.ID = session.UserID
		}

		err = h.service.PostService.PopulatePostData(postID, &data)
		if err != nil {
			h.logger.Errorf("populate post data: %v", err)
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		existLike := h.service.ReactionService.IsPostLikedByUser(postID, data.ID)
		existDisLike := h.service.ReactionService.IsPostDislikedByUser(postID, data.ID)

		data.Post.IsLike = existLike
		data.Post.IsDisLike = existDisLike

		data.Comment = h.service.ReactionService.GetCommentsWithReactions(postID, data.Username, data.ID)

		// default value current page
		data.CurrentPage = AllPostsNavigation

		h.Render(w, "post.html", data)
	}
}

func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(models.SessionKey).(models.Session)

	if r.Method == http.MethodPost {
		if ok {
			user, err := h.service.UserService.GetUserByUserID(session.UserID)
			if err != nil {
				h.logger.Errorf("get user by user_id: %v", err)
				h.ErrorHandler(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			err = r.ParseForm()
			if err != nil {
				h.logger.Errorf("parse form: %v", err)
				h.ErrorHandler(w, http.StatusInternalServerError, "Status Internal Server Error")
				return
			}
			postID := r.Form.Get("post-id")
			commentText := r.Form.Get("comment_text")

			if commentText == "" {
				http.Error(w, "Comment text cannot be empty", http.StatusBadRequest)
				return
			}

			newComment := &models.Comment{
				PostID:      postID,
				Author:      user.Username,
				CommentText: commentText,
			}

			err = h.service.ReactionService.CreateCommentInPost(newComment)
			if err != nil {
				h.logger.Errorf("create comment in post: %v", err)
				h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
				return
			}

			err = h.service.PostService.IncrementCommentCount(postID)
			if err != nil {
				h.logger.Errorf("increment comment count: %v", err)
				h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
				return
			}

			currentPath := strings.TrimPrefix(r.Header.Get("Referer"), r.Header.Get("Origin"))

			http.Redirect(w, r, currentPath, http.StatusSeeOther)
		}
	}
}

func (h *Handler) CommentLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	session, ok := r.Context().Value(models.SessionKey).(models.Session)
	if !ok {
		h.ErrorHandler(w, http.StatusUnauthorized, "Unauthorized\nPlease sign in or sign up")
		return
	}

	err := r.ParseForm()
	if err != nil {
		h.logger.Errorf("parse form: %v", err)
		h.ErrorHandler(w, http.StatusBadRequest, "Bad Request")
	}
	commentID := r.FormValue("CommentID")
	postID := r.FormValue("PostID")

	err = h.service.ReactionService.HandleCommentLike(session.UserID, commentID)
	if err != nil {
		h.logger.Errorf("handle comment like: %v", err)
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/?post-id=%s", postID), http.StatusSeeOther)
}

func (h *Handler) CommentDisLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	session, ok := r.Context().Value(models.SessionKey).(models.Session)
	if !ok {
		h.ErrorHandler(w, http.StatusUnauthorized, "Unauthorized\nPlease sign in or sign up")
		return
	}

	commentID := r.FormValue("CommentID")
	postID := r.FormValue("PostID")

	err := h.service.ReactionService.HandleCommentDislike(session.UserID, commentID)
	if err != nil {
		h.logger.Errorf("handle comment dislike: %v", err)
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/?post-id=%s", postID), http.StatusSeeOther)
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	session, ok := r.Context().Value(models.SessionKey).(models.Session)
	if !ok {
		h.ErrorHandler(w, http.StatusUnauthorized, "Unauthorized\nPlease sign in or sign up")
		return
	}

	_, err := h.service.UserService.GetUserByUserID(session.UserID)
	if err != nil {
		h.logger.Errorf("get user by user_id: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	PostID := r.FormValue("PostID")

	err = h.service.PostService.DeletePostByPostID(PostID)
	if err != nil {
		h.logger.Errorf("delete post by post_id: %v", err)
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal Sever error")
		return
	}

	err = h.service.ReactionService.DeleteReaction(PostID, session.UserID)
	if err != nil {
		h.logger.Errorf("delete reaction: %v", err)
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal Sever error")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) RemoveCommentByCommentID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	session, ok := r.Context().Value(models.SessionKey).(models.Session)
	if !ok {
		h.ErrorHandler(w, http.StatusUnauthorized, "Unauthorized\nPlease sign in or sign up")
		h.logger.Error("session doesn't exist")
		return
	}

	_, err := h.service.UserService.GetUserByUserID(session.UserID)
	if err != nil {
		h.logger.Errorf("get user by user_id: %v", err)
		h.ErrorHandler(w, http.StatusInternalServerError, "Status Internal Server Error")
		return
	}

	err = r.ParseForm()
	if err != nil {
		h.logger.Errorf("parse form: %v", err)
		h.ErrorHandler(w, http.StatusInternalServerError, "Status Internal Server Error")
		return
	}
	commentID := r.FormValue("CommentID")
	postID := r.Form.Get("PostID")

	err = h.service.ReactionService.RemoveCommentByCommentID(commentID)
	if err != nil {
		h.logger.Errorf("remove comment by comment_id: %v", err)
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}

	err = h.service.PostService.DecrementCommentCount(postID)
	if err != nil {
		h.logger.Errorf("decrement comment count: %v", err)
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/?post-id=%s", postID), http.StatusSeeOther)
}
