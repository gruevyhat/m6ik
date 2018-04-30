package chargen

import (
	"encoding/json"
	"fmt"
	. "github.com/kniren/gota/dataframe"
	"io/ioutil"
	"os"
)

const (
	abilityFile string = "assets/abilities.json"
	armorFile   string = "assets/armors.json"
	benefitFile string = "assets/benefits.json"
	careerFile  string = "assets/careers.json"
	skillFile   string = "assets/skills.json"
	spellFile   string = "assets/spells.json"
	weaponFile  string = "assets/weapons.json"
)

type CharacterDatabase struct {
	Abilities DataFrame
	Armors    DataFrame
	Benefits  DataFrame
	Careers   DataFrame
	Skills    DataFrame
	Spells    DataFrame
	Weapons   DataFrame
}

type Ability struct {
	Ability        string `json:"Ability"`
	Cp             string `json:"CP Cost"`
	Prerequisites  string `json:"Prerequisites"`
	Description    string `json:"Description"`
	Notes          string `json:"Notes"`
	OldDescription string `json:"Old IKRPG Description"`
}

type Benefit struct {
	Benefit             string `json:"Benefit"`
	Archetype           string `json:"Archetype"`
	Description         string `json:"Description"`
	OldIKRPGDescription string `json:"Old IKRPG Description"`
}

type Skill struct {
	Skill         string `json:"skill"`
	Attribute     string `json:"attribute"`
	SkillType     string `json:"skillType"`
	Social        string `json:"social"`
	Advanced      string `json:"advanced"`
	Prerequisites string `json:"prerequisites"`
	Description   string `json:"description"`
	IkrpgSkill    string `json:"ikrpgSkill"`
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
	Abilities           string  `json:"Abilities"`
	SkillMaximums       string  `json:"Skill Maximums"`
	Restrictions        string  `json:"Restrictions"`
	Special             string  `json:"Special"`
	StartingAssets      string  `json:"Starting Assets"`
	StartingConnections string  `json:"Starting Connections"`
	StartingMoney       float32 `json:"Starting Money"`
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
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return raw
}

func (c *CharacterDatabase) Build() {
	// Abilities
	var abils = []Ability{}
	json.Unmarshal(readJson(abilityFile), &abils)
	c.Abilities = LoadStructs(abils)
	// Armors
	var armors = []Armor{}
	json.Unmarshal(readJson(armorFile), &armors)
	c.Armors = LoadStructs(armors)
	// Benefits
	var bens = []Benefit{}
	json.Unmarshal(readJson(benefitFile), &bens)
	c.Benefits = LoadStructs(bens)
	// Careers
	var cars = []Career{}
	json.Unmarshal(readJson(careerFile), &cars)
	c.Careers = LoadStructs(cars)
	// Skills
	var skills = []Skill{}
	json.Unmarshal(readJson(skillFile), &skills)
	c.Skills = LoadStructs(skills)
	// Spells
	var spells = []Spell{}
	json.Unmarshal(readJson(spellFile), &spells)
	c.Abilities = LoadStructs(spells)
	// Weapons
	var weaps = []Weapon{}
	json.Unmarshal(readJson(weaponFile), &weaps)
	c.Weapons = LoadStructs(weaps)
}

var CharDB = CharacterDatabase{}

func init() {
	CharDB.Build()
}
