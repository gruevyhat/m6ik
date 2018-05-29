package m6ik

import (
	"sort"
	"strconv"
	"strings"
	"testing"
)

type opt map[string]string

func TestNewCharacterParams(t *testing.T) {
	opts := []opt{
		opt{
			"name":      "Rasputin",
			"gender":    "male",
			"age":       "11",
			"race":      "Human (Khard)",
			"careers":   "Arcanist (Greylord)/Soldier",
			"archetype": "Gifted",
			"n_perks":   "2",
			"seed":      "",
		},
		opt{
			"name":      "Borkenhekenaken",
			"gender":    "male",
			"age":       "22",
			"race":      "Gobber",
			"careers":   "Spy/Thief",
			"archetype": "Skilled",
			"n_perks":   "2",
			"seed":      "",
		},
		opt{
			"name":      "Xev",
			"gender":    "female",
			"age":       "33",
			"race":      "Human (Thurian)",
			"careers":   "Knight/Stormblade",
			"archetype": "Intellectual",
			"n_perks":   "2",
			"seed":      "",
		},
	}
	for _, o := range opts {
		c := NewCharacter(o)
		if c.Name != o["name"] {
			t.Errorf("Incorrect name. Expected '%s'. Found '%s'.", c.Name, o["name"])
		}
		if c.Gender != o["gender"] {
			t.Errorf("Incorrect gender. Expected '%s'. Found '%s'.", c.Gender, o["gender"])
		}
		age, err := strconv.Atoi(o["age"])
		if err != nil || c.Age != age {
			t.Errorf("Incorrect age. Expected '%d'. Found '%s'.", c.Age, o["age"])
		}
		if c.Race != o["race"] {
			t.Errorf("Incorrect race. Expected '%s'. Found '%s'.", c.Race, o["race"])
		}
		sort.Strings(c.Careers)
		car := strings.Join(c.Careers, "/")
		if car != o["careers"] {
			t.Errorf("Incorrect careers. Expected '%s'. Found '%s'.", c.Careers, o["careers"])
		}
	}
}

func TestNewCharacterRandom(t *testing.T) {
	o := opt{
		"name": "", "gender": "", "age": "", "race": "",
		"careers": "", "archetype": "", "n_perks": "",
		"seed": "1532de8a7946fc4c",
	}
	c := NewCharacter(o)
	if c.Race != "Human (Scharde)" {
		t.Errorf("Incorrect race. Expected 'Human (Scharde)'. Found '%s'.", o["race"])
	}
	car := strings.Join(c.Careers, "/")
	if car != "Duelist/Pirate" {
		t.Errorf("Incorrect careers. Expected 'Duelist/Pirate'. Found '%s'.", o["careers"])
	}

}
