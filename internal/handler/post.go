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
		data.Categories = nil
	}

	if r.Method == http.MethodPost {
		err = r.ParseForm()
		if err != nil {
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
	// get cookie
	session, ok := r.Context().Value(models.SessionKey).(models.Session)
	if !ok {
		h.ErrorHandler(w, http.StatusUnauthorized, "Unauthorized\nPlease sign in or sign up")
		return
	}

	user, err := h.service.UserService.GetUserByUserID(session.UserID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	PostID := r.FormValue("PostID")

	newLike := models.Like{
		PostID: PostID,
		UserID: user.ID,
	}

	existLike := h.service.ReactionService.LikeExistInPost(PostID, user.ID)
	if !existLike {
		err = h.service.ReactionService.CreateLikeInPost(&newLike)
		if err != nil {
			fmt.Println(err)
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		err = h.service.PostService.IncrementLikeCount(PostID)
		if err != nil {
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		err = h.service.ReactionService.DeleteDislike(PostID, user.ID)
		if err != nil {
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
			return
		}
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

	user, err := h.service.UserService.GetUserByUserID(session.UserID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	PostID := r.FormValue("PostID")

	dislike := &models.Dislike{
		PostID: PostID,
		UserID: user.ID,
	}
	existDislike := h.service.ReactionService.DislikeExistInPost(PostID, user.ID)
	if !existDislike {
		err := h.service.ReactionService.CreateDislikeInPost(dislike)
		if err != nil {
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		err = h.service.PostService.IncrementDislikeCount(PostID)
		if err != nil {
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		err = h.service.ReactionService.DeleteLikeInPost(PostID, user.ID)
		if err != nil {
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
			return
		}
	}

	currentPath := strings.TrimPrefix(r.Header.Get("Referer"), r.Header.Get("Origin"))

	http.Redirect(w, r, currentPath, http.StatusSeeOther)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(models.SessionKey).(models.Session)

	if r.Method == http.MethodGet {
		query, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			h.logger.Info("parse query", err)
			h.ErrorHandler(w, http.StatusBadRequest, "Status Bad Request")
			return
		}

		post_id := query.Get("post-id")
		if post_id == "" {
			h.ErrorHandler(w, http.StatusBadRequest, "Missing post ID")
			return
		}
		data := models.Login{IsAuth: false} // default user not auth

		if !ok {
			data.IsAuth = false
		} else {
			data.IsAuth = true
			user, err := h.service.UserService.GetUserByUserID(session.UserID)
			if err != nil {
				h.logger.Info("Get user by id | Internal server error", "home page")
				h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			data.Username = user.Username
			data.Id = session.UserID

		}

		categories, err := h.service.PostService.GetAllCategories()
		if err != nil {
			h.logger.Info("Get categories/tags", "home page")
			data.Categories = nil
		}
		data.Categories = *categories

		comments := h.GetComment(post_id, data.Username, data.Id)
		data.Comment = comments

		post, err := h.service.PostService.GetPostByPostID(post_id)
		if err != nil {
			h.ErrorHandler(w, http.StatusBadRequest, "Status Bad Request")
			// h.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		h.CheckPostReaction(post, data.Id)

		data.Post = *post

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
				h.ErrorHandler(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			err = r.ParseForm()
			if err != nil {
				h.ErrorHandler(w, http.StatusInternalServerError, "Status Internal Server Error")
				return
			}
			PostID := r.Form.Get("post-id")
			commentText := r.Form.Get("comment_text")

			if commentText == "" {
				http.Error(w, "Comment text cannot be empty", http.StatusBadRequest)
				return
			}

			newComment := &models.Comment{
				PostID:      PostID,
				Author:      user.Username,
				CommentText: commentText,
			}

			err = h.service.ReactionService.CreateCommentInPost(newComment)
			if err != nil {
				h.logger.Error("Create Comment In Post", err)
				h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
				return
			}

			err = h.service.PostService.IncrementCommentCount(PostID)
			if err != nil {
				h.logger.Error("Increment Comment In Post", err)
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

	user, err := h.service.UserService.GetUserByUserID(session.UserID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = r.ParseForm()
	if err != nil {
		h.ErrorHandler(w, http.StatusBadRequest, "Bad Request")
	}
	CommentID := r.FormValue("CommentID")
	PostID := r.FormValue("PostID")

	existLike := h.service.ReactionService.ExistLikeInComment(user.ID, CommentID)

	if !existLike {
		like := &models.CommentLike{
			CommentID: CommentID,
			UserID:    user.ID,
		}

		err = h.service.ReactionService.CreateLikeInComment(like)
		if err != nil {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}

		err = h.service.ReactionService.IncrementLikeInComment(CommentID)
		if err != nil {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}

		err = h.service.ReactionService.DeleteDisLikeInComment(CommentID, user.ID)
		if err != nil {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/post/?post-id=%s", PostID), http.StatusSeeOther)
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

	user, err := h.service.UserService.GetUserByUserID(session.UserID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	CommentID := r.FormValue("CommentID")
	PostID := r.FormValue("PostID")
	existDisLike := h.service.ReactionService.ExistDisLikeInComment(user.ID, CommentID)

	if !existDisLike {
		disLike := &models.CommentDislike{
			CommentID: CommentID,
			UserID:    user.ID,
		}

		err = h.service.ReactionService.CreateDislikeInComment(disLike)
		if err != nil {
			fmt.Println("CREATE", err)
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}

		err = h.service.ReactionService.IncrementDislikeCountInComment(CommentID)
		if err != nil {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}

		err = h.service.ReactionService.DeleteLikeInComment(CommentID, user.ID)
		if err != nil {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/post/?post-id=%s", PostID), http.StatusSeeOther)
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
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	PostID := r.FormValue("PostID")

	err = h.service.PostService.DeletePostByPostID(PostID)
	if err != nil {
		fmt.Println("error", err)
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal Sever error")
		return
	}

	err = h.service.ReactionService.DeleteReaction(PostID, session.UserID)
	if err != nil {
		fmt.Println("error", err)
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal Sever error")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
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
		h.ErrorHandler(w, http.StatusInternalServerError, "Status Internal Server Error")
		h.logger.Error("get user: %w", err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		h.ErrorHandler(w, http.StatusInternalServerError, "Status Internal Server Error")
		h.logger.Error("parse form: %w", err)
		return
	}
	CommentID := r.FormValue("CommentID")
	PostID := r.Form.Get("PostID")

	err = h.service.ReactionService.DeleteComment(CommentID)
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}

	err = h.service.PostService.DecrementCommentCount(PostID)
	if err != nil {
		fmt.Println("error:", err)
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/?post-id=%s", PostID), http.StatusSeeOther)
}
