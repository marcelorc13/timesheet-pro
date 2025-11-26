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

func (us UserService) List(ctx context.Context) (*[]domain.User, error) {
	res, err := us.repository.List(ctx)

	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, nil
	}

	usuarios, ok := res.Data.([]domain.User)

	if !ok {
		return nil, fmt.Errorf("erro ao converter dados")
	}

	return &usuarios, nil
}

func (us UserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	res, err := us.repository.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, nil
	}

	usuario, ok := res.Data.(domain.User)

	if !ok {
		return nil, fmt.Errorf("erro ao converter dados")
	}

	return &usuario, nil
}

func (us UserService) Delete(ctx context.Context, id int) error {
	res, err := us.repository.Delete(ctx, id)

	if err != nil {
		return err
	}

	if !res.Success {
		return fmt.Errorf("%s", res.Message)
	}

	return nil
}

func (us UserService) Create(ctx context.Context, u domain.User) error {
	validate := validator.New()
	u.ID = uuid.New()
	err := validate.Struct(u)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	res, err := us.repository.Create(ctx, u)

	if err != nil {
		return err
	}

	if !res.Success {
		return fmt.Errorf("%s", res.Message)
	}

	return nil
}

func (us UserService) Login(ctx context.Context, u domain.LoginUser) (*domain.User, error) {
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

	usuario, ok := res.Data.(domain.User)

	if !ok {
		return nil, fmt.Errorf("erro ao converter dados")
	}

	return &usuario, nil
}

// UpdateProfile updates a user's profile information
func (us UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, name, email string) (*domain.User, error) {
	// Validate name is not empty
	if name == "" {
		return nil, fmt.Errorf("nome não pode ser vazio")
	}
	
	// Validate email format
	validate := validator.New()
	if err := validate.Var(email, "required,email"); err != nil {
		return nil, fmt.Errorf("email inválido")
	}
	
	// Check if email is already in use by another user
	existingRes, err := us.repository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	
	if existingRes.Success {
		existingUser, ok := existingRes.Data.(domain.User)
		if ok && existingUser.ID != userID {
			return nil, fmt.Errorf("email já está em uso")
		}
	}
	
	// Update user
	res, err := us.repository.UpdateUser(ctx, userID.String(), name, email)
	if err != nil {
		return nil, err
	}
	
	if !res.Success {
		return nil, fmt.Errorf("%s", res.Message)
	}
	
	updatedUser, ok := res.Data.(domain.User)
	if !ok {
		return nil, fmt.Errorf("erro ao converter dados")
	}
	
	return &updatedUser, nil
}
