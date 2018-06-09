package m6ik

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var reSkill = regexp.MustCompile(`(.*) (\d)D`)

func parseSkillMax(skill string) (string, Die) {
	sk := reSkill.FindStringSubmatch(skill)
	if len(sk) < 2 {
		fmt.Println(skill, sk)
	}
	m, _ := strconv.Atoi(sk[2])
	return sk[1], Die{codeMax: m}
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

func (c *Character) generateSkills() {
	c.Skills = make(map[string]*Die)
	c.SkillWeights = make(map[string]float64)
	skills := strings.Split(
		strings.Join(CharDB.Careers.Col("SkillMaximums").Records(), ", "), ", ")
	// Add skill taxation.
	skills = append(skills, "Dodge 4D")
	skills = append(skills, "Search 4D")
	// Add a number (Int D * 3) of random general skills.
	for n := 0; n <= c.Attributes["Intellect"].code*3; n++ {
		idx := randomInt(0, len(GeneralSkills))
		skills = append(skills, fmt.Sprintf("%s 4D", GeneralSkills[idx]))
	}
	for _, s := range skills {
		n, d := parseSkillMax(s)
		// Deal with wildcard skills.
		if strings.HasSuffix(n, "(*)") {
			pfx := n[0 : len(n)-4]
			sk := filterDf(CharDB.Skills, "Type", "==", pfx).Col("Skill").Records()
			idx := randomInt(0, len(sk))
			n = GeneralSkills[idx]
		}
		// Set max and add points.
		a := filterDf(CharDB.Skills, "Skill", "==", n).Col("Attribute").Records()[0]
		c.Skills[n] = &d
		c.Skills[n].codeMax += c.Attributes[a].code
		c.Skills[n].addP(c.Attributes[a])
		c.SkillWeights[n] = c.AttrWeights[a]
	}
	// TODO: Check skill prereqs.
}

func (c *Character) distributeSkillDice(nSkills string) error {
	if nSkills == "" {
		nSkills = defaultSkillDice
	}
	skillDice, _ := strconv.Atoi(nSkills)
	if skillDice > len(c.Skills)*4 {
		return errors.New("Invalid number of skill dice. Using default.")
	}
	for i := 0; i < skillDice; i++ {
		c.promoteRandomSkill(Die{code: 1})
	}
	return nil
}
