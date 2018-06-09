package m6ik

import (
	"errors"
	"regexp"
	"strings"
)

var reCarType = regexp.MustCompile(" \\(.*\\)")

func (c *Character) generateCareers(careerOpts string) error {
	careers := []string{}
	var firstCareer, secondCareer string
	if careerOpts == "" {
		careers = CharDB.Careers.Col("Career").Records()
		// sample first career
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
		} else {
			firstCareer = randomChoice(careers)
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
