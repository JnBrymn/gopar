TODO: tsbr

* The current implementation will be really slow because data is fetched as it
is needed. Instead prefetch some user configured amount - maybe a block size.
* The current implementation is computationally intense because for every read
it readjusts the buffer to include only the range from min(subreader_offset) 
to max(subreader_offset) - fix this by only adjusting array size when you have
to fetch the next block

TODO: parser

* hide all the structs and only expose functions like so that we can write code
like

```
Seq(
	Str("eggs and"),
	OneOf(
		Str("cheese"),
		Str("bacon"),
	),
)
```
* The methods OneOrMore, ZeroOrMore, ZeroOrOne are all backed by a Rule for
"at least X and no more than Y"
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
* Make rule replacement strategy used to define Rules that contain rules that
haven't been defined yet - including themselves.
```
thing = Seq(
	Bla1(),
	Bla2(),
	Placeholder("thing"),
	Bla3(),
)
thing.Patch("thing",thing) //this would place the concrete rule in the Placeholder
```
Placeholder.Parse is simply Panic with the placeholder name. (See below for
alternative and better idea to this Patch.)
* Make rules have optional names and descriptions - if name is not specified 
default to the Struct name, if description not specified default to ""
* If we use the previous bullet then we can make patch take no args - it would
sweep through the model collecting rule names, making sure there are no dupes
and then patching in the rules. Or even better yet, upon the first use of the 
Rule, patch would run that once (but I don't like having to keep up with an 
"is_patched" member variable).
* Many times whitespace is ignored - but making optional whitespace rules 
everywhere is annoying. Is the best place to handle this in the String rule by
ignoring prefixed whitespace.
* OneOf rule can be done in parallel waiting on the first matching rule (e.g. 
if the first rule matches then it doesn't hve to wait for rule 2 to complete,
but if the frst rule doesn't match then it has to wait for rule 2, etc.). 
Making subrules stop working and close up readers might be hard - on the other
hand maybe we can just patch the reader to send EOF to all children readers
* When an error occurs for in StringRule for unicode, print the unicode 
characters
