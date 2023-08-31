package usecase_test

import (
	"L0_EVRONE/internal/aggregate"
	"L0_EVRONE/internal/usecase"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type test struct {
	name  string
	mock  func()
	input string
	res   interface{}
	err   error
}

func order(t *testing.T) (*usecase.OrderUseCase, *MockOrderRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	repo := NewMockOrderRepository(mockCtl)

	order := usecase.New(repo)

	return order, repo
}

func TestSet(t *testing.T) {
	t.Parallel()
	order, repo := order(t)

	tests := []test{
		{
			name: "empty result",
			mock: func() {
				repo.EXPECT().PutRep(aggregate.Order{ Order_uid: "123456789"}).Return(nil)
			},
			input: "",
			res:   nil,
			err:   nil,
		},
	}
	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.mock()

			err := order.Set(context.Background(), aggregate.Order{Order_uid: "123456789"})

			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()
	order, repo := order(t)

	tests := []test{
		{
			name: "empty result",
			mock: func() {
				repo.EXPECT().GetRep("123456789").Return(aggregate.Order{
					Order_uid: "123456789",
				}, nil)
			},
			input: "123456789",
			res: aggregate.Order{
				Order_uid: "123456789",
			},
			err: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.mock()

			res, err := order.Get(context.Background(), tc.input)

			require.Equal(t, res, tc.res)
			require.ErrorIs(t, err, tc.err)
		})
	}
}
