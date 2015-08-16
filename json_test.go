package gopar

import (
	"testing"
)

func TestJson(t *testing.T) {
	digit := OneOfChars("0123456789").Rename("Digit")

	char := OneOfChars(" \t\nabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*()_+`-={}|[]\\:;'<>?,./'").Rename("Char")

	number := Seq(
		OneOrMoreOf(digit),
		ZeroOrOneOf(Seq(
			S("."),
			OneOrMoreOf(digit),
		)),
	).Rename("Number")

	str := Seq(
		S("\""),
		OneOrMoreOf(char),
		S("\""),
	).Rename("JsonString")

	value := OneOf(
		str,
		number,
		P("Object"),
		P("List"),
	).Rename("Value")

	list := Seq(
		S("["),
		ZeroOrOneOf(Seq(
			value,
			ZeroOrMoreOf(Seq(
				S(","),
				value,
			)),
		)),
		S("]"),
	).Rename("List")

	keyVal := Seq(
		str,
		S(":"),
		value,
	).Rename("KeyValue")

	object := Seq(
		S("{"),
		ZeroOrOneOf(Seq(
			keyVal,
			ZeroOrMoreOf(Seq(
				S(","),
				keyVal,
			)),
		)),
		S("}"),
	).Rename("Object")

	err := Patch(object,list)
	if err != nil {
		t.Fatal(err)
	}

	expectNoErr(t, digit, "1")
	expectErr(t, digit, "a", "error at offset 0 in rule Digit>'0'. expected '0' found 'a'")
	expectNoErr(t, number, "1")
	expectNoErr(t, number, "12")
	expectNoErr(t, number, "12.3")
	expectNoErr(t, str, "\"complex string contents 1234.2345 !@#$%^&*(}{::?\"")
	expectErr(t, str, "\"missing quote", "error at offset 14 in rule JsonString>'\"'. EOF")
	expectNoErr(t, keyVal, "\"apple\":12.32")
	expectNoErr(t, keyVal, "\"apple\":\"banana\"")
	expectErr(t, keyVal, `3.2:"banana"`, "error at offset 0 in rule KeyValue>JsonString>'\"'. expected '\"' found '3'")
	expectNoErr(t, list, "[]")
	expectNoErr(t, list, "[1]")
	expectNoErr(t, list, `["tree"]`)
	expectNoErr(t, list, `["tree",1.2]`)
	expectNoErr(t, object, `{"apple":"red","banana":1.2}`)
	
	//The big test: nested lists and dicts some which are empty
	expectNoErr(t, object, `{"apple":"red","banana":[1,2],"coconut":{"a":1,"b":[],"c":{}}}`)
	expectErr(t, object,   `{"apple":"red","banana":[1,2],"coconut":{"a":1,"b":[,"c":{}}}`,
		"error at offset 29 in rule Object>'}'. expected '}' found ','")
}
