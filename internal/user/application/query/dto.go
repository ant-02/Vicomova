package query

type GetUserQuery struct {
	UserID   int64
	Username string
}

type UserResult struct {
	UserID   int64
	Username string
	Email    string
}

type BatchGetUsersQuery struct {
	UserIDs []int64
}
