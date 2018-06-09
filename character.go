package m6ik

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("m6ik")

const (
	baseAttrDice     int     = 10
	defaultAttrDice  string  = "15"
	defaultSkillDice string  = "7"
	pipsPerDie       int     = 3
	weightFactor     float64 = 3.0
)

type Character struct {
	Race         string
	Archetype    string
	Careers      []string
	Attributes   map[string]*Die
	AttrWeights  map[string]float64
	Abilities    []string
	Perks        []string
	Skills       map[string]*Die
	SkillWeights map[string]float64
	Spells       []string
	Weapons      []string
	Armors       []string
	Seed         string
	Personal
	StaticDefenses
	Experience
}

func getKeys(m map[string]*Die) []string {
	a := []string{}
	for name := range m {
		a = append(a, name)
	}
	//sort.Strings(a)
	return a
}

func contains(arr []string, s string) bool {
	for _, a := range arr {
		if a == s || strings.HasPrefix(a, s) || strings.HasSuffix(s, a) {
			return true
		}
	}
	return false
}

func remove(s string, a []string) []string {
	for i, x := range a {
		if x == "" || x == s {
			a = append(a[:i], a[i+1:]...)
		}
	}
	return a
}

func (c *Character) applyConstraints(restr string) {
	if restr != "" {
		//restr = reCarType.ReplaceAllString(restr, "")
		restrs := strings.Split(restr, ", ")
		for _, r := range restrs {
			CharDB.filter("Archetypes", "Archetype", "!=", r)
			CharDB.filter("Careers", "Career", "!=", r)
		}
	}
}

var reBonus = regexp.MustCompile(`[+-]\d`)

func parseBonus(bonus string) (string, Die) {
	attr := strings.Split(bonus, " ")
	if strings.Index(attr[0], "D") == -1 {
		attr[0] = "+0D" + attr[0]
	}
	mods := reBonus.FindAllString(attr[0], -1)
	modsInt := [2]int{}
	for i, b := range mods {
		modsInt[i], _ = strconv.Atoi(b)
	}
	return attr[1], Die{code: modsInt[0], pips: modsInt[1]}
}

func NewCharacter(opts map[string]string) Character {

	NewCharDB()

	c := Character{}

	// Set seed
	if opts["seed"] == "" {
		c.Seed = seed
		log.Info("NEW SEED:", c.Seed)
	} else {
		c.Seed = opts["seed"]
		random = setSeed(c.Seed)
		log.Info("OLD SEED:", c.Seed)
	}

	// Base stats
	c.generateAttributes()

	// Race
	c.generateRace(opts["race"])

	// Archetype
	if err := c.generateArchetype(opts["archetype"]); err != nil {
		log.Warning(err)
		c = NewCharacter(map[string]string{})
		return c
	}

	// Careers
	if err := c.generateCareers(opts["careers"]); err != nil {
		log.Warning(err)
		c = NewCharacter(opts)
		return c
	}

	// Distribute dice among Attr.
	if err := c.distributeAttrDice(opts["n_attrs"]); err != nil {
		log.Warning(err)
		opts["n_attrs"] = defaultAttrDice
		c = NewCharacter(opts)
		return c
	}

	// Distribute dice among Skills
	c.generateSkills()
	if err := c.distributeSkillDice(opts["n_skills"]); err != nil {
		log.Warning(err)
		opts["n_skills"] = defaultSkillDice
		c = NewCharacter(opts)
		return c
	}

	// Perks
	if err := c.generatePerks(opts["n_perks"]); err != nil {
		log.Warning(err)
		opts["seed"] = ""
		c = NewCharacter(opts)
		return c
	}

	// Static Defenses
	c.generateStaticDefenses()

	// Flavor
	c.generateName(opts["name"])
	c.generateAge(opts["age"])
	c.generateGender(opts["gender"])
	/*
		c.Quote = generateQuote()
		c.Appearance = generateAppearance()
		c.Personality = generatePersonality()
		c.Height = generateHeight()
		c.Weight = generateWeight()
		c.Equipment = generateEquipment()
	*/

	// Misc.
	c.CharPoints = 5
	c.FatePoints = 1

	return c
}
