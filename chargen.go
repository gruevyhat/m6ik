package chargen

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	startingAttrD  int = 5
	startingSkillD int = 7
	pipsPerDie     int = 3
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
	Dodge int
	Block int
	Parry int
	Soak  int
	Sense int
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
	Personal
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

func stringifyStatDef(c Character) string {
	return fmt.Sprintf("%s %d, %s %d, %s %d, %s %d, %s %d",
		"Dodge", c.Dodge, "Block", c.Block, "Parry", c.Parry,
		"Soak", c.Soak, "Sense", c.Sense)
}

func (c Character) Print() {
	fmt.Println("Race\t" + c.Race)
	fmt.Println("Gender\t" + c.Gender)
	fmt.Println("Careers\t" + strings.Join(c.Careers, "/"))
	fmt.Println("Archetype\t" + c.Archetype)
	fmt.Println("Attributes\t", stringifyDice(c.Attributes))
	fmt.Println("Skills\t", stringifyDice(c.Skills))
	fmt.Println("Perks\t", strings.Join(c.Perks, ", "))
	fmt.Println("Static Def.\t", stringifyStatDef(c))
}

func getMapKeys(m *map[string]Die) []string {
	keys := []string{}
	for k := range *m {
		keys = append(keys, k)
	}
	return keys
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
	ltMax := (p.code + c.Attributes[a].code) < c.Attributes[a].codeMax
	if c.AttrWeights[a] > 0.0 && ltMax {
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
	ltMax := (p.code + c.Skills[sk].code) < c.Skills[sk].codeMax
	if c.SkillWeights[sk] > 0.0 && ltMax {
		c.Skills[sk].add(p)
	} else {
		c.promoteRandomSkill(p)
	}
}

func (c *Character) promoteAttribute(attr string, p Die) {
	c.Attributes[attr].add(p)
}

func (c *Character) generateAge() {
	// TODO: Fix for race.
	c.Age = randomInt(18, 40)
}

func (c *Character) generateGender() {
	sexes := []string{"Male", "Female", "Other"}
	weights := []float64{0.45, 0.45, 0.1}
	c.Gender = weightedRandomChoice(sexes, weights)
}

func contains(arr []string, s string) bool {
	for _, a := range arr {
		if a == s || strings.HasPrefix(a, s) || strings.HasPrefix(s, a) {
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
		if attr != "Arcane" {
			c.Attributes[attr] = &Die{2, 0, 5}
		} else {
			c.Attributes[attr] = &Die{0, 0, 5}
		}
		c.AttrWeights[attr] = 1.0
	}
}

func (c *Character) generateRace() {
	// Sample.
	races := CharDB.Races.Col("Race").Records()
	c.Race = randomChoice(races)
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

func (c *Character) generateArchetype() {
	// Sample.
	archs := CharDB.Archetypes.Col("Archetype").Records()
	c.Archetype = randomChoice(archs)
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
	c.AttrWeights[attr] += 5.0
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

func (c *Character) generateCareers() {
	// sample first career
	careers := CharDB.Careers.Col("Career").Records()
	firstCareer := randomChoice(careers)
	// Deal with Gifted
	if c.Archetype == "Gifted" {
		for _, car := range Casters {
			if !contains(careers, car) {
				Casters = remove(car, Casters)
			}
		}
		casters := dropIfNotIn(CharDB.Careers, "Type", Casters).
			Col("Career").Records()
		firstCareer = randomChoice(casters)
	}
	// sample second career
	carRestr := filterDf(CharDB.Careers, "Career", "==", firstCareer).
		Col("Restrictions").Records()[0]
	if carRestr != "" {
		carRestrSplit := strings.Split(carRestr, ", ")
		careers = dropIfNotIn(CharDB.Careers, "Career", carRestrSplit).Col("Career").Records()
	}
	secondCareer := randomChoice(careers)
	c.Careers = []string{firstCareer, secondCareer}
	CharDB.Careers = dropIfNotIn(CharDB.Careers, "Career", c.Careers)
	// Filter Perks
	occPerks := []string{}
	for i := 0; i < len(c.Careers); i++ {
		occPerks = append(occPerks, strings.Split(CharDB.Careers.Col("Perks").Records()[i], ", ")...)
	}
	CharDB.Perks = dropIfNotIn(CharDB.Perks, "Perk", occPerks)
	// TODO: Add special
	// TODO: Add assets
	// TODO: Add connections
	// TODO: Add money
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

func (c *Character) calcStaticDefenses() {
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

func (c *Character) generatePerks() {
	perks := CharDB.Perks.Col("Perk").Records()
	for i := 0; i < 2; i++ {
		c.Perks = append(c.Perks, randomChoice(perks))
	}
}

func NewCharacter() Character {

	c := Character{}

	// Base stats
	c.generateAttributes()
	c.generateRace()
	c.generateArchetype()
	c.generateCareers()

	// Distribute 5D among Attr.
	for i := 0; i < startingAttrD; i++ {
		c.promoteRandomAttribute(Die{code: 1})
	}
	// Distribute 7D among Skills and Perks
	c.generateSkills()
	for i := 0; i < startingSkillD; i++ {
		c.promoteRandomSkill(Die{code: 1})
	}

	// Perks
	c.generatePerks()

	// Static Defenses
	c.calcStaticDefenses()

	// Purchase equipment

	/*
		c.Name = generateName()
		c.Quote = generateQuote()
		c.Appearance = generateAppearance()
		c.Personality = generatePersonality()
		c.Height = generateHeight()
		c.Weight = generateWeight()
	*/

	// Flavor
	c.generateAge()
	c.generateGender()

	// Misc.
	c.CharPoints = 5
	c.FatePoints = 1

	fmt.Println(c)
	return c
}
