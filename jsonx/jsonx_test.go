package jsonx_test

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/docodex/gopkg/jsonx"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

type applyParam struct {
	Type int8   `json:"type"`
	Key  string `json:"key"`
	// From go1.24, when marshaling, a struct field with the new omitzero
	// option in the struct field tag will be omitted if its value is zero.
	// If the field type has an IsZero() bool method, that will be used to
	// determine whether the value is zero. Otherwise, the value is zero if
	// it is the zero value for its type. The omitzero field tag is clearer
	// and less error-prone than omitempty when the intent is to omit zero
	// values. In particular, unlike omitempty, omitzero omits zero-valued
	// time.Time values, which is a common source of friction.
	//
	// If both omitempty and omitzero are specified, the field will be omitted
	// if the value is either empty or zero (or both).
	T1 time.Time `json:"t1,omitempty"`
	T2 time.Time `json:"t2,omitzero"`
	T3 time.Time `json:"t3,omitempty,omitzero"`
}

func TestJSON(t *testing.T) {
	p1 := &applyParam{
		Type: 2,
		Key:  "test_key",
	}
	v, err := jsonx.MarshalToString(p1)
	assert.Nil(t, err)
	fmt.Println("p1:", p1.T1, p1.T2, p1.T3)
	fmt.Println(v)
	var p2 applyParam
	err = jsonx.UnmarshalFromString(v, &p2)
	assert.Nil(t, err)
	fmt.Println("p2:", p2.Type, p2.Key, p2.T1, p2.T2, p2.T3)
}

func TestUnmarshal1(t *testing.T) {
	a := []int64{1, 2, 3}
	v, err := jsonx.MarshalToString(a)
	assert.Nil(t, err)
	fmt.Println("a:", v)
	var b []int64
	err = jsonx.UnmarshalFromString(v, b)
	assert.NotNil(t, err)
	fmt.Println("b:", b)
	var c []int64
	err = jsonx.UnmarshalFromString(v, &c)
	assert.Nil(t, err)
	fmt.Println("c:", c)
}

func TestUnmarshal2(t *testing.T) {
	a := map[int64]string{
		1: "a",
		2: "b",
		3: "c",
	}
	v, err := jsonx.MarshalToString(a)
	assert.Nil(t, err)
	fmt.Println("a:", v)
	var b map[int64]string
	err = jsonx.UnmarshalFromString(v, b)
	assert.NotNil(t, err)
	fmt.Println("b:", b)
	var c map[int64]string
	err = jsonx.UnmarshalFromString(v, &c)
	assert.Nil(t, err)
	fmt.Println("c:", c)
}

func TestUnmarshal3(t *testing.T) {
	type T struct {
		A int64
		B int64
		C int64
	}
	a := &T{
		A: 1,
		B: 2,
		C: 3,
	}
	v, err := jsonx.MarshalToString(a)
	assert.Nil(t, err)
	fmt.Println("a:", v)
	var b *T
	err = jsonx.UnmarshalFromString(v, b)
	assert.NotNil(t, err)
	fmt.Println("b:", b)
	var c *T
	err = jsonx.UnmarshalFromString(v, &c)
	assert.Nil(t, err)
	fmt.Println("c:", c)
	var d T
	err = jsonx.UnmarshalFromString(v, d)
	assert.NotNil(t, err)
	fmt.Println("d:", d)
	var e T
	err = jsonx.UnmarshalFromString(v, &e)
	assert.Nil(t, err)
	fmt.Println("e:", e)
}

func TestBatchVariousPathCounts(t *testing.T) {
	text := `{"a":"a","b":"b","c":"c"}`
	counts := []int{
		3, 4, 7, 8, 9, 15, 16, 17, 31, 32, 33, 63, 64, 65, 127,
		128, 129, 255, 256, 257, 511, 512, 513,
	}
	paths := []string{"a", "b", "c"}
	expects := []string{"a", "b", "c"}
	for _, count := range counts {
		var gpaths []string
		for i := range count {
			if i < len(paths) {
				gpaths = append(gpaths, paths[i])
			} else {
				gpaths = append(gpaths, fmt.Sprintf("not%d", i))
			}
		}
		results := jsonx.MGet(text, gpaths...)
		for i := range paths {
			if results[i].String() != expects[i] {
				t.Fatalf("expected '%v', got '%v'", expects[i],
					results[i].String())
			}
		}
	}
}

func TestBatchRecursion(t *testing.T) {
	var text string
	var path string
	for range 100 {
		text += `{"a":`
		path += ".a"
	}
	text += `"b"`
	for range 100 {
		text += `}`
	}
	path = path[1:]
	assert.True(t, jsonx.MGet(text, path)[0].String() == "b")
}

