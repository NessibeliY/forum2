package handler

import "net/http"

func (h *Handler) Routes() http.Handler {
	router := http.NewServeMux()

	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// router.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.NotFound(w, r) // Отправляет 404 Not Found
	// })

	router.HandleFunc("/", h.CheckSession(http.HandlerFunc(h.IndexRouter))) // home page

	// user
	router.HandleFunc("/sign-up", h.SignUp) // login
	router.HandleFunc("/sign-in", h.SignIn) // register
	router.HandleFunc("/log-out", h.LogOut) // logout

	// todo add middleware for session |-| done
	// todo add logging |-| done
	// todo add error page |-| done
	// todo extract run to main |-| done
	// todo hide http://localhost:4000/static/
	// todo create new design for 'create-post' like page not modal  |-| done
	// todo normal session for user this not work |-| done
	// Try opening two different browsers and login into one of them. Refresh the other browser. |-| done
	// Can you confirm that the browser non logged remains unregistered
	// posts

	// router.HandleFunc("/posts", h.GetAllPosts)                                                // get all post
	router.HandleFunc("/create-post", h.CheckSession(http.HandlerFunc(h.CreatePost)))         // create post
	router.HandleFunc("/like", h.CheckSession(http.HandlerFunc(h.Like)))                      // like in post
	router.HandleFunc("/dislike", h.CheckSession(http.HandlerFunc(h.DisLike)))                // dislike in post
	router.HandleFunc("/post/", h.CheckSession(http.HandlerFunc(h.Post)))                     // get post by ID
	router.HandleFunc("/add-comment/", h.CheckSession(http.HandlerFunc(h.AddComment)))        // add comment in post
	router.HandleFunc("/comment/like", h.CheckSession(http.HandlerFunc(h.CommentLike)))       // like for comment in post
	router.HandleFunc("/comment/dislike", h.CheckSession(http.HandlerFunc(h.CommentDisLike))) // dislike for comment in post

	router.HandleFunc("/user/posts/", h.CheckSession(http.HandlerFunc(h.UserPost)))              // filter post by user
	router.HandleFunc("/user/post/delete", h.CheckSession(http.HandlerFunc(h.DeletePost)))       // delete post
	router.HandleFunc("/user/comment/delete", h.CheckSession(http.HandlerFunc(h.DeleteComment))) // delete comment

	return router
}
