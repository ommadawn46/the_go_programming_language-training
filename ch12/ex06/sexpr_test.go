package sexpr

import (
	"reflect"
	"testing"
)

func TestBool(t *testing.T) {
	for _, test := range []struct {
		input bool
		want  string
	}{
		{true, "t"},
		{false, "nil"},
	} {
		got, err := Marshal(test.input)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}
		if string(got) != test.want {
			t.Errorf("\ngot: \n%s\nwant: \n%s\n", got, test.want)
		}
	}
}

func TestFloat(t *testing.T) {
	for _, test := range []struct {
		input float64
		want  string
	}{
		{0, "0"},
		{-2.71828182846, "-2.71828182846"},
		{3.14159265359, "3.14159265359"},
	} {
		got, err := Marshal(test.input)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}
		if string(got) != test.want {
			t.Errorf("\ngot: \n%s\nwant: \n%s\n", got, test.want)
		}
	}
}

func TestComplex(t *testing.T) {
	for _, test := range []struct {
		input complex128
		want  string
	}{
		{complex(0, 0), "#C(0 0)"},
		{complex(-1, -2), "#C(-1 -2)"},
		{complex(3.14159265359, 2.71828182846), "#C(3.14159265359 2.71828182846)"},
	} {
		got, err := Marshal(test.input)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}
		if string(got) != test.want {
			t.Errorf("\ngot: \n%s\nwant: \n%s\n", got, test.want)
		}
	}
}

func TestInterface(t *testing.T) {
	type TestStruct struct {
		StringField string
		IntField    int
		SliceField  []float64
	}
	for _, test := range []struct {
		input interface{}
		want  string
	}{
		{nil, "nil"},
		{"ABCDEFG", `("string" "ABCDEFG")`},
		{
			TestStruct{"ABC", 123, []float64{3.14, 2.72}},
			`("sexpr.TestStruct" ((StringField "ABC")
                     (IntField 123)
                     (SliceField (3.14
                                  2.72))))`,
		},
	} {
		got, err := Marshal(&test.input)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}
		if string(got) != test.want {
			t.Errorf("\ngot: \n%s\nwant: \n%s\n", got, test.want)
		}
	}
}

func TestPrettyPrint(t *testing.T) {
	type TestStruct struct {
		ComplexField   complex128
		SliceField     []string
		InterfaceField interface{}
	}
	for _, test := range []struct {
		input interface{}
		want  string
	}{
		{
			TestStruct{complex(0, 0), nil, nil},
			`("sexpr.TestStruct" ())`,
		},
		{
			TestStruct{
				complex(3.14, 2.72),
				[]string{"ABC", "DEF", "GHI"},
				&TestStruct{complex(0, 0), nil, nil},
			},
			`("sexpr.TestStruct" ((ComplexField #C(3.14 2.72))
                     (SliceField ("ABC"
                                  "DEF"
                                  "GHI"))
                     (InterfaceField ("*sexpr.TestStruct" ()))))`,
		},
	} {
		got, err := Marshal(&test.input)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}
		if string(got) != test.want {
			t.Errorf("\ngot: \n%s\nwant: \n%s\n", got, test.want)
		}
	}
}

func TestZeroValue(t *testing.T) {
	type TestStruct struct {
		IntField       int
		ComplexField   complex128
		SliceField     []string
		InterfaceField interface{}
	}
	for _, test := range []struct {
		input TestStruct
		want  string
	}{
		{TestStruct{0, complex(0, 0), nil, nil}, "()"},
		{TestStruct{1, complex(0, 0), nil, nil}, "((IntField 1))"},
		{TestStruct{1, complex(1, 1), nil, nil}, `((IntField 1)
 (ComplexField #C(1 1)))`},
		{TestStruct{1, complex(1, 1), []string{}, nil}, `((IntField 1)
 (ComplexField #C(1 1))
 (SliceField ()))`},
		{
			TestStruct{
				1, complex(1, 1), []string{},
				&TestStruct{0, complex(0, 0), nil, nil},
			},
			`((IntField 1)
 (ComplexField #C(1 1))
 (SliceField ())
 (InterfaceField ("*sexpr.TestStruct" ())))`,
		},
	} {
		got, err := Marshal(test.input)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}
		if string(got) != test.want {
			t.Errorf("\ngot: \n%s\nwant: \n%s\n", got, test.want)
		}
	}
}

// Test verifies that encoding and decoding a complex data value
// produces an equal result.
//
// The test does not make direct assertions about the encoded output
// because the output depends on map iteration order, which is
// nondeterministic.  The output of the t.Log statements can be
// inspected by running the test with the -v flag:
//
// 	$ go test -v gopl.io/ch12/sexpr
//
func Test(t *testing.T) {
	type Movie struct {
		Title, Subtitle string
		Year            int
		Actor           map[string]string
		Oscars          []string
		Sequel          *string
	}
	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Actor: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
			"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
			`Maj. T.J. "King" Kong`:      "Slim Pickens",
		},
		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}

	// Encode it
	data, err := Marshal(strangelove)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	t.Logf("Marshal() = \n%s\n", data)

	// Decode it
	var movie Movie
	if err := Unmarshal(data, &movie); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	t.Logf("Unmarshal() = %+v\n", movie)

	// Check equality.
	if !reflect.DeepEqual(movie, strangelove) {
		t.Fatal("not equal")
	}
}
