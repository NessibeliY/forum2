package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"forum/internal/models"
	"forum/internal/validator"
)

// create post
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(models.SessionKey).(models.Session)
	if !ok {
		h.ErrorHandler(w, http.StatusUnauthorized, "Unauthorized\nPlease sign in or sign up")
		return
	}
	// get info about user
	user, err := h.service.User.GetUserUserID(session.UserId)
	if err != nil {
		h.ErrorHandler(w, http.StatusUnauthorized, "Unauthorized\nPlease sign in or sign up")
		return
	}

	// default user not auth
	data := models.Login{IsAuth: true}

	// categories/tags for post
	categories, err := h.service.Post.GetCategoryList()
	if err != nil {
		data.Categories = nil
	}

	if r.Method == http.MethodPost {
		err = r.ParseForm()
		if err != nil {
			h.ErrorHandler(w, http.StatusBadRequest, "Bad request")
			return
		}

		title := r.Form.Get("title")             // title
		description := r.Form.Get("description") // description
		tags := r.Form["tags"]                   // tags

		v := validator.NewValidator() // validator
		post := &models.Post{
			UserID:      user.ID,
			Author:      user.UserName,
			Title:       title,
			Description: description,
			Tags:        tags,
		}

		// validate post input
		if models.ValidatePost(v, post); !v.Valid() {
			if v.Errors["title"] != "" {
				data.Error.Title = v.Errors["title"]
			}

			if v.Errors["description"] != "" {
				data.Error.Description = v.Errors["description"]
			}

			if v.Errors["tags"] != "" {
				data.Error.Tags = v.Errors["tags"]
			}
		}

		data.Post.Title = title             // title
		data.Post.Description = description // description
		data.Post.Tags = tags               // tags

		if data.Error.Title == "" && data.Error.Tags == "" && data.Error.Description == "" {
			err := h.service.Post.CreatePost(post) // create post
			if err != nil {
				h.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}

	data.UserName = user.UserName // set username
	data.Categories = *categories // set categories

	// http.Redirect(w, r, "/", http.StatusSeeOther)
	h.Render(w, "create_post.html", data)
}

// get all post
func (h *Handler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ErrorHandler(w, http.StatusMethodNotAllowed, "Methot not allowed")
		return
	}
	// get all post list
	// data, err := h.service.Post.GetCategoryList()
	// if err != nil {
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	return
	// }

	h.Render(w, "index.html", nil)
}

