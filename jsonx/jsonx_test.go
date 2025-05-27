package jsonx_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/bytedance/gopkg/lang/fastrand"
	"github.com/docodex/gopkg/jsonx"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/tidwall/sjson"
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

// TestRandomData11 is a fuzzing test that throws random data at the Parse
// function looking for panics.
func TestRandomData(t *testing.T) {
	var lstr string
	defer func() {
		if v := recover(); v != nil {
			println("'" + hex.EncodeToString([]byte(lstr)) + "'")
			println("'" + lstr + "'")
			panic(v)
		}
	}()
	b := make([]byte, 200)
	for range 2000000 {
		n, err := fastrand.Read(b[:fastrand.Int()%len(b)])
		if err != nil {
			t.Fatal(err)
		}
		lstr = string(b[:n])
		gjson.GetBytes([]byte(lstr), "zzzz")
		gjson.Parse(lstr)
	}
}

func TestRandomValidStrings(t *testing.T) {
	b := make([]byte, 200)
	for range 100000 {
		n, err := fastrand.Read(b[:fastrand.Int()%len(b)])
		if err != nil {
			t.Fatal(err)
		}
		sm, err := json.Marshal(string(b[:n]))
		if err != nil {
			t.Fatal(err)
		}
		var su string
		if err := json.Unmarshal(sm, &su); err != nil {
			t.Fatal(err)
		}
		token := gjson.Get(`{"str":`+string(sm)+`}`, "str")
		if token.Type != gjson.String || token.Str != su {
			println("["+token.Raw+"]", "["+token.Str+"]", "["+su+"]",
				"["+string(sm)+"]")
			t.Fatal("string mismatch")
		}
	}
}

func TestEmoji(t *testing.T) {
	const input = `{"utf8":"Example emoji, KO: \ud83d\udd13, \ud83c\udfc3 ` +
		`OK: \u2764\ufe0f "}`
	value := gjson.Get(input, "utf8")
	var s string
	_ = json.Unmarshal([]byte(value.Raw), &s)
	if value.String() != s {
		t.Fatalf("expected '%v', got '%v'", s, value.String())
	}
}

func testEscapePath(t *testing.T, json, path, expect string) {
	if gjson.Get(json, path).String() != expect {
		t.Fatalf("expected '%v', got '%v'", expect, gjson.Get(json, path).String())
	}
}

func TestEscapePath(t *testing.T) {
	text := `{
		"test":{
			"*":"valZ",
			"*v":"val0",
			"keyv*":"val1",
			"key*v":"val2",
			"keyv?":"val3",
			"key?v":"val4",
			"keyv.":"val5",
			"key.v":"val6",
			"keyk*":{"key?":"val7"}
		}
	}`

	testEscapePath(t, text, "test.\\*", "valZ")
	testEscapePath(t, text, "test.\\*v", "val0")
	testEscapePath(t, text, "test.keyv\\*", "val1")
	testEscapePath(t, text, "test.key\\*v", "val2")
	testEscapePath(t, text, "test.keyv\\?", "val3")
	testEscapePath(t, text, "test.key\\?v", "val4")
	testEscapePath(t, text, "test.keyv\\.", "val5")
	testEscapePath(t, text, "test.key\\.v", "val6")
	testEscapePath(t, text, "test.keyk\\*.key\\?", "val7")
}

// this json block is poorly formed on purpose.
var basicJSON = `  {"age":100, "name":{"here":"B\\\"R"},
	"noop":{"what is a wren?":"a bird"},
	"happy":true,"immortal":false,
	"items":[1,2,3,{"tags":[1,2,3],"points":[[1,2],[3,4]]},4,5,6,7],
	"arr":["1",2,"3",{"hello":"world"},"4",5],
	"vals":[1,2,3,{"sadf":sdf"asdf"}],"name":{"first":"tom","last":null},
	"created":"2014-05-16T08:28:06.989Z",
	"loggy":{
		"programmers": [
    	    {
    	        "firstName": "Brett",
    	        "lastName": "McLaughlin",
    	        "email": "aaaa",
				"tag": "good"
    	    },
    	    {
    	        "firstName": "Jason",
    	        "lastName": "Hunter",
    	        "email": "bbbb",
				"tag": "bad"
    	    },
    	    {
    	        "firstName": "Elliotte",
    	        "lastName": "Harold",
    	        "email": "cccc",
				"tag":, "good"
    	    },
			{
				"firstName": 1002.3,
				"age": 101
			}
    	]
	},
	"lastly":{"end...ing":"soon","yay":"final"}
}`

func TestPath(t *testing.T) {
	text := basicJSON
	r := gjson.Get(text, "@this")
	path := r.Path(text)
	if path != "@this" {
		t.FailNow()
	}

	r = gjson.Parse(text)
	path = r.Path(text)
	if path != "@this" {
		t.FailNow()
	}

	obj := gjson.Parse(text)
	obj.ForEach(func(key, val gjson.Result) bool {
		kp := key.Path(text)
		assert.True(t, kp == "")
		vp := val.Path(text)
		if vp == "name" {
			// there are two "name" keys
			return true
		}
		val2 := obj.Get(vp)
		assert.True(t, val2.Raw == val.Raw)
		return true
	})
	arr := obj.Get("loggy.programmers")
	arr.ForEach(func(_, val gjson.Result) bool {
		vp := val.Path(text)
		val2 := gjson.Get(text, vp)
		assert.True(t, val2.Raw == val.Raw)
		return true
	})
	get := func(path string) {
		r1 := gjson.Get(text, path)
		path2 := r1.Path(text)
		r2 := gjson.Get(text, path2)
		assert.True(t, r1.Raw == r2.Raw)
	}
	get("age")
	get("name")
	get("name.here")
	get("noop")
	get("noop.what is a wren?")
	get("arr.0")
	get("arr.1")
	get("arr.2")
	get("arr.3")
	get("arr.3.hello")
	get("arr.4")
	get("arr.5")
	get("loggy.programmers.2.email")
	get("lastly.end\\.\\.\\.ing")
	get("lastly.yay")
}

func TestTimeResult(t *testing.T) {
	assert.True(t, gjson.Get(basicJSON, "created").String() ==
		gjson.Get(basicJSON, "created").Time().Format(time.RFC3339Nano))
}

func TestParseAny(t *testing.T) {
	assert.True(t, gjson.Parse("100").Float() == 100)
	assert.True(t, gjson.Parse("true").Bool())
	assert.True(t, gjson.Parse("false").Bool() == false)
	assert.True(t, gjson.Parse("yikes").Exists() == false)
}

func TestBatchVariousPathCounts(t *testing.T) {
	text := `{"a":"a","b":"b","c":"c"}`
	counts := []int{3, 4, 7, 8, 9, 15, 16, 17, 31, 32, 33, 63, 64, 65, 127,
		128, 129, 255, 256, 257, 511, 512, 513}
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
		results := gjson.GetMany(text, gpaths...)
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
	assert.True(t, gjson.GetMany(text, path)[0].String() == "b")
}

func TestByteSafety(t *testing.T) {
	jsonb := []byte(`{"name":"Janet","age":38}`)
	mtok := gjson.GetBytes(jsonb, "name")
	if mtok.String() != "Janet" {
		t.Fatalf("expected %v, got %v", "Jason", mtok.String())
	}
	mtok2 := gjson.GetBytes(jsonb, "age")
	if mtok2.Raw != "38" {
		t.Fatalf("expected %v, got %v", "Jason", mtok2.Raw)
	}
	jsonb[9] = 'T'
	jsonb[12] = 'd'
	jsonb[13] = 'y'
	if mtok.String() != "Janet" {
		t.Fatalf("expected %v, got %v", "Jason", mtok.String())
	}
}

func get(json, path string) gjson.Result {
	return gjson.GetBytes([]byte(json), path)
}

func TestBasic(t *testing.T) {
	var mtok gjson.Result
	mtok = get(basicJSON, `loggy.programmers.#[tag="good"].firstName`)
	if mtok.String() != "Brett" {
		t.Fatalf("expected %v, got %v", "Brett", mtok.String())
	}
	mtok = get(basicJSON, `loggy.programmers.#[tag="good"]#.firstName`)
	if mtok.String() != `["Brett","Elliotte"]` {
		t.Fatalf("expected %v, got %v", `["Brett","Elliotte"]`, mtok.String())
	}
}

func TestIsArrayIsObject(t *testing.T) {
	mtok := get(basicJSON, "loggy")
	assert.True(t, mtok.IsObject())
	assert.True(t, !mtok.IsArray())

	mtok = get(basicJSON, "loggy.programmers")
	assert.True(t, !mtok.IsObject())
	assert.True(t, mtok.IsArray())

	mtok = get(basicJSON, `loggy.programmers.#[tag="good"]#.firstName`)
	assert.True(t, mtok.IsArray())

	mtok = get(basicJSON, `loggy.programmers.0.firstName`)
	assert.True(t, !mtok.IsObject())
	assert.True(t, !mtok.IsArray())
}

func TestPlus53BitInts(t *testing.T) {
	text := `{"IdentityData":{"GameInstanceId":634866135153775564}}`
	value := gjson.Get(text, "IdentityData.GameInstanceId")
	assert.True(t, value.Uint() == 634866135153775564)
	assert.True(t, value.Int() == 634866135153775564)
	assert.True(t, value.Float() == 634866135153775616)

	text = `{"IdentityData":{"GameInstanceId":634866135153775564.88172}}`
	value = gjson.Get(text, "IdentityData.GameInstanceId")
	assert.True(t, value.Uint() == 634866135153775616)
	assert.True(t, value.Int() == 634866135153775616)
	assert.True(t, value.Float() == 634866135153775616.88172)

	text = `{
		"min_uint64": 0,
		"max_uint64": 18446744073709551615,
		"overflow_uint64": 18446744073709551616,
		"min_int64": -9223372036854775808,
		"max_int64": 9223372036854775807,
		"overflow_int64": 9223372036854775808,
		"min_uint53":  0,
		"max_uint53":  4503599627370495,
		"overflow_uint53": 4503599627370496,
		"min_int53": -2251799813685248,
		"max_int53": 2251799813685247,
		"overflow_int53": 2251799813685248
	}`

	assert.True(t, gjson.Get(text, "min_uint53").Uint() == 0)
	assert.True(t, gjson.Get(text, "max_uint53").Uint() == 4503599627370495)
	assert.True(t, gjson.Get(text, "overflow_uint53").Int() == 4503599627370496)
	assert.True(t, gjson.Get(text, "min_int53").Int() == -2251799813685248)
	assert.True(t, gjson.Get(text, "max_int53").Int() == 2251799813685247)
	assert.True(t, gjson.Get(text, "overflow_int53").Int() == 2251799813685248)
	assert.True(t, gjson.Get(text, "min_uint64").Uint() == 0)
	assert.True(t, gjson.Get(text, "max_uint64").Uint() == 18446744073709551615)
	// this next value overflows the max uint64 by one which will just
	// flip the number to zero
	assert.True(t, gjson.Get(text, "overflow_uint64").Int() == 0)
	assert.True(t, gjson.Get(text, "min_int64").Int() == -9223372036854775808)
	assert.True(t, gjson.Get(text, "max_int64").Int() == 9223372036854775807)
	// this next value overflows the max int64 by one which will just
	// flip the number to the negative sign.
	assert.True(t, gjson.Get(text, "overflow_int64").Int() == -9223372036854775808)
}

