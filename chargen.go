package chargen

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kniren/gota/dataframe"
	"github.com/kniren/gota/series"
)

var seed = rand.NewSource(time.Now().Unix())
var random = rand.New(seed)

var Races = []string{
	"Human", "Trollkin", "Nyss", "Ogrun", "Iosan", "Dwarf", "Gobber", "Tharn",
}

var Archetype = []string{
	"Mighty", "Skilled", "Intellectual", "Gifted",
}

type Race struct {
	maxAttributes Attributes
	minAttributes Attributes
	name          string
	restrictions  []string
}

type Die struct {
	code int
	pips int
}

func (d Die) Roll() int {
	return d.code*3 + d.pips
}

type Attributes struct {
	Agility    Die
	Perception Die
	Intellect  Die
	Strength   Die
	Technical  Die
	Arcane     Die
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
	Sex         string
	Height      string
	Weight      string
}

type Experience struct {
	CharPoints int
	FatePoints int
}

type Character struct {
	Personal
	Race      string
	Archetype string
	Careers   []string
	Attributes
	CharDB CharacterDatabase
	StaticDefenses
	Experience
}

func randomChoice(choices []string) string {
	return choices[random.Intn(len(choices))]
}

func randomInt(fromInt int, toInt int) int {
	return random.Intn(toInt) + fromInt
}

func generateAge() int {
	return randomInt(18, 40)
}

func generateSex() string {
	sexes := []string{"male", "female", "other"}
	return randomChoice(sexes)
}

func generateRace() string {
	return randomChoice(Races)
}

func generateArchetype() string {
	return randomChoice(Archetype)
}

func generateCareer() string {
	careers := make([]string, 0)
	for _, row := range CharDB.Careers {
		careers = append(careers, row.Career)
	}
	return randomChoice(careers)
}

/*
func generateName() string {
	return
}

func generateQuote() string {
	return
}

func generateAppearance() string {
	return
}

func generatePersonality() string {
	return
}

func generateHeight() string {
	return
}

func generateWeight() string {
	return
}

func calcDodge() string {
	return
}

func calcBlock() string {
	return
}

func calcParry() string {
	return
}

func calcSoak() string {
	return
}

func calcSense() string {
	return
}
*/

func GenerateCharacter() Character {

	c := Character{}
	c.Age = generateAge()
	c.Sex = generateSex()
	c.Race = generateRace()
	c.Archetype = generateArchetype()
	for i := 0; i < 2; i++ {
		c.Careers = append(c.Careers, generateCareer())
	}

	df := dataframe.LoadStructs(CharDB.Armors)
	fil := df.Filter(
		dataframe.F{"ArmorModifier", series.Eq, "1D+2"},
	)

	fmt.Println(fil)

	/*
		c.Name = generateName()
		c.Quote = generateQuote()
		c.Appearance = generateAppearance()
		c.Personality = generatePersonality()
		c.Height = generateHeight()
		c.Weight = generateWeight()
		c.Dodge = calcDodge()
		c.Block = calcBlock()
		c.Parry = calcParry()
		c.Soak = calcSoak()
		c.Sense = calcSense()
		c.CharPoints = 5
		c.FatePoints = 1
	*/

	return c
}
