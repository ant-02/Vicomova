package mock

import (
	"context"

	commerceEntity "vicomova/internal/commerce/domain/entity"
	commerceRepo "vicomova/internal/commerce/domain/repository"
)

type MockPointsWalletRepository struct {
	Wallets             map[int64]*commerceEntity.PointsWallet
	CreateErr           error
	GetByIDErr          error
	UpdateErr           error
	UpdateVersionErr    error
	UpdateVersionResult bool // defaults to true (success), set to false to simulate failure
}

func NewMockPointsWalletRepository() *MockPointsWalletRepository {
	return &MockPointsWalletRepository{
		Wallets:             make(map[int64]*commerceEntity.PointsWallet),
		UpdateVersionResult: true, // default to success
	}
}

func (m *MockPointsWalletRepository) GetByUserID(ctx context.Context, userID int64) (*commerceEntity.PointsWallet, error) {
	if m.GetByIDErr != nil {
		return nil, m.GetByIDErr
	}
	return m.Wallets[userID], nil
}

func (m *MockPointsWalletRepository) Create(ctx context.Context, wallet *commerceEntity.PointsWallet) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	wallet.ID = int64(len(m.Wallets) + 1)
	m.Wallets[wallet.UserID] = wallet
	return nil
}

func (m *MockPointsWalletRepository) Update(ctx context.Context, wallet *commerceEntity.PointsWallet) error {
	if m.UpdateErr != nil {
		return m.UpdateErr
	}
	m.Wallets[wallet.UserID] = wallet
	return nil
}

func (m *MockPointsWalletRepository) UpdateVersion(ctx context.Context, userID int64, oldVersion int64, newVersion int64, delta int64) (bool, error) {
	if m.UpdateVersionErr != nil {
		return false, m.UpdateVersionErr
	}
	wallet := m.Wallets[userID]
	if wallet == nil {
		return false, nil
	}
	if wallet.Version != oldVersion {
		return false, nil
	}
	wallet.Version = newVersion
	wallet.Points += delta
	return m.UpdateVersionResult, nil
}

var _ commerceRepo.PointsWalletRepository = (*MockPointsWalletRepository)(nil)

type MockSignInRecordRepository struct {
	Records     map[int64][]*commerceEntity.SignInRecord
	CreateErr   error
	GetTodayErr error
	GetLastErr  error
}

func NewMockSignInRecordRepository() *MockSignInRecordRepository {
	return &MockSignInRecordRepository{
		Records: make(map[int64][]*commerceEntity.SignInRecord),
	}
}

func (m *MockSignInRecordRepository) GetTodaySignIn(ctx context.Context, userID int64) (*commerceEntity.SignInRecord, error) {
	if m.GetTodayErr != nil {
		return nil, m.GetTodayErr
	}
	records := m.Records[userID]
	if len(records) == 0 {
		return nil, nil
	}
	return records[len(records)-1], nil
}

func (m *MockSignInRecordRepository) GetLastSignIn(ctx context.Context, userID int64) (*commerceEntity.SignInRecord, error) {
	if m.GetLastErr != nil {
		return nil, m.GetLastErr
	}
	records := m.Records[userID]
	if len(records) == 0 {
		return nil, nil
	}
	return records[len(records)-1], nil
}

func (m *MockSignInRecordRepository) Create(ctx context.Context, record *commerceEntity.SignInRecord) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	record.ID = int64(len(m.Records[record.UserID]) + 1)
	m.Records[record.UserID] = append(m.Records[record.UserID], record)
	return nil
}

func (m *MockSignInRecordRepository) ListByUser(ctx context.Context, userID int64, cursor int64, limit int) ([]*commerceEntity.SignInRecord, bool, error) {
	records := m.Records[userID]
	return records, false, nil
}

var _ commerceRepo.SignInRecordRepository = (*MockSignInRecordRepository)(nil)

type MockTransactionRepository struct {
	Transactions map[int64][]*commerceEntity.Transaction
	CreateErr    error
}

func NewMockTransactionRepository() *MockTransactionRepository {
	return &MockTransactionRepository{
		Transactions: make(map[int64][]*commerceEntity.Transaction),
	}
}

