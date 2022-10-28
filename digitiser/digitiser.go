package digitiser

import (
	"fmt"
	"math"
)

const base64URL = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

type Digitiser struct {
	base           int
	digits         string
	lookup         map[rune]int
	maxInt, strLen int
}

func New(digits string, length int) (d Digitiser, err error) {
	d = Digitiser{
		base:   len(digits),
		digits: digits,
		strLen: length,
	}

	if err = d.makeLookup(); err != nil {
		err = fmt.Errorf("make lookup failed: %v", err)
		return
	}

	_ = d.countMax(length)
	return
}

func (d *Digitiser) countMax(length int) (err error) {
	var maxStr string
	for i := 0; i < length; i++ {
		maxStr += string(d.digits[len(d.digits)-1])
	}

	maxValue, err := d.digitise(maxStr)
	if err != nil {
		err = fmt.Errorf("digitise failed: %v", err)
		return
	}

	d.maxInt = maxValue
	return
}

func (d *Digitiser) makeLookup() (err error) {
	lookup := make(map[rune]int, 0)
	for i, r := range d.digits {
		if _, ok := lookup[r]; ok {
			err = fmt.Errorf("duplicate rune: %v", r)
			return
		}
		lookup[r] = i
	}
	d.lookup = lookup
	return
}

func (d *Digitiser) Base() int {
	return d.base
}

func (d *Digitiser) Max() int {
	return d.maxInt
}

func (d *Digitiser) Length() int {
	return d.strLen
}

func (d *Digitiser) NewID(s string) (id int, err error) {
	if len(s) > d.Length() {
		err = fmt.Errorf("string exceeds the maximum allowed value(%v)", d.maxInt)
		return
	}

	id, err = d.digitise(s)
	if err != nil {
		err = fmt.Errorf("digitise failed: %v", err)
		return
	}

	return
}

func (d *Digitiser) digitise(s string) (dig int, err error) {
	for i, v := range s {
		m, ok := d.lookup[v]
		if !ok {
			err = fmt.Errorf("rune not found: %v", v)
			return
		}
		dig += m * int(math.Pow(float64(d.Base()), float64(i)))
	}

	return
}

func (d *Digitiser) NewString(id int) (s string, err error) {
	if id > d.Max() {
		err = fmt.Errorf("digit exceeds the maximum:(%v)", d.maxInt)
		return
	}

	var n rune
	for {
		n, err = d.LookupIndex(id % d.Base())
		if err != nil {
			err = fmt.Errorf("lookup index failed: %v", err)
			return
		}

		s += string(n)

		id = id / d.Base()
		if id <= d.Base()-1 {
			if id != 0 {
				n, err = d.LookupIndex(id % d.Base())
				if err != nil {
					err = fmt.Errorf("lookup index failed: %v", err)
					return
				}

				s += string(n)
			}
			break
		}
	}

	return
}

func (d *Digitiser) LookupIndex(i int) (rune, error) {
	for k, v := range d.lookup {
		if i == v {
			return k, nil
		}
	}

	return 0, fmt.Errorf("index out of range: %v", i)
}