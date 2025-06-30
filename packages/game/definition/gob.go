package definition

type Ability struct {
	Cost    string   `json:"Cost,omitempty" yaml:"Cost,omitempty"`
	Effects []Effect `json:"Effects,omitempty" yaml:"Effects,omitempty"`
	Name    string   `json:"Name,omitempty" yaml:"Name,omitempty"`
	Speed   string   `json:"Speed,omitempty" yaml:"Speed,omitempty"`
	Zone    string   `json:"Zone,omitempty" yaml:"Zone,omitempty"`
}

type AbilityOnStack struct {
	AbilityID         string             `json:"AbilityID,omitempty" yaml:"AbilityID,omitempty"`
	Constroller       string             `json:"Controller,omitempty" yaml:"Controller,omitempty"`
	EffectWithTargets []EffectWithTarget `json:"EffectWithTargets,omitempty" yaml:"EffectWithTargets,omitempty"`
	ID                string             `json:"ID,omitempty" yaml:"ID,omitempty"`
	Name              string             `json:"Name,omitempty" yaml:"Name,omitempty"`
	Owner             string             `json:"Owner,omitempty" yaml:"Owner,omitempty"`
	SourceID          string             `json:"SourceID,omitempty" yaml:"SourceID,omitempty"`
}

type Card struct {
	ActivatedAbilities []Ability       `json:"ActivatedAbilities,omitempty" yaml:"ActivatedAbilities,omitempty"`
	CardTypes          []string        `json:"CardTypes,omitempty" yaml:"CardTypes,omitempty"`
	Colors             []string        `json:"Color,omitempty" yaml:"Color,omitempty"`
	Controller         string          `json:"Controller,omitempty" yaml:"Controller,omitempty"`
	Mana               string          `json:"Mana,omitempty" yaml:"Mana,omitempty"`
	AdditionalCost     string          `json:"AdditionalCost,omitempty" yaml:"AdditionalCost,omitempty"`
	ID                 string          `json:"ID,omitempty" yaml:"ID,omitempty"`
	Loyalty            int             `json:"Loyalty,omitempty" yaml:"Loyalty,omitempty"`
	Name               string          `json:"Name,omitempty" yaml:"Name,omitempty"`
	Owner              string          `json:"Owner,omitempty" yaml:"Owner,omitempty"`
	Power              int             `json:"Power,omitempty" yaml:"Power,omitempty"`
	RulesText          string          `json:"RulesText,omitempty" yaml:"RulesText,omitempty"`
	SpellAbility       []Effect        `json:"SpellAbility,omitempty" yaml:"SpellAbility,omitempty"`
	StaticAbilities    []StaticAbility `json:"StaticAbilities,omitempty" yaml:"StaticAbilities,omitempty"`
	Subtypes           []string        `json:"Subtypes,omitempty" yaml:"Subtypes,omitempty"`
	Supertypes         []string        `json:"Supertypes,omitempty" yaml:"Supertypes,omitempty"`
	Toughness          int             `json:"Toughness,omitempty" yaml:"Toughness,omitempty"`
}

type Effect struct {
	Modifiers map[string]any `json:"Modifiers,omitempty" yaml:"Modifiers,omitempty"`
	Name      string         `json:"Name,omitempty" yaml:"Name,omitempty"`
}

type EffectWithTarget struct {
	Effect Effect `json:"Effect" yaml:"Effect,omitempty"`
	Target Target `json:"Target" yaml:"Target,omitempty"`
}

