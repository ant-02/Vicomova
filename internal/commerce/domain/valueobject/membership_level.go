package valueobject

// MembershipLevel 会员等级
type MembershipLevel int8

const (
	MembershipLevelOrdinary MembershipLevel = 1 // 普通会员
	MembershipLevelSilver   MembershipLevel = 2 // 银牌会员
	MembershipLevelGold     MembershipLevel = 3 // 金牌会员
)

func (l MembershipLevel) String() string {
	switch l {
	case MembershipLevelOrdinary:
		return "ordinary"
	case MembershipLevelSilver:
		return "silver"
	case MembershipLevelGold:
		return "gold"
	default:
		return "unknown"
	}
}

func ParseMembershipLevel(s string) MembershipLevel {
	switch s {
	case "ordinary":
		return MembershipLevelOrdinary
	case "silver":
		return MembershipLevelSilver
	case "gold":
		return MembershipLevelGold
	default:
		return 0
	}
}

// CanAccess 判断是否满足所需会员等级
func (l MembershipLevel) CanAccess(required MembershipLevel) bool {
	return l >= required
}
