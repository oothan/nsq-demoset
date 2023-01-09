package dto

type RequestPostView struct {
	PostId uint64 `json:"post_id" form:"post_id" binding:"required"`
}
