package mysql

import (
	commerceEntity "vicomova/internal/commerce/domain/entity"
)

func PointsWalletToPO(w *commerceEntity.PointsWallet) *PointsWalletPO {
	if w == nil {
		return nil
	}
	return &PointsWalletPO{
		ID:      w.ID,
		UserID:  w.UserID,
		Points:  w.Points,
		Version: w.Version,
	}
}

func POToPointsWallet(po *PointsWalletPO) *commerceEntity.PointsWallet {
	if po == nil {
		return nil
	}
	return &commerceEntity.PointsWallet{
		ID:      po.ID,
		UserID:  po.UserID,
		Points:  po.Points,
		Version: po.Version,
	}
}

func SignInRecordToPO(r *commerceEntity.SignInRecord) *SignInRecordPO {
	if r == nil {
		return nil
	}
	return &SignInRecordPO{
		ID:              r.ID,
		UserID:          r.UserID,
		SignDate:        r.SignDate,
		ConsecutiveDays: r.ConsecutiveDays,
		BonusPoints:     r.BonusPoints,
	}
}

func POToSignInRecord(po *SignInRecordPO) *commerceEntity.SignInRecord {
	if po == nil {
		return nil
	}
	return &commerceEntity.SignInRecord{
		ID:              po.ID,
		UserID:          po.UserID,
		SignDate:        po.SignDate,
		ConsecutiveDays: po.ConsecutiveDays,
		BonusPoints:     po.BonusPoints,
	}
}

func TransactionToPO(t *commerceEntity.Transaction) *TransactionPO {
	if t == nil {
		return nil
	}
	return &TransactionPO{
		ID:           t.ID,
		UserID:       t.UserID,
		Type:         t.Type,
		Amount:       t.Amount,
		BalanceAfter: t.BalanceAfter,
		Memo:         t.Memo,
		Status:       t.Status,
		ReferenceID:  t.ReferenceID,
	}
}

func POToTransaction(po *TransactionPO) *commerceEntity.Transaction {
	if po == nil {
		return nil
	}
	return &commerceEntity.Transaction{
		ID:           po.ID,
		UserID:       po.UserID,
		Type:         po.Type,
		Amount:       po.Amount,
		BalanceAfter: po.BalanceAfter,
		Memo:         po.Memo,
		Status:       po.Status,
		ReferenceID:  po.ReferenceID,
	}
}

func ProductToPO(p *commerceEntity.Product) *ProductPO {
	if p == nil {
		return nil
	}
	return &ProductPO{
		ID:              p.ID,
		Name:            p.Name,
		Description:     p.Description,
		ProductType:     int8(p.ProductType),
		Price:           p.Price,
		WalletPrice:     p.WalletPrice,
		Stock:           p.Stock,
		VideoID:         p.VideoID,
		MembershipLevel: p.MembershipLevel,
		IsActive:        p.IsActive,
	}
}

func POToProduct(po *ProductPO) *commerceEntity.Product {
	if po == nil {
		return nil
	}
	return &commerceEntity.Product{
		ID:              po.ID,
		Name:            po.Name,
		Description:     po.Description,
		ProductType:     commerceEntity.ProductType(po.ProductType),
		Price:           po.Price,
		WalletPrice:     po.WalletPrice,
		Stock:           po.Stock,
		VideoID:         po.VideoID,
		MembershipLevel: po.MembershipLevel,
		IsActive:        po.IsActive,
	}
}

func OrderToPO(o *commerceEntity.Order) *OrderPO {
	if o == nil {
		return nil
	}
	return &OrderPO{
		ID:        o.ID,
		UserID:    o.UserID,
		ProductID: o.ProductID,
		PayType:   int8(o.PayType),
		PricePaid: o.PricePaid,
		Status:    int8(o.Status),
		PaidAt:    o.PaidAt,
	}
}

func POToOrder(po *OrderPO) *commerceEntity.Order {
	if po == nil {
		return nil
	}
	return &commerceEntity.Order{
		ID:        po.ID,
		UserID:    po.UserID,
		ProductID: po.ProductID,
		PayType:   commerceEntity.PayType(po.PayType),
		PricePaid: po.PricePaid,
		Status:    commerceEntity.OrderStatus(po.Status),
		PaidAt:    po.PaidAt,
	}
}

func MembershipToPO(m *commerceEntity.Membership) *MembershipPO {
	if m == nil {
		return nil
	}
	return &MembershipPO{
		ID:       m.ID,
		UserID:   m.UserID,
		Level:    int8(m.Level),
		StartAt:  m.StartAt,
		ExpireAt: m.ExpireAt,
		IsActive: m.IsActive,
	}
}

