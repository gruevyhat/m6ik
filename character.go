package m6ik

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	startingAttrDice  int     = 5
	startingSkillDice int     = 7
	pipsPerDie        int     = 3
	weightFactor      float64 = 3.0
)

var (
	Archetypes = []string{
		"Mighty", "Skilled", "Intellectual", "Gifted",
	}
	Attributes = []string{
		"Agility", "Perception", "Strength", "Intellect", "Technical", "Arcane",
	}
	GeneralSkills = []string{
		"Animal Handling", "Athletics", "Climbing", "Cultures", "Detection", "Driving", "First Aid", "Gambling", "Hide", "Intimidation", "Lifting", "Lore (*)", "Pilot (*)", "Riding", "Running", "Scholar", "Search", "Stamina", "Swimming", "Willpower",
	}
	reBonus   = regexp.MustCompile(`[+-]\d`)
	reCarType = regexp.MustCompile(" \\(.*\\)")
)

type StaticDefenses struct {
	Dodge int `json:"dodge"`
	Block int `json:"block"`
	Parry int `json:"parry"`
	Soak  int `json:"soak"`
	Sense int `json:"sense"`
}

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

type Experience struct {
	CharPoints int
	FatePoints int
}

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
	Personal
	StaticDefenses
	Experience
}

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
	}
	return j
}

func getKeys(m map[string]*Die) []string {
	a := []string{}
	for name := range m {
		a = append(a, name)
	}
	sort.Strings(a)
	return a
}

func (c *Character) promoteRandomAttribute(p Die) {
	attrs := getKeys(c.Attributes)
	weights := []float64{}
	for _, k := range attrs {
		weights = append(weights, c.AttrWeights[k])
	}
	a := weightedRandomChoice(attrs, weights)
	lessThanMax := p.toPips()+c.Attributes[a].toPips() <= c.Attributes[a].codeMax*pipsPerDie
	if c.AttrWeights[a] > 0.0 && lessThanMax {
		c.Attributes[a].add(p)
	} else {
		c.promoteRandomAttribute(p)
	}
}

func (c *Character) promoteRandomSkill(p Die) {
	sks := getKeys(c.Skills)
	weights := []float64{}
	for _, k := range sks {
		weights = append(weights, c.SkillWeights[k])
	}
	sk := weightedRandomChoice(sks, weights)
	lessThanMax := p.toPips()+c.Skills[sk].toPips() <= c.Skills[sk].codeMax*pipsPerDie
	if c.SkillWeights[sk] > 0.0 && lessThanMax {
		c.Skills[sk].add(p)
	} else {
		c.promoteRandomSkill(p)
	}
}

