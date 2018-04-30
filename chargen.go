package chargen

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var seed = rand.NewSource(time.Now().Unix())
var random = rand.New(seed)

const (
	startingAttributes int = 15
	diePips            int = 3
)

var Races = []string{
	"Human", "Trollkin", "Nyss", "Ogrun", "Iosan", "Dwarf", "Gobber", "Tharn",
}

var Archetype = []string{
	"Mighty", "Skilled", "Intellectual", "Gifted",
}

var Attributes = []string{
	"Agility", "Perception", "Strength", "Intellect", "Technical", "Arcane",
}

type Die struct {
	code int
	pips int
}

func (d Die) toStr() string {
	var dieStr string
	if d.pips > 0 {
		dieStr = strconv.Itoa(d.code) + "D+" + strconv.Itoa(d.pips)
	} else {
		dieStr = strconv.Itoa(d.code) + "D"
	}
	return dieStr
}

func (d Die) Roll() int {
	return d.code*diePips + d.pips
}

func (d *Die) add(e Die) {
	d.code += e.code
	d.pips += e.pips
	if d.pips > 3 {
		d.code = d.code + d.pips/diePips
		d.pips = d.pips % diePips
	}
}

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
	Race       string
	Archetype  string
	Careers    []string
	Attributes map[string]*Die
	Abilities  []string
	Benefits   []string
	Skills     map[string]*Die
	Spells     []string
	Weapons    []string
	Armors     []string
	StaticDefenses
	Experience
}

func (c Character) Print() {
	fmt.Println("Race\t" + c.Race)
	fmt.Println("Gender\t" + c.Gender)
	fmt.Println("Careers\t" + strings.Join(c.Careers, "/"))
	for _, attr := range Attributes {
		fmt.Printf("%s\t%s\n", attr, c.Attributes[attr].toStr())
	}
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
	c.Attributes[attr].add(p)
}

func (c *Character) promoteAttribute(attr string, p Die) {
	c.Attributes[attr].add(p)
}

func randomChoice(choices []string) string {
	return choices[random.Intn(len(choices))]
}

func randomInt(min, max int) int {
	return random.Intn(max-min) + min
}

func generateAge() int {
	// TODO: Fix for race.
	return randomInt(18, 40)
}

func generateGender() string {
	sexes := []string{"Male", "Female", "Other"}
	return randomChoice(sexes)
}

func generateRace() string {
	// TODO: Impose restrictions on Archetype.
	// TODO: Apply bonuses.
	return randomChoice(Races)
}

func generateArchetype() string {
	// TODO: Impose restrictions on Careers.
	// TODO: Adjust Attr. min/max.
	return randomChoice(Archetype)
}

func generateCareer() string {
	// TODO: Impose restrictions on Skills.
	careers := CharDB.Careers.Col("Career").Records()
	return randomChoice(careers)
}

func dieToPips(d Die) int {
	return d.code*diePips + d.pips
}

func pipsToDie(p int) Die {
	d := Die{code: int(p / diePips), pips: p % diePips}
	return d
}

func GenerateCharacter() Character {

	c := Character{}

	// Add personal characteristics
	c.Age = generateAge()
	c.Gender = generateGender()

	// Select a Race, Archetype, and two Careers
	c.Race = generateRace()
	c.Archetype = generateArchetype()
	for i := 0; i < 2; i++ {
		c.Careers = append(c.Careers, generateCareer())
	}

	// Distribute 15D among Attr.
	c.Attributes = make(map[string]*Die)
	for _, attr := range Attributes {
		c.Attributes[attr] = &Die{}
	}
	for i := 0; i < startingAttributes*diePips; i++ {
		c.promoteRandomAttribute(Die{0, 1})
	}

	// Distribute 12D among Skills and Perks
	c.Skills = make(map[string]*Die)

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

	return c
}
