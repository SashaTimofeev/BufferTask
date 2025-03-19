package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"testTask/pkg/api"
	"testTask/pkg/buffer"
	"testTask/pkg/model"
)

func main() {
	// Create new API client
	apiClient := api.NewClient(
		"https://development.kpi-drive.ru/_api/facts/save_fact",
		"https://development.kpi-drive.ru/_api/indicators/get_facts",
		"48ab34464a5573519725deb5865cc74c",
	)

	// Create new buffer
	factBuffer := buffer.NewBuffer(apiClient)

	// Create context with cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start buffer goroutine
	go factBuffer.Run(ctx)

	// Add facts to buffer
	for i := 1; i < 11; i++ {
		fact := model.Fact{
			PeriodStart:         "2024-12-01",
			PeriodEnd:           "2024-12-31",
			PeriodKey:           "month",
			IndicatorToMoID:     227373,
			IndicatorToMoFactID: 0,
			Value:               3,
			FactTime:            "2024-12-31",
			IsPlan:              0,
			AuthUserID:          40,
			Comment:             "buffer Last_name " + strconv.Itoa(i),
		}
		factBuffer.AddFact(fact)
		log.Printf("Fact added to buffer: %s", fact.Comment)
	}

	// Wait for signal to stop
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Press Ctrl+C to stop...")
	<-sigs

	log.Println("Stopping...")
	// Cancel context
	cancel()
	// Wait for 2 seconds to finish all operations
	time.Sleep(2 * time.Second)
}
