package scheduler

import (
	"context"
	"log"
	"time"
)

// OrderScheduler menjalankan background job terjadwal (contoh: auto cancel order pending)
type OrderScheduler struct {
	ticker *time.Ticker
	quit   chan struct{}
}

func NewOrderScheduler() *OrderScheduler {
	return &OrderScheduler{
		quit: make(chan struct{}),
	}
}

func (s *OrderScheduler) Start(ctx context.Context, interval time.Duration) {
	s.ticker = time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.cleanupExpiredOrders()
			case <-s.quit:
				s.ticker.Stop()
				return
			case <-ctx.Done():
				s.ticker.Stop()
				return
			}
		}
	}()
}

func (s *OrderScheduler) Stop() {
	close(s.quit)
}

func (s *OrderScheduler) cleanupExpiredOrders() {
	// Background task logic (contoh scheduler)
	log.Println("[Scheduler] Menjalankan pengecekan background job order berkala...")
}
