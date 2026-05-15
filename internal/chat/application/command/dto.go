package command

type FollowCommand struct {
	UserID   int64
	TargetID int64
}

type UnfollowCommand struct {
	UserID   int64
	TargetID int64
}

type SendMessageCommand struct {
	SenderID   int64
	ReceiverID int64
	Content    string
}

type SendMessageResult struct {
	MessageID      int64
	ConversationID int64
}

type CreateGroupCommand struct {
	CreatorID int64
	Name      string
}

type CreateGroupResult struct {
	GroupID        int64
	ConversationID int64
}

type AddGroupMemberCommand struct {
	GroupID int64
	UserID  int64
	AdminID int64
}

type RemoveGroupMemberCommand struct {
	GroupID int64
	UserID  int64
	AdminID int64
}

type SendGroupMessageCommand struct {
	SenderID int64
	GroupID  int64
	Content  string
}

type SendGroupMessageResult struct {
	MessageID      int64
	ConversationID int64
}
