package m6ik

import (
	"strconv"
)

type Die struct {
	code    int
	pips    int
	codeMax int
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
	return d.code*pipsPerDie + d.pips
}

func (d *Die) addP(e *Die) {
	d.code += e.code
	d.pips += e.pips
	d.codeMax += e.codeMax
	d.recode()
}

func (d *Die) add(e Die) {
	d.code += e.code
	d.pips += e.pips
	d.codeMax += e.codeMax
	d.recode()
}

func (d *Die) recode() {
	p := d.toPips()
	d.code = p / pipsPerDie
	d.pips = p % pipsPerDie
}

func (d *Die) toPips() int {
	return d.code*pipsPerDie + d.pips
}

func (d *Die) toDie() Die {
	return Die{code: int(d.pips / pipsPerDie), pips: d.pips % pipsPerDie}
}
