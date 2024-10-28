package handler

import "net/http"

func (h *Handler) Routes() http.Handler {
	router := http.NewServeMux()

	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	router.HandleFunc("/", h.SessionAuthMiddleware(http.HandlerFunc(h.Home)))

	router.HandleFunc("/sign-up", h.SignUp)
	router.HandleFunc("/sign-in", h.SignIn)
	router.HandleFunc("/log-out", h.LogOut)

	router.HandleFunc("/create-post", h.SessionAuthMiddleware(http.HandlerFunc(h.CreatePost)))
	router.HandleFunc("/like", h.SessionAuthMiddleware(http.HandlerFunc(h.Like)))
	router.HandleFunc("/dislike", h.SessionAuthMiddleware(http.HandlerFunc(h.DisLike)))
	router.HandleFunc("/post/", h.SessionAuthMiddleware(http.HandlerFunc(h.Post))) // get post by ID

	router.HandleFunc("/post/comment", h.SessionAuthMiddleware(http.HandlerFunc(h.AddComment)))
	router.HandleFunc("/comment/like", h.SessionAuthMiddleware(http.HandlerFunc(h.CommentLike)))
	router.HandleFunc("/comment/dislike", h.SessionAuthMiddleware(http.HandlerFunc(h.CommentDisLike)))

	router.HandleFunc("/user/posts/", h.SessionAuthMiddleware(http.HandlerFunc(h.GetUserActivity))) // filter post by user
	router.HandleFunc("/user/post/delete", h.SessionAuthMiddleware(http.HandlerFunc(h.DeletePost)))
	router.HandleFunc("/user/comment/delete", h.SessionAuthMiddleware(http.HandlerFunc(h.RemoveCommentByCommentID)))

	return router
}
