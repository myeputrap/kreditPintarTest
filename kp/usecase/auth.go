package usecase

import (
	"context"
	"encoding/json"
	"goKreditPintar/domain"
)

type APIResponse struct {
	Kode    string `json:"kode"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type APIRequest struct {
	Key         string `json:"key"`
	NomorKartu  string `json:"nokartu"`
	Nama        string `json:"nama"`
	Nik         string `json:"nik"`
	BirthDate   string `json:"tgllahir"`
	Alamat      string `json:"alamat"`
	NumberPhone string `json:"nohp"`
	Email       string `json:"email"`
	HeirName    string `json:"ahliwaris"`
	HeirPhone   string `json:"nohpahliwaris"`
	CalegDpr    string `json:"calegdpr"`
	CalegDpr1   string `json:"calegdpr1"`
	CalegDpr2   string `json:"calegdpr2"`
	Dapil       string `json:"dapil"`
	Plan        string `json:"plan"`
}

type authUsecase struct {
	authMySQLRepo domain.AuthMySQLRepository
	authRedisRepo domain.AuthRedisRepository
}

// NewAuthUsecase is constructor of usecase
func NewAuthUsecase(authMySQLRepo domain.AuthMySQLRepository, authRedisRepo domain.AuthRedisRepository) domain.AuthUsecase {
	return &authUsecase{
		authMySQLRepo: authMySQLRepo,
		authRedisRepo: authRedisRepo,
	}
}

// AuthorizeAuth is usecase for authorize to service be-service-auth
func (au *authUsecase) Authorize(ctx context.Context, token string) (client domain.Client, err error) {
	data, err := au.authRedisRepo.GetRedis(ctx, token)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(data), &client)

	return
}