func TestGet2(t *testing.T) {
	// These should not fail, even though the unicode is invalid.
	gjson.Get(`["S3O PEDRO DO BUTI\udf93"]`, "0")
	gjson.Get(`["S3O PEDRO DO BUTI\udf93asdf"]`, "0")
	gjson.Get(`["S3O PEDRO DO BUTI\udf93\u"]`, "0")
	gjson.Get(`["S3O PEDRO DO BUTI\udf93\u1"]`, "0")
	gjson.Get(`["S3O PEDRO DO BUTI\udf93\u13"]`, "0")
	gjson.Get(`["S3O PEDRO DO BUTI\udf93\u134"]`, "0")
	gjson.Get(`["S3O PEDRO DO BUTI\udf93\u1345"]`, "0")
	gjson.Get(`["S3O PEDRO DO BUTI\udf93\u1345asd"]`, "0")
}

func TestTypes(t *testing.T) {
	assert.True(t, (gjson.Result{Type: gjson.String}).Type.String() == "String")
	assert.True(t, (gjson.Result{Type: gjson.Number}).Type.String() == "Number")
	assert.True(t, (gjson.Result{Type: gjson.Null}).Type.String() == "Null")
	assert.True(t, (gjson.Result{Type: gjson.False}).Type.String() == "False")
	assert.True(t, (gjson.Result{Type: gjson.True}).Type.String() == "True")
	assert.True(t, (gjson.Result{Type: gjson.JSON}).Type.String() == "JSON")
	assert.True(t, (gjson.Result{Type: 100}).Type.String() == "")
	// bool
	assert.True(t, (gjson.Result{Type: gjson.True}).Bool() == true)
	assert.True(t, (gjson.Result{Type: gjson.False}).Bool() == false)
	assert.True(t, (gjson.Result{Type: gjson.Number, Num: 1}).Bool() == true)
	assert.True(t, (gjson.Result{Type: gjson.Number, Num: 0}).Bool() == false)
	assert.True(t, (gjson.Result{Type: gjson.String, Str: "1"}).Bool() == true)
	assert.True(t, (gjson.Result{Type: gjson.String, Str: "T"}).Bool() == true)
	assert.True(t, (gjson.Result{Type: gjson.String, Str: "t"}).Bool() == true)
	assert.True(t, (gjson.Result{Type: gjson.String, Str: "true"}).Bool() == true)
	assert.True(t, (gjson.Result{Type: gjson.String, Str: "True"}).Bool() == true)
	assert.True(t, (gjson.Result{Type: gjson.String, Str: "TRUE"}).Bool() == true)
	assert.True(t, (gjson.Result{Type: gjson.String, Str: "tRuE"}).Bool() == true)
	assert.True(t, (gjson.Result{Type: gjson.String, Str: "0"}).Bool() == false)
	assert.True(t, (gjson.Result{Type: gjson.String, Str: "f"}).Bool() == false)
	assert.True(t, (gjson.Result{Type: gjson.String, Str: "F"}).Bool() == false)
	assert.True(t, (gjson.Result{Type: gjson.String, Str: "false"}).Bool() == false)
	assert.True(t, (gjson.Result{Type: gjson.String, Str: "False"}).Bool() == false)
	assert.True(t, (gjson.Result{Type: gjson.String, Str: "FALSE"}).Bool() == false)
	assert.True(t, (gjson.Result{Type: gjson.String, Str: "fAlSe"}).Bool() == false)
	assert.True(t, (gjson.Result{Type: gjson.String, Str: "random"}).Bool() == false)

	// int
	assert.True(t, (gjson.Result{Type: gjson.String, Str: "1"}).Int() == 1)
	assert.True(t, (gjson.Result{Type: gjson.True}).Int() == 1)
	assert.True(t, (gjson.Result{Type: gjson.False}).Int() == 0)
	assert.True(t, (gjson.Result{Type: gjson.Number, Num: 1}).Int() == 1)
	// uintgjson.
	assert.True(t, (gjson.Result{Type: gjson.String, Str: "1"}).Uint() == 1)
	assert.True(t, (gjson.Result{Type: gjson.True}).Uint() == 1)
	assert.True(t, (gjson.Result{Type: gjson.False}).Uint() == 0)
	assert.True(t, (gjson.Result{Type: gjson.Number, Num: 1}).Uint() == 1)
	// floatgjson.
	assert.True(t, (gjson.Result{Type: gjson.String, Str: "1"}).Float() == 1)
	assert.True(t, (gjson.Result{Type: gjson.True}).Float() == 1)
	assert.True(t, (gjson.Result{Type: gjson.False}).Float() == 0)
	assert.True(t, (gjson.Result{Type: gjson.Number, Num: 1}).Float() == 1)
}

func TestForEach(t *testing.T) {
	gjson.Result{}.ForEach(nil)
	gjson.Result{Type: gjson.String, Str: "Hello"}.ForEach(func(_, value gjson.Result) bool {
		assert.True(t, value.String() == "Hello")
		return false
	})
	gjson.Result{Type: gjson.JSON, Raw: "*invalid*"}.ForEach(nil)

	text := ` {"name": {"first": "Janet","last": "Prichard"},
	"asd\nf":"\ud83d\udd13","age": 47}`
	var count int
	gjson.ParseBytes([]byte(text)).ForEach(func(key, value gjson.Result) bool {
		count++
		return true
	})
	assert.True(t, count == 3)
	gjson.ParseBytes([]byte(`{"bad`)).ForEach(nil)
	gjson.ParseBytes([]byte(`{"ok":"bad`)).ForEach(nil)
}

func TestMap(t *testing.T) {
	assert.True(t, len(gjson.ParseBytes([]byte(`"asdf"`)).Map()) == 0)
	assert.True(t, gjson.ParseBytes([]byte(`{"asdf":"ghjk"`)).Map()["asdf"].String() == "ghjk")
	assert.True(t, len(gjson.Result{Type: gjson.JSON, Raw: "**invalid**"}.Map()) == 0)
	assert.True(t, gjson.Result{Type: gjson.JSON, Raw: "**invalid**"}.Value() == nil)
	assert.True(t, gjson.Result{Type: gjson.JSON, Raw: "{"}.Map() != nil)
}

func TestBasic1(t *testing.T) {
	mtok := get(basicJSON, `loggy.programmers`)
	var count int
	mtok.ForEach(func(key, value gjson.Result) bool {
		assert.True(t, key.Exists())
		assert.True(t, key.String() == fmt.Sprint(count))
		assert.True(t, key.Int() == int64(count))
		count++
		if count == 3 {
			return false
		}
		if count == 1 {
			i := 0
			value.ForEach(func(key, value gjson.Result) bool {
				switch i {
				case 0:
					if key.String() != "firstName" ||
						value.String() != "Brett" {
						t.Fatalf("expected %v/%v got %v/%v", "firstName",
							"Brett", key.String(), value.String())
					}
				case 1:
					if key.String() != "lastName" ||
						value.String() != "McLaughlin" {
						t.Fatalf("expected %v/%v got %v/%v", "lastName",
							"McLaughlin", key.String(), value.String())
					}
				case 2:
					if key.String() != "email" || value.String() != "aaaa" {
						t.Fatalf("expected %v/%v got %v/%v", "email", "aaaa",
							key.String(), value.String())
					}
				}
				i++
				return true
			})
		}
		return true
	})
	if count != 3 {
		t.Fatalf("expected %v, got %v", 3, count)
	}
}

func TestBasic2(t *testing.T) {
	mtok := get(basicJSON, `loggy.programmers.#[age=101].firstName`)
	if mtok.String() != "1002.3" {
		t.Fatalf("expected %v, got %v", "1002.3", mtok.String())
	}
	mtok = get(basicJSON,
		`loggy.programmers.#[firstName != "Brett"].firstName`)
	if mtok.String() != "Jason" {
		t.Fatalf("expected %v, got %v", "Jason", mtok.String())
	}
	mtok = get(basicJSON, `loggy.programmers.#[firstName % "Bre*"].email`)
	if mtok.String() != "aaaa" {
		t.Fatalf("expected %v, got %v", "aaaa", mtok.String())
	}
	mtok = get(basicJSON, `loggy.programmers.#[firstName !% "Bre*"].email`)
	if mtok.String() != "bbbb" {
		t.Fatalf("expected %v, got %v", "bbbb", mtok.String())
	}
	mtok = get(basicJSON, `loggy.programmers.#[firstName == "Brett"].email`)
	if mtok.String() != "aaaa" {
		t.Fatalf("expected %v, got %v", "aaaa", mtok.String())
	}
	mtok = get(basicJSON, "loggy")
	if mtok.Type != gjson.JSON {
		t.Fatalf("expected %v, got %v", gjson.JSON, mtok.Type)
	}
	if len(mtok.Map()) != 1 {
		t.Fatalf("expected %v, got %v", 1, len(mtok.Map()))
	}
	programmers := mtok.Map()["programmers"]
	if programmers.Array()[1].Map()["firstName"].Str != "Jason" {
		t.Fatalf("expected %v, got %v", "Jason",
			mtok.Map()["programmers"].Array()[1].Map()["firstName"].Str)
	}
}

func TestBasic3(t *testing.T) {
	var mtok gjson.Result
	if gjson.Parse(basicJSON).Get("loggy.programmers").Get("1").
		Get("firstName").Str != "Jason" {
		t.Fatalf("expected %v, got %v", "Jason", gjson.Parse(basicJSON).
			Get("loggy.programmers").Get("1").Get("firstName").Str)
	}
	var token gjson.Result
	if token = gjson.Parse("-102"); token.Num != -102 {
		t.Fatalf("expected %v, got %v", -102, token.Num)
	}
	if token = gjson.Parse("102"); token.Num != 102 {
		t.Fatalf("expected %v, got %v", 102, token.Num)
	}
	if token = gjson.Parse("102.2"); token.Num != 102.2 {
		t.Fatalf("expected %v, got %v", 102.2, token.Num)
	}
	if token = gjson.Parse(`"hello"`); token.Str != "hello" {
		t.Fatalf("expected %v, got %v", "hello", token.Str)
	}
	if token = gjson.Parse(`"\"he\nllo\""`); token.Str != "\"he\nllo\"" {
		t.Fatalf("expected %v, got %v", "\"he\nllo\"", token.Str)
	}
	mtok = get(basicJSON, "loggy.programmers.#.firstName")
	if len(mtok.Array()) != 4 {
		t.Fatalf("expected 4, got %v", len(mtok.Array()))
	}
	for i, ex := range []string{"Brett", "Jason", "Elliotte", "1002.3"} {
		if mtok.Array()[i].String() != ex {
			t.Fatalf("expected '%v', got '%v'", ex, mtok.Array()[i].String())
		}
	}
	mtok = get(basicJSON, "loggy.programmers.#.asd")
	if mtok.Type != gjson.JSON {
		t.Fatalf("expected %v, got %v", gjson.JSON, mtok.Type)
	}
	if len(mtok.Array()) != 0 {
		t.Fatalf("expected 0, got %v", len(mtok.Array()))
	}
}