func POToMembership(po *MembershipPO) *commerceEntity.Membership {
	if po == nil {
		return nil
	}
	return &commerceEntity.Membership{
		ID:       po.ID,
		UserID:   po.UserID,
		Level:    commerceEntity.MembershipLevel(po.Level),
		StartAt:  po.StartAt,
		ExpireAt: po.ExpireAt,
		IsActive: po.IsActive,
	}
}

func VideoPermissionToPO(p *commerceEntity.VideoPermission) *VideoPermissionPO {
	if p == nil {
		return nil
	}
	return &VideoPermissionPO{
		ID:       p.ID,
		UserID:   p.UserID,
		VideoID:  p.VideoID,
		Type:     p.Type,
		ExpireAt: p.ExpireAt,
	}
}

func POToVideoPermission(po *VideoPermissionPO) *commerceEntity.VideoPermission {
	if po == nil {
		return nil
	}
	return &commerceEntity.VideoPermission{
		ID:       po.ID,
		UserID:   po.UserID,
		VideoID:  po.VideoID,
		Type:     po.Type,
		ExpireAt: po.ExpireAt,
	}
}

func FlashSaleToPO(f *commerceEntity.FlashSale) *FlashSalePO {
	if f == nil {
		return nil
	}
	return &FlashSalePO{
		ID:        f.ID,
		Name:      f.Name,
		StartTime: f.StartTime,
		EndTime:   f.EndTime,
		Status:    f.Status,
	}
}

func POToFlashSale(po *FlashSalePO) *commerceEntity.FlashSale {
	if po == nil {
		return nil
	}
	return &commerceEntity.FlashSale{
		ID:        po.ID,
		Name:      po.Name,
		StartTime: po.StartTime,
		EndTime:   po.EndTime,
		Status:    po.Status,
	}
}

func FlashSaleStockToPO(s *commerceEntity.FlashSaleStock) *FlashSaleStockPO {
	if s == nil {
		return nil
	}
	return &FlashSaleStockPO{
		ID:          s.ID,
		FlashSaleID: s.FlashSaleID,
		ProductID:   s.ProductID,
		Stock:       s.Stock,
		RemainStock: s.RemainStock,
		FlashPrice:  s.FlashPrice,
	}
}

func POToFlashSaleStock(po *FlashSaleStockPO) *commerceEntity.FlashSaleStock {
	if po == nil {
		return nil
	}
	return &commerceEntity.FlashSaleStock{
		ID:          po.ID,
		FlashSaleID: po.FlashSaleID,
		ProductID:   po.ProductID,
		Stock:       po.Stock,
		RemainStock: po.RemainStock,
		FlashPrice:  po.FlashPrice,
	}
}

func LotteryDrawToPO(l *commerceEntity.LotteryDraw) *LotteryDrawPO {
	if l == nil {
		return nil
	}
	return &LotteryDrawPO{
		ID:            l.ID,
		Name:          l.Name,
		Description:   l.Description,
		EntryPoints:   l.EntryPoints,
		StartTime:     l.StartTime,
		EndTime:       l.EndTime,
		TotalTickets:  l.TotalTickets,
		RemainTickets: l.RemainTickets,
		Prizes:        l.Prizes,
		Status:        l.Status,
	}
}

func POToLotteryDraw(po *LotteryDrawPO) *commerceEntity.LotteryDraw {
	if po == nil {
		return nil
	}
	return &commerceEntity.LotteryDraw{
		ID:            po.ID,
		Name:          po.Name,
		Description:   po.Description,
		EntryPoints:   po.EntryPoints,
		StartTime:     po.StartTime,
		EndTime:       po.EndTime,
		TotalTickets:  po.TotalTickets,
		RemainTickets: po.RemainTickets,
		Prizes:        po.Prizes,
		Status:        po.Status,
	}
}

func LotteryRecordToPO(r *commerceEntity.LotteryRecord) *LotteryRecordPO {
	if r == nil {
		return nil
	}
	return &LotteryRecordPO{
		ID:         r.ID,
		UserID:     r.UserID,
		LotteryID:  r.LotteryID,
		Prize:      r.Prize,
		IsWin:      r.IsWin,
		DrawNumber: r.DrawNumber,
	}
}

func POToLotteryRecord(po *LotteryRecordPO) *commerceEntity.LotteryRecord {
	if po == nil {
		return nil
	}
	return &commerceEntity.LotteryRecord{
		ID:         po.ID,
		UserID:     po.UserID,
		LotteryID:  po.LotteryID,
		Prize:      po.Prize,
		IsWin:      po.IsWin,
		DrawNumber: po.DrawNumber,
	}
}
