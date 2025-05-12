package interactive_test

/*
func TestREPLAgent_DrawAndPass(t *testing.T) {
	input := strings.NewReader("draw\npass\n")
	scanner := bufio.NewScanner(input)
	agent := interactive.NewREPLAgent(scanner)

	state := game.NewGameState()

	action1 := agent.GetNextAction(state)
	if action1.Cheat != game.CheatDraw {
		t.Errorf("Expected ActionDraw, got %v", action1.Type)
	}

	action2 := agent.GetNextAction(state)
	if action2.Type != game.ActionPass {
		t.Errorf("Expected ActionPass, got %v", action2.Type)
	}
}
*/