func TestBasic4(t *testing.T) {
	if get(basicJSON, "items.3.tags.#").Num != 3 {
		t.Fatalf("expected 3, got %v", get(basicJSON, "items.3.tags.#").Num)
	}
	if get(basicJSON, "items.3.points.1.#").Num != 2 {
		t.Fatalf("expected 2, got %v",
			get(basicJSON, "items.3.points.1.#").Num)
	}
	if get(basicJSON, "items.#").Num != 8 {
		t.Fatalf("expected 6, got %v", get(basicJSON, "items.#").Num)
	}
	if get(basicJSON, "vals.#").Num != 4 {
		t.Fatalf("expected 4, got %v", get(basicJSON, "vals.#").Num)
	}
	if !get(basicJSON, "name.last").Exists() {
		t.Fatal("expected true, got false")
	}
	token := get(basicJSON, "name.here")
	if token.String() != "B\\\"R" {
		t.Fatal("expecting 'B\\\"R'", "got", token.String())
	}
	token = get(basicJSON, "arr.#")
	if token.String() != "6" {
		fmt.Printf("%#v\n", token)
		t.Fatal("expecting 6", "got", token.String())
	}
	token = get(basicJSON, "arr.3.hello")
	if token.String() != "world" {
		t.Fatal("expecting 'world'", "got", token.String())
	}
	_ = token.Value().(string)
	token = get(basicJSON, "name.first")
	if token.String() != "tom" {
		t.Fatal("expecting 'tom'", "got", token.String())
	}
	_ = token.Value().(string)
	token = get(basicJSON, "name.last")
	if token.String() != "" {
		t.Fatal("expecting ''", "got", token.String())
	}
	if token.Value() != nil {
		t.Fatal("should be nil")
	}
}

func TestBasic5(t *testing.T) {
	token := get(basicJSON, "age")
	if token.String() != "100" {
		t.Fatal("expecting '100'", "got", token.String())
	}
	_ = token.Value().(float64)
	token = get(basicJSON, "happy")
	if token.String() != "true" {
		t.Fatal("expecting 'true'", "got", token.String())
	}
	_ = token.Value().(bool)
	token = get(basicJSON, "immortal")
	if token.String() != "false" {
		t.Fatal("expecting 'false'", "got", token.String())
	}
	_ = token.Value().(bool)
	token = get(basicJSON, "noop")
	if token.String() != `{"what is a wren?":"a bird"}` {
		t.Fatal("expecting '"+`{"what is a wren?":"a bird"}`+"'", "got",
			token.String())
	}
	_ = token.Value().(map[string]any)

	if get(basicJSON, "").Value() != nil {
		t.Fatal("should be nil")
	}

	get(basicJSON, "vals.hello")

	type msi = map[string]any
	type fi = []any
	mm := gjson.Parse(basicJSON).Value().(msi)
	fn := mm["loggy"].(msi)["programmers"].(fi)[1].(msi)["firstName"].(string)
	if fn != "Jason" {
		t.Fatalf("expecting %v, got %v", "Jason", fn)
	}
}

func TestUnicode(t *testing.T) {
	var text = `{"key":0,"的情况下解":{"key":1,"的情况":2}}`
	if gjson.Get(text, "的情况下解.key").Num != 1 {
		t.Fatal("fail")
	}
	if gjson.Get(text, "的情况下解.的情况").Num != 2 {
		t.Fatal("fail")
	}
	if gjson.Get(text, "的情况下解.的?况").Num != 2 {
		t.Fatal("fail")
	}
	if gjson.Get(text, "的情况下解.的?*").Num != 2 {
		t.Fatal("fail")
	}
	if gjson.Get(text, "的情况下解.*?况").Num != 2 {
		t.Fatal("fail")
	}
	if gjson.Get(text, "的情?下解.*?况").Num != 2 {
		t.Fatal("fail")
	}
	if gjson.Get(text, "的情下解.*?况").Num != 0 {
		t.Fatal("fail")
	}
}

func TestLess(t *testing.T) {
	assert.True(t, !gjson.Result{Type: gjson.Null}.Less(gjson.Result{Type: gjson.Null}, true))
	assert.True(t, gjson.Result{Type: gjson.Null}.Less(gjson.Result{Type: gjson.False}, true))
	assert.True(t, gjson.Result{Type: gjson.Null}.Less(gjson.Result{Type: gjson.True}, true))
	assert.True(t, gjson.Result{Type: gjson.Null}.Less(gjson.Result{Type: gjson.JSON}, true))
	assert.True(t, gjson.Result{Type: gjson.Null}.Less(gjson.Result{Type: gjson.Number}, true))
	assert.True(t, gjson.Result{Type: gjson.Null}.Less(gjson.Result{Type: gjson.String}, true))
	assert.True(t, !gjson.Result{Type: gjson.False}.Less(gjson.Result{Type: gjson.Null}, true))
	assert.True(t, gjson.Result{Type: gjson.False}.Less(gjson.Result{Type: gjson.True}, true))
	assert.True(t, gjson.Result{Type: gjson.String, Str: "abc"}.Less(gjson.Result{
		Type: gjson.String,
		Str:  "bcd",
	}, true))
	assert.True(t, gjson.Result{Type: gjson.String, Str: "ABC"}.Less(gjson.Result{
		Type: gjson.String,
		Str:  "abc",
	}, true))
	assert.True(t, !gjson.Result{Type: gjson.String, Str: "ABC"}.Less(gjson.Result{
		Type: gjson.String,
		Str:  "abc",
	}, false))
	assert.True(t, gjson.Result{Type: gjson.Number, Num: 123}.Less(gjson.Result{
		Type: gjson.Number,
		Num:  456,
	}, true))
	assert.True(t, !gjson.Result{Type: gjson.Number, Num: 456}.Less(gjson.Result{
		Type: gjson.Number,
		Num:  123,
	}, true))
	assert.True(t, !gjson.Result{Type: gjson.Number, Num: 456}.Less(gjson.Result{
		Type: gjson.Number,
		Num:  456,
	}, true))
}

func TestGet3(t *testing.T) {
	data := `{
      "code": 0,
      "msg": "",
      "data": {
        "sz002024": {
          "qfqday": [
            [
              "2014-01-02",
              "8.93",
              "9.03",
              "9.17",
              "8.88",
              "621143.00"
            ],
            [
              "2014-01-03",
              "9.03",
              "9.30",
              "9.47",
              "8.98",
              "1624438.00"
            ]
          ]
        }
      }
    }`

	var num []string
	for _, v := range gjson.Get(data, "data.sz002024.qfqday.0").Array() {
		num = append(num, v.String())
	}
	if fmt.Sprintf("%v", num) != "[2014-01-02 8.93 9.03 9.17 8.88 621143.00]" {
		t.Fatalf("invalid result")
	}
}

var exampleJSON = `{
	"widget": {
		"debug": "on",
		"window": {
			"title": "Sample Konfabulator Widget",
			"name": "main_window",
			"width": 500,
			"height": 500
		},
		"image": {
			"src": "Images/Sun.png",
			"hOffset": 250,
			"vOffset": 250,
			"alignment": "center"
		},
		"text": {
			"data": "Click Here",
			"size": 36,
			"style": "bold",
			"vOffset": 100,
			"alignment": "center",
			"onMouseUp": "sun1.opacity = (sun1.opacity / 100) * 90;"
		}
	}
}`

func TestUnmarshalMap(t *testing.T) {
	var m1 = gjson.Parse(exampleJSON).Value().(map[string]any)
	var m2 map[string]any
	if err := json.Unmarshal([]byte(exampleJSON), &m2); err != nil {
		t.Fatal(err)
	}
	b1, err := json.Marshal(m1)
	if err != nil {
		t.Fatal(err)
	}
	b2, err := json.Marshal(m2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b1, b2) {
		t.Fatal("b1 != b2")
	}
}

