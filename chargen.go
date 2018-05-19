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
	Benefits     []string
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

func (c Character) Print() {
	fmt.Println("Race\t" + c.Race)
	fmt.Println("Gender\t" + c.Gender)
	fmt.Println("Careers\t" + strings.Join(c.Careers, "/"))
	fmt.Println("Archetype\t" + c.Archetype)
	fmt.Println("Attributes\t", stringifyDice(c.Attributes))
	fmt.Println("Skills\t", stringifyDice(c.Skills))
}

func getMapKeys(m *map[string]Die) []string {
	keys := []string{}
	for k := range *m {
		keys = append(keys, k)
	}
	return keys
}

func (c *Character) promoteRandomAttribute(p Die) {
	attr := randomChoice(Attributes)
	if c.AttrWeights[attr] > 0.0 {
		c.Attributes[attr].add(p)
	} else {
		c.promoteRandomAttribute(p)
	}
}

func (c *Character) promoteRandomSkill(p Die) {
	skills := []string{}
	for k, _ := range c.Skills {
		skills = append(skills, k)
	}
	sk := randomChoice(skills)
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
	c.Gender = randomChoice(sexes)

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
	c.AttrWeights[attr] += 1.0
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
	// TODO: Deal with career type
	// TODO: Add Perks
	// TODO: Set skills maximums
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

func NewCharacter() Character {

	c := Character{}

	// Initialize attributes.
	c.Attributes = make(map[string]*Die)
	c.AttrWeights = make(map[string]float64)
	for _, attr := range Attributes {
		if attr != "Arcane" {
			c.Attributes[attr] = &Die{2, 0, 2}
		} else {
			c.Attributes[attr] = &Die{0, 0, 0}
		}
		c.AttrWeights[attr] = 1.0
	}

	// Add personal characteristics
	c.generateAge()
	c.generateGender()

	// Select a Race, Archetype, and two Careers
	c.generateRace()
	c.generateArchetype()
	c.generateCareers()

	// Distribute 5D among Attr.
	for i := 0; i < startingAttrD; i++ {
		c.promoteRandomAttribute(Die{1, 0, 1})
	}

	// Distribute 12D among Skills and Perks
	c.generateSkills()
	for i := 0; i < startingSkillD; i++ {
		c.promoteRandomSkill(Die{1, 0, 0})
	}

	// Purchase equipment

	/*
		c.Name = generateName()
		c.Quote = generateQuote()
		c.Appearance = generateAppearance()
		c.Personality = generatePersonality()
		c.Height = generateHeight()
		c.Weight = generateWeight()
		c.calcStaticDefenses()
	*/

	c.CharPoints = 5
	c.FatePoints = 1

	fmt.Println(c)
	return c
}
