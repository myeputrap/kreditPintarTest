package usecase_test

import (
	"context"
	"errors"
	"goKreditPintar/domain"
	"goKreditPintar/domain/mocks"
	"testing"
	"time"

	// . "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestGetConsumer_Positive(t *testing.T) {
	// Create an instance of the generated mock
	mock := &mocks.ActionUsecase{}

	expectedResponse := domain.GetAllConsumerResponse{
		MetaData: domain.MetaData{
			TotalData: 100,
			TotalPage: 5,
			Page:      1,
			Limit:     20,
			Sort:      "name",
			Order:     "asc",
		},
		Data: []domain.Consumer{
			{
				ID:           1,
				Name:         "John Doe",
				Nik:          "1234567890",
				PhoneNumber:  "123456789",
				BirthDate:    "1990-01-01",
				PlaceOfBirth: "New York",
				Salary:       "50000",
				Email:        "john@example.com",
				DtmCrt:       time.Now(),
				DtmUpd:       time.Now(),
			},
		},
	}

	request := domain.GetAllConsumerRequest{
		Page:   1,
		Limit:  10,
		Sort:   "name",
		Order:  "asc",
		Name:   "John Doe",
		Salary: "50000",
		Email:  "john@example.com",
	}
	var errExpected error = nil

	mock.On("GetConsumer", context.Background(), request).Return(expectedResponse, errExpected)

	_, err := mock.GetConsumer(context.Background(), request)

	assert.NoError(t, err)
}

func TestGetConsumer_Negative(t *testing.T) {
	mock := &mocks.ActionUsecase{}

	expectedResponse := domain.GetAllConsumerResponse{}
	request := domain.GetAllConsumerRequest{
		Page:   1,
		Limit:  10,
		Sort:   "name",
		Order:  "asc",
		Name:   "John Doe",
		Salary: "50000",
		Email:  "john@example.com",
	}

	expectedError := errors.New("something went wrong")

	mock.On("GetConsumer", context.Background(), request).Return(expectedResponse, expectedError)

	_, err := mock.GetConsumer(context.Background(), request)

	assert.Error(t, err)
	assert.EqualError(t, err, "something went wrong")
}