func (m *MockTransactionRepository) Create(ctx context.Context, tx *commerceEntity.Transaction) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	tx.ID = int64(len(m.Transactions[tx.UserID]) + 1)
	m.Transactions[tx.UserID] = append(m.Transactions[tx.UserID], tx)
	return nil
}

func (m *MockTransactionRepository) GetByID(ctx context.Context, id int64) (*commerceEntity.Transaction, error) {
	return nil, nil
}

func (m *MockTransactionRepository) ListByUser(ctx context.Context, userID int64, cursor int64, limit int) ([]*commerceEntity.Transaction, bool, error) {
	return m.Transactions[userID], false, nil
}

var _ commerceRepo.TransactionRepository = (*MockTransactionRepository)(nil)

type MockProductRepository struct {
	Products   map[int64]*commerceEntity.Product
	CreateErr  error
	GetByIDErr error
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{
		Products: make(map[int64]*commerceEntity.Product),
	}
}

func (m *MockProductRepository) GetByID(ctx context.Context, id int64) (*commerceEntity.Product, error) {
	if m.GetByIDErr != nil {
		return nil, m.GetByIDErr
	}
	return m.Products[id], nil
}

func (m *MockProductRepository) List(ctx context.Context, productType commerceEntity.ProductType, cursor string, limit int) ([]*commerceEntity.Product, string, bool, error) {
	var result []*commerceEntity.Product
	for _, p := range m.Products {
		if productType == 0 || p.ProductType == productType {
			result = append(result, p)
		}
	}
	return result, "", false, nil
}

func (m *MockProductRepository) ListByIDs(ctx context.Context, ids []int64) ([]*commerceEntity.Product, error) {
	var result []*commerceEntity.Product
	for _, id := range ids {
		if p := m.Products[id]; p != nil {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *MockProductRepository) ListFlashSaleProducts(ctx context.Context) ([]*commerceEntity.Product, error) {
	return nil, nil
}

func (m *MockProductRepository) ListActive(ctx context.Context) ([]*commerceEntity.Product, error) {
	var result []*commerceEntity.Product
	for _, p := range m.Products {
		if p.IsActive {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *MockProductRepository) Create(ctx context.Context, product *commerceEntity.Product) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	product.ID = int64(len(m.Products) + 1)
	m.Products[product.ID] = product
	return nil
}

func (m *MockProductRepository) UpdateStock(ctx context.Context, id int64, delta int) error {
	if p := m.Products[id]; p != nil {
		p.Stock += delta
	}
	return nil
}

var _ commerceRepo.ProductRepository = (*MockProductRepository)(nil)

type MockOrderRepository struct {
	Orders     map[int64]*commerceEntity.Order
	CreateErr  error
	GetByIDErr error
}

func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{
		Orders: make(map[int64]*commerceEntity.Order),
	}
}

func (m *MockOrderRepository) Create(ctx context.Context, order *commerceEntity.Order) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	order.ID = int64(len(m.Orders) + 1)
	m.Orders[order.ID] = order
	return nil
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id int64) (*commerceEntity.Order, error) {
	if m.GetByIDErr != nil {
		return nil, m.GetByIDErr
	}
	return m.Orders[id], nil
}

func (m *MockOrderRepository) GetByUserAndProduct(ctx context.Context, userID, productID int64) (*commerceEntity.Order, error) {
	return nil, nil
}

func (m *MockOrderRepository) UpdateStatus(ctx context.Context, id int64, status commerceEntity.OrderStatus) error {
	if o := m.Orders[id]; o != nil {
		o.Status = status
	}
	return nil
}

func (m *MockOrderRepository) ListByUser(ctx context.Context, userID int64, cursor string, limit int) ([]*commerceEntity.Order, bool, error) {
	var result []*commerceEntity.Order
	for _, o := range m.Orders {
		if o.UserID == userID {
			result = append(result, o)
		}
	}
	return result, false, nil
}

var _ commerceRepo.OrderRepository = (*MockOrderRepository)(nil)

type MockMembershipRepository struct {
	Memberships map[int64]*commerceEntity.Membership
	CreateErr   error
}

func NewMockMembershipRepository() *MockMembershipRepository {
	return &MockMembershipRepository{
		Memberships: make(map[int64]*commerceEntity.Membership),
	}
}

func (m *MockMembershipRepository) GetActiveMembership(ctx context.Context, userID int64) (*commerceEntity.Membership, error) {
	return m.Memberships[userID], nil
}

func (m *MockMembershipRepository) GetByID(ctx context.Context, id int64) (*commerceEntity.Membership, error) {
	return nil, nil
}

func (m *MockMembershipRepository) Create(ctx context.Context, membership *commerceEntity.Membership) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	membership.ID = int64(len(m.Memberships) + 1)
	m.Memberships[membership.UserID] = membership
	return nil
}

