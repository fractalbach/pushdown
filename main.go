/*
Pushdown

Input: (string, grammar), Output: (xml).  The grammar is converted
into a pushdown automata, which is used to parse the input string.
The output will be in XML format.


Extended Backus Naur Form

The grammar should be in EBNF:
https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form

*/
package main

import (
	"fmt"
	//"regex"
)

// tokenKind is used to quickly identify the token, and decide how it
// should be processed.  It's like a "token descriptor".
type tokenKind int

const (
	Terminal tokenKind = iota
	Concat
	Union
	Variable
	EndVariable
)

// the data string of a token will depend on the kind of token it is.
// For example, Terminals will use the data string literally, and
// variables will use data string for the name.
type token struct {
	kind         tokenKind
	list         []*token
	data, output string
}

func (t *token) String() string {
	return t.data
}

// Term creates a new terminal token.
func Term(data string) *token {
	return &token{
		kind: Terminal,
		data: data,
	}
}

// And creates a concatenation of multiple tokens.
func And(tokens ...*token) *token {
	return &token{
		kind: Concat,
		list: tokens,
	}
}

// Or creates a union of multiple tokens.
func Or(tokens ...*token) *token {
	return &token{
		kind: Union,
		list: tokens,
	}
}

// Var creates a new variable token.  The name should match one of the
// rules in the defined grammar.
func Var(name string) *token {
	return &token{
		kind:   Variable,
		data:   name,
		output: "<" + name + ">",
	}
}

var (
	stack              []*token
	exampleInputString = "x x .(x x.)"
	ex2                = "021300211"
)

var exampleGrammarEBNF = `
$ = 'a', 'b';
varA = '0', varA, '1'  |  '2';
varB = '1', varB | '3', varA;
`
var exMap = map[string](*token){
	"$": And(Var("φ")),
	"φ": Or(
		And(
			Or(Term("("), Term("["), Term("{")),
			Var("φ"),
			Var("3"),
		),
		And(Term("x"), Var("ζ")),
		And(Var("ε"), Var("φ"), Var("3")),
		And(Var("λ"), Var("3")),
	),
	"ε": Or(
		Term("("), Term("["), Term("{"),
	),
	"3": Or(
		Term(")"), Term("]"), Term("}"),
	),
	"λ": Or(
		And(Term("x"), Var("ζ")),
	),
	"ζ": Or(
		Term(";"),
		Term("."),
		And(Term(" "), Var("φ")),
	),
}
var exMap2 = map[string](*token){
	"$": And(Var("a"), Var("b")),
	"a": Or(
		And(Term("0"), Var("a"), Term("1")),
		Term("2"),
	),
	"b": Or(
		And(Term("1"), Var("b")),
		And(Term("3"), Var("a")),
	),
}
var exMap3 = map[string](*token){
	"$":    Var("Rule"),
	"rule": And(Var("lhs"), Term("="), Var("rhs"), Term(";")),
	"lhs":  Var("Identifier"),
	"rhs": Or(
		Var("Identifier"),
		And(Term("["), Var("rhs"), Term("]")),
		And(Term("{"), Var("rhs"), Term("}")),
		And(Term("("), Var("rhs"), Term(")")),
		And(Var("rhs"), Term("|"), Var("rhs")), // note: potential bug here.
		And(Var("rhs"), Term(","), Var("rhs")), // <- might not be reachable.
	),
}

/*
special sequences are used for regular expressions
`[a-z]|[a-Z]`
regex.Match(`[a-z]+\`)
*/

// var rulemap = map[rune]([]string){
// 	'$': []string{"AB"},
// 	'A': []string{"0A1", "2"},
// 	'B': []string{"1B", "3A"},
// }

func init() {
	// initialize by pushing the starting production.
	toks, _ := parseToken("$", exMap["$"])
	push(toks...)
}

