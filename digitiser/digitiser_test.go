package digitiser

import (
	"fmt"
	"reflect"
	"sync"
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

func testIDInRange(t *testing.T, d *Digitiser, from, till int) {
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

func TestDigitiser_Results(t *testing.T) {
	t.Parallel()

	digits, length := base64URL, 3
	sotka := 100000

	d, err := New(digits, length)
	if err != nil {
		t.Fatal("new digitiser failed in testing:", err)
	}

	if length < 3 {
		testIDInRange(t, &d, 0, d.Max())
		//if err != nil {
		//	t.Fatal(err)
		//}
	}

	//var testErr error
	//wg := sync.WaitGroup{}
	fmt.Println("max", d.Max())

	for i, j := 0, d.Max(); i < j; {

		fmt.Println(i, i+sotka, j-sotka, j)

		if (j-sotka)-(i+sotka) < sotka {
			fmt.Println(i+sotka, j-sotka)
			break
		}

		//if (i + 100000) < d.Max() {
		//	wg.Add(1)
		//	go func() {
		//		defer wg.Done()
		//		fmt.Println(i, i+100000)
		//		testIDInRange(t, &d, i, i+100000)
		//	}()
		//
		//} else {
		//	wg.Add(1)
		//	go func() {
		//		defer wg.Done()
		//		fmt.Println(i, d.Max())
		//
		//		testIDInRange(t, &d, i, d.Max())
		//	}()
		//}
		//wg.Wait()
		//
		i += 100000
		j -= 100000
	}
	//if testErr != nil {
	//	t.Fatal(testErr)
	//}
}
