package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/lastbyte32/gofemart/internal/domain"
	"github.com/lastbyte32/gofemart/internal/storage"
)

type accrualGetter interface {
	GetOrder(ctx context.Context, number string) (*domain.AccrualOrderInfo, error)
}

type work struct {
	client accrualGetter
	store  storage.Order
}

func (s *work) orderProcessing(ctx context.Context, wg *sync.WaitGroup, number string) {
	fmt.Printf("Processing order: %s\n", number)

	info, err := s.client.GetOrder(ctx, number)
	if err != nil {
		fmt.Println("accrual request failed: " + err.Error())
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

func (s *work) Run(ctx context.Context) {
	fmt.Println("start worker")
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	var wg sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			fmt.Println("stop worker")
			return
		case <-ticker.C:
			orders, err := s.store.GetOrdersUnprocessed(ctx)
			if err != nil || len(orders) == 0 {
				continue
			}
			wg.Add(len(orders))
			for _, order := range orders {
				go s.orderProcessing(ctx, &wg, order.Number)
			}
			wg.Wait()
		}
	}
}

func New(store storage.Order, client accrualGetter) *work {
	return &work{
		store:  store,
		client: client,
	}
}
