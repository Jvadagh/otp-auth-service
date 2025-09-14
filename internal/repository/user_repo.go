package repository

import (
	"github.com/jvadagh/otp-auth-service/internal/model"
	"gorm.io/gorm"
	"time"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByPhone(phone string) (*model.User, error) {
	var u model.User
	if err := r.db.Where("phone_number = ?", phone).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByID(id uint) (*model.User, error) {
	var u model.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) Create(phone string) (*model.User, error) {
	u := &model.User{
		PhoneNumber: phone,
		CreatedAt:   time.Now(),
	}
	if err := r.db.Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) List(search string, page, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	q := r.db.Model(&model.User{})
	if search != "" {
		q = q.Where("phone_number ILIKE ?", "%"+search+"%")
	}
	q.Count(&total)
	offset := (page - 1) * limit
	if err := q.Order("created_at desc").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