var manyJSON = `  {
	"a":{"a":{"a":{"a":{"a":{"a":{"a":{"a":{"a":{"a":{
	"a":{"a":{"a":{"a":{"a":{"a":{"a":{"a":{"a":{"a":{
	"a":{"a":{"a":{"a":{"a":{"a":{"a":{"a":{"a":{"a":{
	"a":{"a":{"a":{"a":{"a":{"a":{"a":{"a":{"a":{"a":{
	"a":{"a":{"a":{"a":{"a":{"a":{"a":{"a":{"a":{"a":{
	"a":{"a":{"a":{"a":{"a":{"a":{"a":{"a":{"a":{"a":{
	"a":{"a":{"a":{"a":{"a":{"a":{"a":{"a":{"a":{"a":{"hello":"world"
	}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}
	"position":{"type":"Point","coordinates":[-115.24,33.09]},
	"loves":["world peace"],
	"name":{"last":"Anderson","first":"Nancy"},
	"age":31
	"":{"a":"emptya","b":"emptyb"},
	"name.last":"Yellow",
	"name.first":"Cat",
}`

func TestManyBasic(t *testing.T) {
	testMany := func(expect string, paths ...string) {
		results := jsonx.MGetBytes(
			[]byte(manyJSON),
			paths...,
		)
		if len(results) != len(paths) {
			t.Fatalf("expected %v, got %v", len(paths), len(results))
		}
		if fmt.Sprintf("%v", results) != expect {
			fmt.Printf("%v\n", paths)
			t.Fatalf("expected %v, got %v", expect, results)
		}
	}
	testMany("[Point]", "position.type")
	testMany(`[emptya ["world peace"] 31]`, ".a", "loves", "age")
	testMany(`[["world peace"]]`, "loves")
	testMany(`[{"last":"Anderson","first":"Nancy"} Nancy]`, "name",
		"name.first")
	testMany(`[]`, strings.Repeat("a.", 40)+"hello")
	res := gjson.Get(manyJSON, strings.Repeat("a.", 48)+"a")
	testMany(`[`+res.String()+`]`, strings.Repeat("a.", 48)+"a")
	// these should fallback
	testMany(`[Cat Nancy]`, "name\\.first", "name.first")
	testMany(`[world]`, strings.Repeat("a.", 70)+"hello")
}

func testManyAny(t *testing.T, json string, paths, expected []string, bytes bool) {
	var result []gjson.Result
	for i := range 2 {
		var which string
		if i == 0 {
			which = "Get"
			result = nil
			for j := 0; j < len(expected); j++ {
				if bytes {
					result = append(result, gjson.GetBytes([]byte(json), paths[j]))
				} else {
					result = append(result, gjson.Get(json, paths[j]))
				}
			}
		} else if i == 1 {
			which = "GetMany"
			if bytes {
				result = jsonx.MGetBytes([]byte(json), paths...)
			} else {
				result = jsonx.MGet(json, paths...)
			}
		}
		for j := range expected {
			if result[j].String() != expected[j] {
				t.Fatalf("Using key '%s' for '%s'\nexpected '%v', got '%v'",
					paths[j], which, expected[j], result[j].String())
			}
		}
	}
}

func TestRandomMany(t *testing.T) {
	var lstr string
	defer func() {
		if v := recover(); v != nil {
			println("'" + hex.EncodeToString([]byte(lstr)) + "'")
			println("'" + lstr + "'")
			panic(v)
		}
	}()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 512)
	for range 50000 {
		n, err := r.Read(b[:rand.Int()%len(b)])
		if err != nil {
			t.Fatal(err)
		}
		lstr = string(b[:n])
		paths := make([]string, rand.Int()%64)
		for i := range paths {
			var b []byte
			n := rand.Int() % 5
			for j := range n {
				if j > 0 {
					b = append(b, '.')
				}
				nn := rand.Int() % 10
				for range nn {
					b = append(b, 'a'+byte(rand.Int()%26))
				}
			}
			paths[i] = string(b)
		}
		jsonx.MGet(lstr, paths...)
	}
}

