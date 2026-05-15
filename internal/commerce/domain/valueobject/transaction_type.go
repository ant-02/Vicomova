package valueobject

// TransactionType 交易类型
type TransactionType int8

const (
	TransactionTypeEarn     TransactionType = 1 // 获得积分
	TransactionTypeSpend    TransactionType = 2 // 消费积分
	TransactionTypeRecharge TransactionType = 3 // 充值
	TransactionTypeRefund   TransactionType = 4 // 退款
)

func (t TransactionType) String() string {
	switch t {
	case TransactionTypeEarn:
		return "earn"
	case TransactionTypeSpend:
		return "spend"
	case TransactionTypeRecharge:
		return "recharge"
	case TransactionTypeRefund:
		return "refund"
	default:
		return "unknown"
	}
}

func ParseTransactionType(s string) TransactionType {
	switch s {
	case "earn":
		return TransactionTypeEarn
	case "spend":
		return TransactionTypeSpend
	case "recharge":
		return TransactionTypeRecharge
	case "refund":
		return TransactionTypeRefund
	default:
		return 0
	}
}
