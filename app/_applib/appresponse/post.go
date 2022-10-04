package appresponse

type PostResponse struct {
	Id           uint64 `json:"id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	Status       string `json:"status"`
	UserId       uint64 `json:"user_id"`
	ReadCount    int64  `json:"read_count"`
	LikeCount    int64  `json:"like_count"`
	CommentCount int64  `json:"comment_count"`
}
