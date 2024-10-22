package handler

import (
	"net/http"

	"forum/internal/models"
)

func (h *Handler) Render(w http.ResponseWriter, page string, data interface{}) {
	ts, ok := h.cache[page]
	if !ok {
		return
	}

	err := ts.ExecuteTemplate(w, page, data)
	if err != nil {
		h.logger.Errorf("execute template: %v", err)
		return
	}
}

func (h *Handler) ResponseData(posts []models.Post, username, userID string) []models.UserPostsResponse {
	res := make([]models.UserPostsResponse, 0, len(posts))
	for _, post := range posts {
		existLike := h.service.ReactionService.IsPostLikedByUser(post.PostID, userID)
		existDisLike := h.service.ReactionService.IsPostDislikedByUser(post.PostID, userID)

		post.IsLike = existLike
		post.IsDisLike = existDisLike

		comments, _ := h.service.ReactionService.GetCommentsByPostID(post.PostID)

		for el := range comments {
			comments[el].OwnerID = username
			like := h.service.ReactionService.IsCommentLikedByUser(userID, comments[el].CommentID)
			dislike := h.service.ReactionService.IsCommentDislikedByUser(userID, comments[el].CommentID)

			comments[el].IsLiked = like
			comments[el].DisLiked = dislike

		}

		// response/result data
		posts := models.UserPostsResponse{
			Posts:    post,
			Comments: comments,
			OwnerID:  userID,
		}
		res = append(res, posts)
	}
	return res
}
