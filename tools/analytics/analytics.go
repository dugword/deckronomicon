package main

import (
	"deckronomicon/packages/engine/store"
	"encoding/json"
	"fmt"
	"os"
)

type RunResult struct {
	CardDrawn int `json:"CardDrawn"`
	RunID     int `json:"RunID"`
}

type Summary struct {
	NumberOfTurns                 int       `json:"NumberOfTurns"`
	TotalRuns                     int       `json:"TotalRuns"`
	TotalCardDrawn                int       `json:"TotalCardsDrawn"`
	AverageCardDrawn              float64   `json:"AverageCardDrawn"`
	AvgCumulativeCardDrawnPerTurn []float64 `json:"AvgCumulativeCardDrawnPerTurn"`
}

func main() {
	if err := Run(); err != nil {
		fmt.Printf("Error running the analytics: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Analytics analysis completed successfully!")
}

func Run() error {
	results := []store.RunResult{}
	file, err := os.Open("results/results.json")
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&results); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}
	numTurns := 0
	for _, run := range results {
		if numTurns == 0 {
			numTurns = run.Turns
		} else if numTurns != run.Turns {
			return fmt.Errorf("inconsistent number of turns across runs: expected %d, got %d", numTurns, run.Turns)
		}
	}
	summary := Summary{
		NumberOfTurns:                 numTurns,
		AvgCumulativeCardDrawnPerTurn: make([]float64, numTurns),
	}
	for _, run := range results {
		summary.TotalRuns++
		summary.TotalCardDrawn += run.Totals["CardDrawn"]
		for i := range numTurns {
			summary.AvgCumulativeCardDrawnPerTurn[i] += float64(run.Cumulatives["CardDrawn"][i])
		}
	}
	for i := range numTurns {
		summary.AvgCumulativeCardDrawnPerTurn[i] /= float64(summary.TotalRuns)
	}
	summary.AverageCardDrawn = float64(summary.TotalCardDrawn) / float64(summary.TotalRuns)
	out, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal summary: %w", err)
	}
	if err := os.WriteFile("results/summary.json", out, 0644); err != nil {
		return fmt.Errorf("failed to write summary to file: %w", err)
	}
	return nil
}