func TestGetMany47(t *testing.T) {
	text := `{"bar": {"id": 99, "mybar": "my mybar" }, "foo": ` +
		`{"myfoo": [605]}}`
	paths := []string{"foo.myfoo", "bar.id", "bar.mybar", "bar.mybarx"}
	expected := []string{"[605]", "99", "my mybar", ""}
	results := jsonx.MGet(text, paths...)
	if len(expected) != len(results) {
		t.Fatalf("expected %v, got %v", len(expected), len(results))
	}
	for i, path := range paths {
		if results[i].String() != expected[i] {
			t.Fatalf("expected '%v', got '%v' for path '%v'", expected[i],
				results[i].String(), path)
		}
	}
}

func TestGetMany48(t *testing.T) {
	text := `{"bar": {"id": 99, "xyz": "my xyz"}, "foo": {"myfoo": [605]}}`
	paths := []string{"foo.myfoo", "bar.id", "bar.xyz", "bar.abc"}
	expected := []string{"[605]", "99", "my xyz", ""}
	results := jsonx.MGet(text, paths...)
	if len(expected) != len(results) {
		t.Fatalf("expected %v, got %v", len(expected), len(results))
	}
	for i, path := range paths {
		if results[i].String() != expected[i] {
			t.Fatalf("expected '%v', got '%v' for path '%v'", expected[i],
				results[i].String(), path)
		}
	}
}

func TestIssue54(t *testing.T) {
	var r []gjson.Result
	text := `{"MarketName":null,"Nounce":6115}`
	r = jsonx.MGet(text, "Nounce", "Buys", "Sells", "Fills")
	if strings.Replace(fmt.Sprintf("%v", r), " ", "", -1) != "[6115]" {
		t.Fatalf("expected '%v', got '%v'", "[6115]",
			strings.Replace(fmt.Sprintf("%v", r), " ", "", -1))
	}
	r = jsonx.MGet(text, "Nounce", "Buys", "Sells")
	if strings.Replace(fmt.Sprintf("%v", r), " ", "", -1) != "[6115]" {
		t.Fatalf("expected '%v', got '%v'", "[6115]",
			strings.Replace(fmt.Sprintf("%v", r), " ", "", -1))
	}
	r = jsonx.MGet(text, "Nounce")
	if strings.Replace(fmt.Sprintf("%v", r), " ", "", -1) != "[6115]" {
		t.Fatalf("expected '%v', got '%v'", "[6115]",
			strings.Replace(fmt.Sprintf("%v", r), " ", "", -1))
	}
}

func TestIssue55(t *testing.T) {
	text := `{"one": {"two": 2, "three": 3}, "four": 4, "five": 5}`
	results := jsonx.MGet(text, "four", "five", "one.two", "one.six")
	expected := []string{"4", "5", "2", ""}
	for i, r := range results {
		if r.String() != expected[i] {
			t.Fatalf("expected %v, got %v", expected[i], r.String())
		}
	}
}

func TestSet(t *testing.T) {
	got, err := jsonx.Set(`{"name":{"first":"Tom"}}`, "name.last", "Anderson")
	assert.NoError(t, err)
	assert.Equal(t, `{"name":{"first":"Tom","last":"Anderson"}}`, got)

	// invalid path returns error
	_, err = jsonx.Set(`{}`, "", "x")
	assert.Error(t, err)
}

func TestSetBytes(t *testing.T) {
	got, err := jsonx.SetBytes([]byte(`{"age":36}`), "age", 37)
	assert.NoError(t, err)
	assert.Equal(t, `{"age":37}`, string(got))
}

func TestSetRaw(t *testing.T) {
	got, err := jsonx.SetRaw(`{}`, "children", `["Sara","Alex","Jack"]`)
	assert.NoError(t, err)
	assert.Equal(t, `{"children":["Sara","Alex","Jack"]}`, got)
}

func TestSetRawBytes(t *testing.T) {
	got, err := jsonx.SetRawBytes([]byte(`{}`), "x", []byte(`{"a":1}`))
	assert.NoError(t, err)
	assert.Equal(t, `{"x":{"a":1}}`, string(got))
}

func TestSetWithBeOptimistic(t *testing.T) {
	// BeOptimistic hits that the path exists; should produce same result as plain Set.
	got, err := jsonx.Set(`{"a":1}`, "a", 2, jsonx.BeOptimistic())
	assert.NoError(t, err)
	assert.Equal(t, `{"a":2}`, got)
}

func TestSetBytesWithDoReplaceInPlace(t *testing.T) {
	// With DoReplaceInPlace, sjson may modify the underlying buffer; we only
	// assert the returned byte slice has the expected value.
	in := []byte(`{"a":1}`)
	got, err := jsonx.SetBytes(in, "a", 2, jsonx.DoReplaceInPlace())
	assert.NoError(t, err)
	assert.Equal(t, `{"a":2}`, string(got))
}