func (m *MockMembershipRepository) Update(ctx context.Context, membership *commerceEntity.Membership) error {
	m.Memberships[membership.UserID] = membership
	return nil
}

func (m *MockMembershipRepository) ListByUser(ctx context.Context, userID int64) ([]*commerceEntity.Membership, error) {
	return nil, nil
}

var _ commerceRepo.MembershipRepository = (*MockMembershipRepository)(nil)

type MockVideoPermissionRepository struct {
	Permissions map[int64]map[int64]*commerceEntity.VideoPermission
	CreateErr   error
}

func NewMockVideoPermissionRepository() *MockVideoPermissionRepository {
	return &MockVideoPermissionRepository{
		Permissions: make(map[int64]map[int64]*commerceEntity.VideoPermission),
	}
}

func (m *MockVideoPermissionRepository) Check(ctx context.Context, userID, videoID int64) (bool, error) {
	if perms, ok := m.Permissions[userID]; ok {
		_, exists := perms[videoID]
		return exists, nil
	}
	return false, nil
}

func (m *MockVideoPermissionRepository) Get(ctx context.Context, userID, videoID int64) (*commerceEntity.VideoPermission, error) {
	if perms, ok := m.Permissions[userID]; ok {
		return perms[videoID], nil
	}
	return nil, nil
}

func (m *MockVideoPermissionRepository) Create(ctx context.Context, p *commerceEntity.VideoPermission) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	p.ID = int64(len(m.Permissions[p.UserID]) + 1)
	if m.Permissions[p.UserID] == nil {
		m.Permissions[p.UserID] = make(map[int64]*commerceEntity.VideoPermission)
	}
	m.Permissions[p.UserID][p.VideoID] = p
	return nil
}

func (m *MockVideoPermissionRepository) ListByUser(ctx context.Context, userID int64) ([]*commerceEntity.VideoPermission, error) {
	return nil, nil
}

func (m *MockVideoPermissionRepository) ListByVideo(ctx context.Context, videoID int64) ([]*commerceEntity.VideoPermission, error) {
	return nil, nil
}

var _ commerceRepo.VideoPermissionRepository = (*MockVideoPermissionRepository)(nil)

type MockFlashSaleRepository struct {
	FlashSales map[int64]*commerceEntity.FlashSale
	CreateErr  error
}

func NewMockFlashSaleRepository() *MockFlashSaleRepository {
	return &MockFlashSaleRepository{
		FlashSales: make(map[int64]*commerceEntity.FlashSale),
	}
}

func (m *MockFlashSaleRepository) GetActiveFlashSale(ctx context.Context) (*commerceEntity.FlashSale, error) {
	return nil, nil
}

func (m *MockFlashSaleRepository) GetByID(ctx context.Context, id int64) (*commerceEntity.FlashSale, error) {
	return m.FlashSales[id], nil
}

func (m *MockFlashSaleRepository) Create(ctx context.Context, fs *commerceEntity.FlashSale) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	fs.ID = int64(len(m.FlashSales) + 1)
	m.FlashSales[fs.ID] = fs
	return nil
}

func (m *MockFlashSaleRepository) ListActive(ctx context.Context) ([]*commerceEntity.FlashSale, error) {
	var result []*commerceEntity.FlashSale
	for _, fs := range m.FlashSales {
		result = append(result, fs)
	}
	return result, nil
}

var _ commerceRepo.FlashSaleRepository = (*MockFlashSaleRepository)(nil)