type Permanent struct {
	ActivatedAbilities []Ability       `json:"ActivatedAbilities,omitempty" yaml:"ActivatedAbilities,omitempty"`
	Card               Card            `json:"Card" yaml:"Card,omitempty"`
	CardTypes          []string        `json:"CardTypes,omitempty" yaml:"CardTypes,omitempty"`
	Colors             []string        `json:"Colors,omitempty" yaml:"Colors,omitempty"`
	Controller         string          `json:"Controller,omitempty" yaml:"Controller,omitempty"`
	ID                 string          `json:"ID,omitempty" yaml:"ID,omitempty"`
	Loyalty            int             `json:"Loyalty,omitempty" yaml:"Loyalty,omitempty"`
	Cost               string          `json:"Cost,omitempty" yaml:"Cost,omitempty"`
	Name               string          `json:"Name,omitempty" yaml:"Name,omitempty"`
	Owner              string          `json:"Owner,omitempty" yaml:"Owner,omitempty"`
	Power              int             `json:"Power,omitempty" yaml:"Power,omitempty"`
	RulesText          string          `json:"RulesText,omitempty" yaml:"RulesText,omitempty"`
	StaticAbilities    []StaticAbility `json:"StaticAbilities,omitempty" yaml:"StaticAbilities,omitempty"`
	Subtypes           []string        `json:"Subtypes,omitempty" yaml:"Subtypes,omitempty"`
	SummoningSickness  bool            `json:"SummoningSickness,omitempty" yaml:"SummoningSickness,omitempty"`
	Supertypes         []string        `json:"Supertypes,omitempty" yaml:"Supertypes,omitempty"`
	Tapped             bool            `json:"Tapped,omitempty" yaml:"Tapped,omitempty"`
	Toughness          int             `json:"Toughness,omitempty" yaml:"Toughness,omitempty"`
	TriggeredAbilities []Ability       `json:"TriggeredAbilities,omitempty" yaml:"TriggeredAbilities,omitempty"`
}

type Spell struct {
	Card              Card               `json:"Card,omitempty" yaml:"Card,omitempty"`
	CardTypes         []string           `json:"CardTypes,omitempty" yaml:"CardTypes,omitempty"`
	Colors            []string           `json:"Colors,omitempty" yaml:"Colors,omitempty"`
	Controller        string             `json:"Controller,omitempty" yaml:"Controller,omitempty"`
	EffectWithTargets []EffectWithTarget `json:"EffectWithTargets,omitempty" yaml:"EffectWithTargets,omitempty"`
	Flashback         bool               `json:"Flashback,omitempty" yaml:"Flashback,omitempty"`
	ID                string             `json:"ID,omitempty" yaml:"ID,omitempty"`
	Copy              bool               `json:"Copy,omitempty" yaml:"Copy,omitempty"`
	Loyalty           int                `json:"Loyalty,omitempty" yaml:"Loyalty,omitempty"`
	Cost              string             `json:"Cost,omitempty" yaml:"Cost,omitempty"`
	Name              string             `json:"Name,omitempty" yaml:"Name,omitempty"`
	Owner             string             `json:"Owner,omitempty" yaml:"Owner,omitempty"`
	Power             int                `json:"Power,omitempty" yaml:"Power,omitempty"`
	RulesText         string             `json:"RulesText,omitempty" yaml:"RulesText,omitempty"`
	StaticAbilities   []StaticAbility    `json:"StaticAbilities,omitempty" yaml:"StaticAbilities,omitempty"`
	Subtypes          []string           `json:"Subtypes,omitempty" yaml:"Subtypes,omitempty"`
	Supertypes        []string           `json:"Supertypes,omitempty" yaml:"Supertypes,omitempty"`
	Toughness         int                `json:"Toughness,omitempty" yaml:"Toughness,omitempty"`
}

type StaticAbility struct {
	Cost      string         `json:"Cost,omitempty" yaml:"Cost,omitempty"`
	Modifiers map[string]any `json:"Modifiers,omitempty" yaml:"Modifiers,omitempty"`
	Name      string         `json:"Name,omitempty" yaml:"Name,omitempty"`
}

type Target struct {
	Type string `json:"Type,omitempty" yaml:"Type,omitempty"`
	ID   string `json:"ID,omitempty" yaml:"ID,omitempty"`
}
