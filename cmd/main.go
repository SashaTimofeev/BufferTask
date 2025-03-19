package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	client "testTask/pkg/api"
	"testTask/pkg/model"
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

	// Save fact
	log.Println("Saving fact...")
	err := apiClient.SaveFact(ctx, model.Fact{
		PeriodStart:         "2024-12-01",
		PeriodEnd:           "2024-12-31",
		PeriodKey:           "month",
		IndicatorToMoID:     227373,
		IndicatorToMoFactID: 0,
		Value:               3,
		FactTime:            "2024-12-31",
		IsPlan:              0,
		AuthUserID:          40,
		Comment:             "buffer Last_name ",
	})
	if err != nil {
		log.Fatalf("Error saving fact: %v", err)
	}

	// Get facts for period
	log.Println("Getting facts for period...")
	facts, err := apiClient.GetFacts(ctx, "2024-12-01", "2024-12-31", "month", 227373)
	if err != nil {
		log.Fatalf("Error getting facts: %v", err)
	}

	// Create or open the file facts.txt to write the facts
	file, err := os.Create("facts.txt")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	// Write the facts to the file
	_, err = file.WriteString(facts)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	// Parse the facts into a slice of Fact objects
	var factsSlice []model.Fact
	err = json.Unmarshal([]byte(facts), &factsSlice)
	if err != nil {
		log.Fatalf("Error unmarshalling facts: %v", err)
	}

	// Check if the fact exists in the slice
	found := false
	for _, fact := range factsSlice {
		if fact.Value == 3 && fact.PeriodStart == "2024-12-01" {
			found = true
			break
		}
	}

	if found {
		log.Println("Fact successfully saved and found in GetFacts response!")
	} else {
		log.Println("Fact not found! Possible issue with saving.")
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