// like post
func (h *Handler) Like(w http.ResponseWriter, r *http.Request) {
	// CHECK METHOD
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

	user, err := h.service.User.GetUserUserID(session.UserId)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// post_id
	postID := r.FormValue("postId")

	// new like model
	newLike := models.Likes{
		PostId: postID,
		UserId: user.ID,
	}

	existLike := h.service.Reaction.LikeExistInPost(postID, user.ID)
	if !existLike {
		err = h.service.Reaction.CreateLikeInPost(&newLike)
		if err != nil {
			fmt.Println(err)
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		err = h.service.Post.IncrementLike(postID)
		if err != nil {
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		err = h.service.Reaction.DeleteDisLike(postID, user.ID)
		if err != nil {
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// dislike post
func (h *Handler) DisLike(w http.ResponseWriter, r *http.Request) {
	// CHECK METHOD
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	session, ok := r.Context().Value(models.SessionKey).(models.Session)
	if !ok {
		h.ErrorHandler(w, http.StatusUnauthorized, "Unauthorized\nPlease sign in or sign up")
		return
	}

	user, err := h.service.User.GetUserUserID(session.UserId)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// post_id
	postID := r.FormValue("postId")

	// new dislike
	dislike := &models.DisLike{
		PostId: postID,
		UserId: user.ID,
	}
	existDislike := h.service.Reaction.DisLikeExistInPost(postID, user.ID)
	if !existDislike {
		err := h.service.Reaction.CreateDislikeInPost(dislike)
		if err != nil {
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		err = h.service.Post.IncrementDisLike(postID)
		if err != nil {
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		err = h.service.Reaction.DeleteLikeInPost(postID, user.ID)
		if err != nil {
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(models.SessionKey).(models.Session)

	if r.Method == http.MethodGet {
		query, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			fmt.Println("ERR", err)
			return
		}

		post_id := query.Get("post-id")
		data := models.Login{IsAuth: false} // default user not auth

		if !ok {
			data.IsAuth = false
		} else {
			data.IsAuth = true
			// get user info
			user, err := h.service.User.GetUserUserID(session.UserId)
			if err != nil {
				h.logger.Info("Get user by id | Internal server error", "home page")
				h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			data.UserName = user.UserName
			data.Id = session.UserId

		}

		categories, err := h.service.Post.GetCategoryList()
		if err != nil {
			h.logger.Info("Get categories/tags", "home page")
			data.Categories = nil
		}
		data.Categories = *categories

		comments := h.GetComment(post_id, data.UserName, data.Id)
		data.Comment = comments

		post, err := h.service.Post.GetPostByID(post_id)
		if err != nil {
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		data.Post = *post

		// defalut value current page
		data.CurrentPage = "all"

		h.Render(w, "post.html", data)
	}
}

// add new comment
func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(models.SessionKey).(models.Session)

	if r.Method == http.MethodPost {
		if ok {
			user, err := h.service.User.GetUserUserID(session.UserId)
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

			err = h.service.Reaction.CreateCommentInPost(newComment)
			if err != nil {
				h.logger.Error("Create Comment In Post", err)
				h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
				return
			}

			err = h.service.Post.IncrementComment(postID)
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

// like for comment
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

	user, err := h.service.User.GetUserUserID(session.UserId)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	commentId := r.FormValue("commentId")
	postID := r.FormValue("PostId")

	existLike := h.service.Reaction.ExistLikeInComment(user.ID, commentId)

	if !existLike {
		like := &models.CommentLike{
			CommentId: commentId,
			UserId:    user.ID,
		}

		err = h.service.Reaction.CreateLikeInComment(like)
		if err != nil {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}

		err = h.service.Reaction.IncrementLikeInComment(commentId)
		if err != nil {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}

		err = h.service.Reaction.DeleteDisLikeInComment(commentId, user.ID)
		if err != nil {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/post/?post-id=%s", postID), http.StatusSeeOther)
}

// dislike for comment
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

	user, err := h.service.User.GetUserUserID(session.UserId)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	commentId := r.FormValue("commentId")
	postID := r.FormValue("PostId")
	existDisLike := h.service.Reaction.ExistDisLikeInComment(user.ID, commentId)

	if !existDisLike {
		disLike := &models.CommentDisLike{
			CommentId: commentId,
			UserId:    user.ID,
		}

		err = h.service.Reaction.CreateDisLikeInComment(disLike)
		if err != nil {
			fmt.Println("CREATE", err)
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}

		err = h.service.Reaction.IncrementDisLikeInComment(commentId)
		if err != nil {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}

		err = h.service.Reaction.DeleteLikeInComment(commentId, user.ID)
		if err != nil {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/post/?post-id=%s", postID), http.StatusSeeOther)
}

// delete post
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

	_, err := h.service.User.GetUserUserID(session.UserId)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postId := r.FormValue("postId")

	err = h.service.Post.Delete(postId)
	if err != nil {
		fmt.Println("error", err)
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal Sever error")
		return
	}

	err = h.service.Reaction.DeleteReaction(postId, session.UserId)
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

	_, err := h.service.User.GetUserUserID(session.UserId)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	r.ParseForm()
	commentId := r.FormValue("CommentId")
	postID := r.Form.Get("PostId")

	err = h.service.Reaction.DeleteComment(commentId)
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}

	err = h.service.Post.DecrementComment(postID)
	if err != nil {
		fmt.Println("error:", err)
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/?post-id=%s", postID), http.StatusSeeOther)
}
