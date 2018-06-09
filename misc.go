package m6ik

import (
	"errors"
	"strconv"
	"strings"
)

var (
	Archetypes = []string{
		"Mighty", "Skilled", "Intellectual", "Gifted",
	}
)

type StaticDefenses struct {
	Dodge int `json:"dodge"`
	Block int `json:"block"`
	Parry int `json:"parry"`
	Soak  int `json:"soak"`
	Sense int `json:"sense"`
}

type Experience struct {
	CharPoints int
	FatePoints int
}

func (c *Character) generateRace(race string) {
	// Sample.
	if race == "" {
		races := CharDB.Races.Col("Race").Records()
		c.Race = randomChoice(races)
	} else {
		c.Race = race
	}
	CharDB.filter("Races", "Race", "==", c.Race)
	// Apply restrictions
	raceCon := CharDB.Races.Col("Proscriptions").Records()[0]
	c.applyConstraints(raceCon)
	// Apply bonuses.
	raceAttr := strings.Split(CharDB.Races.Col("Attributes").Records()[0], ", ")
	for _, r := range raceAttr {
		a, d := parseBonus(r)
		d.codeMax += d.code
		c.promoteAttribute(a, d)
	}
}

func (c *Character) generateArchetype(archetype string) error {
	// Sample.
	if archetype == "" {
		archs := CharDB.Archetypes.Col("Archetype").Records()
		c.Archetype = randomChoice(archs)
	} else {
		c.Archetype = archetype
	}
	CharDB.Archetypes = filterDf(CharDB.Archetypes, "Archetype", "==", c.Archetype)
	if len(CharDB.Archetypes.Col("Archetype").Records()) == 0 {
		return errors.New("Bad race/archetype combination. Randomizing.")
	}
	// Weight Attribute
	var attr string
	switch c.Archetype {
	case "Mighty":
		attr = "Strength"
	case "Gifted":
		attr = "Arcane"
	case "Intellectual":
		attr = randomChoice([]string{"Intellect", "Technical"})
	case "Skilled":
		attr = randomChoice([]string{"Agility", "Perception"})
	}
	c.AttrWeights[attr] += weightFactor
	// Handle Gifted
	if c.Archetype != "Gifted" {
		for _, c := range Casters {
			CharDB.Careers = filterDf(CharDB.Careers, "Type", "!=", c)
		}
		c.AttrWeights["Arcane"] = 0.0
	}
	// Apply restrictions
	archCon := CharDB.Archetypes.Col("Proscriptions").Records()[0]
	c.applyConstraints(archCon)
	// Apply bonuses
	b := strings.Split(CharDB.Archetypes.Col("Bonus").Records()[0], ", ")[0]
	if strings.Contains(b, " or ") {
		b = randomChoice(strings.Split(b, " or "))
	}
	a, d := parseBonus(b)
	d.codeMax += d.code
	c.promoteAttribute(a, d)
	return nil
}

func (c *Character) getDefVal(s, a string) int {
	p := 0
	if _, ok := c.Skills[s]; ok {
		p = c.Skills[s].toPips()
	} else {
		p = c.Attributes[a].toPips()
	}
	return p
}

func (c *Character) generateStaticDefenses() {
	c.Dodge = c.Skills["Dodge"].toPips()
	c.Block = c.getDefVal("Unarmed Combat", "Agility")
	c.Parry = 0
	for _, v := range []string{
		"Hand Weapon", "Great Weapon", "Shield Weapon",
	} {
		if _, ok := c.Skills[v]; ok {
			p := c.getDefVal(v, "Agility")
			if p > c.Parry {
				c.Parry = p
			}
		}
	}
	c.Soak = c.getDefVal("Lifting", "Strength")
	c.Sense = c.getDefVal("Search", "Perception")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (c *Character) generatePerks(n_perks string) error {
	n, _ := strconv.Atoi(n_perks)
	if n == 0 {
		n = 2
	}
	for {
		lvl := randomInt(0, min(4, n+1)) + 1
		p := filterDf(CharDB.Perks, "Level", "==", strconv.Itoa(lvl)).
			Col("Perk").Records()
		if len(p) == 0 {
			return errors.New("No available perks. Retrying")
		}
		idx := randomInt(0, len(p))
		CharDB.Perks = filterDf(CharDB.Perks, "Perk", "!=", p[idx])
		c.Perks = append(c.Perks, p[idx])
		if n -= lvl; n <= 0 {
			break
		}
	}
	return nil
}
