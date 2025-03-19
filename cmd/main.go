package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	client "testTask/pkg/api"
	"time"
)

func main() {
	// Create new API client
	apiClient := client.NewClient(
		"https://development.kpi-drive.ru/_api/facts/save_fact",
		"https://development.kpi-drive.ru/_api/indicators/get_facts",
		"48ab34464a5573519725deb5865cc74c",
	)

	// Create context with cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get facts for period
	log.Println("Getting facts for period...")
	_, err := apiClient.GetFacts(ctx, "2024-12-01", "2024-12-31", "month", 227373)
	if err != nil {
		log.Fatalf("Error getting facts: %v", err)
	}

	// Wait for signal to stop
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Waiting for signal to stop...")
	<-sigs

	log.Println("Stopping...")
	// Cancel context
	cancel()
	// Wait for 2 seconds to finish all operations
	time.Sleep(2 * time.Second)
}