func TestSingleArrayValue(t *testing.T) {
	var text = `{"key": "value","key2":[1,2,3,4,"A"]}`
	var result = gjson.Get(text, "key")
	var array = result.Array()
	if len(array) != 1 {
		t.Fatal("array is empty")
	}
	if array[0].String() != "value" {
		t.Fatalf("got %s, should be %s", array[0].String(), "value")
	}

	array = gjson.Get(text, "key2.#").Array()
	if len(array) != 1 {
		t.Fatalf("got '%v', expected '%v'", len(array), 1)
	}

	array = gjson.Get(text, "key3").Array()
	if len(array) != 0 {
		t.Fatalf("got '%v', expected '%v'", len(array), 0)
	}

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
		results := gjson.GetManyBytes(
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

func testMany(t *testing.T, json string, paths, expected []string) {
	testManyAny(t, json, paths, expected, true)
	testManyAny(t, json, paths, expected, false)
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
				result = gjson.GetManyBytes([]byte(json), paths...)
			} else {
				result = gjson.GetMany(json, paths...)
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

func TestIssue20(t *testing.T) {
	text := `{ "name": "FirstName", "name1": "FirstName1", ` +
		`"address": "address1", "addressDetails": "address2", }`
	paths := []string{"name", "name1", "address", "addressDetails"}
	expected := []string{"FirstName", "FirstName1", "address1", "address2"}
	t.Run("SingleMany", func(t *testing.T) {
		testMany(t, text, paths,
			expected)
	})
}

func TestIssue21(t *testing.T) {
	text := `{ "Level1Field1":3,
	           "Level1Field4":4,
			   "Level1Field2":{ "Level2Field1":[ "value1", "value2" ],
			   "Level2Field2":{ "Level3Field1":[ { "key1":"value1" } ] } } }`
	paths := []string{"Level1Field1", "Level1Field2.Level2Field1",
		"Level1Field2.Level2Field2.Level3Field1", "Level1Field4"}
	expected := []string{"3", `[ "value1", "value2" ]`,
		`[ { "key1":"value1" } ]`, "4"}
	t.Run("SingleMany", func(t *testing.T) {
		testMany(t, text, paths,
			expected)
	})
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
		gjson.GetMany(lstr, paths...)
	}
}

var complicatedJSON = `
{
	"tagged": "OK",
	"Tagged": "KO",
	"NotTagged": true,
	"unsettable": 101,
	"Nested": {
		"Yellow": "Green",
		"yellow": "yellow"
	},
	"nestedTagged": {
		"Green": "Green",
		"Map": {
			"this": "that",
			"and": "the other thing"
		},
		"Ints": {
			"Uint": 99,
			"Uint16": 16,
			"Uint32": 32,
			"Uint64": 65
		},
		"Uints": {
			"int": -99,
			"Int": -98,
			"Int16": -16,
			"Int32": -32,
			"int64": -64,
			"Int64": -65
		},
		"Uints": {
			"Float32": 32.32,
			"Float64": 64.64
		},
		"Byte": 254,
		"Bool": true
	},
	"LeftOut": "you shouldn't be here",
	"SelfPtr": {"tagged":"OK","nestedTagged":{"Ints":{"Uint32":32}}},
	"SelfSlice": [{"tagged":"OK","nestedTagged":{"Ints":{"Uint32":32}}}],
	"SelfSlicePtr": [{"tagged":"OK","nestedTagged":{"Ints":{"Uint32":32}}}],
	"SelfPtrSlice": [{"tagged":"OK","nestedTagged":{"Ints":{"Uint32":32}}}],
	"interface": "Tile38 Rocks!",
	"Interface": "Please Download",
	"Array": [0,2,3,4,5],
	"time": "2017-05-07T13:24:43-07:00",
	"Binary": "R0lGODlhPQBEAPeo",
	"NonBinary": [9,3,100,115]
}
`

func TestGetMany47(t *testing.T) {
	text := `{"bar": {"id": 99, "mybar": "my mybar" }, "foo": ` +
		`{"myfoo": [605]}}`
	paths := []string{"foo.myfoo", "bar.id", "bar.mybar", "bar.mybarx"}
	expected := []string{"[605]", "99", "my mybar", ""}
	results := gjson.GetMany(text, paths...)
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
	results := gjson.GetMany(text, paths...)
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

func TestResultRawForLiteral(t *testing.T) {
	for _, lit := range []string{"null", "true", "false"} {
		result := gjson.Parse(lit)
		if result.Raw != lit {
			t.Fatalf("expected '%v', got '%v'", lit, result.Raw)
		}
	}
}

func TestNullArray(t *testing.T) {
	n := len(gjson.Get(`{"data":null}`, "data").Array())
	if n != 0 {
		t.Fatalf("expected '%v', got '%v'", 0, n)
	}
	n = len(gjson.Get(`{}`, "data").Array())
	if n != 0 {
		t.Fatalf("expected '%v', got '%v'", 0, n)
	}
	n = len(gjson.Get(`{"data":[]}`, "data").Array())
	if n != 0 {
		t.Fatalf("expected '%v', got '%v'", 0, n)
	}
	n = len(gjson.Get(`{"data":[null]}`, "data").Array())
	if n != 1 {
		t.Fatalf("expected '%v', got '%v'", 1, n)
	}
}

func TestIssue54(t *testing.T) {
	var r []gjson.Result
	text := `{"MarketName":null,"Nounce":6115}`
	r = gjson.GetMany(text, "Nounce", "Buys", "Sells", "Fills")
	if strings.Replace(fmt.Sprintf("%v", r), " ", "", -1) != "[6115]" {
		t.Fatalf("expected '%v', got '%v'", "[6115]",
			strings.Replace(fmt.Sprintf("%v", r), " ", "", -1))
	}
	r = gjson.GetMany(text, "Nounce", "Buys", "Sells")
	if strings.Replace(fmt.Sprintf("%v", r), " ", "", -1) != "[6115]" {
		t.Fatalf("expected '%v', got '%v'", "[6115]",
			strings.Replace(fmt.Sprintf("%v", r), " ", "", -1))
	}
	r = gjson.GetMany(text, "Nounce")
	if strings.Replace(fmt.Sprintf("%v", r), " ", "", -1) != "[6115]" {
		t.Fatalf("expected '%v', got '%v'", "[6115]",
			strings.Replace(fmt.Sprintf("%v", r), " ", "", -1))
	}
}

func TestIssue55(t *testing.T) {
	text := `{"one": {"two": 2, "three": 3}, "four": 4, "five": 5}`
	results := gjson.GetMany(text, "four", "five", "one.two", "one.six")
	expected := []string{"4", "5", "2", ""}
	for i, r := range results {
		if r.String() != expected[i] {
			t.Fatalf("expected %v, got %v", expected[i], r.String())
		}
	}
}

func TestIssue58(t *testing.T) {
	text := `{"data":[{"uid": 1},{"uid": 2}]}`
	res := gjson.Get(text, `data.#[uid!=1]`).Raw
	if res != `{"uid": 2}` {
		t.Fatalf("expected '%v', got '%v'", `{"uid": 1}`, res)
	}
}

func TestObjectGrouping(t *testing.T) {
	text := `
[
	true,
	{"name":"tom"},
	false,
	{"name":"janet"},
	null
]
`
	res := gjson.Get(text, "#.name")
	if res.String() != `["tom","janet"]` {
		t.Fatalf("expected '%v', got '%v'", `["tom","janet"]`, res.String())
	}
}

func TestJSONLines(t *testing.T) {
	text := `
true
false
{"name":"tom"}
[1,2,3,4,5]
{"name":"janet"}
null
12930.1203
	`
	paths := []string{"..#", "..0", "..2.name", "..#.name", "..6", "..7"}
	ress := []string{"7", "true", "tom", `["tom","janet"]`, "12930.1203", ""}
	for i, path := range paths {
		res := gjson.Get(text, path)
		if res.String() != ress[i] {
			t.Fatalf("expected '%v', got '%v'", ress[i], res.String())
		}
	}

	text = `
{"name": "Gilbert", "wins": [["straight", "7♣"], ["one pair", "10♥"]]}
{"name": "Alexa", "wins": [["two pair", "4♠"], ["two pair", "9♠"]]}
{"name": "May", "wins": []}
{"name": "Deloise", "wins": [["three of a kind", "5♣"]]}
`

	var i int
	lines := strings.Split(strings.TrimSpace(text), "\n")
	gjson.ForEachLine(text, func(line gjson.Result) bool {
		if line.Raw != lines[i] {
			t.Fatalf("expected '%v', got '%v'", lines[i], line.Raw)
		}
		i++
		return true
	})
	if i != 4 {
		t.Fatalf("expected '%v', got '%v'", 4, i)
	}

}

func TestNumUint64String(t *testing.T) {
	var i int64 = 9007199254740993 //2^53 + 1
	j := fmt.Sprintf(`{"data":  [  %d, "hello" ] }`, i)
	res := gjson.Get(j, "data.0")
	if res.String() != "9007199254740993" {
		t.Fatalf("expected '%v', got '%v'", "9007199254740993", res.String())
	}
}

func TestNumInt64String(t *testing.T) {
	var i int64 = -9007199254740993
	j := fmt.Sprintf(`{"data":[ "hello", %d ]}`, i)
	res := gjson.Get(j, "data.1")
	if res.String() != "-9007199254740993" {
		t.Fatalf("expected '%v', got '%v'", "-9007199254740993", res.String())
	}
}

func TestNumBigString(t *testing.T) {
	i := "900719925474099301239109123101" // very big
	j := fmt.Sprintf(`{"data":[ "hello", "%s" ]}`, i)
	res := gjson.Get(j, "data.1")
	if res.String() != "900719925474099301239109123101" {
		t.Fatalf("expected '%v', got '%v'", "900719925474099301239109123101",
			res.String())
	}
}

func TestNumFloatString(t *testing.T) {
	var i int64 = -9007199254740993
	j := fmt.Sprintf(`{"data":[ "hello", %d ]}`, i) //No quotes around value!!
	res := gjson.Get(j, "data.1")
	if res.String() != "-9007199254740993" {
		t.Fatalf("expected '%v', got '%v'", "-9007199254740993", res.String())
	}
}

func TestDuplicateKeys(t *testing.T) {
	var text = `{"name": "Alex","name": "Peter"}`
	if gjson.Parse(text).Get("name").String() !=
		gjson.Parse(text).Map()["name"].String() {
		t.Fatalf("expected '%v', got '%v'",
			gjson.Parse(text).Get("name").String(),
			gjson.Parse(text).Map()["name"].String(),
		)
	}
	if !gjson.Valid(text) {
		t.Fatal("should be valid")
	}
}

func TestArrayValues(t *testing.T) {
	var text = `{"array": ["PERSON1","PERSON2",0],}`
	values := gjson.Get(text, "array").Array()
	var output string
	for i, val := range values {
		if i > 0 {
			output += "\n"
		}
		output += fmt.Sprintf("%#v", val)
	}
	expect := strings.Join([]string{
		`gjson.Result{Type:3, Raw:"\"PERSON1\"", Str:"PERSON1", Num:0, ` +
			`Index:11, Indexes:[]int(nil)}`,
		`gjson.Result{Type:3, Raw:"\"PERSON2\"", Str:"PERSON2", Num:0, ` +
			`Index:21, Indexes:[]int(nil)}`,
		`gjson.Result{Type:2, Raw:"0", Str:"", Num:0, Index:31, Indexes:[]int(nil)}`,
	}, "\n")
	if output != expect {
		t.Fatalf("expected '%v', got '%v'", expect, output)
	}
}

func BenchmarkValid(b *testing.B) {
	for b.Loop() {
		gjson.Valid(complicatedJSON)
	}
}

func BenchmarkValidBytes(b *testing.B) {
	complicatedJSON := []byte(complicatedJSON)
	for b.Loop() {
		gjson.ValidBytes(complicatedJSON)
	}
}

func BenchmarkGoStdlibValidBytes(b *testing.B) {
	complicatedJSON := []byte(complicatedJSON)
	for b.Loop() {
		json.Valid(complicatedJSON)
	}
}

func TestChaining(t *testing.T) {
	text := `{
		"info": {
			"friends": [
				{"first": "Dale", "last": "Murphy", "age": 44},
				{"first": "Roger", "last": "Craig", "age": 68},
				{"first": "Jane", "last": "Murphy", "age": 47}
			]
		}
	  }`
	res := gjson.Get(text, "info.friends|0|first").String()
	if res != "Dale" {
		t.Fatalf("expected '%v', got '%v'", "Dale", res)
	}
	res = gjson.Get(text, "info.friends|@reverse|0|age").String()
	if res != "47" {
		t.Fatalf("expected '%v', got '%v'", "47", res)
	}
	res = gjson.Get(text, "@ugly|i\\nfo|friends.0.first").String()
	if res != "Dale" {
		t.Fatalf("expected '%v', got '%v'", "Dale", res)
	}
}

func TestArrayEx(t *testing.T) {
	text := `
	[
		{
			"c":[
				{"a":10.11}
			]
		}, {
			"c":[
				{"a":11.11}
			]
		}
	]`
	res := gjson.Get(text, "@ugly|#.c.#[a=10.11]").String()
	if res != `[{"a":10.11}]` {
		t.Fatalf("expected '%v', got '%v'", `[{"a":10.11}]`, res)
	}
	res = gjson.Get(text, "@ugly|#.c.#").String()
	if res != `[1,1]` {
		t.Fatalf("expected '%v', got '%v'", `[1,1]`, res)
	}
	res = gjson.Get(text, "@reverse|0|c|0|a").String()
	if res != "11.11" {
		t.Fatalf("expected '%v', got '%v'", "11.11", res)
	}
	res = gjson.Get(text, "#.c|#").String()
	if res != "2" {
		t.Fatalf("expected '%v', got '%v'", "2", res)
	}
}

func TestPipeDotMixing(t *testing.T) {
	text := `{
		"info": {
			"friends": [
				{"first": "Dale", "last": "Murphy", "age": 44},
				{"first": "Roger", "last": "Craig", "age": 68},
				{"first": "Jane", "last": "Murphy", "age": 47}
			]
		}
	  }`
	var res string
	res = gjson.Get(text, `info.friends.#[first="Dale"].last`).String()
	if res != "Murphy" {
		t.Fatalf("expected '%v', got '%v'", "Murphy", res)
	}
	res = gjson.Get(text, `info|friends.#[first="Dale"].last`).String()
	if res != "Murphy" {
		t.Fatalf("expected '%v', got '%v'", "Murphy", res)
	}
	res = gjson.Get(text, `info|friends.#[first="Dale"]|last`).String()
	if res != "Murphy" {
		t.Fatalf("expected '%v', got '%v'", "Murphy", res)
	}
	res = gjson.Get(text, `info|friends|#[first="Dale"]|last`).String()
	if res != "Murphy" {
		t.Fatalf("expected '%v', got '%v'", "Murphy", res)
	}
	res = gjson.Get(text, `@ugly|info|friends|#[first="Dale"]|last`).String()
	if res != "Murphy" {
		t.Fatalf("expected '%v', got '%v'", "Murphy", res)
	}
	res = gjson.Get(text, `@ugly|info.@ugly|friends|#[first="Dale"]|last`).String()
	if res != "Murphy" {
		t.Fatalf("expected '%v', got '%v'", "Murphy", res)
	}
	res = gjson.Get(text, `@ugly.info|@ugly.friends|#[first="Dale"]|last`).String()
	if res != "Murphy" {
		t.Fatalf("expected '%v', got '%v'", "Murphy", res)
	}
}

func TestDeepSelectors(t *testing.T) {
	text := `{
		"info": {
			"friends": [
				{
					"first": "Dale", "last": "Murphy",
					"extra": [10,20,30],
					"details": {
						"city": "Tempe",
						"state": "Arizona"
					}
				},
				{
					"first": "Roger", "last": "Craig",
					"extra": [40,50,60],
					"details": {
						"city": "Phoenix",
						"state": "Arizona"
					}
				}
			]
		}
	  }`
	var res string
	res = gjson.Get(text, `info.friends.#[first="Dale"].extra.0`).String()
	if res != "10" {
		t.Fatalf("expected '%v', got '%v'", "10", res)
	}
	res = gjson.Get(text, `info.friends.#[first="Dale"].extra|0`).String()
	if res != "10" {
		t.Fatalf("expected '%v', got '%v'", "10", res)
	}
	res = gjson.Get(text, `info.friends.#[first="Dale"]|extra|0`).String()
	if res != "10" {
		t.Fatalf("expected '%v', got '%v'", "10", res)
	}
	res = gjson.Get(text, `info.friends.#[details.city="Tempe"].last`).String()
	if res != "Murphy" {
		t.Fatalf("expected '%v', got '%v'", "Murphy", res)
	}
	res = gjson.Get(text, `info.friends.#[details.city="Phoenix"].last`).String()
	if res != "Craig" {
		t.Fatalf("expected '%v', got '%v'", "Craig", res)
	}
	res = gjson.Get(text, `info.friends.#[details.state="Arizona"].last`).String()
	if res != "Murphy" {
		t.Fatalf("expected '%v', got '%v'", "Murphy", res)
	}
}

func TestMultiArrayEx(t *testing.T) {
	text := `{
		"info": {
			"friends": [
				{
					"first": "Dale", "last": "Murphy", "kind": "Person",
					"cust1": true,
					"extra": [10,20,30],
					"details": {
						"city": "Tempe",
						"state": "Arizona"
					}
				},
				{
					"first": "Roger", "last": "Craig", "kind": "Person",
					"cust2": false,
					"extra": [40,50,60],
					"details": {
						"city": "Phoenix",
						"state": "Arizona"
					}
				}
			]
		}
	  }`

	var res string

	res = gjson.Get(text, `info.friends.#[kind="Person"]#.kind|0`).String()
	if res != "Person" {
		t.Fatalf("expected '%v', got '%v'", "Person", res)
	}
	res = gjson.Get(text, `info.friends.#.kind|0`).String()
	if res != "Person" {
		t.Fatalf("expected '%v', got '%v'", "Person", res)
	}

	res = gjson.Get(text, `info.friends.#[kind="Person"]#.kind`).String()
	if res != `["Person","Person"]` {
		t.Fatalf("expected '%v', got '%v'", `["Person","Person"]`, res)
	}
	res = gjson.Get(text, `info.friends.#.kind`).String()
	if res != `["Person","Person"]` {
		t.Fatalf("expected '%v', got '%v'", `["Person","Person"]`, res)
	}

	res = gjson.Get(text, `info.friends.#[kind="Person"]#|kind`).String()
	if res != `` {
		t.Fatalf("expected '%v', got '%v'", ``, res)
	}
	res = gjson.Get(text, `info.friends.#|kind`).String()
	if res != `` {
		t.Fatalf("expected '%v', got '%v'", ``, res)
	}

	res = gjson.Get(text, `i*.f*.#[kind="Other"]#`).String()
	if res != `[]` {
		t.Fatalf("expected '%v', got '%v'", `[]`, res)
	}
}

func TestQueries(t *testing.T) {
	text := `{
		"info": {
			"friends": [
				{
					"first": "Dale", "last": "Murphy", "kind": "Person",
					"cust1": true,
					"extra": [10,20,30],
					"details": {
						"city": "Tempe",
						"state": "Arizona"
					}
				},
				{
					"first": "Roger", "last": "Craig", "kind": "Person",
					"cust2": false,
					"extra": [40,50,60],
					"details": {
						"city": "Phoenix",
						"state": "Arizona"
					}
				}
			]
		}
	  }`

	// numbers
	assert.True(t, gjson.Get(text, "i*.f*.#[extra.0<11].first").Exists())
	assert.True(t, gjson.Get(text, "i*.f*.#[extra.0<=11].first").Exists())
	assert.True(t, !gjson.Get(text, "i*.f*.#[extra.0<10].first").Exists())
	assert.True(t, gjson.Get(text, "i*.f*.#[extra.0<=10].first").Exists())
	assert.True(t, gjson.Get(text, "i*.f*.#[extra.0=10].first").Exists())
	assert.True(t, !gjson.Get(text, "i*.f*.#[extra.0=11].first").Exists())
	assert.True(t, gjson.Get(text, "i*.f*.#[extra.0!=10].first").String() == "Roger")
	assert.True(t, gjson.Get(text, "i*.f*.#[extra.0>10].first").String() == "Roger")
	assert.True(t, gjson.Get(text, "i*.f*.#[extra.0>=10].first").String() == "Dale")

	// strings
	assert.True(t, gjson.Get(text, `i*.f*.#[extra.0<"11"].first`).Exists())
	assert.True(t, gjson.Get(text, `i*.f*.#[first>"Dale"].last`).String() == "Craig")
	assert.True(t, gjson.Get(text, `i*.f*.#[first>="Dale"].last`).String() == "Murphy")
	assert.True(t, gjson.Get(text, `i*.f*.#[first="Dale"].last`).String() == "Murphy")
	assert.True(t, gjson.Get(text, `i*.f*.#[first!="Dale"].last`).String() == "Craig")
	assert.True(t, !gjson.Get(text, `i*.f*.#[first<"Dale"].last`).Exists())
	assert.True(t, gjson.Get(text, `i*.f*.#[first<="Dale"].last`).Exists())
	assert.True(t, gjson.Get(text, `i*.f*.#[first%"Da*"].last`).Exists())
	assert.True(t, gjson.Get(text, `i*.f*.#[first%"Dale"].last`).Exists())
	assert.True(t, gjson.Get(text, `i*.f*.#[first%"*a*"]#|#`).String() == "1")
	assert.True(t, gjson.Get(text, `i*.f*.#[first%"*e*"]#|#`).String() == "2")
	assert.True(t, gjson.Get(text, `i*.f*.#[first!%"*e*"]#|#`).String() == "0")

	// trues
	assert.True(t, gjson.Get(text, `i*.f*.#[cust1=true].first`).String() == "Dale")
	assert.True(t, gjson.Get(text, `i*.f*.#[cust2=false].first`).String() == "Roger")
	assert.True(t, gjson.Get(text, `i*.f*.#[cust1!=false].first`).String() == "Dale")
	assert.True(t, gjson.Get(text, `i*.f*.#[cust2!=true].first`).String() == "Roger")
	assert.True(t, !gjson.Get(text, `i*.f*.#[cust1>true].first`).Exists())
	assert.True(t, gjson.Get(text, `i*.f*.#[cust1>=true].first`).Exists())
	assert.True(t, !gjson.Get(text, `i*.f*.#[cust2<false].first`).Exists())
	assert.True(t, gjson.Get(text, `i*.f*.#[cust2<=false].first`).Exists())

}

func TestQueryArrayValues(t *testing.T) {
	text := `{
		"artists": [
			["Bob Dylan"],
			"John Lennon",
			"Mick Jagger",
			"Elton John",
			"Michael Jackson",
			"John Smith",
			true,
			123,
			456,
			false,
			null
		]
	}`
	assert.True(t, gjson.Get(text, `a*.#[0="Bob Dylan"]#|#`).String() == "1")
	assert.True(t, gjson.Get(text, `a*.#[0="Bob Dylan 2"]#|#`).String() == "0")
	assert.True(t, gjson.Get(text, `a*.#[%"John*"]#|#`).String() == "2")
	assert.True(t, gjson.Get(text, `a*.#[_%"John*"]#|#`).String() == "0")
	assert.True(t, gjson.Get(text, `a*.#[="123"]#|#`).String() == "1")
}

func TestParenQueries(t *testing.T) {
	text := `{
		"friends": [{"a":10},{"a":20},{"a":30},{"a":40}]
	}`
	assert.True(t, gjson.Get(text, "friends.#(a>9)#|#").Int() == 4)
	assert.True(t, gjson.Get(text, "friends.#(a>10)#|#").Int() == 3)
	assert.True(t, gjson.Get(text, "friends.#(a>40)#|#").Int() == 0)
}

func TestSubSelectors(t *testing.T) {
	text := `{
		"info": {
			"friends": [
				{
					"first": "Dale", "last": "Murphy", "kind": "Person",
					"cust1": true,
					"extra": [10,20,30],
					"details": {
						"city": "Tempe",
						"state": "Arizona"
					}
				},
				{
					"first": "Roger", "last": "Craig", "kind": "Person",
					"cust2": false,
					"extra": [40,50,60],
					"details": {
						"city": "Phoenix",
						"state": "Arizona"
					}
				}
			]
		}
	  }`
	assert.True(t, gjson.Get(text, "[]").String() == "[]")
	assert.True(t, gjson.Get(text, "{}").String() == "{}")
	res := gjson.Get(text, `{`+
		`abc:info.friends.0.first,`+
		`info.friends.1.last,`+
		`"a`+"\r"+`a":info.friends.0.kind,`+
		`"abc":info.friends.1.kind,`+
		`{123:info.friends.1.cust2},`+
		`[info.friends.#[details.city="Phoenix"]#|#]`+
		`}.@pretty.@ugly`).String()
	// println(res)
	// {"abc":"Dale","last":"Craig","\"a\ra\"":"Person","_":{"123":false},"_":[1]}
	assert.True(t, gjson.Get(res, "abc").String() == "Dale")
	assert.True(t, gjson.Get(res, "last").String() == "Craig")
	assert.True(t, gjson.Get(res, "\"a\ra\"").String() == "Person")
	assert.True(t, gjson.Get(res, "@reverse.abc").String() == "Person")
	assert.True(t, gjson.Get(res, "_.123").String() == "false")
	assert.True(t, gjson.Get(res, "@reverse._.0").String() == "1")
	assert.True(t, gjson.Get(text, "info.friends.[0.first,1.extra.0]").String() ==
		`["Dale",40]`)
	assert.True(t, gjson.Get(text, "info.friends.#.[first,extra.0]").String() ==
		`[["Dale",10],["Roger",40]]`)
}

func TestArrayCountRawOutput(t *testing.T) {
	assert.True(t, gjson.Get(`[1,2,3,4]`, "#").Raw == "4")
}

func TestParentSubQuery(t *testing.T) {
	var text = `{
		"topology": {
		  "instances": [
			{
			  "service_version": "1.2.3",
			  "service_locale": {"lang": "en"},
			  "service_roles": ["one", "two"]
			},
			{
			  "service_version": "1.2.4",
			  "service_locale": {"lang": "th"},
			  "service_roles": ["three", "four"]
			},
			{
			  "service_version": "1.2.2",
			  "service_locale": {"lang": "en"},
			  "service_roles": ["one"]
			}
		  ]
		}
	  }`
	res := gjson.Get(text, `topology.instances.#( service_roles.#(=="one"))#.service_version`)
	// should return two instances
	assert.True(t, res.String() == `["1.2.3","1.2.2"]`)
}

func TestSingleModifier(t *testing.T) {
	var data = `{"@key": "value"}`
	assert.True(t, gjson.Get(data, "@key").String() == "value")
	assert.True(t, gjson.Get(data, "\\@key").String() == "value")
}

func TestIssue141(t *testing.T) {
	text := `{"data": [{"q": 11, "w": 12}, {"q": 21, "w": 22}, {"q": 31, "w": 32} ], "sql": "some stuff here"}`
	assert.True(t, gjson.Get(text, "data.#").Int() == 3)
	assert.True(t, gjson.Get(text, "data.#.{q}|@ugly").Raw == `[{"q":11},{"q":21},{"q":31}]`)
	assert.True(t, gjson.Get(text, "data.#.q|@ugly").Raw == `[11,21,31]`)
}

func TestFlatten(t *testing.T) {
	text := `[1,[2],[3,4],[5,[6,[7]]],{"hi":"there"},8,[9]]`
	assert.True(t, gjson.Get(text, "@flatten").String() == `[1,2,3,4,5,[6,[7]],{"hi":"there"},8,9]`)
	assert.True(t, gjson.Get(text, `@flatten:{"deep":true}`).String() == `[1,2,3,4,5,6,7,{"hi":"there"},8,9]`)
	assert.True(t, gjson.Get(`{"9999":1234}`, "@flatten").String() == `{"9999":1234}`)
}

func TestJoin(t *testing.T) {
	assert.True(t, gjson.Get(`[{},{}]`, "@join").String() == `{}`)
	assert.True(t, gjson.Get(`[{"a":1},{"b":2}]`, "@join").String() == `{"a":1,"b":2}`)
	assert.True(t, gjson.Get(`[{"a":1,"b":1},{"b":2}]`, "@join").String() == `{"a":1,"b":2}`)
	assert.True(t, gjson.Get(`[{"a":1,"b":1},{"b":2},5,{"c":3}]`, "@join").String() == `{"a":1,"b":2,"c":3}`)
	assert.True(t, gjson.Get(`[{"a":1,"b":1},{"b":2},5,{"c":3}]`, `@join:{"preserve":true}`).String() == `{"a":1,"b":1,"b":2,"c":3}`)
	assert.True(t, gjson.Get(`[{"a":1,"b":1},{"b":2},5,{"c":3}]`, `@join:{"preserve":true}.b`).String() == `1`)
	assert.True(t, gjson.Get(`{"9999":1234}`, "@join").String() == `{"9999":1234}`)
}

func TestValid(t *testing.T) {
	assert.True(t, gjson.Get("[{}", "@valid").Exists() == false)
	assert.True(t, gjson.Get("[{}]", "@valid").Exists() == true)
}

// https://github.com/tidwall/gjson/issues/152
func TestJoin152(t *testing.T) {
	var text = `{
		"distance": 1374.0,
		"validFrom": "2005-11-14",
		"historical": {
		  "type": "Day",
		  "name": "last25Hours",
		  "summary": {
			"units": {
			  "temperature": "C",
			  "wind": "m/s",
			  "snow": "cm",
			  "precipitation": "mm"
			},
			"days": [
			  {
				"time": "2020-02-08",
				"hours": [
				  {
					"temperature": {
					  "min": -2.0,
					  "max": -1.6,
					  "value": -1.6
					},
					"wind": {},
					"precipitation": {},
					"humidity": {
					  "value": 92.0
					},
					"snow": {
					  "depth": 49.0
					},
					"time": "2020-02-08T16:00:00+01:00"
				  },
				  {
					"temperature": {
					  "min": -1.7,
					  "max": -1.3,
					  "value": -1.3
					},
					"wind": {},
					"precipitation": {},
					"humidity": {
					  "value": 92.0
					},
					"snow": {
					  "depth": 49.0
					},
					"time": "2020-02-08T17:00:00+01:00"
				  },
				  {
					"temperature": {
					  "min": -1.3,
					  "max": -0.9,
					  "value": -1.2
					},
					"wind": {},
					"precipitation": {},
					"humidity": {
					  "value": 91.0
					},
					"snow": {
					  "depth": 49.0
					},
					"time": "2020-02-08T18:00:00+01:00"
				  }
				]
			  },
			  {
				"time": "2020-02-09",
				"hours": [
				  {
					"temperature": {
					  "min": -1.7,
					  "max": -0.9,
					  "value": -1.5
					},
					"wind": {},
					"precipitation": {},
					"humidity": {
					  "value": 91.0
					},
					"snow": {
					  "depth": 49.0
					},
					"time": "2020-02-09T00:00:00+01:00"
				  },
				  {
					"temperature": {
					  "min": -1.5,
					  "max": 0.9,
					  "value": 0.2
					},
					"wind": {},
					"precipitation": {},
					"humidity": {
					  "value": 67.0
					},
					"snow": {
					  "depth": 49.0
					},
					"time": "2020-02-09T01:00:00+01:00"
				  }
				]
			  }
			]
		  }
		}
	  }`

	res := gjson.Get(text, "historical.summary.days.#.hours|@flatten|#.humidity.value")
	assert.True(t, res.Raw == `[92.0,92.0,91.0,91.0,67.0]`)
}

func TestSubpathsWithMultipaths(t *testing.T) {
	const text = `
[
  {"a": 1},
  {"a": 2, "values": ["a", "b", "c", "d", "e"]},
  true,
  ["a", "b", "c", "d", "e"],
  4
]
`
	assert.True(t, gjson.Get(text, `1.values.@ugly`).Raw == `["a","b","c","d","e"]`)
	assert.True(t, gjson.Get(text, `1.values.[0,3]`).Raw == `["a","d"]`)
	assert.True(t, gjson.Get(text, `3.@ugly`).Raw == `["a","b","c","d","e"]`)
	assert.True(t, gjson.Get(text, `3.[0,3]`).Raw == `["a","d"]`)
	assert.True(t, gjson.Get(text, `#.@ugly`).Raw == `[{"a":1},{"a":2,"values":["a","b","c","d","e"]},true,["a","b","c","d","e"],4]`)
	assert.True(t, gjson.Get(text, `#.[0,3]`).Raw == `[[],[],[],["a","d"],[]]`)
}

func TestFlattenRemoveNonExist(t *testing.T) {
	raw := gjson.Get("[[1],[2,[[],[3]],[4,[5],[],[[[6]]]]]]", `@flatten:{"deep":true}`).Raw
	assert.True(t, raw == "[1,2,3,4,5,6]")
}

func TestPipeEmptyArray(t *testing.T) {
	raw := gjson.Get("[]", `#(hello)#`).Raw
	assert.True(t, raw == "[]")
}

func TestEncodedQueryString(t *testing.T) {
	text := `{
		"friends": [
			{"first": "Dale", "last": "Mur\nphy", "age": 44},
			{"first": "Roger", "last": "Craig", "age": 68},
			{"first": "Jane", "last": "Murphy", "age": 47}
		]
	}`
	assert.True(t, gjson.Get(text, `friends.#(last=="Mur\nphy").age`).Int() == 44)
	assert.True(t, gjson.Get(text, `friends.#(last=="Murphy").age`).Int() == 47)
}

func TestBoolConvertQuery(t *testing.T) {
	text := `{
		"vals": [
			{ "a": 1, "b": true },
			{ "a": 2, "b": true },
			{ "a": 3, "b": false },
			{ "a": 4, "b": "0" },
			{ "a": 5, "b": 0 },
			{ "a": 6, "b": "1" },
			{ "a": 7, "b": 1 },
			{ "a": 8, "b": "true" },
			{ "a": 9, "b": false },
			{ "a": 10, "b": null },
			{ "a": 11 }
		]
	}`
	ts := gjson.Get(text, `vals.#(b==~true)#.a`).Raw
	fs := gjson.Get(text, `vals.#(b==~false)#.a`).Raw
	assert.True(t, ts == "[1,2,6,7,8]")
	assert.True(t, fs == "[3,4,5,9,10,11]")
}

func TestIndexes(t *testing.T) {
	var exampleJSON = `{
		"vals": [
			[1,66,{test: 3}],
			[4,5,[6]]
		],
		"objectArray":[
			{"first": "Dale", "age": 44},
			{"first": "Roger", "age": 68},
		]
	}`

	testCases := []struct {
		path     string
		expected []string
	}{
		{
			`vals.#.1`,
			[]string{`6`, "5"},
		},
		{
			`vals.#.2`,
			[]string{"{", "["},
		},
		{
			`objectArray.#(age>43)#.first`,
			[]string{`"`, `"`},
		},
		{
			`objectArray.@reverse.#.first`,
			nil,
		},
	}

	for _, tc := range testCases {
		r := gjson.Get(exampleJSON, tc.path)

		assert.True(t, len(r.Indexes) == len(tc.expected))

		for i, a := range r.Indexes {
			assert.True(t, string(exampleJSON[a]) == tc.expected[i])
		}
	}
}

func TestIndexesMatchesRaw(t *testing.T) {
	var exampleJSON = `{
		"objectArray":[
			{"first": "Jason", "age": 41},
			{"first": "Dale", "age": 44},
			{"first": "Roger", "age": 68},
			{"first": "Mandy", "age": 32}
		]
	}`
	r := gjson.Get(exampleJSON, `objectArray.#(age>43)#.first`)
	assert.True(t, len(r.Indexes) == 2)
	assert.True(t, gjson.Parse(exampleJSON[r.Indexes[0]:]).String() == "Dale")
	assert.True(t, gjson.Parse(exampleJSON[r.Indexes[1]:]).String() == "Roger")
	r = gjson.Get(exampleJSON, `objectArray.#(age>43)#`)
	assert.True(t, gjson.Parse(exampleJSON[r.Indexes[0]:]).Get("first").String() == "Dale")
	assert.True(t, gjson.Parse(exampleJSON[r.Indexes[1]:]).Get("first").String() == "Roger")
}

func TestIssue240(t *testing.T) {
	nonArrayData := `{"jsonrpc":"2.0","method":"subscription","params":
		{"channel":"funny_channel","data":
			{"name":"Jason","company":"good_company","number":12345}
		}
	}`
	parsed := gjson.Parse(nonArrayData)
	assert.True(t, len(parsed.Get("params.data").Array()) == 1)

	arrayData := `{"jsonrpc":"2.0","method":"subscription","params":
		{"channel":"funny_channel","data":[
			{"name":"Jason","company":"good_company","number":12345}
		]}
	}`
	parsed = gjson.Parse(arrayData)
	assert.True(t, len(parsed.Get("params.data").Array()) == 1)
}

func TestKeysValuesModifier(t *testing.T) {
	var text = `{
		"1300014": {
		  "code": "1300014",
		  "price": 59.18,
		  "symbol": "300014",
		  "update": "2020/04/15 15:59:54",
		},
		"1300015": {
		  "code": "1300015",
		  "price": 43.31,
		  "symbol": "300015",
		  "update": "2020/04/15 15:59:54",
		}
	  }`
	assert.True(t, gjson.Get(text, `@keys`).String() == `["1300014","1300015"]`)
	assert.True(t, gjson.Get(``, `@keys`).String() == `[]`)
	assert.True(t, gjson.Get(`"hello"`, `@keys`).String() == `[null]`)
	assert.True(t, gjson.Get(`[]`, `@keys`).String() == `[]`)
	assert.True(t, gjson.Get(`[1,2,3]`, `@keys`).String() == `[null,null,null]`)

	assert.True(t, gjson.Get(text, `@values.#.code`).String() == `["1300014","1300015"]`)
	assert.True(t, gjson.Get(``, `@values`).String() == `[]`)
	assert.True(t, gjson.Get(`"hello"`, `@values`).String() == `["hello"]`)
	assert.True(t, gjson.Get(`[]`, `@values`).String() == `[]`)
	assert.True(t, gjson.Get(`[1,2,3]`, `@values`).String() == `[1,2,3]`)
}

func TestNaNInf(t *testing.T) {
	text := `[+Inf,-Inf,Inf,iNF,-iNF,+iNF,NaN,nan,nAn,-0,+0]`
	raws := []string{"+Inf", "-Inf", "Inf", "iNF", "-iNF", "+iNF", "NaN", "nan",
		"nAn", "-0", "+0"}
	nums := []float64{math.Inf(+1), math.Inf(-1), math.Inf(0), math.Inf(0),
		math.Inf(-1), math.Inf(+1), math.NaN(), math.NaN(), math.NaN(),
		math.Copysign(0, -1), 0}

	assert.True(t, int(gjson.Get(text, `#`).Int()) == len(raws))
	for i := 0; i < len(raws); i++ {
		r := gjson.Get(text, fmt.Sprintf("%d", i))
		assert.True(t, r.Raw == raws[i])
		assert.True(t, r.Num == nums[i] || (math.IsNaN(r.Num) && math.IsNaN(nums[i])))
		assert.True(t, r.Type == gjson.Number)
	}

	var i int
	gjson.Parse(text).ForEach(func(_, r gjson.Result) bool {
		assert.True(t, r.Raw == raws[i])
		assert.True(t, r.Num == nums[i] || (math.IsNaN(r.Num) && math.IsNaN(nums[i])))
		assert.True(t, r.Type == gjson.Number)
		i++
		return true
	})

	// Parse should also return valid numbers
	assert.True(t, math.IsNaN(gjson.Parse("nan").Float()))
	assert.True(t, math.IsNaN(gjson.Parse("NaN").Float()))
	assert.True(t, math.IsNaN(gjson.Parse(" NaN").Float()))
	assert.True(t, math.IsInf(gjson.Parse("+inf").Float(), +1))
	assert.True(t, math.IsInf(gjson.Parse("-inf").Float(), -1))
	assert.True(t, math.IsInf(gjson.Parse("+INF").Float(), +1))
	assert.True(t, math.IsInf(gjson.Parse("-INF").Float(), -1))
	assert.True(t, math.IsInf(gjson.Parse(" +INF").Float(), +1))
	assert.True(t, math.IsInf(gjson.Parse(" -INF").Float(), -1))
}

func TestEmptyValueQuery(t *testing.T) {
	// issue: https://github.com/tidwall/gjson/issues/246
	assert.True(t, gjson.Get(
		`["ig","","tw","fb","tw","ig","tw"]`,
		`#(!="")#`).Raw ==
		`["ig","tw","fb","tw","ig","tw"]`)
	assert.True(t, gjson.Get(
		`["ig","","tw","fb","tw","ig","tw"]`,
		`#(!=)#`).Raw ==
		`["ig","tw","fb","tw","ig","tw"]`)
}

func TestParseIndex(t *testing.T) {
	assert.True(t, gjson.Parse(`{}`).Index == 0)
	assert.True(t, gjson.Parse(` {}`).Index == 1)
	assert.True(t, gjson.Parse(` []`).Index == 1)
	assert.True(t, gjson.Parse(` true`).Index == 1)
	assert.True(t, gjson.Parse(` false`).Index == 1)
	assert.True(t, gjson.Parse(` null`).Index == 1)
	assert.True(t, gjson.Parse(` +inf`).Index == 1)
	assert.True(t, gjson.Parse(` -inf`).Index == 1)
}

const readmeJSON = `
{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}
`

func TestQueryGetPath(t *testing.T) {
	assert.True(t, strings.Join(
		gjson.Get(readmeJSON, "friends.#.first").Paths(readmeJSON), " ") ==
		"friends.0.first friends.1.first friends.2.first")
	assert.True(t, strings.Join(
		gjson.Get(readmeJSON, "friends.#(last=Murphy)").Paths(readmeJSON), " ") ==
		"")
	assert.True(t, gjson.Get(readmeJSON, "friends.#(last=Murphy)").Path(readmeJSON) ==
		"friends.0")
	assert.True(t, strings.Join(
		gjson.Get(readmeJSON, "friends.#(last=Murphy)#").Paths(readmeJSON), " ") ==
		"friends.0 friends.2")
	arr := gjson.Get(readmeJSON, "friends.#.first").Array()
	for i := 0; i < len(arr); i++ {
		assert.True(t, arr[i].Path(readmeJSON) == fmt.Sprintf("friends.%d.first", i))
	}
}

func TestStaticJSON(t *testing.T) {
	text := `{
		"name": {"first": "Tom", "last": "Anderson"}
	}`
	assert.True(t, gjson.Get(text,
		`"bar"`).Raw ==
		``)
	assert.True(t, gjson.Get(text,
		`!"bar"`).Raw ==
		`"bar"`)
	assert.True(t, gjson.Get(text,
		`!{"name":{"first":"Tom"}}.{name.first}.first`).Raw ==
		`"Tom"`)
	assert.True(t, gjson.Get(text,
		`{name.last,"foo":!"bar"}`).Raw ==
		`{"last":"Anderson","foo":"bar"}`)
	assert.True(t, gjson.Get(text,
		`{name.last,"foo":!{"a":"b"},"that"}`).Raw ==
		`{"last":"Anderson","foo":{"a":"b"}}`)
	assert.True(t, gjson.Get(text,
		`{name.last,"foo":!{"c":"d"},!"that"}`).Raw ==
		`{"last":"Anderson","foo":{"c":"d"},"_":"that"}`)
	assert.True(t, gjson.Get(text,
		`[!true,!false,!null,!inf,!nan,!hello,{"name":!"andy",name.last},+inf,!["any","thing"]]`).Raw ==
		`[true,false,null,inf,nan,{"name":"andy","last":"Anderson"},["any","thing"]]`,
	)
}

func TestArrayKeys(t *testing.T) {
	N := 100
	text := "["
	for i := 0; i < N; i++ {
		if i > 0 {
			text += ","
		}
		text += fmt.Sprint(i)
	}
	text += "]"
	var i int
	gjson.Parse(text).ForEach(func(key, value gjson.Result) bool {
		assert.True(t, key.String() == fmt.Sprint(i))
		assert.True(t, key.Int() == int64(i))
		i++
		return true
	})
	assert.True(t, i == N)
}

func TestToFromStr(t *testing.T) {
	text := `{"Message":"{\"Records\":[{\"eventVersion\":\"2.1\"}]"}`
	res := gjson.Get(text, "Message.@fromstr.Records.#.eventVersion.@tostr").Raw
	assert.True(t, res == `["\"2.1\""]`)
}

func TestGroup(t *testing.T) {
	text := `{"id":["123","456","789"],"val":[2,1]}`
	res := gjson.Get(text, "@group").Raw
	assert.True(t, res == `[{"id":"123","val":2},{"id":"456","val":1},{"id":"789"}]`)

	text = `
{
	"issues": [
	  {
		"fields": {
		  "labels": [
			"milestone_1",
			"group:foo",
			"plan:a",
			"plan:b"
		  ]
		},
		"id": "123"
	  },{
		"fields": {
		  "labels": [
			"milestone_1",
			"group:foo",
			"plan:a",
			"plan"
		  ]
		},
		"id": "456"
	  }
	]
  }
  `
	res = gjson.Get(text, `{"id":issues.#.id,"plans":issues.#.fields.labels.#(%"plan:*")#|#.#}|@group|#(plans>=2)#.id`).Raw
	assert.True(t, res == `["123"]`)
}

const (
	setRaw    = 1
	setBool   = 2
	setInt    = 3
	setFloat  = 4
	setString = 5
	setDelete = 6
)

func sortJSON(json string) string {
	opts := pretty.Options{SortKeys: true}
	return string(pretty.Ugly(pretty.PrettyOptions([]byte(json), &opts)))
}

func testRaw(t *testing.T, kind int, expect, json, path string, value any) {
	t.Helper()
	expect = sortJSON(expect)
	var json2 string
	var err error
	switch kind {
	default:
		json2, err = sjson.Set(json, path, value)
	case setRaw:
		json2, err = sjson.SetRaw(json, path, value.(string))
	case setDelete:
		json2, err = sjson.Delete(json, path)
	}

	if err != nil {
		t.Fatal(err)
	}
	json2 = sortJSON(json2)
	if json2 != expect {
		t.Fatalf("expected '%v', got '%v'", expect, json2)
	}
	var json3 []byte
	switch kind {
	default:
		json3, err = sjson.SetBytes([]byte(json), path, value)
	case setRaw:
		json3, err = sjson.SetRawBytes([]byte(json), path, []byte(value.(string)))
	case setDelete:
		json3, err = sjson.DeleteBytes([]byte(json), path)
	}
	json3 = []byte(sortJSON(string(json3)))
	if err != nil {
		t.Fatal(err)
	} else if string(json3) != expect {
		t.Fatalf("expected '%v', got '%v'", expect, string(json3))
	}
}

func TestBasic11(t *testing.T) {
	testRaw(t, setRaw, `[{"hiw":"planet","hi":"world"}]`, `[{"hi":"world"}]`, "0.hiw", `"planet"`)
	testRaw(t, setRaw, `[true]`, ``, "0", `true`)
	testRaw(t, setRaw, `[null,true]`, ``, "1", `true`)
	testRaw(t, setRaw, `[1,null,true]`, `[1]`, "2", `true`)
	testRaw(t, setRaw, `[1,true,false]`, `[1,null,false]`, "1", `true`)
	testRaw(t, setRaw,
		`[1,{"hello":"when","this":[0,null,2]},false]`,
		`[1,{"hello":"when","this":[0,1,2]},false]`,
		"1.this.1", `null`)
	testRaw(t, setRaw,
		`{"a":1,"b":{"hello":"when","this":[0,null,2]},"c":false}`,
		`{"a":1,"b":{"hello":"when","this":[0,1,2]},"c":false}`,
		"b.this.1", `null`)
	testRaw(t, setRaw,
		`{"a":1,"b":{"hello":"when","this":[0,null,2,null,4]},"c":false}`,
		`{"a":1,"b":{"hello":"when","this":[0,null,2]},"c":false}`,
		"b.this.4", `4`)
	testRaw(t, setRaw,
		`{"b":{"this":[null,null,null,null,4]}}`,
		``,
		"b.this.4", `4`)
	testRaw(t, setRaw,
		`[null,{"this":[null,null,null,null,4]}]`,
		``,
		"1.this.4", `4`)
	testRaw(t, setRaw,
		`{"1":{"this":[null,null,null,null,4]}}`,
		``,
		":1.this.4", `4`)
	testRaw(t, setRaw,
		`{":1":{"this":[null,null,null,null,4]}}`,
		``,
		"\\:1.this.4", `4`)
	testRaw(t, setRaw,
		`{":\\1":{"this":[null,null,null,null,{".HI":4}]}}`,
		``,
		"\\:\\\\1.this.4.\\.HI", `4`)
	testRaw(t, setRaw,
		`{"app.token":"cde"}`,
		`{"app.token":"abc"}`,
		"app\\.token", `"cde"`)
	testRaw(t, setRaw,
		`{"b":{"this":{"😇":""}}}`,
		``,
		"b.this.😇", `""`)
	testRaw(t, setRaw,
		`[ 1,2  ,3]`,
		`  [ 1,2  ] `,
		"-1", `3`)
	testRaw(t, setInt, `[1234]`, ``, `0`, int64(1234))
	testRaw(t, setFloat, `[1234.5]`, ``, `0`, 1234.5)
	testRaw(t, setString, `["1234.5"]`, ``, `0`, "1234.5")
	testRaw(t, setBool, `[true]`, ``, `0`, true)
	testRaw(t, setBool, `[null]`, ``, `0`, nil)
	testRaw(t, setString, `{"arr":[1]}`, ``, `arr.-1`, 1)
	testRaw(t, setString, `{"a":"\\"}`, ``, `a`, "\\")
	testRaw(t, setString, `{"a":"C:\\Windows\\System32"}`, ``, `a`, `C:\Windows\System32`)
}

func TestDelete(t *testing.T) {
	testRaw(t, setDelete, `[456]`, `[123,456]`, `0`, nil)
	testRaw(t, setDelete, `[123,789]`, `[123,456,789]`, `1`, nil)
	testRaw(t, setDelete, `[123,456]`, `[123,456,789]`, `-1`, nil)
	testRaw(t, setDelete, `{"a":[123,456]}`, `{"a":[123,456,789]}`, `a.-1`, nil)
	testRaw(t, setDelete, `{"and":"another"}`, `{"this":"that","and":"another"}`, `this`, nil)
	testRaw(t, setDelete, `{"this":"that"}`, `{"this":"that","and":"another"}`, `and`, nil)
	testRaw(t, setDelete, `{}`, `{"and":"another"}`, `and`, nil)
	testRaw(t, setDelete, `{"1":"2"}`, `{"1":"2"}`, `3`, nil)
}

// TestRandomData11 is a fuzzing test that throws random data at SetRaw
// function looking for panics.
func TestRandomData11(t *testing.T) {
	var lstr string
	defer func() {
		if v := recover(); v != nil {
			println("'" + hex.EncodeToString([]byte(lstr)) + "'")
			println("'" + lstr + "'")
			panic(v)
		}
	}()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 200)
	for range 2000000 {
		n, err := r.Read(b[:rand.Int()%len(b)])
		if err != nil {
			t.Fatal(err)
		}
		lstr = string(b[:n])
		_, _ = sjson.SetRaw(lstr, "zzzz.zzzz.zzzz", "123")
	}
}

func TestDelete11(t *testing.T) {
	text := `{"country_code_from":"NZ","country_code_to":"SA","date_created":"2018-09-13T02:56:11.25783Z","date_updated":"2018-09-14T03:15:16.67356Z","disabled":false,"last_edited_by":"Developers","id":"a3e...bc454","merchant_id":"f2b...b91abf","signed_date":"2018-02-01T00:00:00Z","start_date":"2018-03-01T00:00:00Z","url":"https://www.google.com"}`
	res1 := gjson.Get(text, "date_updated")
	var err error
	text, err = sjson.Delete(text, "date_updated")
	if err != nil {
		t.Fatal(err)
	}
	res2 := gjson.Get(text, "date_updated")
	res3 := gjson.Get(text, "date_created")
	if !res1.Exists() || res2.Exists() || !res3.Exists() {
		t.Fatal("bad news")
	}

	// We change the number of characters in this to make the section of the string before the section that we want to delete a certain length

	//---------------------------
	lenBeforeToDeleteIs307AsBytes := `{"1":"","0":"012345678901234567890123456789012345678901234567890123456789012345678901234567","to_delete":"0","2":""}`

	expectedForLenBefore307AsBytes := `{"1":"","0":"012345678901234567890123456789012345678901234567890123456789012345678901234567","2":""}`
	//---------------------------

	//---------------------------
	lenBeforeToDeleteIs308AsBytes := `{"1":"","0":"0123456789012345678901234567890123456789012345678901234567890123456789012345678","to_delete":"0","2":""}`

	expectedForLenBefore308AsBytes := `{"1":"","0":"0123456789012345678901234567890123456789012345678901234567890123456789012345678","2":""}`
	//---------------------------

	//---------------------------
	lenBeforeToDeleteIs309AsBytes := `{"1":"","0":"01234567890123456789012345678901234567890123456789012345678901234567890123456","to_delete":"0","2":""}`

	expectedForLenBefore309AsBytes := `{"1":"","0":"01234567890123456789012345678901234567890123456789012345678901234567890123456","2":""}`
	//---------------------------

	var data = []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "len before \"to_delete\"... = 307",
			input:    lenBeforeToDeleteIs307AsBytes,
			expected: expectedForLenBefore307AsBytes,
		},
		{
			desc:     "len before \"to_delete\"... = 308",
			input:    lenBeforeToDeleteIs308AsBytes,
			expected: expectedForLenBefore308AsBytes,
		},
		{
			desc:     "len before \"to_delete\"... = 309",
			input:    lenBeforeToDeleteIs309AsBytes,
			expected: expectedForLenBefore309AsBytes,
		},
	}

	for i, d := range data {
		result, err := sjson.Delete(d.input, "to_delete")

		if err != nil {
			t.Error(fmtErrorf(testError{
				unexpected: "error",
				desc:       d.desc,
				i:          i,
				lenInput:   len(d.input),
				input:      d.input,
				expected:   d.expected,
				result:     result,
			}))
		}
		if result != d.expected {
			t.Error(fmtErrorf(testError{
				unexpected: "result",
				desc:       d.desc,
				i:          i,
				lenInput:   len(d.input),
				input:      d.input,
				expected:   d.expected,
				result:     result,
			}))
		}
	}
}

