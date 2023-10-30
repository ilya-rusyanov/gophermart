package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
	"github.com/ilya-rusyanov/gophermart/internal/usecases/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLogin(t *testing.T) {
	var dummyLogger DummyLogger

	ctx := context.Background()

	someError := errors.New("some error")

	tests := []struct {
		name           string
		storageReturns error
		wantUCReturn   error
	}{
		{
			name:           "correct login",
			storageReturns: nil,
			wantUCReturn:   nil,
		},
		{
			name:           "not found",
			storageReturns: entities.ErrNotFound,
			wantUCReturn:   entities.ErrLoginIncorrect,
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

			m := mocks.NewMockLoginStorage(ctrl)
			m.EXPECT().FindCredentials(ctx, creds).Return(testCase.storageReturns)

			l := NewLogin(&dummyLogger, time.Hour, "signature", m)

			_, err := l.Auth(ctx, creds)

			if testCase.wantUCReturn != nil {
				assert.ErrorIs(t, err, testCase.wantUCReturn)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
