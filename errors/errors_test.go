package errors

import (
	"errors"
	"testing"
)

var (
	errTest = errors.New("test")
)

func testDrop1() (int, error) {
	return 0, errTest
}

func testDrop2() (float32, string, error) {
	return 0.0, "", errTest
}

func testMustV() (int, error) {
	return 42, nil
}

func TestDrop1(t *testing.T) {
	if err := Drop(testDrop1()); err != errTest {
		t.Error("Failed")
	}
}

func TestDrop12(t *testing.T) {
	if err := Drop2(testDrop2()); err != errTest {
		t.Error("Failed")
	}
}

func TestMustV(t *testing.T) {
	val := MustV(testMustV())
	if val != 42 {
		t.Error("Failed")
	}
}
