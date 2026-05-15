package valueobject

// OrderStatus 订单状态
type OrderStatus int8

const (
	OrderStatusPending   OrderStatus = 1 // 待支付
	OrderStatusPaid      OrderStatus = 2 // 已支付
	OrderStatusCompleted OrderStatus = 3 // 已完成
	OrderStatusCancelled OrderStatus = 4 // 已取消
	OrderStatusRefunded  OrderStatus = 5 // 已退款
)

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusPending:
		return "pending"
	case OrderStatusPaid:
		return "paid"
	case OrderStatusCompleted:
		return "completed"
	case OrderStatusCancelled:
		return "cancelled"
	case OrderStatusRefunded:
		return "refunded"
	default:
		return "unknown"
	}
}

func ParseOrderStatus(s string) OrderStatus {
	switch s {
	case "pending":
		return OrderStatusPending
	case "paid":
		return OrderStatusPaid
	case "completed":
		return OrderStatusCompleted
	case "cancelled":
		return OrderStatusCancelled
	case "refunded":
		return OrderStatusRefunded
	default:
		return 0
	}
}

// CanPay 判断是否可以支付
func (s OrderStatus) CanPay() bool {
	return s == OrderStatusPending
}

// CanCancel 判断是否可以取消
func (s OrderStatus) CanCancel() bool {
	return s == OrderStatusPending
}
