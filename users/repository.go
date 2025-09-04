package users

import (
	"fmt"
	"log"
	"strings"

	"github.com/juanrojas09/gocourse_domain/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *domain.User) (*domain.User, error)
	GetAll(filters Filters, offset int, limit int) ([]domain.User, error)
	GetById(id string) (*domain.User, error)
	Delete(id string) error
	Update(req *UpdateReq) error
	Count(filters Filters) (int, error)
}

type repository struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewRepository(db *gorm.DB, logger *log.Logger) UserRepository {
	return &repository{
		db:     db,
		logger: logger,
	}
}

func (r *repository) buildErrorResponse(err error) error {
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Count(filters Filters) (int, error) {
	var users domain.User
	var count int64
	tx := r.db.Model(&users)
	tx = applyFilters(filters, tx)
	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *repository) Create(user *domain.User) (*domain.User, error) {

	if res := r.db.Create(user); res.Error != nil {
		r.logger.Println(res.Error)
		return nil, res.Error
	}
	log.Println("User created successfully", user.ID)
	return user, nil
}

func (r *repository) GetAll(filters Filters, offset int, limit int) ([]domain.User, error) {

	var users []domain.User
	tx := r.db.Model(&users)
	tx = applyFilters(filters, tx)
	tx = tx.Offset(offset).Limit(limit)
	res := tx.Order("created_at desc").Find(&users)
	if err := r.buildErrorResponse(res.Error); err != nil {
		return nil, err
	}
	r.logger.Println(users)
	return users, nil

}

func (r *repository) GetById(Id string) (*domain.User, error) {

	// res := r.db.Find(&user, User{ID: Id})
	// otra opcion
	var user = domain.User{ID: Id}
	// r.db.First(&user) e implica que hace el where con el id de la instancia de arriba
	// db.First(&user) suponiendo que el objeto ya tiene el id seteado que llega por param
	res := r.db.First(&user)
	if err := r.buildErrorResponse(res.Error); err != nil {
		return nil, err
	}
	return &user, nil

}

func (r *repository) Update(req *UpdateReq) error {
	values := make(map[string]interface{})
	if req.FirstName != nil {
		values["first_name"] = req.FirstName
	}
	if req.LastName != nil {
		values["last_name"] = req.LastName
	}
	if req.Email != nil {
		values["email"] = req.Email
	}
	if req.Phone != nil {
		values["phone"] = req.Phone
	}
	res := r.db.Model(&domain.User{}).Where("id = ?", req.ID).Updates(values)
	if err := r.buildErrorResponse(res.Error); err != nil {
		return err
	}
	return nil
}

func (r *repository) Delete(id string) error {
	res := r.db.Delete(&domain.User{ID: id})

	if err := r.buildErrorResponse(res.Error); err != nil {
		return err
	}

	return nil
}

func applyFilters(filters Filters, tx *gorm.DB) *gorm.DB {
	if filters.FirstName != "" {
		filters.FirstName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.FirstName))
		tx = tx.Where("lower(first_name) like ?", filters.FirstName)
	}

	if filters.LastName != "" {
		filters.FirstName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.LastName))
		tx = tx.Where("lower(last_name) like ?", filters.LastName)
	}
	return tx

}
