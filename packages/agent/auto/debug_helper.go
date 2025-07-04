package auto

import (
	"bufio"
	"deckronomicon/packages/choose"
	"fmt"
	"os"
)

func (a *RuleBasedAgent) EnterToContinue(message string) {
	a.uiBuffer.UpdateMessage([]string{message})
	if err := a.uiBuffer.Render(); err != nil {
		panic(fmt.Errorf("failed to render UI buffer: %w", err))
	}
	fmt.Print("Press Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func (a *RuleBasedAgent) EnterToContinueOnChoices(message string, title string, choices []choose.Choice) {
	a.uiBuffer.UpdateMessage([]string{message})
	a.uiBuffer.UpdateChoices(title, choices)
}
