package main

import (
	"os"

	"github.com/docopt/docopt-go"
	"github.com/gruevyhat/m6ik"
)

var usage = `M6IK Character Generator

Usage: m6ikgen [options]

Options:
  --name	The character's full name.
  --gender	The character's gender.
  --race	The race and ethnicity.
  --archetype	Mighty, Skilled, Intellectual, or Gifted.  
  --careers	Slash-delimited career list (e.g., Soldier/Spy).
  --perks	Number of random perks to assign.
	--attrs	Number of starting attribute dice. [default: 15]
	--skills	Number of starting skill dice. [default: 7]
  --seed	Character generation signature.
	--log-level	One of {INFO, WARNING, ERROR}. [default: ERROR]
  -h --help
  --version
`

var Opts struct {
	Name       string `docopt:"--name"`
	Gender     string `docopt:"--gender"`
	Age        string `docopt:"--age"`
	Race       string `docopt:"--race"`
	Careers    string `docopt:"--careers"`
	Archetype  string `docopt:"--archetype"`
	NPerks     string `docopt:"--perks"`
	NAttrDice  string `docopt:"--attrs"`
	NSkillDice string `docopt:"--skills"`
	Seed       string `docopt:"--seed"`
	LogLevel   string `docopt:"--log-level"`
}

func main() {

	optFlags, _ := docopt.ParseArgs(usage, os.Args[1:], "0.0.1")
	optFlags.Bind(&Opts)

	opts := map[string]string{
		"name":      Opts.Name,
		"gender":    Opts.Gender,
		"age":       Opts.Age,
		"race":      Opts.Race,
		"careers":   Opts.Careers,
		"archetype": Opts.Archetype,
		"perks":     Opts.NPerks,
		"attrs":     Opts.NAttrDice,
		"skills":    Opts.NSkillDice,
		"seed":      Opts.Seed,
		"log-level": Opts.LogLevel,
	}

	c := m6ik.NewCharacter(opts)
	c.Print()
}
