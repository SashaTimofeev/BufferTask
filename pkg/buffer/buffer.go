package buffer

import (
	"context"
	"log"
	"testTask/pkg/api"
	"testTask/pkg/model"
)

// Buffer for saving facts
type Buffer struct {
	apiClient *api.Client
	factsChan chan model.Fact
}

// Create new buffer with capacity 1000
func NewBuffer(apiClient *api.Client) *Buffer {
	return &Buffer{
		apiClient: apiClient,
		factsChan: make(chan model.Fact, 1000),
	}
}

// Add fact to buffer
func (b *Buffer) AddFact(fact model.Fact) {
	b.factsChan <- fact
}

// Run buffer
func (b *Buffer) Run(ctx context.Context) {
	// Run buffer in separate goroutine
	log.Println("Starting to save facts from buffer...")
	for {
		select {
		case <-ctx.Done():
			// Stop buffer if context was canceled
			log.Println("Stopping buffer...")
			return
		case fact := <-b.factsChan:
			// Save fact to API
			if err := b.apiClient.SaveFact(ctx, fact); err != nil {
				log.Printf("Error saving fact: %v", err)
			}
		}
	}
}
