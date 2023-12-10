package domain

import (
	"context"
	"time"
)

// Client is table client data
type Client struct {
	ID          int64
	InsuranceNo string
	NIK         string
	Name        string
	BirthDate   time.Time
	ProvID      int
	KabID       int
	KecID       int
	KelID       int
	Address     string
	PhoneNumber string
	Email       *string
	Active      bool
	Log         *string
	DtmCrt      time.Time
	DtmUpd      time.Time
}

// AuthUsecase is auth usecase
type AuthUsecase interface {
	Authorize(ctx context.Context, token string) (client Client, err error)
}

// AuthMySQLRepository is auth repository in MySQL
type AuthMySQLRepository interface {
	GetAllClient(ctx context.Context, page, limit int64, sort, order string, calegId ...string) (client []Client, err error)
	GetClientByID(ctx context.Context, id int64) (client Client, err error)
	InsertClient(ctx context.Context, client Client) (err error)
	UpdateClient(ctx context.Context, client Client, audit string) (err error)
}

// AuthRedisRepository is auth repository in Redis
type AuthRedisRepository interface {
	SetRedis(ctx context.Context, key string, value string) (err error)
	GetRedis(ctx context.Context, key string) (value string, err error)
}
