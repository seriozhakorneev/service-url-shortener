package digitiser

import (
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"testing"
)

const (
	// changing this parameters will affect tests performance
	length         = 5
	postgresMaxInt = 2147483647
	digits         = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
)

var (
	base, max int
)

func TestMain(m *testing.M) {
	d := Digitiser{digits: digits, strLen: length, digBase: len(digits)}

	err := d.makeLookup()
	if err != nil {
		log.Fatalf("Failed to make lookup in test: %s", err)
	}

	err = d.countMax(length)
	if err != nil {
		log.Fatalf("Failed to count Max in test: %s", err)
	}

	base, max = d.base(), d.Max()

	code := m.Run()
	os.Exit(code)
}

func BenchmarkName(b *testing.B) {

	//for i := 1; i <= length; i++ {

	d, err := New(digits, 6, postgresMaxInt)
	if err != nil {
		log.Fatal("new digitiser failed in testing:", err)
	}

	fmt.Println(d)

	//for i := 0; i < b.N; i++ {
	//	fmt.Println(i)
	//}
}

func FuzzDigitiserNew(f *testing.F) {
	f.Add(digits)
	f.Fuzz(func(t *testing.T, fuzzDigits string) {
		t.Parallel()
		expected := Digitiser{
			digBase: len(fuzzDigits),
			digits:  fuzzDigits,
			strLen:  length,
		}

		err := expected.makeLookup()
		if err != nil {
			return
		}

		result, err := New(fuzzDigits, length, postgresMaxInt)
		if err != nil {
			return
		}
		expected.maxInt = result.Max()
		if !reflect.DeepEqual(expected, result) {
			t.Fatalf("expected result: %v, got: %v", expected, result)
		}
	})
}

func TestDigitiserNew(t *testing.T) {
	t.Parallel()

	expError := fmt.Errorf(
		"impossible configuration: digits len(%d) less than 1",
		len(""),
	)
	_, err := New("", 0, 0)
	if !reflect.DeepEqual(err, expError) {
		t.Fatalf("Expected error: %s\nGot:%s", expError, err)
	}

	expected := Digitiser{
		digBase: base,
		digits:  digits,
		lookup:  map[rune]int{},
		maxInt:  max,
		strLen:  length,
	}

	err = expected.makeLookup()
	if err != nil {
		t.Fatal("make lookup failed in testing:", err)
	}

	result, err := New(digits, length, postgresMaxInt)
	if err != nil {
		t.Fatal("unexpected failed new in testing:", err)
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("expected result: %v, got: %v", expected, result)
	}

	duplicateDigits := "AA"
	expectedRes := Digitiser{}
	expectedErr := fmt.Errorf("make lookup failed: %w",
		fmt.Errorf("duplicate rune: %d", duplicateDigits[0]))

	result, err = New(duplicateDigits, 0, postgresMaxInt)
	if !reflect.DeepEqual(expectedErr, err) {
		t.Fatalf("expected err: %v, got: %v", expectedErr, err)
	}

	if !reflect.DeepEqual(expectedRes, result) {
		t.Fatalf("expected result: %v, got: %v", expectedRes, result)
	}

	maxRepo := 10
	expectedErr = fmt.Errorf(
		"impossible configuration: "+
			"maximum digit(%d) exceeds maximum repository integer(%d), "+
			"should shorten maxLength or base",
		max,
		maxRepo,
	)

	_, err = New(digits, length, maxRepo)
	if !reflect.DeepEqual(expectedErr, err) {
		t.Fatalf("Expected error: %s\nGot:%s", expectedErr, err)
	}
}

func TestDigitiserCountMax(t *testing.T) {
	t.Parallel()

	d := Digitiser{digits: digits, digBase: len(digits)}
	expectedErr := fmt.Errorf("digitise failed: %w",
		fmt.Errorf("rune not found: %v", digits[len(digits)-1]))

	err := d.countMax(length)
	if !reflect.DeepEqual(expectedErr, err) {
		t.Fatalf("expected err: %v\ngot: %v", expectedErr, err)
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
		t.Fatalf("expected Max int: %v, got: %v", max, d.Max())
	}
}

