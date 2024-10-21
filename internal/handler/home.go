package handler

import (
	"net/http"
	"strconv"

	"forum/internal/models"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.logger.Info("Page not found", "home page")
		h.ErrorHandler(w, http.StatusNotFound, "Whoops...Page not found")
		return
	}

	if r.Method != http.MethodGet {
		h.logger.Info("Method not allowed", "home page")
		h.ErrorHandler(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	data := models.Login{IsAuth: false} // default user not auth
	pageStr := r.URL.Query().Get("page")
	var page int
	if pageStr == "" {
		page = 1
		data.Page = page
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
		data.Page = page
	}

	postsPerPage := 10

	totalPosts, err := h.service.PostService.GetPostsCount()
	if err != nil {
		h.logger.Info("Error fetching total post count", "home page")
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	totalPages := (totalPosts + postsPerPage - 1) / postsPerPage
	data.TotalPages = totalPages

	offset := (page - 1) * postsPerPage

	posts, err := h.service.PostService.GetAllPosts(postsPerPage*totalPages, offset)
	if err != nil {
		posts = []models.Post{}
	}

	// check session
	s, ok := r.Context().Value(models.SessionKey).(models.Session)
	if !ok {
		data.IsAuth = false
	} else {
		data.IsAuth = true
		user, err := h.service.UserService.GetUserByUserID(s.UserID)
		if err != nil {
			h.logger.Info("Get user by id | Internal server error", "home page")
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		data.UserName = user.UserName
		data.Id = s.UserID
	}

	// get tags list
	categories, err := h.service.PostService.GetAllCategories()
	if err != nil {
		h.logger.Info("Get categories/tags", "home page")
		data.Categories = nil
	}
	data.Categories = *categories

	// get post and comment
	data.Posts = h.ResponseData(posts, data.UserName, data.Id)
	// defalut value current page
	data.CurrentPage = "all"
	data.ShowCreatePostForm = false
	if r.URL.Path == "/create-post" {
		data.ShowCreatePostForm = true
	}

	h.Render(w, "index.html", data)
}

func (h *Handler) GetUserActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Info("Method not allowed", "GetUserActivity handler")
		h.ErrorHandler(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	data := models.Login{IsAuth: false}

	filter := r.URL.Query().Get("filter")
	var posts []models.Post
	posts, err := h.service.PostService.GetPostsByCategory(filter)
	if err != nil {
		h.logger.Info("Get post by tags", "GetUserActivity handler")
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	s, ok := r.Context().Value(models.SessionKey).(models.Session)

	if !ok {
		data.IsAuth = false
	} else {
		data.IsAuth = true

		user, _ := h.service.UserService.GetUserByUserID(s.UserID)
		data.UserName = user.UserName
		data.Id = s.UserID

		switch filter {
		case "all":
			http.Redirect(w, r, "/", http.StatusSeeOther)
		case "my":
			posts, err = h.service.PostService.GetPostsByUsername(user.UserName)
			data.CurrentPage = "my"
			if err != nil {
				h.logger.Info("Get post by user name - 'my'", "GetUserActivity Handler", err)
				h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
				return
			}
		case "liked":
			posts, err = h.service.PostService.GetPostsLikedByUser(user.ID)
			data.CurrentPage = "liked"
			if err != nil {
				h.logger.Info("Get post by user name - 'liked'", "GetUserActivity Handler", err)
				h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
				return
			}
		case "disliked":
			posts, err = h.service.PostService.GetPostsDislikedByUser(user.ID)
			data.CurrentPage = "disliked"
			if err != nil {
				h.logger.Info("Get post by user name - 'disliked'", "GetUserActivity Handler", err)
				h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
				return
			}

		}

	}

	categories, _ := h.service.PostService.GetAllCategories()
	data.Categories = *categories

	// get post and comment
	data.Posts = h.ResponseData(posts, data.UserName, data.Id)

	h.Render(w, "index.html", data)
}
