package m6ik

import (
	"errors"
	"strconv"
)

var (
	Attributes = []string{
		"Agility", "Perception", "Strength", "Intellect", "Technical", "Arcane",
	}
)

func (c *Character) promoteRandomAttribute(p Die) {
	attrs := getKeys(c.Attributes)
	weights := []float64{}
	for _, k := range attrs {
		weights = append(weights, c.AttrWeights[k])
	}
	a := weightedRandomChoice(attrs, weights)
	lessThanMax := p.toPips()+c.Attributes[a].toPips() <= c.Attributes[a].codeMax*pipsPerDie+2
	if c.AttrWeights[a] > 0.0 && lessThanMax {
		c.Attributes[a].add(p)
	} else {
		c.promoteRandomAttribute(p)
	}
}

func (c *Character) promoteAttribute(attr string, p Die) {
	c.Attributes[attr].add(p)
}

func (c *Character) generateAttributes() {
	c.Attributes = make(map[string]*Die)
	c.AttrWeights = make(map[string]float64)
	for _, attr := range Attributes {
		if attr == "Arcane" {
			c.Attributes[attr] = &Die{0, 0, 4}
		} else {
			c.Attributes[attr] = &Die{2, 0, 4}
		}
		c.AttrWeights[attr] = 1.0
	}
}

func (c *Character) distributeAttrDice(nAttrs string) error {
	if nAttrs == "" {
		nAttrs = defaultAttrDice
	}
	attrDice, _ := strconv.Atoi(nAttrs)
	if attrDice > baseAttrDice {
		attrDice -= baseAttrDice
	}
	if attrDice > len(Attributes)*2 {
		return errors.New("Invalid number of attribute dice. Using default.")
	}
	for i := 0; i < attrDice; i++ {
		c.promoteRandomAttribute(Die{code: 1})
	}
	return nil
}
