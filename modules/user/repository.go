package user

import (
	"context"
	"database/sql"

	"lince/datastore"
	"lince/entities"

	"golang.org/x/crypto/bcrypt"
)

type userRepository struct {
	conn func(company entities.CompanyDatabaseConfig) *sql.DB
}

func NewUserRepository(settings datastore.SettingsRepository) datastore.UserRepository {
	return userRepository{conn: settings.Connection}
}

func (r userRepository) GetByEmail(ctx context.Context, company entities.CompanyDatabaseConfig, email string) (*entities.Usuario, error) {
	db := r.conn(company)
	var u entities.Usuario

	query := `
	SELECT id,
	nome, 
	email, 
	senha 
	FROM usuario 
	WHERE email = ?
	`

	err := db.QueryRowContext(ctx,
		query,
		email,
	).Scan(&u.ID, &u.Nome, &u.Email, &u.Senha)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r userRepository) Create(ctx context.Context, company entities.CompanyDatabaseConfig, nome, email, senha string) (*entities.Usuario, error) {
	db := r.conn(company)

	query := `
	INSERT INTO usuario (nome, email, senha) VALUES (?, ?, ?)
	`

	hash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(ctx, query, nome, email, string(hash))
	if err != nil {
		return nil, err
	}

	user, err := r.GetByEmail(ctx, company, email)
	if err != nil {
		return nil, err
	}

	return user, nil

}
