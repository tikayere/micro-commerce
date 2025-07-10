package repository

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       string `gorm:"primary_key"`
	Email    string `gorm:"unique_index"`
	Password string
	Name     string
	Role     string
	TenantID string `gorm:"index"`
}

type Role struct {
	gorm.Model
	Name        string
	Permissions []Permission `gorm:"many2many:role_permissions"`
}

type Permission struct {
	gorm.Model
	Name string `gorm:"unique_index"`
}

type RolePermission struct {
	gorm.Model
	RoleID       uint
	PermissionID uint
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return r.db.Create(user).Error
}

func (r *UserRepository) FindUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindUserByID(id, tenantID string) (*User, error) {
	var user User
	err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(user *User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) CheckPermission(userID, tenantID, permission string) (bool, error) {
	var count int64
	err := r.db.Raw(`
	SELECT COUNT (*)
	FROM users u
	JOIN roles r ON u.role = r.name
	JOIN role_permissions rp ON r.id = rp.role_id
	JOIN permissions p ON rp.permission_id = p.id
	WHERE u.id = ? AND u.tenant_id = ? AND p.name = ?
	`, userID, tenantID, permission).Count(&count).Error
	return count > 0, err
}
