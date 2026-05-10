package errors

var (
	ErrVideoNotFound     = NewBizError(404, "Video not found")
	ErrVideoStatusInvalid = NewBizError(400, "Invalid video status")
	ErrForbidden         = NewBizError(403, "Forbidden")
)
