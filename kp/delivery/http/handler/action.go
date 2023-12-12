package handler

import (
	"errors"
	"goKreditPintar/domain"
	"goKreditPintar/helper"
	"math/rand"
	"strconv"
	"time"

	//"goKreditPintar/helper"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

type ActionHandler struct {
	ActionUsecase domain.ActionUsecase
}

func (ah *ActionHandler) GetConsumer(c *fiber.Ctx) (err error) {
	var input domain.GetAllConsumerRequest

	input.Limit, err = strconv.ParseInt(c.Query("limit", strconv.FormatInt(int64(viper.GetInt("dafault_limit_query")), 10)), 10, 64)
	if err != nil {
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}

	input.Page, err = strconv.ParseInt(c.Query("page", "1"), 10, 64)
	if err != nil {
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}
	if c.Query("name") != "" {
		input.Name = c.Query("name")
	}

	if c.Query("min_salary") != "" {
		input.Salary = c.Query("min_salary")
	}
	input.Sort = c.Query("sort", "id")
	input.Order = c.Query("order", "asc")
	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		log.Errorf("error validator PostLogin %s", err.Error())
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}

	res, err := ah.ActionUsecase.GetConsumer(c.Context(), input)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(fasthttp.StatusBadRequest).SendString("The data you filled is incorrect.")
		}
		if strings.Contains(err.Error(), "not found") {
			return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
		}
		log.Errorf("error  PostLogin %s", err.Error())
		return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
	}

	err = c.JSON(res)

	return err
}

func (ah *ActionHandler) GetBilling(c *fiber.Ctx) (err error) {
	var input domain.GetAllBillingRequest

	input.Limit, err = strconv.ParseInt(c.Query("limit", strconv.FormatInt(int64(viper.GetInt("dafault_limit_query")), 10)), 10, 64)
	if err != nil {
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}

	input.Page, err = strconv.ParseInt(c.Query("page", "1"), 10, 64)
	if err != nil {
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}
	if c.Query("bill_amount") != "" {
		input.BillAmount = c.Query("bill_amount")
	}

	if c.Query("due_date") != "" {
		input.DueDate = c.Query("due_date")
	}
	if c.Query("status") != "" {
		input.Status = c.Query("status")
	}
	input.Sort = c.Query("sort", "id")
	input.Order = c.Query("order", "asc")
	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		log.Errorf("error validator PostLogin %s", err.Error())
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}

	res, err := ah.ActionUsecase.GetBilling(c.Context(), input)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(fasthttp.StatusBadRequest).SendString("The data you filled is incorrect.")
		}
		if strings.Contains(err.Error(), "not found") {
			return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
		}
		log.Errorf("error  PostLogin %s", err.Error())
		return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
	}

	err = c.JSON(res)

	return err
}

func (ah *ActionHandler) PostConsumer(c *fiber.Ctx) (err error) {
	var input domain.Consumer
	err = c.BodyParser(&input)
	if err != nil {
		log.Errorf("error bodyparser Post %s", err.Error())
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}
	input.PlaceOfBirth = c.FormValue("place_of_birth")
	input.PhoneNumber = c.FormValue("phone_number")
	input.BirthDate = c.FormValue("birth_date")
	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		log.Errorf("error validator Post %s", err.Error())
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}
	input.PhoneNumber, err = helper.NumberUniformity(input.PhoneNumber)
	if err != nil {
		log.Error("error in NumberUniformity PostLogin err : ", err.Error())
		err = errors.New("phone number is not valid")
		return
	}
	err = ah.ActionUsecase.PostConsumer(c.Context(), input)
	if err != nil {
		if strings.Contains(err.Error(), "data exist") {
			return c.Status(fasthttp.StatusBadRequest).SendString("The data you input is exist.")
		}
		if strings.Contains(err.Error(), "not found") {
			return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
		}
		log.Errorf("error  Post %s", err.Error())
		return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
	}

	return helper.HTTPSimpleResponse(c, fasthttp.StatusCreated)
}