type MockFlashSaleStockRepository struct {
	Stocks map[int64]map[int64]*commerceEntity.FlashSaleStock
}

func NewMockFlashSaleStockRepository() *MockFlashSaleStockRepository {
	return &MockFlashSaleStockRepository{
		Stocks: make(map[int64]map[int64]*commerceEntity.FlashSaleStock),
	}
}

func (m *MockFlashSaleStockRepository) GetByFlashSaleID(ctx context.Context, flashSaleID int64) ([]*commerceEntity.FlashSaleStock, error) {
	return nil, nil
}

func (m *MockFlashSaleStockRepository) GetStock(ctx context.Context, flashSaleID, productID int64) (*commerceEntity.FlashSaleStock, error) {
	if stocks, ok := m.Stocks[flashSaleID]; ok {
		return stocks[productID], nil
	}
	return nil, nil
}

func (m *MockFlashSaleStockRepository) Create(ctx context.Context, stock *commerceEntity.FlashSaleStock) error {
	stock.ID = int64(len(m.Stocks[stock.FlashSaleID]) + 1)
	if m.Stocks[stock.FlashSaleID] == nil {
		m.Stocks[stock.FlashSaleID] = make(map[int64]*commerceEntity.FlashSaleStock)
	}
	m.Stocks[stock.FlashSaleID][stock.ProductID] = stock
	return nil
}

func (m *MockFlashSaleStockRepository) DecrementStock(ctx context.Context, flashSaleID, productID int64) (bool, error) {
	if stocks, ok := m.Stocks[flashSaleID]; ok {
		if stock := stocks[productID]; stock != nil && stock.RemainStock > 0 {
			stock.RemainStock--
			return true, nil
		}
	}
	return false, nil
}

var _ commerceRepo.FlashSaleStockRepository = (*MockFlashSaleStockRepository)(nil)

type MockLotteryRepository struct {
	Lotteries map[int64]*commerceEntity.LotteryDraw
	CreateErr error
}

func NewMockLotteryRepository() *MockLotteryRepository {
	return &MockLotteryRepository{
		Lotteries: make(map[int64]*commerceEntity.LotteryDraw),
	}
}

func (m *MockLotteryRepository) GetActiveLottery(ctx context.Context) (*commerceEntity.LotteryDraw, error) {
	return nil, nil
}

func (m *MockLotteryRepository) GetByID(ctx context.Context, id int64) (*commerceEntity.LotteryDraw, error) {
	return m.Lotteries[id], nil
}

func (m *MockLotteryRepository) Create(ctx context.Context, l *commerceEntity.LotteryDraw) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	l.ID = int64(len(m.Lotteries) + 1)
	m.Lotteries[l.ID] = l
	return nil
}

func (m *MockLotteryRepository) ListActive(ctx context.Context) ([]*commerceEntity.LotteryDraw, error) {
	var result []*commerceEntity.LotteryDraw
	for _, l := range m.Lotteries {
		result = append(result, l)
	}
	return result, nil
}

var _ commerceRepo.LotteryRepository = (*MockLotteryRepository)(nil)

type MockLotteryRecordRepository struct {
	Records map[int64][]*commerceEntity.LotteryRecord
}

func NewMockLotteryRecordRepository() *MockLotteryRecordRepository {
	return &MockLotteryRecordRepository{
		Records: make(map[int64][]*commerceEntity.LotteryRecord),
	}
}

func (m *MockLotteryRecordRepository) RecordParticipation(ctx context.Context, record *commerceEntity.LotteryRecord) error {
	record.ID = int64(len(m.Records[record.UserID]) + 1)
	m.Records[record.UserID] = append(m.Records[record.UserID], record)
	return nil
}

func (m *MockLotteryRecordRepository) GetUserParticipation(ctx context.Context, userID, lotteryID int64) (*commerceEntity.LotteryRecord, error) {
	return nil, nil
}

func (m *MockLotteryRecordRepository) CountByUser(ctx context.Context, userID, lotteryID int64) (int, error) {
	return len(m.Records[userID]), nil
}

var _ commerceRepo.LotteryRecordRepository = (*MockLotteryRecordRepository)(nil)
