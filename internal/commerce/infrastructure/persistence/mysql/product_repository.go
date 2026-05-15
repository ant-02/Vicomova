package mysql

import (
	"context"
	"strconv"

	commerceEntity "vicomova/internal/commerce/domain/entity"
	commerceRepo "vicomova/internal/commerce/domain/repository"
	sharedMysql "vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/log"

	"gorm.io/gorm"
)

type ProductRepository struct {
	mysql *sharedMysql.Client
}

func NewProductRepository(mysqlClient *sharedMysql.Client) commerceRepo.ProductRepository {
	return &ProductRepository{mysql: mysqlClient}
}

func (r *ProductRepository) Create(ctx context.Context, product *commerceEntity.Product) error {
	po := ProductToPO(product)
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("ProductRepository.Create: failed: %v", err)
		return err
	}
	product.ID = po.ID
	return nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*commerceEntity.Product, error) {
	var po ProductPO
	err := r.mysql.WithContext(ctx).First(&po, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error.Printf("ProductRepository.GetByID: failed: %v", err)
		return nil, err
	}
	return POToProduct(&po), nil
}

func (r *ProductRepository) List(ctx context.Context, productType commerceEntity.ProductType, cursor string, limit int) ([]*commerceEntity.Product, string, bool, error) {
	var pos []ProductPO
	query := r.mysql.WithContext(ctx).Where("is_active = ? AND stock > 0", true)
	if productType > 0 {
		query = query.Where("product_type = ?", int8(productType))
	}
	if cursor != "" {
		if id, err := strconv.ParseInt(cursor, 10, 64); err == nil {
			query = query.Where("id < ?", id)
		}
	}
	if err := query.Order("id DESC").Limit(limit + 1).Find(&pos).Error; err != nil {
		log.Error.Printf("ProductRepository.List: failed: %v", err)
		return nil, "", false, err
	}
	hasMore := len(pos) > limit
	if hasMore {
		pos = pos[:limit]
	}
	nextCursor := ""
	if len(pos) > 0 {
		nextCursor = strconv.FormatInt(pos[len(pos)-1].ID, 10)
	}
	products := make([]*commerceEntity.Product, len(pos))
	for i := range pos {
		products[i] = POToProduct(&pos[i])
	}
	return products, nextCursor, hasMore, nil
}

func (r *ProductRepository) ListByIDs(ctx context.Context, ids []int64) ([]*commerceEntity.Product, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var pos []ProductPO
	if err := r.mysql.WithContext(ctx).Where("id IN ?", ids).Find(&pos).Error; err != nil {
		log.Error.Printf("ProductRepository.ListByIDs: failed: %v", err)
		return nil, err
	}
	products := make([]*commerceEntity.Product, len(pos))
	for i := range pos {
		products[i] = POToProduct(&pos[i])
	}
	return products, nil
}

func (r *ProductRepository) ListFlashSaleProducts(ctx context.Context) ([]*commerceEntity.Product, error) {
	var pos []ProductPO
	if err := r.mysql.WithContext(ctx).Where("product_type = ? AND is_active = ? AND stock > 0", int8(commerceEntity.ProductTypeFlashSale), true).Find(&pos).Error; err != nil {
		log.Error.Printf("ProductRepository.ListFlashSaleProducts: failed: %v", err)
		return nil, err
	}
	products := make([]*commerceEntity.Product, len(pos))
	for i := range pos {
		products[i] = POToProduct(&pos[i])
	}
	return products, nil
}

func (r *ProductRepository) ListActive(ctx context.Context) ([]*commerceEntity.Product, error) {
	var pos []ProductPO
	if err := r.mysql.WithContext(ctx).Where("is_active = ? AND stock > 0", true).Find(&pos).Error; err != nil {
		log.Error.Printf("ProductRepository.ListActive: failed: %v", err)
		return nil, err
	}
	products := make([]*commerceEntity.Product, len(pos))
	for i := range pos {
		products[i] = POToProduct(&pos[i])
	}
	return products, nil
}

func (r *ProductRepository) UpdateStock(ctx context.Context, id int64, delta int) error {
	if err := r.mysql.WithContext(ctx).Model(&ProductPO{}).Where("id = ?", id).
		Update("stock", gorm.Expr("stock + ?", delta)).Error; err != nil {
		log.Error.Printf("ProductRepository.UpdateStock: failed: %v", err)
		return err
	}
	return nil
}
