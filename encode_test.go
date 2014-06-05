package logfmt

import (
	"bytes"
	"testing"
)

func TestEncodeBools(t *testing.T) {
	want := []byte("one=true two=false")
	got, err := Marshal("one", true, "two", false)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, want) {
		t.Errorf("\nwant %s\n got %s", want, got)
	}
}

func TestEncodeInts(t *testing.T) {
	want := []byte("zero=0 one=-1 two=2 three=3 four=4")
	got, err := Marshal("zero", 0, "one", -1, "two", 2, "three", 3, "four", 4)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, want) {
		t.Errorf("\nwant %s\n got %s", want, got)
	}
}

func TestEncodeFloat(t *testing.T) {
	// NOTE: There will be rounding truncation of numbers with large decimal places
	want := []byte("zero=0 one=-1.1 two=2 three=33333333333.333332 four=4 five=5.555555555555555")
	got, err := Marshal(
		"zero", 0.0, "one", -1.1, "two", 2.0000000000,
		"three", 33333333333.333333, "four", 4.0,
		"five", 5.5555555555555555555555555)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, want) {
		t.Errorf("\nwant %s\n got %s", want, got)
	}
}

func TestEncodeStrings(t *testing.T) {
	// NOTE: There will be rounding truncation of numbers with large decimal places
	want := []byte(`zero=nothing one=-1.1 three="three things" crazytext="!@#$%^&*()_+{}|:\"<>?;',./\
 some other stuff"`)
	got, err := Marshal(
		"zero", "nothing", "one", "-1.1", "three", "three things",
		"crazytext", `!@#$%^&*()_+{}|:"<>?;',./\
 some other stuff`)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, want) {
		t.Errorf("\nwant %s\n got %s", want, got)
	}
}