// pop() splits the stack into two seperate strings, with one of those
// strings containing a single character (the last character).  A type
// conversion needs to be done in order to return a type rune.
func pop() (*token, error) {
	if len(stack) <= 0 {
		return nil, fmt.Errorf("stack empty.")
	}
	L := len(stack) - 1
	t := stack[L]
	stack = stack[:L]
	return t, nil
}

func push(t ...*token) {
	t = reverse(t)
	stack = append(stack, t...)
}

// In formal theory, a transition takes the 3-tuple (state,
// inputSymbol, stackSymbol), but in this program, we are only going
// to use the 2-tuple (inputSymbol, stackSymbol).  Similarily, the
// result will omit the state, and additionally will include a
// bubble-up error message.
//
// type transition func(string, *token) ([]string, error)

var functionMap = map[tokenKind](func(string, *token) ([]*token, error)){
	Terminal: parseTerminal,
	Concat:   parseConcat,
	Union:    parseUnion,
}

// proccess always pops a symbol from the stack.  This token is
// examined alongside the input symbol to decide what kind of
// transition will happen.  If the transition results in values that
// need to pushed to the stack, this will push them.  It's like a
// "stack manager" function.
func process(a string) error {
	X, err := pop()
	if err != nil {
		return err
	}

	// check for the special stack token "EndVar", which only
	// exists in the stack language. This transition is unlike
	// others, because it does NOT consume an input token.  To
	// achieve this affect, call "process" again, using the same
	// input.
	switch X.kind {

	case EndVariable:
		fmt.Println(X.data)
		return process(a)

	case Variable:
		X = exMap[X.data]

	}
	//δ := functionMap[X.kind]
	// results, err := δ(a, X)

	results, err := parseToken(a, X)
	if err != nil {
		return err
	}
	push(results...)

	return nil
}

func parseToken(s string, t *token) ([]*token, error) {

	var δ (func(string, *token) ([]*token, error))

	switch t.kind {
	case Terminal:
		δ = parseTerminal
	case Union:
		δ = parseUnion
	case Concat:
		δ = parseConcat
	case Variable:
		δ = parseVariable
	default:
		panic(fmt.Sprint("unknown token kind(%v", t.kind))
	}

	return δ(s, t)
}

func parseTerminal(s string, t *token) ([]*token, error) {
	if s != t.data {
		// return fmt.Errorf(
		// 	"Invalid symbol. expected(%s), got:(%s)",
		// 	input.data, t.data,
		// )
		return nil, fmt.Errorf("invalid symbol")
	}
	return nil, nil
}

// In a concat, we only need to match the first token.  If the match
// is successful, the the remaning tokens are pushed to the stack.
func parseConcat(s string, t *token) ([]*token, error) {
	// beware of infinite loops caused by the grammar definition.
	result, err := parseToken(s, t.list[0])
	if err != nil {
		return nil, err
	}
	if len(t.list) > 1 {
		result = append(result, t.list[1:]...)
		return result, nil
	}
	return result, nil
}

// parsing a union looks through all of the possibilties, and does not
// push any symbols.  Returns an error only if ALL other parses return
// an error.
func parseUnion(s string, t *token) ([]*token, error) {
	for _, tok := range t.list {
		result, err := parseToken(s, tok)
		if err == nil {
			return result, nil
		}
	}
	// No matches found?! We aren't in the right context!
	return nil, fmt.Errorf("symbol(%s) not expected in token:(%v) ", s, t)
}

func parseVariable(s string, t *token) ([]*token, error) {
	return []*token{t}, nil
}

// prints the reverse of the stack; used for displaying output of the
// stack in a pushdown automata style.
func reverse(a []*token) []*token {
	var rev []*token
	for _, t := range a {
		rev = append([]*token{t}, rev...)
	}
	return rev
}

func main() {
	// ⊢
	fmt.Printf("- %q  %s \n", '$', reverse(stack))
	for i, c := range exampleInputString {
		err := process(string(c))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%d %q  %s \n", i, c, reverse(stack))
	}
	if len(stack) == 0 {
		fmt.Println("string accepted!")
	} else {
		fmt.Println("Unexpected End: stack should be empty.")
	}
}