func (ah *ActionHandler) GetConsumerDetail(c *fiber.Ctx) (err error) {
	id, err := c.ParamsInt("id")
	if err != nil {
		log.Error(err)
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}
	res, err := ah.ActionUsecase.GetConsumerDetail(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(fasthttp.StatusBadRequest).SendString("The data you filled is incorrect.")
		}
		if strings.Contains(err.Error(), "not found") {
			return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
		}
		log.Errorf("error  PostLogin %s", err.Error())
		return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
	}

	err = c.JSON(res)

	return err
}
func (ah *ActionHandler) PatchBilling(c *fiber.Ctx) (err error) {
	var token string
	auth := c.GetReqHeaders()
	authorization := auth["Authorization"]
	if len(authorization) != 0 {
		token = authorization[0][7:]
	}
	if token == "" {
		return helper.HTTPSimpleResponse(c, fasthttp.StatusUnauthorized)
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		log.Error(err)
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}
	err = ah.ActionUsecase.PatchBilling(c.Context(), id, token)
	if err != nil {
		if strings.Contains(err.Error(), "data exist") {
			return c.Status(fasthttp.StatusBadRequest).SendString("The data you input is exist.")
		}
		if strings.Contains(err.Error(), "consumer id not valid") {
			return c.Status(fasthttp.StatusBadRequest).SendString(err.Error())
		}
		if strings.Contains(err.Error(), "not found") {
			return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
		}
		log.Errorf("error  Post %s", err.Error())
		return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
	}

	return helper.HTTPSimpleResponse(c, fasthttp.StatusOK)

}
func (ah *ActionHandler) PostCreditCards(c *fiber.Ctx) (err error) {
	var input domain.ConsumerCreditCard
	input.Nik = c.FormValue("nik")
	input.RequestLimit = c.FormValue("request_limit")

	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		log.Errorf("error validator Post %s", err.Error())
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}
	input.CardNumber = helper.GenerateValidCardNumber()
	currentDate := time.Now()
	expirationDate := currentDate.AddDate(3, 0, 0) //3 default expire date creidt card
	input.ExpirationDate = expirationDate.Format("2006-01-02")
	input.CurrentBalance = "0"
	numericBytes := "0123456789"
	b := make([]byte, 3)
	for i := range b {
		b[i] = numericBytes[rand.Intn(len(numericBytes))]
	}
	input.CVV = string(b)
	err = ah.ActionUsecase.PostConsumerCreditCard(c.Context(), input)
	if err != nil {
		if strings.Contains(err.Error(), "data exist") {
			return c.Status(fasthttp.StatusBadRequest).SendString("The data you input is exist.")
		}
		if strings.Contains(err.Error(), "consumer id not valid") {
			return c.Status(fasthttp.StatusBadRequest).SendString(err.Error())
		}
		if strings.Contains(err.Error(), "not found") {
			return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
		}
		log.Errorf("error  Post %s", err.Error())
		return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
	}

	return helper.HTTPSimpleResponse(c, fasthttp.StatusCreated)

}

func (ah *ActionHandler) PostTransactionDetail(c *fiber.Ctx) (err error) {
	var token string
	auth := c.GetReqHeaders()
	authorization := auth["Authorization"]
	if len(authorization) != 0 {
		token = authorization[0][7:]
	}
	if token == "" {
		return helper.HTTPSimpleResponse(c, fasthttp.StatusUnauthorized)
	}

	var input domain.TransactionCreditCard
	input.ProductType = c.FormValue("product_type")
	input.PurchaseAmount = c.FormValue("purchase_amount")
	input.InstallmentCount = c.FormValue("installment_count")
	input.AssetName = c.FormValue("asset_name")
	input.CVV = c.FormValue("cvv")

	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		log.Errorf("error validator Post %s", err.Error())
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}

	input.InterestAmount = 0.15

	err = ah.ActionUsecase.PostTransactionCredit(c.Context(), input, token)
	if err != nil {
		if strings.Contains(err.Error(), "data exist") {
			return c.Status(fasthttp.StatusBadRequest).SendString("The data you input is exist.")
		}
		if strings.Contains(err.Error(), "consumer id not valid") {
			return c.Status(fasthttp.StatusBadRequest).SendString(err.Error())
		}
		if strings.Contains(err.Error(), "not found") {
			return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
		}
		log.Errorf("error  Post %s", err.Error())
		return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
	}

	return helper.HTTPSimpleResponse(c, fasthttp.StatusCreated)

}

func (ah *ActionHandler) GetCreditCards(c *fiber.Ctx) (err error) {
	var input domain.GetAllCreditCardRequest
	input.Limit, err = strconv.ParseInt(c.Query("limit", strconv.FormatInt(int64(viper.GetInt("dafault_limit_query")), 10)), 10, 64)
	if err != nil {
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}

	input.Page, err = strconv.ParseInt(c.Query("page", "1"), 10, 64)
	if err != nil {
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}
	if c.Query("card_number") != "" {
		input.CardNumber = c.Query("card_number")
	}

	if c.Query("expiration_date") != "" {
		input.ExpirationDate = c.Query("expiration_date")
	}

	if c.Query("credit_limit") != "" {
		input.CreditLimit = c.Query("credit_limit")
	}

	if c.Query("current_balance") != "" {
		input.CurrentBalance = c.Query("current_balance")
	}
	input.Sort = c.Query("sort", "id")
	input.Order = c.Query("order", "asc")
	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		log.Errorf("error validator PostLogin %s", err.Error())
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}

	res, err := ah.ActionUsecase.GetCreditCard(c.Context(), input)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(fasthttp.StatusBadRequest).SendString("The data you filled is incorrect.")
		}
		if strings.Contains(err.Error(), "not found") {
			return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
		}
		log.Errorf("error   %s", err.Error())
		return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
	}

	err = c.JSON(res)

	return err
}

func (ah *ActionHandler) GetCreditCardDetail(c *fiber.Ctx) (err error) {
	id, err := c.ParamsInt("id")
	if err != nil {
		log.Error(err)
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}
	res, err := ah.ActionUsecase.GetCreditCardDetail(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(fasthttp.StatusBadRequest).SendString("The data you filled is incorrect.")
		}
		if strings.Contains(err.Error(), "not found") {
			return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
		}
		log.Errorf("error   %s", err.Error())
		return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
	}

	err = c.JSON(res)

	return err
}
