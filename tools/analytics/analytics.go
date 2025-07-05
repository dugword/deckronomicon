package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type RunResult struct {
	CardDrawn                int `json:"CardDrawn"`
	FlushOutDrawn            int `json:"FlushOutDrawn"`
	ThrillOfPossibilityDrawn int `json:"ThrillOfPossibilityDrawn"`
	PearledUnicornDrawn      int `json:"PearledUnicornDrawn"`
	RunID                    int `json:"RunID"`
}

type Summary struct {
	TotalRuns                                int     `json:"TotalRuns"`
	TotalCardDrawn                           int     `json:"TotalCardsDrawn"`
	TotalFlushOutDrawn                       int     `json:"TotalFlushOutDrawn"`
	TotalThrillOfPossibilityDrawn            int     `json:"TotalThrillOfPossibilityDrawn"`
	TotalPearledUnicornDrawn                 int     `json:"TotalPearledUnicornDrawn"`
	RunsWithAtLeastOneFlushOut               int     `json:"RunsWithAtLeastOneFlushOut"`
	RunsWithAtLeastOneThrillOfPossibility    int     `json:"RunsWithAtLeastOneThrillOfPossibility"`
	RunsWithAtLeastOnePearledUnicorn         int     `json:"RunsWithAtLeastOnePearledUnicorn"`
	AverageCardDrawn                         float64 `json:"AverageCardDrawn"`
	AverageFlushOutDrawn                     float64 `json:"AverageFlushOutDrawn"`
	AverageThrillOfPossibilityDrawn          float64 `json:"AverageThrillOfPossibilityDrawn"`
	AveragePearledUnicornDrawn               float64 `json:"AveragePearledUnicornDrawn"`
	ChanceAtLeastOneFlushOutDrawn            float64 `json:"ChanceAtLeastOneFlushOut"`
	ChanceAtLeastOneThrillOfPossibilityDrawn float64 `json:"ChanceAtLeastOneThrillOfPossibilityDrawn"`
	ChanceAtLeastOnePearledUnicornDrawn      float64 `json:"ChanceAtLeastOnePearledUnicornDrawn"`
}

func main() {
	if err := Run(); err != nil {
		fmt.Printf("Error running the analytics: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Analytics analysis completed successfully!")
}

func Run() error {
	var results []RunResult
	var summary Summary
	file, err := os.Open("results/did_this_work.json")
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&results); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}
	for _, run := range results {
		summary.TotalRuns++
		summary.TotalCardDrawn += run.CardDrawn
		summary.TotalFlushOutDrawn += run.FlushOutDrawn
		summary.TotalThrillOfPossibilityDrawn += run.ThrillOfPossibilityDrawn
		summary.TotalPearledUnicornDrawn += run.PearledUnicornDrawn
		if run.FlushOutDrawn > 0 {
			summary.RunsWithAtLeastOneFlushOut++
		}
		if run.ThrillOfPossibilityDrawn > 0 {
			summary.RunsWithAtLeastOneThrillOfPossibility++
		}
		if run.PearledUnicornDrawn > 0 {
			summary.RunsWithAtLeastOnePearledUnicorn++
		}
	}
	summary.AverageCardDrawn = float64(summary.TotalCardDrawn) / float64(summary.TotalRuns)
	summary.AverageFlushOutDrawn = float64(summary.TotalFlushOutDrawn) / float64(summary.TotalRuns)
	summary.AverageThrillOfPossibilityDrawn = float64(summary.TotalThrillOfPossibilityDrawn) / float64(summary.TotalRuns)
	summary.AveragePearledUnicornDrawn = float64(summary.TotalPearledUnicornDrawn) / float64(summary.TotalRuns)

	summary.ChanceAtLeastOneFlushOutDrawn = float64(summary.RunsWithAtLeastOneFlushOut) / float64(summary.TotalRuns) * 100
	summary.ChanceAtLeastOneThrillOfPossibilityDrawn = float64(summary.RunsWithAtLeastOneThrillOfPossibility) / float64(summary.TotalRuns) * 100
	summary.ChanceAtLeastOnePearledUnicornDrawn = float64(summary.RunsWithAtLeastOnePearledUnicorn) / float64(summary.TotalRuns) * 100

	out, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal summary: %w", err)
	}
	if err := os.WriteFile("results/summary.json", out, 0644); err != nil {
		return fmt.Errorf("failed to write summary to file: %w", err)
	}
	return nil
}
