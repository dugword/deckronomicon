package state

import (
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"strconv"
)

func (g *Game) WithGetNextID() (id string, game *Game) {
	newGame := *g
	newGame.nextID++
	return strconv.Itoa(newGame.nextID), &newGame
}

func (g *Game) WithCheatsEnabled(enabled bool) *Game {
	newGame := *g
	newGame.cheatsEnabled = enabled
	return &newGame
}

func (g *Game) WithPhase(phase mtg.Phase) *Game {
	newGame := *g
	newGame.phase = phase
	return &newGame
}

func (g *Game) WithStep(step mtg.Step) *Game {
	newGame := *g
	newGame.step = step
	return &newGame
}

func (g *Game) WithGameOver(winnerID string) *Game {
	newGame := *g
	newGame.winnerID = winnerID
	return &newGame
}

func (g *Game) WithPlayers(players []*Player) *Game {
	newGame := *g
	newGame.players = players
	return &newGame
}

func (g *Game) WithActivePlayer(playerID string) *Game {
	newGame := *g
	var idx int
	for i, p := range newGame.players {
		if p.id == playerID {
			idx = i
			break
		}
	}
	newGame.activePlayerIdx = idx
	return &newGame
}

func (g *Game) WithResetPriorityPasses() *Game {
	newGame := *g
	newGame.playersPassedPriority = map[string]bool{}
	return &newGame
}

func (g *Game) WithPlayerPassedPriority(playerID string) *Game {
	newGame := *g
	playersPassedPriority := map[string]bool{}
	for pID := range newGame.playersPassedPriority {
		playersPassedPriority[pID] = newGame.playersPassedPriority[pID]
	}
	playersPassedPriority[playerID] = true
	newGame.playersPassedPriority = playersPassedPriority
	return &newGame
}

func (g *Game) WithUpdatedPlayer(player *Player) *Game {
	newGame := *g
	var players []*Player
	for _, p := range newGame.players {
		if p.id == player.id {
			players = append(players, player)
			continue
		}
		players = append(players, p)
	}
	newGame.players = players
	return &newGame
}

func (g *Game) WithBattlefield(battlefield *Battlefield) *Game {
	newGame := *g
	newGame.battlefield = battlefield
	return &newGame
}

// TODO: Maybe manage all cards moving zones like this so I can register and deregister triggered abilities?
func (g *Game) WithPutCardInGraveyard(playerID string, card *gob.Card) *Game {
	player := g.GetPlayer(playerID)
	player, ok := player.WithAddCardToZone(card, mtg.ZoneGraveyard)
	if !ok {
		return g
	}
	game := g.WithUpdatedPlayer(player)
	for _, triggeredAbility := range card.TriggeredAbilities() {
		if triggeredAbility.Zone() != mtg.ZoneGraveyard {
			continue
		}
		game = game.WithRegisteredTriggeredAbility(
			playerID,
			triggeredAbility.Name(),
			card.Name(),
			card.ID(),
			triggeredAbility.Trigger(),
			triggeredAbility.Effects(),
			"",
			false,
		)
	}
	return game
}

func (g *Game) WithRemoveCardFromGraveyard(playerID, cardID string) (*gob.Card, *Game, bool) {
	player := g.GetPlayer(playerID)
	card, player, ok := player.TakeCardFromZone(cardID, mtg.ZoneGraveyard)
	if !ok {
		return nil, nil, false
	}
	game := g.WithUpdatedPlayer(player)
	for _, registeredTriggeredAbilities := range game.RegisteredTriggeredAbilities() {
		if registeredTriggeredAbilities.SourceID != card.ID() {
			continue
		}
		game = game.WithRemoveTriggeredAbility(registeredTriggeredAbilities.ID)
	}
	return card, game, true
}

func (g *Game) WithRemovePermanentFromBattlefield(permanentID string) (*gob.Permanent, *Game, bool) {
	permanent, battlefield, ok := g.battlefield.Take(permanentID)
	if !ok {
		return nil, nil, false
	}
	newGame := g.WithBattlefield(battlefield)
	for _, registeredTriggeredAbilities := range newGame.RegisteredTriggeredAbilities() {
		if registeredTriggeredAbilities.SourceID != permanent.ID() {
			continue
		}
		newGame = newGame.WithRemoveTriggeredAbility(registeredTriggeredAbilities.ID)
	}
	return permanent, newGame, true
}

