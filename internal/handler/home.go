package handler

import (
	"net/http"
	"strconv"

	"forum/internal/models"
)

// home page
func (h *Handler) IndexRouter(w http.ResponseWriter, r *http.Request) {
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

	totalPosts, err := h.service.PostService.GetCountPost()
	if err != nil {
		h.logger.Info("Error fetching total post count", "home page")
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	totalPages := (totalPosts + postsPerPage - 1) / postsPerPage
	data.TotalPages = totalPages

	offset := (page - 1) * postsPerPage
	var posts []models.Post
	// get all post
	posts, err = h.service.PostService.GetPostList(postsPerPage*totalPages, offset)
	if err != nil {
		posts = []models.Post{}
	}

	// check session
	s, ok := r.Context().Value(models.SessionKey).(models.Session)
	if !ok {
		data.IsAuth = false
	} else {

		data.IsAuth = true
		// get user info
		user, err := h.service.UserService.GetUserUserID(s.UserId)
		if err != nil {
			h.logger.Info("Get user by id | Internal server error", "home page")
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		data.UserName = user.UserName
		data.Id = s.UserId

	}

	// get tags list
	categories, err := h.service.PostService.GetCategoryList()
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
	// log.Println("Posts data:", posts)
	h.Render(w, "index.html", data)
}

// filter user post
func (h *Handler) UserPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Info("Method not allowed", "UserPost handler")
		h.ErrorHandler(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	data := models.Login{IsAuth: false} // default user not auth

	filter := r.URL.Query().Get("filter")
	var posts []models.Post
	posts, err := h.service.PostService.GetPostByTags(filter)
	if err != nil {
		h.logger.Info("Get post by tags", "UserPost handler")
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	s, ok := r.Context().Value(models.SessionKey).(models.Session)

	if !ok {
		data.IsAuth = false
	} else {
		data.IsAuth = true

		user, _ := h.service.UserService.GetUserUserID(s.UserId) // get user info
		data.UserName = user.UserName                            // set username
		data.Id = s.UserId                                       // set user id

		switch filter {
		case "all":
			http.Redirect(w, r, "/", http.StatusSeeOther)
		case "my":
			posts, err = h.service.PostService.GetPostByName(user.UserName)
			data.CurrentPage = "my"
			if err != nil {
				h.logger.Info("Get post by user name - 'my'", "UserPost Handler", err)
				h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
				return
				// fmt.Println(err)
			}
		case "liked":
			posts, err = h.service.PostService.GetPostByLiked(user.ID)
			data.CurrentPage = "liked"
			if err != nil {
				h.logger.Info("Get post by user name - 'liked'", "UserPost Handler", err)
				h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
				return
			}
		case "disliked":
			posts, err = h.service.PostService.GetPostByDisLike(user.ID)
			data.CurrentPage = "disliked"
			if err != nil {
				h.logger.Info("Get post by user name - 'disliked'", "UserPost Handler", err)
				h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
				return
			}

		}

	}

	categories, _ := h.service.PostService.GetCategoryList() // get tags list
	data.Categories = *categories                            // set tags

	// get post and comment
	data.Posts = h.ResponseData(posts, data.UserName, data.Id)

	h.Render(w, "index.html", data)
}
