package digitiser

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDigitiser_New(t *testing.T) {
	t.Parallel()
	length, digits := 5, base64URL
	expected := Digitiser{
		base:   64,
		digits: digits,
		lookup: map[rune]int{},
		maxInt: 1073741823,
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

	length, digits, expected := 5, base64URL, 1073741823

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

	if d.Max() != expected {
		t.Fatalf("expected max int: %v, got: %v", expected, d.Max())
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

	digits, length := base64URL, 5
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

func testIDInRange(d *Digitiser, from, till int) error {
	for i, j := from, till; i < j; i, j = i+1, j-1 {
		err := testID(d, i)
		if err != nil {
			return err
		}

		err = testID(d, j)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestDigitiser_Results(t *testing.T) {
	t.Parallel()

	digits, length := base64URL, 4

	d, err := New(digits, length)
	if err != nil {
		t.Fatal("new digitiser failed in testing:", err)
	}

	if length < 3 {
		err = testIDInRange(&d, 0, d.Max())
		if err != nil {
			t.Fatal(err)
		}
	}

	// TODO начиная с length=3 d.Max()/100000 будет давать результат > 0
	//fmt.Println("max:", d.Max(), "sotka counter:", d.Max()/100000)

	for i := 0; i <= d.Max(); {

		// TODO: WAITGROUP n Err as var
		if (i + 100000) < d.Max() {

			err = testIDInRange(&d, i, i+100000)
			if err != nil {
				t.Fatal(err)
			}
		} else {

			err = testIDInRange(&d, i, d.Max())
			if err != nil {
				t.Fatal(err)
			}

		}
		i += 100000
	}
}

//func TestDigitiserCountMax(t *testing.T) {
//	t.Parallel()
//
//	d := Digitiser{}
//	err := d.countMax()
//	if err != nil {
//		return
//	}
//}

//func TestNewID(t *testing.T) {
//	err := Init(base64URL, 5)
//	if err != nil {
//		t.Fatalf("Failed to init Digit")
//	}
//	var result int
//	testStr := "s"
//	expectedIndex := strings.Index(base64URL, testStr)
//	expectedUniqueIds := []string{"0", "c", "Z", "01", "b1"}
//	testCounters := []int{0, 12, 61, 62, 73}
//
//	for i, v := range []rune(testStr) {
//		index, _ := Digits.LookupRune(v)
//		if index != expectedIndex {
//			t.Fatalf("Expected index: %v. Got: %v", expectedIndex, index)
//		}
//		result += index * int(math.Pow(62, float64(i)))
//	}
//
//	for i, v := range expectedUniqueIds {
//		id, _ := Digits.NewID(v)
//		if id != testCounters[i] {
//			t.Fatalf("Expected ID: %v. Got: %v", testCounters[i], id)
//		}
//	}
//}
//
//func Test_NewString(t *testing.T) {
//	err := Init(base64URL, 5)
//	if err != nil {
//		t.Fatalf("Failed to init Digit")
//	}
//
//	testCounters := []int{0, 12, 61, 62, 73}
//	expectedUniqueIds := []string{"0", "c", "Z", "01", "b1"}
//
//	for i, v := range testCounters {
//		x, _ := Digits.NewString(v)
//		if x != expectedUniqueIds[i] {
//			t.Fatalf("Expected id: %v. Got: %v", expectedUniqueIds[i], x)
//		}
//	}
//}
