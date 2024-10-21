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

	user, err := h.service.UserService.GetUserByUserID(session.UserId)
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
			Author:      user.UserName,
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

	data.UserName = user.UserName
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

	user, err := h.service.UserService.GetUserByUserID(session.UserId)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postID := r.FormValue("postId")

	newLike := models.Like{
		PostId: postID,
		UserId: user.ID,
	}

	existLike := h.service.ReactionService.LikeExistInPost(postID, user.ID)
	if !existLike {
		err = h.service.ReactionService.CreateLikeInPost(&newLike)
		if err != nil {
			fmt.Println(err)
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		err = h.service.PostService.IncrementLikeCount(postID)
		if err != nil {
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		err = h.service.ReactionService.DeleteDisLike(postID, user.ID)
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

	user, err := h.service.UserService.GetUserByUserID(session.UserId)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postID := r.FormValue("postId")

	dislike := &models.Dislike{
		PostId: postID,
		UserId: user.ID,
	}
	existDislike := h.service.ReactionService.DisLikeExistInPost(postID, user.ID)
	if !existDislike {
		err := h.service.ReactionService.CreateDislikeInPost(dislike)
		if err != nil {
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		err = h.service.PostService.IncrementDislikeCount(postID)
		if err != nil {
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		err = h.service.ReactionService.DeleteLikeInPost(postID, user.ID)
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
			user, err := h.service.UserService.GetUserByUserID(session.UserId)
			if err != nil {
				h.logger.Info("Get user by id | Internal server error", "home page")
				h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			data.UserName = user.UserName
			data.Id = session.UserId

		}

		categories, err := h.service.PostService.GetAllCategories()
		if err != nil {
			h.logger.Info("Get categories/tags", "home page")
			data.Categories = nil
		}
		data.Categories = *categories

		comments := h.GetComment(post_id, data.UserName, data.Id)
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
		data.CurrentPage = "all"

		h.Render(w, "post.html", data)
	}
}

func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(models.SessionKey).(models.Session)

	if r.Method == http.MethodPost {
		if ok {
			user, err := h.service.UserService.GetUserByUserID(session.UserId)
			if err != nil {
				h.ErrorHandler(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			r.ParseForm()
			postID := r.Form.Get("post-id")
			commentText := r.Form.Get("comment_text")

			if commentText == "" {
				http.Error(w, "Comment text cannot be empty", http.StatusBadRequest)
				return
			}

			newComment := &models.Comment{
				PostId:      postID,
				Author:      user.UserName,
				CommentText: commentText,
			}

			err = h.service.ReactionService.CreateCommentInPost(newComment)
			if err != nil {
				h.logger.Error("Create Comment In Post", err)
				h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
				return
			}

			err = h.service.PostService.IncrementCommentCount(postID)
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

	user, err := h.service.UserService.GetUserByUserID(session.UserId)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = r.ParseForm()
	if err != nil {
		h.ErrorHandler(w, http.StatusBadRequest, "Bad Request")
	}
	commentId := r.FormValue("commentId")
	postID := r.FormValue("PostId")

	existLike := h.service.ReactionService.ExistLikeInComment(user.ID, commentId)

	if !existLike {
		like := &models.CommentLike{
			CommentId: commentId,
			UserId:    user.ID,
		}

		err = h.service.ReactionService.CreateLikeInComment(like)
		if err != nil {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}

		err = h.service.ReactionService.IncrementLikeInComment(commentId)
		if err != nil {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}

		err = h.service.ReactionService.DeleteDisLikeInComment(commentId, user.ID)
		if err != nil {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}
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

	user, err := h.service.UserService.GetUserByUserID(session.UserId)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	commentId := r.FormValue("commentId")
	postID := r.FormValue("PostId")
	existDisLike := h.service.ReactionService.ExistDisLikeInComment(user.ID, commentId)

	if !existDisLike {
		disLike := &models.CommentDislike{
			CommentId: commentId,
			UserId:    user.ID,
		}

		err = h.service.ReactionService.CreateDisLikeInComment(disLike)
		if err != nil {
			fmt.Println("CREATE", err)
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}

		err = h.service.ReactionService.IncrementDisLikeInComment(commentId)
		if err != nil {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}

		err = h.service.ReactionService.DeleteLikeInComment(commentId, user.ID)
		if err != nil {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}
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

	_, err := h.service.UserService.GetUserByUserID(session.UserId)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postId := r.FormValue("postId")

	err = h.service.PostService.DeletePostByPostID(postId)
	if err != nil {
		fmt.Println("error", err)
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal Sever error")
		return
	}

	err = h.service.ReactionService.DeleteReaction(postId, session.UserId)
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
		return
	}

	_, err := h.service.UserService.GetUserByUserID(session.UserId)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	r.ParseForm()
	commentId := r.FormValue("commentId")
	postID := r.Form.Get("PostId")

	err = h.service.ReactionService.DeleteComment(commentId)
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}

	err = h.service.PostService.DecrementCommentCount(postID)
	if err != nil {
		fmt.Println("error:", err)
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/?post-id=%s", postID), http.StatusSeeOther)
}
