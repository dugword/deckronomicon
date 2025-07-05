package main

import (
	"deckronomicon/packages/engine/store"
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
	NumberOfTurns                                int       `json:"NumberOfTurns"`
	TotalRuns                                    int       `json:"TotalRuns"`
	TotalCardDrawn                               int       `json:"TotalCardsDrawn"`
	TotalFlushOutDrawn                           int       `json:"TotalFlushOutDrawn"`
	TotalThrillOfPossibilityDrawn                int       `json:"TotalThrillOfPossibilityDrawn"`
	TotalPearledUnicornDrawn                     int       `json:"TotalPearledUnicornDrawn"`
	RunsWithAtLeastOneFlushOut                   int       `json:"RunsWithAtLeastOneFlushOut"`
	RunsWithAtLeastOneThrillOfPossibility        int       `json:"RunsWithAtLeastOneThrillOfPossibility"`
	RunsWithAtLeastOnePearledUnicorn             int       `json:"RunsWithAtLeastOnePearledUnicorn"`
	AverageCardDrawn                             float64   `json:"AverageCardDrawn"`
	AverageFlushOutDrawn                         float64   `json:"AverageFlushOutDrawn"`
	AverageThrillOfPossibilityDrawn              float64   `json:"AverageThrillOfPossibilityDrawn"`
	AveragePearledUnicornDrawn                   float64   `json:"AveragePearledUnicornDrawn"`
	ChanceAtLeastOneFlushOutDrawn                float64   `json:"ChanceAtLeastOneFlushOut"`
	ChanceAtLeastOneThrillOfPossibilityDrawn     float64   `json:"ChanceAtLeastOneThrillOfPossibilityDrawn"`
	ChanceAtLeastOnePearledUnicornDrawn          float64   `json:"ChanceAtLeastOnePearledUnicornDrawn"`
	AvgCumulativeCardDrawnPerTurn                []float64 `json:"AvgCumulativeCardDrawnPerTurn"`
	AvgCumulativeFlushOutDrawnPerTurn            []float64 `json:"AvgCumulativeFlushOutDrawnPerTurn"`
	AvgCumulativeThrillOfPossibilityDrawnPerTurn []float64 `json:"AvgCumulativeThrillOfPossibilityDrawnPerTurn"`
	AvgCumulativePearledUnicornDrawnPerTurn      []float64 `json:"AvgCumulativePearledUnicornDrawnPerTurn"`
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
	file, err := os.Open("results/did_this_work.json")
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
		NumberOfTurns:                                numTurns,
		AvgCumulativeCardDrawnPerTurn:                make([]float64, numTurns),
		AvgCumulativeFlushOutDrawnPerTurn:            make([]float64, numTurns),
		AvgCumulativeThrillOfPossibilityDrawnPerTurn: make([]float64, numTurns),
		AvgCumulativePearledUnicornDrawnPerTurn:      make([]float64, numTurns),
	}
	for _, run := range results {
		summary.TotalRuns++
		summary.TotalCardDrawn += run.Totals["CardDrawn"]
		summary.TotalFlushOutDrawn += run.Totals["FlushOutDrawn"]
		summary.TotalThrillOfPossibilityDrawn += run.Totals["ThrillOfPossibilityDrawn"]
		summary.TotalPearledUnicornDrawn += run.Totals["PearledUnicornDrawn"]
		if run.Totals["FlushOutDrawn"] > 0 {
			summary.RunsWithAtLeastOneFlushOut++
		}
		if run.Totals["ThrillOfPossibilityDrawn"] > 0 {
			summary.RunsWithAtLeastOneThrillOfPossibility++
		}
		if run.Totals["PearledUnicornDrawn"] > 0 {
			summary.RunsWithAtLeastOnePearledUnicorn++
		}
		for i := range numTurns {
			summary.AvgCumulativeCardDrawnPerTurn[i] += float64(run.Cumulatives["CardDrawn"][i])
			summary.AvgCumulativeFlushOutDrawnPerTurn[i] += float64(run.Cumulatives["FlushOutDrawn"][i])
			summary.AvgCumulativeThrillOfPossibilityDrawnPerTurn[i] += float64(run.Cumulatives["ThrillOfPossibilityDrawn"][i])
			summary.AvgCumulativePearledUnicornDrawnPerTurn[i] += float64(run.Cumulatives["PearledUnicornDrawn"][i])
		}
	}
	for i := range numTurns {
		summary.AvgCumulativeCardDrawnPerTurn[i] /= float64(summary.TotalRuns)
		summary.AvgCumulativeFlushOutDrawnPerTurn[i] /= float64(summary.TotalRuns)
		summary.AvgCumulativeThrillOfPossibilityDrawnPerTurn[i] /= float64(summary.TotalRuns)
		summary.AvgCumulativePearledUnicornDrawnPerTurn[i] /= float64(summary.TotalRuns)
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
