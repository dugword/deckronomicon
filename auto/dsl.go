package auto

import "fmt"

// These helper functions create strategy condition trees in Go code.
// Great for testing or writing built-in behaviors without JSON.

func When(conds ...ConditionNode) ConditionNode {
	return &AndCondition{Conditions: conds}
}

func And(conds ...ConditionNode) ConditionNode {
	return &AndCondition{Conditions: conds}
}

func Or(conds ...ConditionNode) ConditionNode {
	return &OrCondition{Conditions: conds}
}

func Not(cond ConditionNode) ConditionNode {
	return &NotCondition{Condition: cond}
}

/*
func InHand(cards ...string) ConditionNode {
	return &InZoneCondition{Zone: "Hand", Cards: cards}
}

func OnBattlefield(cards ...string) ConditionNode {
	return &InZoneCondition{Zone: "Battlefield", Cards: cards}
}

func InGraveyard(cards ...string) ConditionNode {
	return &InZoneCondition{Zone: "Graveyard", Cards: cards}
}

func OnStack(cards ...string) ConditionNode {
	return &InZoneCondition{Zone: "Stack", Cards: cards}
}

/*
func LifeTotal(op string, value int) ConditionNode {
	return &PlayerStatCondition{Stat: "lifeTotal", Op: op, Value: value}
}
*/

// TODO Think about how to handle expressions vs op/values
func LifeTotal(expr string) ConditionNode {
	stat, op, value, err := parseStatShortcut("lifeTotal", expr)
	if err != nil {
		panic(fmt.Sprintf("Invalid life total expression: %s", expr))
	}
	return &PlayerStatCondition{Stat: stat, Op: op, Value: value}
}

func HandSize(op string, value int) ConditionNode {
	return &PlayerStatCondition{Stat: "handSize", Op: op, Value: value}
}

func GraveyardSize(op string, value int) ConditionNode {
	return &PlayerStatCondition{Stat: "graveyardSize", Op: op, Value: value}
}

/*
func Cast(card string) ActionNode {
	return ActionNode{
		ActionType: "Cast",
		CardName:   card,
	}
}
*/
