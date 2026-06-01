package auction

import (
	"context"
	"os"
	"testing"
	"time"

	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
)

func TestStartAuctionCloserClosesAuctionAfterDuration(t *testing.T) {
	originalTimer := auctionTimer
	defer func() {
		auctionTimer = originalTimer
	}()

	if err := os.Setenv("AUCTION_DURATION", "1ms"); err != nil {
		t.Fatalf("falha ao definir AUCTION_DURATION: %v", err)
	}

	closed := make(chan struct{}, 1)
	auctionTimer = func(d time.Duration) *time.Timer {
		return time.NewTimer(10 * time.Millisecond)
	}

	mockUpdater := &MockAuctionUpdater{
		closed: closed,
		t:      t,
	}

	repo := &AuctionRepository{
		Updater: mockUpdater,
	}
	repo.StartAuctionCloser(context.Background(), "test-auction-id")

	select {
	case <-closed:
		// sucesso
	case <-time.After(200 * time.Millisecond):
		t.Fatal("o leilão não foi fechado dentro do tempo esperado")
	}
}

type MockAuctionUpdater struct {
	closed chan struct{}
	t      *testing.T
}

func (m *MockAuctionUpdater) UpdateAuctionStatus(ctx context.Context, auctionId string, status auction_entity.AuctionStatus) *internal_error.InternalError {
	if auctionId != "test-auction-id" {
		m.t.Errorf("esperado auctionId test-auction-id, recebeu %s", auctionId)
	}
	if status != auction_entity.Completed {
		m.t.Errorf("esperado status Completed, recebeu %v", status)
	}
	m.closed <- struct{}{}
	return nil
}
