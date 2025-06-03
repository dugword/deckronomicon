package agent

import "fmt"

type Agent struct {
	id string
}

func NewAgent(id string) *Agent {
	agent := Agent{id: id}
	return &agent
}

func (a *Agent) GetNextAction() string {
	// press enter to continue
	fmt.Println("Press Enter to continue for agent:", a.id)
	fmt.Scanln() // wait for Enter Key
	return "pass"
}