func TestDigitiserMakeLookup(t *testing.T) {
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

func FuzzDigitiserDigit(f *testing.F) {
	f.Add("abcd")
	f.Fuzz(func(t *testing.T, short string) {
		t.Parallel()
		digitiser, err := New(digits, length, postgresMaxInt)
		if err != nil {
			return
		}

		id, err := digitiser.Digit(short)
		if err != nil {
			return
		}
		if id > digitiser.Max() {
			t.Fatalf("expected result id <= (%d), got: %v", digitiser.Max(), id)
		}
	})
}

func TestDigitiserDigitErrors(t *testing.T) {
	t.Parallel()

	d, err := New(digits, length, postgresMaxInt)
	if err != nil {
		t.Fatal("new digitiser failed in testing:", err)
	}

	expectedErr := fmt.Errorf("string exceeds the maximum allowed value(%v)", d.maxInt)
	_, err = d.Digit("Heelloo8")

	if !reflect.DeepEqual(expectedErr, err) {
		t.Fatalf("expected err: %v, got: %v", expectedErr, err)
	}

	expectedErr = fmt.Errorf("digitise failed: %w",
		fmt.Errorf("rune not found: %v", '&'))

	_, err = d.Digit("Hell&")
	if !reflect.DeepEqual(expectedErr, err) {
		t.Fatalf("expected err: %v\ngot: %v", expectedErr, err)
	}
}

func TestDigitiserLookupIndex(t *testing.T) {
	t.Parallel()

	d := Digitiser{
		digBase: base,
		digits:  digits,
		lookup:  nil,
		maxInt:  max,
		strLen:  length,
	}

	expectedErr := fmt.Errorf("index out of range: %v", 1)

	_, err := d.lookupIndex(1)
	if !reflect.DeepEqual(expectedErr, err) {
		t.Fatalf("expected err: %v, got: %v", expectedErr, err)
	}
}

func FuzzDigitiserString(f *testing.F) {
	f.Add(224)
	f.Fuzz(func(t *testing.T, id int) {
		t.Parallel()
		digitiser, err := New(digits, length, postgresMaxInt)
		if err != nil {
			return
		}

		short, err := digitiser.String(id)
		if err != nil {
			return
		}
		if len(short) > length {
			t.Fatalf("expected short length %d, got: %d", len(short), length)
		}
		for _, r := range short {
			if v, ok := digitiser.lookup[r]; !ok {
				t.Fatalf("expected r(%d) in lookup, got: %d", r, v)
			}
		}
	})
}

func TestDigitiserString(t *testing.T) {
	t.Parallel()

	d, err := New(digits, length, postgresMaxInt)
	if err != nil {
		t.Fatal("new digitiser failed in testing:", err)
	}

	expectedErr := fmt.Errorf("digit exceeds the maximum:(%v)", d.Max())

	_, err = d.String(d.Max() + 1)
	if !reflect.DeepEqual(expectedErr, err) {
		t.Fatalf("expected err: %v, got: %v", expectedErr, err)
	}

	d.lookup = nil
	expectedErr = fmt.Errorf("lookup index failed: %w", fmt.Errorf("index out of range: %v", 0))

	_, err = d.String(0)
	if !reflect.DeepEqual(expectedErr, err) {
		t.Fatalf("expected err: %v, got: %v", expectedErr, err)
	}
}

func TestDigitiser_Results(t *testing.T) {
	t.Parallel()

	var (
		notMax int
		err    error
		d      Digitiser
	)

	for i := 1; i <= length; i++ {
		d, err = New(digits, i, postgresMaxInt)
		if err != nil {
			t.Fatal("new digitiser failed in testing:", err)
		}

		// shorten Max value to speed up test
		notMax = d.Max() / int(math.Pow(float64(i), float64(i)))

		for from, till := 0, notMax; from < till; from, till = from+1, till-1 {
			err = testID(&d, from)
			if err != nil {
				t.Fatal(err)
			}

			err = testID(&d, till)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func testID(d *Digitiser, expected int) error {
	str, err := d.String(expected)
	if err != nil {
		return fmt.Errorf("new string failed in testing: %w, id(%v)", err, expected)
	}

	result, err := d.Digit(str)
	if err != nil {
		return fmt.Errorf("new id failed in testing: %w, id(%v)", err, expected)
	}

	if expected != result {
		return fmt.Errorf("expected id: %v, got: %v, string: '%s'", expected, result, str)
	}

	return nil
}
