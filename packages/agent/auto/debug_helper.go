package auto

import (
	"bufio"
	"deckronomicon/packages/game/mtg"
	"fmt"
	"os"
	"slices"
)

func (a *RuleBasedAgent) enterToContinue() {
	if !a.interactive {
		return
	}
	if err := a.uiBuffer.Render(); err != nil {
		panic(fmt.Errorf("failed to render UI buffer: %w", err))
	}
	fmt.Print("Press Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func (a *RuleBasedAgent) enterToContinueOnSteps(step mtg.Step) {
	if !a.interactive {
		return
	}
	if slices.Contains(a.stops, step) {
		if err := a.uiBuffer.Render(); err != nil {
			panic(fmt.Errorf("failed to render UI buffer: %w", err))
		}
		fmt.Print("Press Enter to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}
