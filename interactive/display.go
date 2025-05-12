package interactive

import (
	"strings"

	"deckronomicon/game"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func ClearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// TODO Maybe text/templates would be good
func DisplayGameState(g *game.GameState) {
	DisplayGameStatus(g)
	fmt.Println()
	DisplayBattlefield(g.Battlefield)
	DisplayGraveyard(g.Graveyard)
	fmt.Println()
	DisplayManaPool(g)
	DisplayHand(g.Hand)
	fmt.Println()
	DisplayTurnMessageLog(g.TurnMessageLog)
	fmt.Println()
}

func DisplayGameStatus(g *game.GameState) {
	fmt.Println("=== GAME STATE ===")
	fmt.Printf("Current Player: Player %d\n", g.CurrentPlayer)
	fmt.Printf("Turn: %d\n", g.Turn)
	fmt.Printf("Phase: %s\n", g.CurrentPhase)
	fmt.Printf("Player %d: %d life\n", 1, g.Life)
	fmt.Printf("Hand: %d cards\n", len(g.Hand))
	fmt.Printf("Library: %d cards\n", len(g.Deck))
	fmt.Printf("Graveyard: %d cards\n", len(g.Graveyard))
	fmt.Printf("Message: %s\n", g.Message)
}

func DisplayTurnMessageLog(messageLog []string) {
	fmt.Println("Message log:")
	for _, message := range messageLog {
		fmt.Println(message)
	}
}
func DisplayHand(hand []*game.Card) {
	names := []string{}
	for _, card := range hand {
		names = append(names, card.Name)
	}
	fmt.Printf("Hand: [%s]\n", strings.Join(names, ", "))
}

func DisplayManaPool(g *game.GameState) {
	potential := g.PotentialMana
	potentialOut := "Potential Mana: "
	for color, amt := range potential {
		potentialOut += fmt.Sprintf("{%s}=%d ", color, amt)
	}
	fmt.Println(potentialOut)
	pool := g.ManaPool
	if len(pool) == 0 {
		fmt.Println("Mana Pool: (empty)")
		return
	}
	poolOut := "Mana Pool: "
	for color, amt := range pool {
		poolOut += fmt.Sprintf("{%s}=%d ", color, amt)
	}
	fmt.Println(poolOut)
}

func DisplayGraveyard(graveyard []*game.Card) {
	cutoff := 3
	var summary []string
	if len(graveyard) <= cutoff {
		for _, card := range graveyard {
			summary = append(summary, card.Name)
		}
	} else {
		summary = append(summary, fmt.Sprintf("...%d more", len(graveyard)-cutoff))
		for _, card := range graveyard[len(graveyard)-cutoff:] {
			summary = append(summary, card.Name)
		}
	}
	fmt.Printf("Graveyard:")
	fmt.Println(strings.Join(summary, ", "))
}

func DisplayBattlefield(perms []*game.Permanent) {
	fmt.Println("Battlefield:")
	for i, p := range perms {
		status := []string{}
		if p.Tapped {
			status = append(status, "🔒 Tapped")
		}
		if p.SummoningSick {
			status = append(status, "💤 Summoning Sick")
		}
		icon := "✨"
		if p.Object.HasType(game.CardTypeCreature) {
			icon = "🧍"
		} else if p.Object.HasType(game.CardTypeArtifact) {
			icon = "⚙️"
		} else if p.Object.HasType(game.CardTypeLand) {
			icon = "🌿"
		}
		desc := icon + " " + p.Object.Name
		if p.Power != 0 || p.Toughness != 0 {
			desc += fmt.Sprintf(" [%d/%d]", p.Power, p.Toughness)
		}
		if len(status) > 0 {
			desc += " (" + strings.Join(status, ", ") + ")"
		}
		fmt.Printf(" [%d] %s\n", i, desc)
	}
}
