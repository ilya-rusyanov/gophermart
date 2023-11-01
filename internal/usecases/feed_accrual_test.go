package usecases

import (
	"context"
	"testing"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
	"github.com/ilya-rusyanov/gophermart/internal/usecases/mocks"
	"go.uber.org/mock/gomock"
)

func TestFeedAccrual(t *testing.T) {
	t.Run("increasing user balance", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		storage := mocks.NewMockAccrualStorage(ctrl)
		storage.EXPECT().GetUnfinishedOrders(ctx).
			Return(
				entities.OrderList{
					entities.Order{Status: entities.OrderStatusNew},
				},
				nil,
			)
	})
}
