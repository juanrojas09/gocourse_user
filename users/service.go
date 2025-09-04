package users

import (
	"errors"
	"log"
	"strings"

	"github.com/juanrojas09/gocourse_domain/domain"
)

type (
	Service interface {
		Create(firstName, lastName, email, phone string) (*domain.User, error)
		Update(req *UpdateReq) (*UpdateReq, error)
		Delete(id string) (string, error)
		GetById(id string) (*domain.User, error)
		GetAll(filters Filters, offset int, limit int) ([]domain.User, error)
		Count(filters Filters) (int, error)
	} //interface del servicio
	service struct { //representa las propiedades de la estructura del servicio "clase"
		r      UserRepository
		logger *log.Logger
	}

	Filters struct {
		FirstName string
		LastName  string
		Page      int
		Offset    int
	}
)

func NewService(r UserRepository, logger *log.Logger) Service {
	return &service{
		r:      r,
		logger: logger,
	}
}

func (s *service) Create(firstName, lastName, email, phone string) (*domain.User, error) {
	user := domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}
	s.logger.Printf("Usuario a crear: %s", user)
	usr, err := s.r.Create(&user)
	if err != nil {
		return nil, err
	}

	return usr, nil

}

func (s *service) Update(req *UpdateReq) (*UpdateReq, error) {

	if strings.TrimSpace(*req.ID) == "" {
		return nil, errors.New("User id cannot be null on update")
	}

	s.logger.Printf("Usuario a crear: %s", req)
	err := s.r.Update(req)
	if err != nil {
		return nil, err
	}

	return req, nil

}
func (s *service) Delete(id string) (string, error) {

	s.logger.Printf("Usuario con id a borrar: %s", id)
	err := s.r.Delete(id)
	if err != nil {
		return "", err
	}

	return id, nil

}
func (s *service) GetById(Id string) (*domain.User, error) {

	s.logger.Printf("Usuario a obtener con el id: %s", Id)
	usr, err := s.r.GetById(Id)
	if err != nil {
		return nil, err
	}

	return usr, nil

}
func (s *service) GetAll(filters Filters, offset int, limit int) ([]domain.User, error) {

	s.logger.Printf("Obteniendo todos los usuarios")

	usr, err := s.r.GetAll(filters, offset, limit)
	if err != nil {
		return nil, err
	}

	return usr, nil

}

func (s *service) Count(filters Filters) (int, error) {
	return s.r.Count(filters)
}
