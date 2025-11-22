package repository

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marcelorc13/timesheet-pro/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) List(ctx context.Context) (domain.DBResponse, error) {
	results, err := r.DB.Query(ctx, "SELECT id, name, email, password FROM users")
	if err != nil {
		return domain.DBResponse{Message: "Ocorreu um erro na query"}, err
	}

	res := []domain.User{}

	for results.Next() {
		var user domain.User

		err = results.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			panic(err)
		}

		res = append(res, user)
	}

	if len(res) == 0 {
		return domain.DBResponse{Message: "O banco ainda não possui usuários"}, nil
	}

	return domain.DBResponse{Success: true, Data: res}, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (domain.DBResponse, error) {
	var user domain.User

	err := r.DB.QueryRow(ctx, "SELECT id, name, email, password FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		return domain.DBResponse{Message: "usuário não encontrado"}, nil
	} else if err != nil {
		return domain.DBResponse{Message: err.Error()}, err
	}
	return domain.DBResponse{Success: true, Data: user}, nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) (domain.DBResponse, error) {
	res, err := r.DB.Exec(ctx, "DELETE FROM usuarios WHERE id = $1", id)
	if err != nil {
		return domain.DBResponse{Message: "ocorreu um erro na query"}, err
	}

	rows := res.RowsAffected()
	if rows != 1 {
		return domain.DBResponse{Message: "usuário não encontrado"}, nil
	}

	return domain.DBResponse{Success: true}, nil
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) (domain.DBResponse, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return domain.DBResponse{Message: "erro ao hashear password"}, nil
	}

	res, err := r.DB.Exec(ctx, `
		INSERT INTO users(name, email, password)
		VALUES($1, $2, $3);
	`, u.Name, u.Email, string(passwordBytes))
	if err != nil {
		return domain.DBResponse{Message: "erro ao criar usuario"}, err
	}

	rows := res.RowsAffected()
	if rows != 1 {
		return domain.DBResponse{Message: "erro ao criar usuario"}, err
	}

	return domain.DBResponse{Success: true}, nil
}

func (r *UserRepository) Login(ctx context.Context, u domain.LoginUser) (domain.DBResponse, error) {
	var usuario domain.LoginUser
	err := r.DB.QueryRow(ctx, "SELECT id, email, password FROM users WHERE email = $1", u.Email).
		Scan(&usuario.ID, &usuario.Email, &usuario.Password)

	if err == sql.ErrNoRows {
		return domain.DBResponse{Message: "usuário não encontrado"}, nil
	} else if err != nil {
		return domain.DBResponse{Message: err.Error()}, err
	}

	errpassword := bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(u.Password))

	if errpassword != nil {
		return domain.DBResponse{Message: "password incorreta"}, nil
	}
	return domain.DBResponse{Success: true, Data: usuario}, nil
}
