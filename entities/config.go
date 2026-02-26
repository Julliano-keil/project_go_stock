package entities

type Config struct {
	DSN       string
	JWTSecret string // chave para assinatura dos tokens JWT
}

type CompanyDatabaseConfig struct{}
