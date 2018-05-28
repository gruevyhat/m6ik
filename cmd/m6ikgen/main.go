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
  --n_perks	Number of random perks to assign.
  --seed	Character generation signature.
  -h --help
  --version
`

var Opts struct {
	Name      string `docopt:"--name"`
	Gender    string `docopt:"--gender"`
	Age       string `docopt:"--age"`
	Race      string `docopt:"--race"`
	Careers   string `docopt:"--careers"`
	Archetype string `docopt:"--archetype"`
	NPerks    string `docopt:"--n_perks"`
	Seed      string `docopt:"--seed"`
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
		"n_perks":   Opts.NPerks,
		"seed":      Opts.Seed,
	}

	c := m6ik.NewCharacter(opts)
	c.Print()
}
