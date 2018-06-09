package m6ik

import (
	"strconv"
)

type Personal struct {
	Name        string
	Quote       string
	Appearance  string
	Personality string
	Age         int
	Gender      string
	Height      string
	Weight      string
}

func (c *Character) generateAge(age string) {
	// TODO: Fix for race.
	if age == "" {
		c.Age = randomInt(15, 35)
	} else {
		c.Age, _ = strconv.Atoi(age)
	}
}

func (c *Character) generateGender(gender string) {
	if gender == "" {
		sexes := []string{"Male", "Female", "Other"}
		weights := []float64{0.45, 0.45, 0.1}
		c.Gender = weightedRandomChoice(sexes, weights)
	} else {
		c.Gender = gender
	}
}

func (c *Character) generateName(name string) {
	if name == "" {
		c.Name = "Nameless"
	} else {
		c.Name = name
	}
}
