package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"goKreditPintar/domain"
	"goKreditPintar/helper"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

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
	data, err := au.authRedisRepo.GetRedis(ctx, "id"+token)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(data), &client)

	return
}

func (au *authUsecase) PostLogin(ctx context.Context, req domain.LoginRequest) (res domain.LoginResponse, err error) {
	var isAdmin, isConsumer bool
	var outAdm domain.IsAdminResponse
	var outCust domain.Client
	if req.IsAdmin {
		outAdm, err = au.authMySQLRepo.LoginAdmin(ctx, req)
		if err != nil {
			log.Error("error usecase LoginAdmin ", err.Error())
			return domain.LoginResponse{}, err
		}
		err = bcrypt.CompareHashAndPassword([]byte(outAdm.Password), []byte(req.Password))
		if err != nil {
			err = errors.New("answer request isn't correct")
			return
		}
		isAdmin = true
		isConsumer = false
	} else {
		if req.PhoneNumber == "" {
			err = errors.New("phone number is empty")
			return
		}

		req.PhoneNumber, err = helper.NumberUniformity(req.PhoneNumber)
		if err != nil {
			log.Error("error in NumberUniformity PostLogin err : ", err.Error())
			err = errors.New("phone number is not valid")
			return
		}
		outCust, err = au.authMySQLRepo.LoginConsumer(ctx, req)
		if err != nil {
			log.Error("error usecase LoginConsumer ", err.Error())
			return domain.LoginResponse{}, err
		}
		isAdmin = false
		isConsumer = true
	}
	var variable string
	if req.IsAdmin {
		variable = outAdm.Username
	} else {
		variable = outCust.PhoneNumber
	}
	expTime := viper.GetInt64("expire_time_redis")
	date := time.Now()
	dateStr := date.Format(time.RFC3339)
	tokenByte := sha256.Sum256([]byte(variable + "_" + dateStr))
	token := base64.URLEncoding.EncodeToString(tokenByte[:])
	res.Token = token

	fields := map[string]interface{}{
		"ID":          outCust.ID,
		"NIK":         outCust.NIK,
		"Name":        outCust.Name,
		"BirthDate":   outCust.BirthDate.Format("2006-01-02"),
		"PhoneNumber": outCust.PhoneNumber,
		"Email":       outCust.Email,
		"Username":    outAdm.Username,
		"IsAdmin":     isAdmin,
		"IsConsumer":  isConsumer,
	}
	err = au.authRedisRepo.HSetRedis(ctx, "id:"+token, fields, expTime)
	if err != nil {
		log.Error("error setRedis HSetRedis, err is ", err.Error())
		return
	}

	return
}
func (au *authUsecase) PostLogout(ctx context.Context, req string) (err error) {
	exists, err := au.authRedisRepo.HExistsRedis(ctx, "id:"+req)
	if err != nil {
		log.Error("error checking Redis key: ", err.Error())
	}

	if !exists {
		err = errors.New("unauthorize token")
		return
	}

	err = au.authRedisRepo.HDelRedis(ctx, "id:"+req)
	if err != nil {
		log.Error("error delete HDelRedis , err is ", err.Error())
		return
	}
	return
}

func (au *authUsecase) IsAdmin(ctx context.Context, token string) (err error) {
	var resRedis domain.Session
	resRedis, err = au.authRedisRepo.HGetRedis(ctx, "id:"+token)
	if err != nil {
		err = errors.New("redis not found")
		return
	}
	if !resRedis.IsAdmin {
		err = errors.New("not valid")
	}
	return
}
func (au *authUsecase) IsConsumer(ctx context.Context, token string) (err error) {
	var resRedis domain.Session
	resRedis, err = au.authRedisRepo.HGetRedis(ctx, "id:"+token)
	if err != nil {
		err = errors.New("redis not found")
		return
	}
	if !resRedis.IsConsumer {
		err = errors.New("not valid")
	}
	return
}
