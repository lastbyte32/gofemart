package order

import (
	"context"
	"fmt"
	"sync"

	"github.com/lastbyte32/gofemart/internal/domain"
	"github.com/lastbyte32/gofemart/internal/storage"
	"github.com/lastbyte32/gofemart/internal/util/luhn"
)

type accrualGetter interface {
	GetOrder(ctx context.Context, number string) (*domain.AccrualOrderInfo, error)
}

type Order interface {
	IsOrderNumberValid(number string) bool
	UploadOrder(ctx context.Context, number, userId string) error
	GetOrdersByUserID(ctx context.Context, userID string) ([]*domain.Order, error)
}

type order struct {
	store  storage.Order
	client accrualGetter
}

func (s *order) GetOrdersByUserID(ctx context.Context, userID string) ([]*domain.Order, error) {
	return s.store.GetOrdersByUserID(ctx, userID)
}

func (s *order) UploadOrder(ctx context.Context, number, userId string) error {
	newOrder := domain.Order{
		Number: number,
		UserID: userId,
		Status: domain.OrderNew,
	}
	order, err := s.store.GetByNumber(ctx, number)
	if err != nil {
		return err
	}

	if order != nil && order.UserID == userId {
		return domain.ErrDuplicateOrderUploadSameUser
	}

	if order != nil && order.UserID != userId {
		return domain.ErrDuplicateOrderUploadOtherUser
	}

	if _, err := s.store.Create(ctx, newOrder); err != nil {
		return err
	}
	return nil
}

func (s *order) IsOrderNumberValid(number string) bool {
	return luhn.Validation(number)
}

func (s *order) orderProcessing(ctx context.Context, wg *sync.WaitGroup, accrual accrualGetter, number string) {
	fmt.Printf("Processing order: %s\n", number)

	info, err := accrual.GetOrder(ctx, number)
	if err != nil {
		fmt.Println("Failed to request order status from accrual system" + err.Error())
		wg.Done()
		return
	}

	err = s.store.UpdateOrder(ctx, info)
	if err != nil {
		fmt.Println("updateOrder failed: " + err.Error())
		wg.Done()
		return
	}

}
func (s *order) WorkerAccrual(ctx context.Context, client accrualGetter) {
	var wg sync.WaitGroup
	fmt.Println("start worker")
	for {
		orders, err := s.store.GetOrdersUnprocessed(ctx)
		if err != nil || len(orders) == 0 {
			continue
		}
		wg.Add(len(orders))
		for _, order := range orders {
			go s.orderProcessing(ctx, &wg, client, order.Number)
		}
		wg.Wait()
	}
}

func NewService(ctx context.Context, store storage.Order, client accrualGetter) Order {
	o := &order{
		client: client,
		store:  store,
	}

	o.WorkerAccrual(ctx, client)

	return o
}
