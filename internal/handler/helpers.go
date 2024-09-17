package handler

import (
	"net/http"

	"forum/internal/models"
)

func (h *Handler) Render(w http.ResponseWriter, page string, data interface{}) {
	ts, ok := h.cache[page]
	if !ok {
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	err := ts.ExecuteTemplate(w, page, data)
	if err != nil {
		h.ErrorHandler(w, http.StatusInternalServerError, "Internal server error")
		return
	}
}

func (h *Handler) ResponseData(posts []models.Post, user_name, user_id string) []models.ResPostModel {
	var res []models.ResPostModel
	for _, post := range posts {
		// check exist like in post from current user
		existLike := h.service.Reaction.LikeExistInPost(post.PostID, user_id)
		// check exist dislike in post from current user
		existDisLike := h.service.Reaction.DisLikeExistInPost(post.PostID, user_id)

		post.IsLike = existLike
		post.IsDisLike = existDisLike

		// comments in post
		comments, _ := h.service.Reaction.GetCommentsByID(post.PostID)

		for el := range comments {
			comments[el].OwnerId = user_name
			// check if comment is liked
			like := h.service.Reaction.ExistLikeInComment(user_id, comments[el].CommentId)
			dislike := h.service.Reaction.ExistDisLikeInComment(user_id, comments[el].CommentId)

			comments[el].IsLiked = like
			comments[el].DisLiked = dislike

		}

		// response/result data
		posts := models.ResPostModel{
			Posts:    post,
			Comments: comments,
			OwnerId:  user_id,
		}
		res = append(res, posts)
	}
	return res
}

func (h *Handler) GetComment(post_id, user_name, user_id string) []models.Comment {
	comments, err := h.service.Reaction.GetCommentsByID(post_id)
	if err != nil {
		h.logger.Error("GET COMMENT", "ADD COMMENT HANDLER")
		return []models.Comment{}
	}

	for el := range comments {
		comments[el].OwnerId = user_name
		// check if comment is liked
		like := h.service.Reaction.ExistLikeInComment(user_id, comments[el].CommentId)
		dislike := h.service.Reaction.ExistDisLikeInComment(user_id, comments[el].CommentId)

		comments[el].IsLiked = like
		comments[el].DisLiked = dislike

	}

	return comments
}
