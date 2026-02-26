package user

import (
	"context"
	"database/sql"

	"lince/datastore"
	"lince/entities"
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
	err := db.QueryRowContext(ctx,
		"SELECT id, nome, email, senha, salt FROM usuario WHERE email = ?",
		email,
	).Scan(&u.ID, &u.Nome, &u.Email, &u.Senha, &u.Salt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
