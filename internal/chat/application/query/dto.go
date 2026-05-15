package query

type ListResult struct {
	Items interface{}
	Total int64
}

type UserInfo struct {
	ID       int64
	Username string
	Avatar   string
}

type GroupInfo struct {
	ID             int64
	ConversationID int64
	Name           string
	Avatar         string
	MemberCount    int64
}

type GroupMemberInfo struct {
	UserID   int64
	Username string
	Avatar   string
	Role     int8
}
