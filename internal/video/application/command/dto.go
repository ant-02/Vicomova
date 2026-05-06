package command

// Command DTOs

type PublishVideoCommand struct {
	UserID      int64
	Title       string
	Description string
	CategoryID  int
	CoverURL    string
	VideoURL    string
	Duration    int
}

// Result DTOs

type PublishVideoResult struct {
	VideoID int64
}