func (c *Character) promoteAttribute(attr string, p Die) {
	c.Attributes[attr].add(p)
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

func contains(arr []string, s string) bool {
	for _, a := range arr {
		if a == s { // || strings.HasPrefix(a, s) || strings.HasPrefix(s, a) {
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

func (c *Character) generateAttributes() {
	c.Attributes = make(map[string]*Die)
	c.AttrWeights = make(map[string]float64)
	for _, attr := range Attributes {
		if attr == "Arcane" {
			c.Attributes[attr] = &Die{0, 0, 5}
		} else {
			c.Attributes[attr] = &Die{2, 0, 5}
		}
		c.AttrWeights[attr] = 1.0
	}
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
		c.promoteAttribute(a, d)
	}
}

func (c *Character) generateArchetype(archetype string) {
	// Sample.
	if archetype == "" {
		archs := CharDB.Archetypes.Col("Archetype").Records()
		c.Archetype = randomChoice(archs)
	} else {
		c.Archetype = archetype
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
		c.AttrWeights["Arcane"] = 0.0
		for _, caster := range Casters {
			CharDB.filter("Careers", "Type", "!=", caster)
		}
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
	c.promoteAttribute(a, d)
}

func (c *Character) generateCareers(careerOpts string) error {
	careers := []string{}
	var (
		firstCareer, secondCareer string
	)
	if careerOpts == "" {
		careers = CharDB.Careers.Col("Career").Records()
		// sample first career
		firstCareer = randomChoice(careers)
		// Deal with Gifted
		if c.Archetype == "Gifted" {
			for _, car := range Casters {
				if !contains(careers, car) {
					Casters = remove(car, Casters)
				}
			}
			casters := dropIfNotIn(CharDB.Careers, "Type", Casters).
				Col("Career").Records()
			if len(casters) == 0 {
				return errors.New("Gifted character needs a magical career.")
			}
			firstCareer = randomChoice(casters)
		}
		// sample second career
		carRestr := filterDf(CharDB.Careers, "Career", "==", firstCareer).
			Col("Restrictions").Records()[0]
		if carRestr != "" {
			carRestrSplit := strings.Split(carRestr, ", ")
			careers = dropIfNotIn(CharDB.Careers, "Career", carRestrSplit).Col("Career").Records()
		}
		secondCareer = randomChoice(careers)
	} else {
		// Handle specified parameters
		careers = strings.Split(careerOpts, "/")
		firstCareer, secondCareer = careers[0], careers[1]
	}
	// Assign careers and filter db
	c.Careers = []string{firstCareer, secondCareer}
	CharDB.Careers = dropIfNotIn(CharDB.Careers, "Career", c.Careers)
	// Filter Perks
	occPerks := []string{}
	perks := CharDB.Careers.Col("Perks").Records()
	if len(perks) != len(c.Careers) {
		return errors.New("Impossible character generated. Retrying.")
	}
	for i := 0; i < len(c.Careers); i++ {
		occPerks = append(occPerks, strings.Split(perks[i], ", ")...)
	}
	CharDB.Perks = dropIfNotIn(CharDB.Perks, "Perk", occPerks)
	// TODO: Add special
	// TODO: Add assets
	// TODO: Add connections
	// TODO: Add money
	return nil
}

func parseSkillMax(skill string) (string, Die) {
	r, _ := regexp.Compile(`(.*) (\d)D`)
	sk := r.FindStringSubmatch(skill)
	m, _ := strconv.Atoi(sk[2])
	return sk[1], Die{codeMax: m}
}

func (c *Character) generateSkills() {
	c.Skills = make(map[string]*Die)
	c.SkillWeights = make(map[string]float64)
	skills := strings.Split(
		strings.Join(CharDB.Careers.Col("SkillMaximums").Records(), ", "), ", ")
	skills = append(skills, "Dodge 4D")
	skills = append(skills, "Search 4D")
	for _, skill := range GeneralSkills {
		if randomInt(0, 10) > 5 {
			skills = append(skills, fmt.Sprintf("%s 4D", skill))
		}
	}
	for _, s := range skills {
		n, d := parseSkillMax(s)
		r := filterDf(CharDB.Skills, "Skill", "==", n).Col("Attribute").Records()
		if len(r) == 0 {
			fmt.Println(n)
		}
		a := filterDf(CharDB.Skills, "Skill", "==", n).Col("Attribute").Records()[0]
		c.Skills[n] = &d
		c.Skills[n].addP(c.Attributes[a])
		c.SkillWeights[n] = c.AttrWeights[a]
	}
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

func (c *Character) generatePerks(n_perks string) {
	n, _ := strconv.Atoi(n_perks)
	if n == 0 {
		n = 2
	}
	perks := CharDB.Perks.Col("Perk").Records()
	c.Perks = sampleWithoutReplacement(perks, n)
}

func (c *Character) generateName(name string) {
	if name == "" {
		c.Name = "Nameless"
	} else {
		c.Name = name
	}
}

func NewCharacter(opts map[string]string) Character {

	NewCharDB()
	c := Character{}

	// Base stats
	c.generateAttributes()
	c.generateRace(opts["race"])
	c.generateArchetype(opts["archetype"])
	err := c.generateCareers(opts["careers"])
	if err != nil {
		fmt.Println(err)
		NewCharDB()
		c = NewCharacter(opts)
		return c
	}

	// Distribute 5D among Attr.
	for i := 0; i < startingAttrDice; i++ {
		c.promoteRandomAttribute(Die{code: 1})
	}
	// Distribute 7D among Skills and Perks
	c.generateSkills()
	for i := 0; i < startingSkillDice; i++ {
		c.promoteRandomSkill(Die{code: 1})
	}

	// Perks
	c.generatePerks(opts["n_perks"])

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
