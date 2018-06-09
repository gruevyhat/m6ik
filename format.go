package m6ik

import (
	"fmt"
	"sort"
	"strings"
)

func stringifyDice(dice map[string]*Die) string {
	keys := []string{}
	for k := range dice {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	outString := []string{}
	for _, k := range keys {
		if dice[k].code > 0 || dice[k].pips > 0 {
			outString = append(outString, fmt.Sprintf("%s %s", k, dice[k].toStr()))
		}
	}
	return strings.Join(outString, ", ")
}

func stringifyStatDef(sd StaticDefenses) string {
	return fmt.Sprintf("%s %d, %s %d, %s %d, %s %d, %s %d",
		"Dodge", sd.Dodge, "Block", sd.Block, "Parry", sd.Parry,
		"Soak", sd.Soak, "Sense", sd.Sense)
}

func stringifyPerks(perks []string) map[string]string {
	perkMap := map[string]string{}
	for _, p := range perks {
		perkMap[p] = filterDf(CharDB.Perks, "Perk", "==", p).
			Col("Description").Records()[0]
	}
	return perkMap
}

func (c Character) Print() {
	fmt.Println("Name\t" + c.Name)
	fmt.Println("Race\t" + c.Race)
	fmt.Println("Gender\t" + c.Gender)
	fmt.Println("Careers\t" + strings.Join(c.Careers, "/"))
	fmt.Println("Archetype\t" + c.Archetype)
	fmt.Println("Attributes\t", stringifyDice(c.Attributes))
	fmt.Println("Skills\t", stringifyDice(c.Skills))
	fmt.Println("Perks\t", strings.Join(c.Perks, ", "))
	fmt.Println("Static Def.\t", stringifyStatDef(c.StaticDefenses))
	fmt.Println("Random Seed\t", c.Seed)
}

type charJSON struct {
	Name       string            `json:"name"`
	Race       string            `json:"race"`
	Gender     string            `json:"gender"`
	Age        int               `json:"age"`
	Careers    []string          `json:"careers"`
	Archetype  string            `json:"archetype"`
	Attributes []string          `json:"attributes"`
	Skills     []string          `json:"skills"`
	Perks      map[string]string `json:"perks"`
	StaticDef  StaticDefenses    `json:"statdefs"`
	Seed       string            `json:"seed"`
}

func (c Character) ToJSON() charJSON {
	j := charJSON{
		Name:       c.Name,
		Age:        c.Age,
		Race:       c.Race,
		Gender:     c.Gender,
		Careers:    c.Careers,
		Archetype:  c.Archetype,
		Attributes: strings.Split(stringifyDice(c.Attributes), ", "),
		Skills:     strings.Split(stringifyDice(c.Skills), ", "),
		Perks:      stringifyPerks(c.Perks),
		StaticDef:  c.StaticDefenses,
		Seed:       c.Seed,
	}
	return j
}
