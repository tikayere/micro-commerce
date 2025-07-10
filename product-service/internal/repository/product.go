package repository

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ProductID   string `gorm:"primary_key"`
	TenantID    string `gorm:"index:idx_product_tenant_id"`
	Name        string `gorm:"index:idx_product_name"`
	Description string
	Price       float64
	Stock       int32
	Category    string `gorm:"index:idx_product_category"`
}

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(product *Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) UpdateProduct(product *Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) DeleteProduct(productID, tenantID string) error {
	return r.db.Where("product_id = ? AND tenant_id = ?", productID, tenantID).Delete(&Product{}).Error
}

func (r *ProductRepository) GetProduct(produdctID, tenantID string) (*Product, error) {
	var product Product
	err := r.db.Where("product_id = ? AND tenant_id = ?", produdctID, tenantID).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) ListProducts(tenantID, category string, page, pageSize int32) ([]Product, int, error) {
	var products []Product
	query := r.db.Where("tenant_id = ?", tenantID)
	if category != "" {
		query = query.Where("category = ?", category)
	}
	var total int64
	query.Model(&Product{}).Count(&total)

	err := query.Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}
	return products, int(total), nil
}

func (r *ProductRepository) SearchProducts(tenantID, query string, page, pageSize int32) ([]Product, int, error) {
	var products []Product
	searchQuery := "%" + query + "%"
	queryDB := r.db.Where("tenant_id = ? AND (name ILIKE ? OR description ILIKE ?)", tenantID, searchQuery, searchQuery)
	var total int64
	queryDB.Model(&Product{}).Count(&total)

	err := queryDB.Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}
	return products, int(total), nil
}
