package logfmt

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

// Starts with unmarshalling first
func testRoundtripU(t *testing.T, start *[]byte, want *[]byte) {

	c := new(coll)
	if err := Unmarshal(*start, c); err != nil {
		t.Fatal(err)
	}

	r := []interface{}{}
	for _, p := range c.a {
		r = append(r, p.k, p.v)
	}

	got, err := Marshal(r...)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, *want) {
		t.Errorf("\nwant %s\n got %s", *want, got)
	}

}

// Starts with marshalling first
func testRoundtripM(t *testing.T, start *[]interface{}, want *[]interface{}) {

	s, err := Marshal(*start...)
	if err != nil {
		t.Fatal(err)
	}

	c := new(coll)
	if err := Unmarshal(s, c); err != nil {
		t.Fatal(err)
	}

	got := []interface{}{}
	for _, p := range c.a {
		got = append(got, p.k, p.v)
	}

	if !reflect.DeepEqual(got, *want) {
		t.Errorf("\nwant %v\n got %v", *want, got)
	}

}

func TestRoundtripBoolU(t *testing.T) {
	start := []byte(`d foo= emp=`) // '=' will always be added in
	want := []byte(`d= foo= emp=`)

	testRoundtripU(t, &start, &want)

}

// Booleans get turned into "true" and "false"
func TestRoundtripBoolM(t *testing.T) {
	start := []interface{}{"first", true, "second", false}
	want := []interface{}{"first", "true", "second", "false"}

	testRoundtripM(t, &start, &want)
}

func TestRoundtripIntU(t *testing.T) {
	start := []byte(`a=1 b=2 c=3`)
	want := []byte(`a=1 b=2 c=3`)

	testRoundtripU(t, &start, &want)

}

func TestRoundtripIntM(t *testing.T) {

	start := []interface{}{"one", -1, "two", 2, "three", 3333333333333333}
	want := []interface{}{"one", "-1", "two", "2", "three", "3333333333333333"}

	testRoundtripM(t, &start, &want)

}

func TestRoundtripFloatU(t *testing.T) {
	start := []byte(`a=1.0 b=77.777777777777 c=555555555.55555555`)

	testRoundtripU(t, &start, &start)

}

// There can be rounding on occasion and elimination of unnecessary characters
func TestRoundtripFloatM(t *testing.T) {

	start := []interface{}{"a", -1.0, "b", 77.777777777777, "three", 555555555.5555555555}
	want := []interface{}{"a", "-1", "b", "77.777777777777", "three", "555555555.5555556"} // rounding

	testRoundtripM(t, &start, &want)

}

func TestRoundtripStringU(t *testing.T) {
	start := []byte(`a=coolness b="all your base" c="are
belong to us"`)

	fmt.Println("zero")
	fmt.Println(string(start))
	testRoundtripU(t, &start, &start)

	fmt.Println("one")
	start = []byte("newlines=\"several\nlines\nof\ncode\"")
	fmt.Println(string(start))
	testRoundtripU(t, &start, &start)

	fmt.Println("two")
	start = []byte(`newlines="several\nlines\nof\ncode"`)
	fmt.Println(string(start))
	testRoundtripU(t, &start, &start)

}

// There can be rounding on occasion and elimination of unnecessary characters
func TestRoundtripStringM(t *testing.T) {

	start := []interface{}{"a", "several\nlines\nof stuff", "b", "日本国", "c", "汉语", "last", "!@#$%^&*()_+{}|;':<>?,./\\`~"}

	testRoundtripM(t, &start, &start)

}
