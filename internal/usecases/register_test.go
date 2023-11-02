package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
	"github.com/ilya-rusyanov/gophermart/internal/usecases/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRegister(t *testing.T) {
	var dummyLogger DummyLogger

	ctx := context.Background()

	tests := []struct {
		name         string
		findReturns  error
		storeReturns error
		wantUCReturn error
	}{
		{
			name:         "successfull registration",
			findReturns:  entities.ErrNotFound,
			storeReturns: nil,
			wantUCReturn: nil,
		},
		{
			name:         "user already exists",
			findReturns:  nil,
			storeReturns: nil,
			wantUCReturn: entities.ErrLoginConflict,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			creds := entities.AuthCredentials{
				Login:    "user",
				Password: "password",
			}

			m := mocks.NewMockRegisterStorage(ctrl)
			m.EXPECT().
				FindUser(ctx, creds.Login).
				Return(testCase.findReturns)

			m.EXPECT().
				AddCredentials(ctx, creds).
				Return(testCase.storeReturns).
				AnyTimes()

			r := NewRegister(&dummyLogger, time.Hour, "signature", m)

			_, err := r.Auth(ctx, creds)

			switch {
			case testCase.wantUCReturn == nil:
				assert.NoError(t, err)
			default:
				assert.ErrorIs(t, err, testCase.wantUCReturn)
			}
		})
	}
}
