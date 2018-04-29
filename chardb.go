package chargen

import (
	"encoding/json"
	"fmt"
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

type Characteristic []interface{}

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

type CharacterDatabase struct {
	Abilities []Ability
	Benefits  []Benefit
	Skills    []Skill
	Armors    []Armor
	Careers   []Career
	Spells    []Spell
	Weapons   []Weapon
}

func readJson(filename string) []byte {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return raw
}

func (c *CharacterDatabase) getAbilities() {
	raw := readJson(abilityFile)
	json.Unmarshal(raw, &c.Abilities)
}

func (c *CharacterDatabase) getSkills() {
	raw := readJson(skillFile)
	json.Unmarshal(raw, &c.Skills)
}

func (c *CharacterDatabase) getArmors() {
	raw := readJson(armorFile)
	json.Unmarshal(raw, &c.Armors)
}

func (c *CharacterDatabase) getCareers() {
	raw := readJson(careerFile)
	json.Unmarshal(raw, &c.Careers)
}

func (c *CharacterDatabase) getSpells() {
	raw := readJson(spellFile)
	json.Unmarshal(raw, &c.Spells)
}

func (c *CharacterDatabase) getBenefits() {
	raw := readJson(benefitFile)
	json.Unmarshal(raw, &c.Benefits)
}

func (c *CharacterDatabase) getWeapons() {
	raw := readJson(weaponFile)
	json.Unmarshal(raw, &c.Weapons)
}

var CharDB = CharacterDatabase{}

func init() {
	CharDB.getAbilities()
	CharDB.getArmors()
	CharDB.getBenefits()
	CharDB.getCareers()
	CharDB.getSkills()
	CharDB.getSpells()
	CharDB.getWeapons()
}
