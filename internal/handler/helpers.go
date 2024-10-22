package handler

import (
	"net/http"

	"forum/internal/models"
)

func (h *Handler) Render(w http.ResponseWriter, page string, data interface{}) {
	ts, ok := h.cache[page]
	if !ok {
		// h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	err := ts.ExecuteTemplate(w, page, data)
	if err != nil {
		// h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
		return
	}
}

func (h *Handler) ResponseData(posts []models.Post, username, UserID string) []models.UserPostsResponse {
	var res []models.UserPostsResponse
	for _, post := range posts {
		// check exist like in post from current user
		existLike := h.service.ReactionService.LikeExistInPost(post.PostID, UserID)
		// check exist dislike in post from current user
		existDisLike := h.service.ReactionService.DislikeExistInPost(post.PostID, UserID)

		post.IsLike = existLike
		post.IsDisLike = existDisLike

		// comments in post
		comments, _ := h.service.ReactionService.GetCommentsByID(post.PostID)

		for el := range comments {
			comments[el].OwnerId = username
			// check if comment is liked
			like := h.service.ReactionService.ExistLikeInComment(UserID, comments[el].CommentID)
			dislike := h.service.ReactionService.ExistDisLikeInComment(UserID, comments[el].CommentID)

			comments[el].IsLiked = like
			comments[el].DisLiked = dislike

		}

		// response/result data
		posts := models.UserPostsResponse{
			Posts:    post,
			Comments: comments,
			OwnerId:  UserID,
		}
		res = append(res, posts)
	}
	return res
}

func (h *Handler) GetComment(post_id, user_name, user_id string) []models.Comment {
	comments, err := h.service.ReactionService.GetCommentsByID(post_id)
	if err != nil {
		h.logger.Error("get comments by id: %w", err)
		return []models.Comment{}
	}

	for el := range comments {
		comments[el].OwnerId = user_name
		// check if comment is liked
		like := h.service.ReactionService.ExistLikeInComment(user_id, comments[el].CommentID)
		dislike := h.service.ReactionService.ExistDisLikeInComment(user_id, comments[el].CommentID)

		comments[el].IsLiked = like
		comments[el].DisLiked = dislike

	}

	return comments
}

func (h *Handler) CheckPostReaction(post *models.Post, UserID string) {
	existLike := h.service.ReactionService.LikeExistInPost(post.PostID, UserID)
	existDisLike := h.service.ReactionService.DislikeExistInPost(post.PostID, UserID)

	if existLike {
		post.IsLike = existLike
	}

	if existDisLike {
		post.IsDisLike = existDisLike
	}

	// return post
}
