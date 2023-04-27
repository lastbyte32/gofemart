package accrual

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/lastbyte32/gofemart/internal/domain"
)

const defaultTimeout = time.Second * 10

type service struct {
	client *resty.Client
}

func New(accrualUrl string) *service {
	client := resty.New().SetTimeout(defaultTimeout).SetBaseURL(accrualUrl)

	return &service{
		client: client,
	}
}

func (s *service) GetOrder(ctx context.Context, number string) (*domain.AccrualOrderInfo, error) {
	url := fmt.Sprintf("/api/orders/%s", number)
	var order domain.AccrualOrderInfo

	response, err := s.client.R().SetResult(order).SetContext(ctx).Get(url)
	if err != nil {
		return nil, errors.New("accrual is not available")
	}
	if response.StatusCode() != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("accrual api call failed: %s", response.Status()))
	}

	return &order, nil
}
