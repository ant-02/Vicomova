package valueobject

// ProductType 商品类型
type ProductType int8

const (
	ProductTypePointsItem  ProductType = 1 // 积分商品
	ProductTypeVIP         ProductType = 2 // VIP会员
	ProductTypeVideoAccess ProductType = 3 // 视频权限
	ProductTypeFlashSale   ProductType = 4 // 秒杀商品
)

func (p ProductType) String() string {
	switch p {
	case ProductTypePointsItem:
		return "points_item"
	case ProductTypeVIP:
		return "vip"
	case ProductTypeVideoAccess:
		return "video_access"
	case ProductTypeFlashSale:
		return "flash_sale"
	default:
		return "unknown"
	}
}

func ParseProductType(s string) ProductType {
	switch s {
	case "points_item":
		return ProductTypePointsItem
	case "vip":
		return ProductTypeVIP
	case "video_access":
		return ProductTypeVideoAccess
	case "flash_sale":
		return ProductTypeFlashSale
	default:
		return 0
	}
}
