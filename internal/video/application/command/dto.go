package command

// Command DTOs

type SaveVideoCommand struct {
	VideoID     int64 // 可选，有则更新，无则创建
	UserID      int64
	Title       string
	Description string
	CategoryID  int
	CoverURL    string
	VideoURL    string
	Duration    int
}

type SubmitVideoCommand struct {
	VideoID int64
	UserID  int64
}

type PublishVideoCommand struct {
	VideoID     int64 // 可选，有则更新，无则创建
	UserID      int64
	Title       string
	Description string
	CategoryID  int
	CoverURL    string
	VideoURL    string
	Duration    int
}

// Result DTOs

type SaveVideoResult struct {
	VideoID int64
}

type SubmitVideoResult struct {
	VideoID int64
}

type PublishVideoResult struct {
	VideoID int64
}
