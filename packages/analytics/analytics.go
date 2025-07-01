package analytics

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
	"deckronomicon/packages/view"
	"encoding/json"
	"fmt"
	"os"
)

func WriteGameRecordToFile(record *engine.GameRecord, dirname string) error {
	filename := dirname + "/game_record.json"
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create game record file: %w", err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(record.Export()); err != nil {
		return fmt.Errorf("failed to write game record: %w", err)
	}
	return nil
}

func WriteAnalyticsEventsToFile(record *engine.GameRecord, dirname string) error {
	filename := dirname + "/analytics.json"
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create game record file: %w", err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(record.ExportAnalytics()); err != nil {
		return fmt.Errorf("failed to write game record: %w", err)
	}
	return nil
}

func WriteGameStateToFile(game state.Game, dirname string) error {
	gameView := view.NewGameViewFromState(game)
	filename := dirname + "/game_view.json"
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create game view file: %w", err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(gameView); err != nil {
		return fmt.Errorf("failed to write game view: %w", err)
	}
	for _, player := range game.Players() {
		filename := dirname + fmt.Sprintf("/player_%s_view.json", player.ID())
		playerFile, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create player view file for %s: %w", player.ID(), err)
		}
		defer playerFile.Close()
		playerView := view.NewPlayerViewFromState(game, player, "")
		encoder := json.NewEncoder(playerFile)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(playerView); err != nil {
			return fmt.Errorf("failed to write player view for %s: %w", player.ID(), err)
		}
	}
	return nil
}
