package digitiser

import (
	"fmt"
	"math"
)

type Digitiser struct {
	digits         string
	digBase        int
	lookup         map[rune]int
	maxInt, strLen int
}

func New(digits string, length, maxRepoInt int) (Digitiser, error) {
	if len(digits) < 1 {
		return Digitiser{}, fmt.Errorf(
			"impossible configuration: digits len(%d) less than 1",
			len(digits),
		)
	}

	d := Digitiser{
		digBase: len(digits),
		digits:  digits,
		strLen:  length,
	}

	if err := d.makeLookup(); err != nil {
		return Digitiser{}, fmt.Errorf("make lookup failed: %w", err)
	}

	_ = d.countMax(length)

	if d.maxInt > maxRepoInt {
		return Digitiser{}, fmt.Errorf(
			"impossible configuration: "+
				"maximum digit(%d) exceeds maximum repository integer(%d), "+
				"should shorten length or base",
			d.maxInt,
			maxRepoInt,
		)
	}

	return d, nil
}

func (d *Digitiser) countMax(length int) (err error) {
	var maxStr string
	for i := 0; i < length; i++ {
		maxStr += string(d.digits[len(d.digits)-1])
	}

	maxValue, err := d.digitise(maxStr)
	if err != nil {
		err = fmt.Errorf("digitise failed: %w", err)
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

func (d *Digitiser) base() int {
	return d.digBase
}

func (d *Digitiser) Max() int {
	return d.maxInt
}

func (d *Digitiser) Length() int {
	return d.strLen
}

func (d *Digitiser) Digit(s string) (id int, err error) {
	if len(s) > d.Length() {
		err = fmt.Errorf("string exceeds the maximum allowed value(%v)", d.maxInt)
		return
	}

	id, err = d.digitise(s)
	if err != nil {
		err = fmt.Errorf("digitise failed: %w", err)
		return
	}

	return
}

func (d *Digitiser) digitise(s string) (digit int, err error) {
	for i, v := range s {
		m, ok := d.lookup[v]
		if !ok {
			err = fmt.Errorf("rune not found: %v", v)
			return
		}

		digit += m * int(math.Pow(float64(d.base()), float64(i)))
	}

	return
}

func (d *Digitiser) String(id int) (string, error) {
	if id > d.Max() {
		return "", fmt.Errorf("digit exceeds the maximum:(%v)", d.maxInt)
	}

	var s string

	for {
		n, err := d.lookupIndex(id % d.base())
		if err != nil {
			return "", fmt.Errorf("lookup index failed: %w", err)
		}

		s += string(n)
		id /= d.base()

		if id <= d.base()-1 {
			if id != 0 {
				n, err = d.lookupIndex(id % d.base())
				if err != nil {
					return "", fmt.Errorf("lookup index failed: %w", err)
				}
				s += string(n)
			}
			break
		}
	}

	return s, nil
}

func (d *Digitiser) lookupIndex(i int) (rune, error) {
	for k, v := range d.lookup {
		if i == v {
			return k, nil
		}
	}
	return 0, fmt.Errorf("index out of range: %v", i)
}
