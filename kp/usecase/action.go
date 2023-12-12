package usecase

import (
	"context"
	"errors"
	"goKreditPintar/domain"
	"goKreditPintar/helper"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type actionUsecase struct {
	actionMySQLRepo domain.ActionMySQLRepository
	authRedisRepo   domain.AuthRedisRepository
}

// NewActionUsecase is constructor of usecase
func NewActionUsecase(actionMySQLRepo domain.ActionMySQLRepository, authRedisRepo domain.AuthRedisRepository) domain.ActionUsecase {
	return &actionUsecase{
		actionMySQLRepo: actionMySQLRepo,
		authRedisRepo:   authRedisRepo,
	}
}

func (u *actionUsecase) GetConsumer(ctx context.Context, req domain.GetAllConsumerRequest) (res domain.GetAllConsumerResponse, err error) {
	log.Debug("start of getcust")
	total, err := u.actionMySQLRepo.CountConsumer(ctx, domain.GetAllConsumerRequest{
		Salary: req.Salary,
		Name:   req.Name,
		Email:  req.Email,
	})
	if err != nil {
		return
	}

	res.MetaData = domain.MetaData{
		TotalData: uint(total),
		TotalPage: (uint(total) + uint(req.Limit) - 1) / uint(req.Limit),
		Page:      uint(req.Page),
		Limit:     uint(req.Limit),
		Sort:      req.Sort,
		Order:     req.Order,
	}

	offset := (req.Page - 1) * req.Limit
	res.Data, err = u.actionMySQLRepo.GetConsumer(ctx, domain.GetAllConsumerRequest{
		Salary: req.Salary,
		Name:   req.Name,
		Email:  req.Email,
		Page:   offset,
		Limit:  req.Limit,
		Sort:   req.Sort,
		Order:  req.Order,
	})

	return
}
func (u *actionUsecase) GetCreditCard(ctx context.Context, req domain.GetAllCreditCardRequest) (res domain.GetAllCreditCardResponse, err error) {
	log.Debug("start of getcust")
	total, err := u.actionMySQLRepo.CountCreditCard(ctx, domain.GetAllCreditCardRequest{
		CardNumber:     req.CardNumber,
		ExpirationDate: req.ExpirationDate,
		CreditLimit:    req.CreditLimit,
		CurrentBalance: req.CurrentBalance,
	})
	if err != nil {
		return
	}

	res.MetaData = domain.MetaData{
		TotalData: uint(total),
		TotalPage: (uint(total) + uint(req.Limit) - 1) / uint(req.Limit),
		Page:      uint(req.Page),
		Limit:     uint(req.Limit),
		Sort:      req.Sort,
		Order:     req.Order,
	}

	offset := (req.Page - 1) * req.Limit
	res.Data, err = u.actionMySQLRepo.GetCreditCard(ctx, domain.GetAllCreditCardRequest{
		CardNumber:     req.CardNumber,
		ExpirationDate: req.ExpirationDate,
		CreditLimit:    req.CreditLimit,
		CurrentBalance: req.CurrentBalance,
		Page:           offset,
		Limit:          req.Limit,
		Sort:           req.Sort,
		Order:          req.Order,
	})

	for i, v := range res.Data {
		number, _ := strconv.Atoi(v.ConsumerID)
		resDetail, err := u.actionMySQLRepo.GetConsumerDetail(ctx, number)
		if err != nil {
			return domain.GetAllCreditCardResponse{}, err
		}
		res.Data[i].Name = resDetail.Name
		res.Data[i].Nik = resDetail.Nik
		res.Data[i].CardNumber = helper.FormatCardNumber(res.Data[i].CardNumber)
	}
	return
}

func (u *actionUsecase) PostConsumer(ctx context.Context, req domain.Consumer) (err error) {
	log.Debug("start of postcust")
	err = u.actionMySQLRepo.CheckValidate(ctx, "nik", req.Nik)
	if err != nil {
		return
	}
	err = u.actionMySQLRepo.CheckValidate(ctx, "phone_number", req.PhoneNumber)
	if err != nil {
		return
	}
	err = u.actionMySQLRepo.PostConsumer(ctx, req)
	if err != nil {
		return
	}
	return
}

func (u *actionUsecase) GetConsumerDetail(ctx context.Context, req int) (res domain.Consumer, err error) {

	res, err = u.actionMySQLRepo.GetConsumerDetail(ctx, req)
	if err != nil {
		return
	}
	return
}

func (u *actionUsecase) GetCreditCardDetail(ctx context.Context, req int) (res domain.ConsumerCreditCard, err error) {

	res, err = u.actionMySQLRepo.GetCreditCardDetail(ctx, req)
	if err != nil {
		return
	}
	number, _ := strconv.Atoi(res.ConsumerID)
	resDetail, err := u.actionMySQLRepo.GetConsumerDetail(ctx, number)
	if err != nil {
		return domain.ConsumerCreditCard{}, err
	}
	res.Name = resDetail.Name
	res.Nik = resDetail.Nik
	res.CardNumber = helper.FormatCardNumber(res.CardNumber)
	return
}

func (u *actionUsecase) PostConsumerCreditCard(ctx context.Context, req domain.ConsumerCreditCard) (err error) {
	log.Debug("start of postcust")

	out, err := u.actionMySQLRepo.GetConsumerByParameter(ctx, "nik", req.Nik)
	if err != nil {
		return
	}
	req.ConsumerID = strconv.Itoa(int(out.ID))

	count, err := u.actionMySQLRepo.CheckValidateCC(ctx, "consumer_id", strconv.Itoa(int(out.ID)))
	if err != nil {
		return
	}

	if count > 0 {
		err = errors.New("data exist")
		return
	}

	requestLimitFloat, err := strconv.ParseFloat(req.RequestLimit, 64)
	if err != nil {
		return
	}

	salaryFloat, err := strconv.ParseFloat(out.Salary, 64)
	if err != nil {
		return
	}

	if requestLimitFloat < 0.5*salaryFloat {
		req.CreditLimit = req.RequestLimit
	} else {
		valCL := 0.5 * salaryFloat
		req.CreditLimit = strconv.Itoa(int(valCL))
	}
	err = u.actionMySQLRepo.PostConsumerCreditCard(ctx, req)
	if err != nil {
		return
	}
	return
}
func (u *actionUsecase) PostTransactionCredit(ctx context.Context, req domain.TransactionCreditCard, token string) (err error) {

	resRedis, err := u.authRedisRepo.HGetRedis(ctx, "id:"+token)
	if err != nil {
		err = errors.New("redis not found")
		return
	}
	out, err := u.actionMySQLRepo.GetConsumerByParameter(ctx, "id", strconv.Itoa(int(resRedis.ID)))
	if err != nil {
		return
	}
	req.ConsumerID = strconv.Itoa(int(out.ID))
	purchaseAmount, err := strconv.ParseFloat(req.PurchaseAmount, 32)
	if err != nil {
		return
	}
	otr := purchaseAmount + (purchaseAmount * req.InterestAmount)
	req.OTR = strconv.FormatFloat(otr, 'f', -1, 64)

	id, _ := strconv.Atoi(req.ConsumerID)
	res, err := u.actionMySQLRepo.GetCreditCardDetail(ctx, id)
	if err != nil {
		return
	}
	if req.CVV == "" || req.CVV != res.CVV {
		err = errors.New("wrong cvv")
		return
	}
	limit, _ := strconv.ParseFloat(res.CreditLimit, 64)
	balance, _ := strconv.ParseFloat(res.CurrentBalance, 64)
	if otr > (limit - balance) {
		err = errors.New("overlimit")
		return
	}
	if otr < 250000 && req.InstallmentCount != "1" {
		err = errors.New("if otr below certain amount, installment have to be 1")
		return
	} else if req.InstallmentCount != "3" && req.InstallmentCount != "6" && req.InstallmentCount != "12" && req.InstallmentCount != "1" {
		err = errors.New("please choose available count")
		return
	}
	/*
		if otr < 250000 {
				req.InstallmentCount = "1"
			} else if req.InstallmentCount != "3" && req.InstallmentCount != "6" && req.InstallmentCount != "12" {
				err = errors.New("please choose available count")
				return
			}*/
	req.ContractNumber = "AS24" + res.CardNumber[0:4] + out.PhoneNumber[2:6]

	countID, err := u.actionMySQLRepo.PostTransactionCredit(ctx, req)
	if err != nil {
		return
	}
	req.TransactionID = strconv.Itoa(int(countID))
	installments, _ := strconv.Atoi(req.InstallmentCount)
	billings := make([]domain.Billing, installments)
	equalInstallment := otr / float64(installments)
	remainingAmount := otr

	for i := 0; i < installments; i++ {
		dueDate := time.Now().AddDate(0, i+1, 0)
		var amount float64
		if i == installments-1 {
			amount = remainingAmount
		} else {
			amount = equalInstallment
			remainingAmount -= equalInstallment
		}
		consID, _ := strconv.Atoi(req.ConsumerID)
		transacID, _ := strconv.Atoi(req.TransactionID)
		billings[i] = domain.Billing{
			ConsumerID:     consID,
			TransactionID:  transacID,
			ContractNumber: req.ContractNumber,
			DueDate:        dueDate,
			BillAmount:     amount,
			Status:         "Pending",
		}
		err = u.actionMySQLRepo.PostBillingCredit(ctx, billings[i])
		if err != nil {
			return
		}
		err = u.actionMySQLRepo.UpdateCreditBalance(ctx, otr, req.ConsumerID)
		if err != nil {
			return
		}
	}

	return
}
func (u *actionUsecase) PatchBilling(ctx context.Context, req int, token string) (err error) {
	resRedis, err := u.authRedisRepo.HGetRedis(ctx, "id:"+token)
	if err != nil {
		err = errors.New("redis not found")
		return
	}
	outConsumer, err := u.actionMySQLRepo.GetConsumerByParameter(ctx, "id", strconv.Itoa(int(resRedis.ID)))
	if err != nil {
		return
	}
	out, err := u.actionMySQLRepo.CheckValidateBilling(ctx, "consumer_id", strconv.Itoa(int(outConsumer.ID)))
	if err != nil {
		return
	}
	if out < 1 {
		err = errors.New("not found")
		return
	}

	outs, err := u.actionMySQLRepo.CheckValidateBilling(ctx, "id", strconv.Itoa(int(req)))
	if err != nil {
		return
	}
	if outs < 1 {
		err = errors.New("not found")
		return
	}

	err = u.actionMySQLRepo.PatchBilling(ctx, strconv.Itoa(int(req)))
	if err != nil {
		return
	}
	otr, err := u.actionMySQLRepo.CountCreditBalance(ctx, strconv.Itoa(int(outConsumer.ID)))
	if err != nil {
		return
	}
	err = u.actionMySQLRepo.UpdateCreditBalance(ctx, otr, strconv.Itoa(int(outConsumer.ID)))
	if err != nil {
		return
	}
	return
}

func (u *actionUsecase) GetBilling(ctx context.Context, req domain.GetAllBillingRequest) (res domain.GetAllBillingResponse, err error) {
	log.Debug("start of getcust")
	total, err := u.actionMySQLRepo.CountBilling(ctx, domain.GetAllBillingRequest{
		BillAmount: req.BillAmount,
		DueDate:    req.DueDate,
		Status:     req.Status,
	})
	if err != nil {
		return
	}

	res.MetaData = domain.MetaData{
		TotalData: uint(total),
		TotalPage: (uint(total) + uint(req.Limit) - 1) / uint(req.Limit),
		Page:      uint(req.Page),
		Limit:     uint(req.Limit),
		Sort:      req.Sort,
		Order:     req.Order,
	}

	offset := (req.Page - 1) * req.Limit
	res.Data, err = u.actionMySQLRepo.GetBilling(ctx, domain.GetAllBillingRequest{
		BillAmount: req.BillAmount,
		DueDate:    req.DueDate,
		Status:     req.Status,
		Page:       offset,
		Limit:      req.Limit,
		Sort:       req.Sort,
		Order:      req.Order,
	})
	if err != nil {
		return
	}

	return
}
