# Gopar - "The go parser that needs a better name"™

So lookey here: I'm using gopar to build up a symple 
```
	digit := OneOfChars("0123456789").Rename("Digit")

	char := OneOfChars(" \t\nabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*()_+`-={}|[]\\:;'<>?,./'").Rename("Char")

	// numbers are a sequence of digits optionally with a '.' and then some
	// more digits
	number := Seq(
		OneOrMoreOf(digit),
		ZeroOrOneOf(Seq(
			S("."),
			OneOrMoreOf(digit),
		)),
	).Rename("Number")

	// strings are a bunch of characters surrounded by " (I'm lazy, I did
	// include ' strings and quote characters in strings)
	str := Seq(
		S("\""),
		OneOrMoreOf(char),
		S("\""),
	).Rename("JsonString")

	// values can be strings or numbers or Objects or Lists ... hey wait,
	// we haven't defined Objects or Lists yet. No problem, `P(string)`
	// creates a placeholderRule that will later be patched with the rule it
	// names
	value := OneOf(
		str,
		number,
		P("Object"),
		P("List"),
	).Rename("Value")

	// a list is, well, a list of values
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

	// keyVal has a string key and a value val
	keyVal := Seq(
		str,
		S(":"),
		value,
	).Rename("KeyValue")

	// an object is a bunch of keyVal pairs
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

	// an object is a bunch of keyVal pairs
	err := Patch(object,list)
	if err != nil {
		t.Fatal(err)
	}
		
	//The big test: nested lists and dicts some which are empty
	expectNoErr(t, object, `{"apple":"red","banana":[1,2],"coconut":{"a":1,"b":[],"c":{}}}`)
	expectErr(t, object, `{"apple":"red","banana":[1,2],"coconut":{"a":1,"b":[],"c":{}}`,
		"error at offset 61 in rule Object>'}'. EOF")
```


# Notes for later

TODO: tsbr

* The current implementation will be really slow because data is fetched as it
is needed. Instead prefetch some user configured amount - maybe a block size.
* The current implementation is computationally intense because for every read
it readjusts the buffer to include only the range from min(subreader_offset) 
to max(subreader_offset) - fix this by only adjusting array size when you have
to fetch the next block

TODO: parser

* There should be a delimited sequence rule. Because I commonly use this pattern
```
Seq(
	Str("["),
	ZeroOrOne(
		Value(),
		ZeroOrMore(
			Str(","),
			Value(),
		),
	),
	Str("]"),
)
```
With DelimitedSeq this would become simply
```
Seq(
	Str("["),
	DelimitedSeq(Val(),Str(","))
	Str("]"),
)
```
but be cautious that this will possible cause parsing ambiguity - should the
elements of the array have affinity to the left or right side?
* Many times whitespace is ignored - but making optional whitespace rules 
everywhere is annoying. Is the best place to handle this in the String rule by
ignoring prefixed whitespace.
* When an error occurs for in StringRule for unicode, print the unicode 
characters
* There is a bug in reporting error offset for rules with ZeroOrMore or 
ZeroOrOne because it is perfectly valid for these not to match and therefore
the error character gets reported right before these rules start. But, if the 
entire parse was an error, then it's probably better to report the furthest
recorded error
* Patch is a little sloppy, if probably patches and re-patches the same 
placeholderRule multiple times. Additionally, since I'm specifying all the
patchRules, I don't need to crawl the syntax tree to build the rule lookup.
