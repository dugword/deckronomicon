package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const ScryfallURLPattern = "https://api.scryfall.com/cards/named?exact=%s"

type ScryfallCard struct {
	ColorIdentity []string `json:"color_identity"`
	Colors        []string `json:"colors"`
	ManaCost      string   `json:"mana_cost"`
	Name          string   `json:"name"`
	OracleText    string   `json:"oracle_text"`
	Power         string   `json:"power"`
	Toughness     string   `json:"toughness"`
	TypeLine      string   `json:"type_line"`
}

type CardImport struct {
	ActivatedAbilities []string `json:"ActivatedAbilities,omitempty"`
	CardTypes          []string `json:"CardTypes,omitempty"`
	Colors             []string `json:"Color,omitempty"`
	Loyalty            int      `json:"Loyalty,omitempty"`
	ManaCost           string   `json:"ManaCost,omitempty"`
	Name               string   `json:"Name,omitempty"`
	Power              int      `json:"Power,omitempty"`
	RulesText          string   `json:"RulesText,omitempty"`
	SpellAbility       string   `json:"SpellAbility,omitempty"`
	Subtypes           []string `json:"Subtypes,omitempty"`
	Supertypes         []string `json:"Supertypes,omitempty"`
	Toughness          int      `json:"Toughness,omitempty"`
}

// main is the entry point for the application.
func main() {
	if err := Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// Run is an abtraction for the main function to enable testing.
func Run(args []string) error {
	flags := flag.NewFlagSet("deckronomicon", flag.ContinueOnError)
	cardName := flags.String("card", "", "Name of the card to fetch")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	if *cardName == "" {
		return fmt.Errorf("card name is required")
	}
	path, err := fetchAndWriteCard(*cardName)
	if err != nil {
		return fmt.Errorf("failed to fetch card %q: %w", *cardName, err)
	}
	fmt.Println("imported card to:", path)
	return nil
}

func fetchAndWriteCard(name string) (string, error) {
	url := fmt.Sprintf(ScryfallURLPattern, url.PathEscape(name))
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch card %q: %w", name, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("response error: %d", resp.StatusCode)
	}
	var s ScryfallCard
	if err := json.Unmarshal(body, &s); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}
	color := "colorless"
	// TODO: use a different field here, we are putting lands in the wrong
	// folder.
	if len(s.ColorIdentity) == 1 {
		color = colorFromIdentity(s.ColorIdentity[0])
	} else {
		color = "multicolor"
	}
	cardTypes, subtypes, supertypes := parseTypes(s.TypeLine)
	dc := CardImport{
		CardTypes:  cardTypes,
		Colors:     s.Colors,
		ManaCost:   s.ManaCost,
		Name:       s.Name,
		RulesText:  s.OracleText,
		Subtypes:   subtypes,
		Supertypes: supertypes,
	}
	if s.Power != "" {
		power, err := strconv.Atoi(s.Power)
		if err != nil {
			return "", fmt.Errorf("failed to parse power %q: %w", s.Power, err)
		}
		dc.Power = power
	}
	if s.Toughness != "" {
		toughness, err := strconv.Atoi(s.Toughness)
		if err != nil {
			return "", fmt.Errorf("failed to parse toughness %q: %w", s.Toughness, err)
		}
		dc.Toughness = toughness
	}
	path, err := writeCardJSON(color, dc)
	if err != nil {
		return "", fmt.Errorf("failed to write card JSON: %w", err)
	}
	return path, nil
}

func writeCardJSON(color string, card CardImport) (string, error) {
	nameSlug := strings.ToLower(strings.ReplaceAll(card.Name, " ", "_"))
	dir := filepath.Join("cards", color)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory %q: %w", dir, err)
	}
	path := filepath.Join(dir, nameSlug+".json")
	if _, err := os.Stat(path); err == nil {
		return "", fmt.Errorf("file already exists: %s", path)
	}
	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create file %q: %w", path, err)
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(card); err != nil {
		return "", fmt.Errorf("failed to encode JSON: %w", err)
	}
	return path, nil
}

func colorFromIdentity(c string) string {
	switch c {
	case "W":
		return "white"
	case "U":
		return "blue"
	case "B":
		return "black"
	case "R":
		return "red"
	case "G":
		return "green"
	default:
		return "colorless"
	}
}

func parseTypes(typeLine string) (cardTypes, subtypes, supertypes []string) {
	// TODO: use the definitions in mtg.
	knownSupertypes := map[string]bool{
		"Basic":     true,
		"Legendary": true,
		"Ongoing":   true,
		"Snow":      true,
		"World":     true,
	}
	knownTypes := map[string]bool{
		"Artifact":     true,
		"Battle":       true,
		"Creature":     true,
		"Enchantment":  true,
		"Instant":      true,
		"Land":         true,
		"Planeswalker": true,
		"Sorcery":      true,
		"Tribal":       true,
	}
	parts := strings.Split(typeLine, "—")
	left := strings.Fields(strings.TrimSpace(parts[0]))
	for _, word := range left {
		if knownSupertypes[word] {
			supertypes = append(supertypes, word)
		} else if knownTypes[word] {
			cardTypes = append(cardTypes, word)
		}
	}
	if len(parts) > 1 {
		s := strings.Fields(strings.TrimSpace(parts[1]))
		subtypes = append(subtypes, s...)
	}
	return cardTypes, subtypes, supertypes
}
