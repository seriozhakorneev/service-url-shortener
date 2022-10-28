package digitiser

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

const (
	length = 5
	digits = base64URL
	base   = 64
	max    = 1073741823
)

func TestDigitiser_New(t *testing.T) {
	t.Parallel()
	expected := Digitiser{
		base:   base,
		digits: digits,
		lookup: map[rune]int{},
		maxInt: max,
		strLen: length,
	}

	err := expected.makeLookup()
	if err != nil {
		t.Fatal("make lookup failed in testing:", err)
	}

	result, err := New(digits, length)
	if err != nil {
		t.Fatal("new failed in testing:", err)
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("expected result: %v, got: %v", expected, result)
	}

	duplicateDigits := "AA"
	expectedErr := fmt.Errorf("make lookup failed: duplicate rune: %d", duplicateDigits[0])
	expectedRes := Digitiser{base: 2, digits: duplicateDigits}
	result, err = New(duplicateDigits, 0)

	if !reflect.DeepEqual(expectedErr, err) {
		t.Fatalf("expected err: %v, got: %v", expectedErr, err)
	}
	if !reflect.DeepEqual(expectedRes, result) {
		t.Fatalf("expected err: %v, got: %v", expectedRes, result)
	}
}

func TestDigitiser_countMax(t *testing.T) {
	t.Parallel()

	d := Digitiser{digits: digits, base: len(digits)}
	expectedErr := fmt.Errorf("digitise failed: rune not found: %v", digits[len(digits)-1])

	err := d.countMax(length)
	if !reflect.DeepEqual(expectedErr, err) {
		t.Fatalf("expected err: %v, got: %v", expectedErr, err)
	}

	err = d.makeLookup()
	if err != nil {
		t.Fatal("make lookup error in testing:", err)
	}

	err = d.countMax(length)
	if err != nil {
		t.Fatal("expected nil error, got:", err)
	}

	if d.Max() != max {
		t.Fatalf("expected max int: %v, got: %v", max, d.Max())
	}
}

func TestDigitiser_makeLookup(t *testing.T) {
	t.Parallel()

	expectedLookup := map[rune]int{'A': 0, 'B': 1}
	d := Digitiser{digits: "AB"}

	err := d.makeLookup()
	if err != nil {
		t.Fatal("expected nil error, got:", err)
	}

	if !reflect.DeepEqual(expectedLookup, d.lookup) {
		t.Fatalf("expected lookup: %v, got: %v", expectedLookup, d.lookup)
	}

	d = Digitiser{digits: "AA"}
	expectedErr := fmt.Errorf("duplicate rune: %v", d.digits[0])

	err = d.makeLookup()
	if !reflect.DeepEqual(expectedErr, err) {
		t.Fatalf("expected error: %v, got: %v", expectedErr, err)
	}

}

func TestDigitiser_NewID_Errors(t *testing.T) {
	t.Parallel()

	d, err := New(digits, length)
	if err != nil {
		t.Fatal("new digitiser failed in testing:", err)
	}

	expectedErr := fmt.Errorf("string exceeds the maximum allowed value(%v)", d.maxInt)
	_, err = d.NewID("Helloo")
	if !reflect.DeepEqual(expectedErr, err) {
		t.Fatalf("expected err: %v, got: %v", expectedErr, err)
	}

	expectedErr = fmt.Errorf("digitise failed: rune not found: %v", '&')
	_, err = d.NewID("Hell&")
	if !reflect.DeepEqual(expectedErr, err) {
		t.Fatalf("expected err: %v, got: %v", expectedErr, err)
	}
}

func TestDigitiser_LookupIndex(t *testing.T) {
	t.Parallel()

	d := Digitiser{
		base:   base,
		digits: digits,
		lookup: nil,
		maxInt: max,
		strLen: length,
	}

	expectedErr := fmt.Errorf("index out of range: %v", 1)
	_, err := d.lookupIndex(1)
	if !reflect.DeepEqual(expectedErr, err) {
		t.Fatalf("expected err: %v, got: %v", expectedErr, err)
	}
}

func TestDigitiser_NewString(t *testing.T) {
	t.Parallel()

	d, err := New(digits, length)
	if err != nil {
		t.Fatal("new digitiser failed in testing:", err)
	}

	expectedErr := fmt.Errorf("digit exceeds the maximum:(%v)", d.Max())
	_, err = d.NewString(d.Max() + 1)
	if !reflect.DeepEqual(expectedErr, err) {
		t.Fatalf("expected err: %v, got: %v", expectedErr, err)
	}

	d.lookup = nil
	expectedErr = fmt.Errorf("lookup index failed: %v", fmt.Errorf("index out of range: %v", 0))
	_, err = d.NewString(0)
	if !reflect.DeepEqual(expectedErr, err) {
		t.Fatalf("expected err: %v, got: %v", expectedErr, err)
	}
}

func TestDigitiser_ResultsTillLength3(t *testing.T) {
	t.Parallel()

	for i := 0; i <= 3; i++ {
		d, err := New(digits, i)
		if err != nil {
			t.Fatal("new digitiser failed in testing:", err)
		}
		testIDsInRange(t, &d, 0, d.Max())
	}
}

func TestDigitiser_ResultsLength4(t *testing.T) {
	t.Parallel()

	hundredT := 100000
	d, err := New(digits, 4)
	if err != nil {
		t.Fatal("new digitiser failed in testing:", err)
	}

	for i, j := 0, d.Max(); i < j; i, j = i+hundredT, j-hundredT {
		testIDsInRange(t, &d, i, i+hundredT)
		testIDsInRange(t, &d, j-hundredT, j)

		if (j-hundredT)-(i+hundredT) < hundredT {
			testIDsInRange(t, &d, i+hundredT, j-hundredT)
			break
		}
	}
}

func testID(d *Digitiser, expected int) error {
	str, err := d.NewString(expected)
	if err != nil {
		return fmt.Errorf("new string failed in testing: %v, id(%v)", err, expected)
	}

	result, err := d.NewID(str)
	if err != nil {
		return fmt.Errorf("new id failed in testing: %v, id(%v)", err, expected)
	}

	if expected != result {
		return fmt.Errorf("expected id: %v, got: %v, string: '%s'", expected, result, str)
	}
	return nil
}

func testIDsInRange(t *testing.T, d *Digitiser, from, till int) {
	var (
		err error
		wg  sync.WaitGroup
	)

	for i, j := from, till; i < j; i, j = i+1, j-1 {
		wg.Add(2)

		go func(a int) {
			defer wg.Done()
			err = testID(d, a)
		}(i)

		go func(a int) {
			defer wg.Done()
			err = testID(d, a)
		}(j)
	}

	wg.Wait()
	if err != nil {
		t.Fatal(err)
	}
}