type testError struct {
	unexpected string
	desc       string
	i          int
	lenInput   int
	input      any
	expected   any
	result     any
}

func fmtErrorf(e testError) string {
	return fmt.Sprintf(
		"Unexpected %s:\n\t"+
			"for=%q\n\t"+
			"i=%d\n\t"+
			"len(input)=%d\n\t"+
			"input=%v\n\t"+
			"expected=%v\n\t"+
			"result=%v",
		e.unexpected, e.desc, e.i, e.lenInput, e.input, e.expected, e.result,
	)
}

func TestSetDotKey(t *testing.T) {
	text := `{"app.token":"abc"}`
	text, _ = sjson.Set(text, `app\.token`, "cde")
	if text != `{"app.token":"cde"}` {
		t.Fatalf("expected '%v', got '%v'", `{"app.token":"cde"}`, text)
	}
}

func TestDeleteDotKey2(t *testing.T) {
	bytes_ := []byte(`{"data":{"key1":"value1","key2.something":"value2"}}`)
	bytes_, _ = sjson.DeleteBytes(bytes_, `data.key2\.something`)
	if string(bytes_) != `{"data":{"key1":"value1"}}` {
		t.Fatalf("expected '%v', got '%v'", `{"data":{"key1":"value1"}}`, bytes_)
	}
}

