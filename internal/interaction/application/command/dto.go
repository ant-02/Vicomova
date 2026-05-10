package command

type LikeCommand struct {
	UserID     int64
	TargetType string
	TargetID   int64
}

type CommentCommand struct {
	UserID   int64
	VideoID  int64
	ParentID int64
	Content  string
}

type FavoriteCommand struct {
	UserID  int64
	VideoID int64
}

type CommentResult struct {
	CommentID int64
}
