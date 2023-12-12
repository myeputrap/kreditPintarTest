package domain

import (
	"context"
	"time"
)

// Client is table client data
type Client struct {
	ID           int64
	NIK          string
	Name         string
	BirthDate    time.Time
	BirthOfPlace string
	PhoneNumber  string
	Email        string
	DtmCrt       time.Time
	DtmUpd       time.Time
}

type LoginRequest struct {
	IsAdmin     bool   `json:"is_admin"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Username    string `json:"user_name"`
	Password    string `json:"password"`
}

type IsAdminResponse struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type Session struct {
	ID           int64  `json:"ID" redis:"ID"`
	NIK          string `json:"NIK" redis:"NIK"`
	Name         string `json:"Name" redis:"Name"`
	BirthDate    string `json:"BirthDate" redis:"BirthDate"`
	Address      string `json:"Address" redis:"Address"`
	PhoneNumber  string `json:"PhoneNumber" redis:"PhoneNumber"`
	Email        string `json:"Email" redis:"Email"`
	Username     string `json:"UserName" redis:"UserName"`
	IsAdmin      bool   `json:"IsAdmin" redis:"IsAdmin"`
	IsConsumer   bool   `json:"IsConsumer" redis:"IsConsumer"`
	IsCreditUser bool   `json:"IsCreditUser" redis:"IsCreditUser"`
}

// AuthUsecase is auth usecase
type AuthUsecase interface {
	Authorize(ctx context.Context, token string) (client Client, err error)
	PostLogin(ctx context.Context, req LoginRequest) (res LoginResponse, err error)
	PostLogout(ctx context.Context, req string) (err error)
	IsAdmin(ctx context.Context, token string) (err error)
	IsConsumer(ctx context.Context, token string) (err error)
}

// AuthMySQLRepository is auth repository in MySQL
type AuthMySQLRepository interface {
	GetAllClient(ctx context.Context, page, limit int64, sort, order string, calegId ...string) (client []Client, err error)
	GetClientByID(ctx context.Context, id int64) (client Client, err error)
	InsertClient(ctx context.Context, client Client) (err error)
	LoginAdmin(ctx context.Context, req LoginRequest) (res IsAdminResponse, err error)
	LoginConsumer(ctx context.Context, req LoginRequest) (res Client, err error)
}

// AuthRedisRepository is auth repository in Redis
type AuthRedisRepository interface {
	SetRedis(ctx context.Context, key string, value string) (err error)
	GetRedis(ctx context.Context, key string) (value string, err error)
	HSetRedis(ctx context.Context, key string, fields map[string]interface{}, expTime int64) (err error)
	HGetRedis(ctx context.Context, key string) (res Session, err error)
	HDelRedis(ctx context.Context, key string) error
	HExistsRedis(ctx context.Context, key string) (bool, error)
}
