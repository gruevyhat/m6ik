package chargen

import (
	"encoding/json"
	"fmt"
	"github.com/kniren/gota/dataframe"
	"github.com/kniren/gota/series"
	"io/ioutil"
	"os"
	"strings"
)

func getDataDir() string {
	dir := os.Getenv("M6IK")
	fmt.Println(dir)
	if dir != "" {
		return dir + "/"
	}
	return "./"
}

var (
	dataDir       = getDataDir()
	perkFile      = dataDir + "assets/perks.json"
	armorFile     = dataDir + "assets/armors.json"
	archetypeFile = dataDir + "assets/archetypes.json"
	careerFile    = dataDir + "assets/careers.json"
	raceFile      = dataDir + "assets/races.json"
	skillFile     = dataDir + "assets/skills.json"
	spellFile     = dataDir + "assets/spells.json"
	weaponFile    = dataDir + "assets/weapons.json"
)

const (
	startErr = "\033[31m"
	endErr   = "\033[0m"
)

type CharacterDatabase struct {
	Perks      dataframe.DataFrame
	Armors     dataframe.DataFrame
	Archetypes dataframe.DataFrame
	Careers    dataframe.DataFrame
	Races      dataframe.DataFrame
	Skills     dataframe.DataFrame
	Spells     dataframe.DataFrame
	Weapons    dataframe.DataFrame
}

type Perk struct {
	Perk          string `json:"Perk"`
	Level         int    `json:"Level"`
	Prerequisites string `json:"Prerequisites"`
	Description   string `json:"Description"`
}

type Archetype struct {
	Archetype     string `json:"Archetype"`
	Bonus         string `json:"Bonus"`
	Proscriptions string `json:"Proscriptions"`
}

type Skill struct {
	Skill         string `json:"Skill"`
	Type          string `json:"Type"`
	Attribute     string `json:"Attribute"`
	SkillType     string `json:"Skill Type"`
	Social        string `json:"Social"`
	Advanced      string `json:"Advanced"`
	Prerequisites string `json:"Prerequisites"`
	Description   string `json:"Description"`
}

type Armor struct {
	Armor           string  `json:"Armor"`
	AgilityModifier string  `json:"Agility Modifier"`
	ArmorModifier   string  `json:"Armor Modifier"`
	PowerField      string  `json:"Power Field"`
	Cost            float32 `json:"Cost"`
	SpecialRules    string  `json:"Special Rules"`
}

type Career struct {
	Career              string  `json:"Career"`
	Type                string  `json:"Type"`
	Perks               string  `json:"Occupational Perks"`
	SkillMaximums       string  `json:"Skill Maximums"`
	Restrictions        string  `json:"Restrictions"`
	Special             string  `json:"Special"`
	StartingAssets      string  `json:"Starting Assets"`
	StartingConnections string  `json:"Starting Connections"`
	StartingMoney       float32 `json:"Starting Money"`
}

type Race struct {
	Race          string `json:"Race"`
	Type          string `json:"Type"`
	Attributes    string `json:"Attributes"`
	Skills        string `json:"Skills"`
	Perks         string `json:"Perks"`
	Special       string `json:"Special"`
	Proscriptions string `json:"Proscriptions"`
	Description   string `json:"Description"`
}

type Spell struct {
	Spell     string `json:"Abbrev."`
	Technique string `json:"Technique"`
	Form      string `json:"Form"`
	TN        string `json:"TN"`
	Range     string `json:"Range"`
	Duration  string `json:"Duration"`
	Target    string `json:"Target"`
	Effect    string `json:"Effect"`
}

type Weapon struct {
	Weapon   string  `json:"Weapon"`
	Skill    string  `json:"Skill"`
	Modifier string  `json:"Modifier"`
	Damage   string  `json:"Damage"`
	Range    string  `json:"Range"`
	AOE      string  `json:"AOE (m)"`
	Ammo     string  `json:"Ammo"`
	Scale    string  `json:"Scale"`
	Cost     float32 `json:"Cost"`
	Special  string  `json:"Special"`
}

func readJson(filename string) []byte {
	raw := logger(ioutil.ReadFile)(filename)
	return raw
}

func (c *CharacterDatabase) Build() {
	// Perks
	var perks = []Perk{}
	json.Unmarshal(readJson(perkFile), &perks)
	c.Perks = dataframe.LoadStructs(perks)
	// Armors
	var armors = []Armor{}
	json.Unmarshal(readJson(armorFile), &armors)
	c.Armors = dataframe.LoadStructs(armors)
	// Archetypes
	var archs = []Archetype{}
	json.Unmarshal(readJson(archetypeFile), &archs)
	c.Archetypes = dataframe.LoadStructs(archs)
	// Careers
	var cars = []Career{}
	json.Unmarshal(readJson(careerFile), &cars)
	c.Careers = dataframe.LoadStructs(cars)
	// Races
	var races = []Race{}
	json.Unmarshal(readJson(raceFile), &races)
	c.Races = dataframe.LoadStructs(races)
	// Skills
	var skills = []Skill{}
	json.Unmarshal(readJson(skillFile), &skills)
	c.Skills = dataframe.LoadStructs(skills)
	// Spells
	var spells = []Spell{}
	json.Unmarshal(readJson(spellFile), &spells)
	c.Spells = dataframe.LoadStructs(spells)
	// Weapons
	var weaps = []Weapon{}
	json.Unmarshal(readJson(weaponFile), &weaps)
	c.Weapons = dataframe.LoadStructs(weaps)
}

var CharDB = CharacterDatabase{}

var ops = map[string]series.Comparator{
	"==": series.Eq,
	"!=": series.Neq,
}

func dropIfNotIn(df dataframe.DataFrame, col string, vals []string) dataframe.DataFrame {
	filters := make([]dataframe.F, len(vals))
	for i, val := range vals {
		filters[i] = dataframe.F{col, series.Eq, val}
	}
	newdf := df.Filter(filters...)
	return newdf
}

func filterDf(df dataframe.DataFrame, col, op, val string) dataframe.DataFrame {
	newdf := df.Filter(
		dataframe.F{col, ops[op], val},
	)
	return newdf
}

func (db *CharacterDatabase) filter(table, col, op, val string) {
	switch table {
	case "Perks":
		db.Perks = filterDf(db.Perks, col, op, val)
	case "Armors":
		db.Armors = filterDf(db.Armors, col, op, val)
	case "Archetypes":
		db.Archetypes = filterDf(db.Archetypes, col, op, val)
	case "Careers":
		db.Careers = filterDf(db.Careers, col, op, val)
	case "Races":
		db.Races = filterDf(db.Races, col, op, val)
	case "Skills":
		db.Skills = filterDf(db.Skills, col, op, val)
	case "Spells":
		db.Spells = filterDf(db.Spells, col, op, val)
	case "Weapons":
		db.Weapons = filterDf(db.Weapons, col, op, val)
	}
}

var Casters []string

func init() {
	CharDB.Build()
	Casters = strings.Split(CharDB.Archetypes.Col("Proscriptions").Records()[0], ", ")
}