func (g *Game) WithPutPermanentOnBattlefield(card *gob.Card, playerID string) (*Game, error) {
	id, gameWithID := g.WithGetNextID()
	permanent, err := gob.NewPermanent(id, card, playerID)
	if err != nil {
		return nil, err
	}
	battlefield := gameWithID.battlefield.Add(permanent)
	gameWithBattlefield := gameWithID.WithBattlefield(battlefield)
	for _, triggeredAbility := range permanent.TriggeredAbilities() {
		if triggeredAbility.Zone() != mtg.ZoneBattlefield {
			continue
		}
		gameWithBattlefield = gameWithBattlefield.WithRegisteredTriggeredAbility(
			playerID,
			triggeredAbility.Name(),
			permanent.Name(),
			permanent.ID(),
			triggeredAbility.Trigger(),
			triggeredAbility.Effects(),
			"",
			false,
		)
	}
	return gameWithBattlefield, nil
}

func (g *Game) WithStack(stack *Stack) *Game {
	newGame := *g
	newGame.stack = stack
	return &newGame
}

func (g *Game) WithPutSpellOnStack(
	card *gob.Card,
	playerID string,
	effectWithTargets []*effect.EffectWithTarget,
	flashback bool,
) (*Game, error) {
	id, gameWithNextID := g.WithGetNextID()
	spell, err := gob.NewSpell(
		id,
		card,
		playerID,
		effectWithTargets,
		flashback,
	)
	if err != nil {
		return nil, err
	}
	stack := gameWithNextID.stack.AddTop(spell)
	gameWithStack := gameWithNextID.WithStack(stack)
	return gameWithStack, nil
}

func (g *Game) WithPutCopiedSpellOnStack(
	spell *gob.Spell,
	playerID string,
	effectWithTargets []*effect.EffectWithTarget,
) (*Game, error) {
	id, gameWithNextID := g.WithGetNextID()
	spell, err := gob.CopySpell(
		id,
		spell,
		playerID,
		effectWithTargets,
	)
	if err != nil {
		return nil, err
	}
	stack := gameWithNextID.stack.AddTop(spell)
	gameWithStack := gameWithNextID.WithStack(stack)
	return gameWithStack, nil
}

func (g *Game) WithPutAbilityOnStack(
	playerID,
	sourceID,
	abilityID,
	abilityName string,
	effectWithTargets []*effect.EffectWithTarget,
) (*Game, error) {
	id, gameWithNextID := g.WithGetNextID()
	abilityOnStack := gob.NewAbilityOnStack(
		id,
		playerID,
		sourceID,
		abilityID,
		abilityName,
		effectWithTargets,
	)
	stack := gameWithNextID.stack.AddTop(abilityOnStack)
	gameWithStack := gameWithNextID.WithStack(stack)
	return gameWithStack, nil
}

func (g *Game) WithRegisteredTriggeredAbility(
	playerID string,
	abilityName string,
	sourceName string,
	sourceID string,
	trigger gob.Trigger,
	effects []effect.Effect,
	duration mtg.Duration,
	oneShot bool,
) *Game {
	id, newGame := g.WithGetNextID()
	triggeredEffect := gob.RegisteredTriggeredAbility{
		ID:       id,
		Name:     abilityName,
		SourceID: sourceID,
		PlayerID: playerID,
		Trigger:  trigger,
		Effects:  effects,
		Duration: duration,
		OneShot:  oneShot,
	}
	newGame.registeredTriggeredAbilities = append(newGame.registeredTriggeredAbilities[:], triggeredEffect)
	return newGame
}

func (g *Game) WithRemoveTriggeredAbility(triggeredAbilityID string) *Game {
	newGame := *g
	var newTriggeredAbilities []gob.RegisteredTriggeredAbility
	for _, triggeredAbility := range newGame.registeredTriggeredAbilities {
		if triggeredAbility.ID != triggeredAbilityID {
			newTriggeredAbilities = append(newTriggeredAbilities, triggeredAbility)
		}
	}
	newGame.registeredTriggeredAbilities = newTriggeredAbilities
	return &newGame
}
