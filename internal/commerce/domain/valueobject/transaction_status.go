package valueobject

// TransactionStatus 交易状态
type TransactionStatus int8

const (
	TransactionStatusPending   TransactionStatus = 1 // 待处理
	TransactionStatusCompleted TransactionStatus = 2 // 已完成
	TransactionStatusFailed    TransactionStatus = 3 // 失败
)

func (s TransactionStatus) String() string {
	switch s {
	case TransactionStatusPending:
		return "pending"
	case TransactionStatusCompleted:
		return "completed"
	case TransactionStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

func ParseTransactionStatus(s string) TransactionStatus {
	switch s {
	case "pending":
		return TransactionStatusPending
	case "completed":
		return TransactionStatusCompleted
	case "failed":
		return TransactionStatusFailed
	default:
		return 0
	}
}