func TestSetRaw11(t *testing.T) {
	var text = `
	{
	    "size": 1000
    }
`
	var raw = `
	{
	    "sample": "hello"
	}
`
	_ = raw
	if true {
		text, _ = sjson.SetRaw(text, "aggs", raw)
	}
	if !gjson.Valid(text) {
		t.Fatal("invalid json text")
	}
	res := gjson.Get(text, "aggs.sample").String()
	if res != "hello" {
		t.Fatal("unexpected result")
	}
}

var example = `
{
	"name": {"first": "Tom", "last": "Anderson"},
	"age":37,
	"children": ["Sara","Alex","Jack"],
	"fav.movie": "Deer Hunter",
	"friends": [
	  {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
	  {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
	  {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
	]
  }
  `

func TestIndex(t *testing.T) {
	path := `friends.#(last="Murphy").last`
	text, err := sjson.Set(example, path, "Johnson")
	if err != nil {
		t.Fatal(err)
	}
	if gjson.Get(text, "friends.#.last").String() != `["Johnson","Craig","Murphy"]` {
		t.Fatal("mismatch")
	}
}

func TestIndexes11(t *testing.T) {
	path := `friends.#(last="Murphy")#.last`
	text, err := sjson.Set(example, path, "Johnson")
	if err != nil {
		t.Fatal(err)
	}
	if gjson.Get(text, "friends.#.last").String() != `["Johnson","Craig","Johnson"]` {
		t.Fatal("mismatch")
	}
}
