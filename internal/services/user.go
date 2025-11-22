package service

import (
	"context"
	"fmt"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/marcelorc13/timesheet-pro/internal/domain"
	"github.com/marcelorc13/timesheet-pro/internal/repository"
)

type UserService struct {
	repository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{repository: userRepository}
}

func (us UserService) GetUsuarios(ctx context.Context) (*[]domain.Usuario, error) {
	res, err := us.repository.GetUsuarios(ctx)

	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, nil
	}

	usuarios, ok := res.Data.([]domain.Usuario)

	if !ok {
		return nil, fmt.Errorf("erro ao converter dados")
	}

	return &usuarios, nil
}

func (us UserService) GetUsuario(ctx context.Context, id string) (*domain.Usuario, error) {
	res, err := us.repository.GetUsuario(ctx, id)

	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, nil
	}

	usuario, ok := res.Data.(domain.Usuario)

	if !ok {
		return nil, fmt.Errorf("erro ao converter dados")
	}

	return &usuario, nil
}

func (us UserService) DeleteUsuario(ctx context.Context, id int) error {
	res, err := us.repository.DeleteUsuario(ctx, id)

	if err != nil {
		return err
	}

	if !res.Success {
		return fmt.Errorf("%s", res.Message)
	}

	return nil
}

func (us UserService) CreateUsuario(ctx context.Context, u domain.Usuario) error {
	validate := validator.New()
	u.ID = uuid.New()
	err := validate.Struct(u)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	res, err := us.repository.CreateUsuario(ctx, u)

	if err != nil {
		return err
	}

	if !res.Success {
		return fmt.Errorf("%s", res.Message)
	}

	return nil
}

func (us UserService) Login(ctx context.Context, u domain.LoginUsuario) (*domain.LoginUsuario, error) {
	// validate := validator.New()
	// err := validate.Struct(u)
	// if err != nil {
	// 	return nil, err.(validator.ValidationErrors)
	// }

	res, err := us.repository.Login(ctx, u)

	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, fmt.Errorf("%s", res.Message)
	}

	usuario, ok := res.Data.(domain.LoginUsuario)

	if !ok {
		return nil, fmt.Errorf("erro ao converter dados")
	}

	return &usuario, nil
}
