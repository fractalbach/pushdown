/*
CFG to PushDown Automata

This is a first example of a program that automatically generates
parsers on the fly.  What "on the fly" means is that the program is
given an string and a Context-Free Grammar, and the parse tree will be
output.


The Goal:

Designing and Parsing strings should be simple.  Domain-specific
languages should be casual and easy to use, no matter what they are.
Heavy-weight or confusing tools are hard to get started with, and/or
hard to use and share.

Thus, the goal of this tool is to make the creation of a parser very
simple, without the bells and whistles.  If you can define the
encoding with a context-free grammar, then its structure can be
analyzed using this tool.  Major goals include: making the definitions
simple, and make the error messages helpful.


Currently:

Hard-coded example with an input alphabet of Ʃ={0,1,2,3}, and variable
alphabet V={A,B}.  The stack alphabet equivalent to the union(input
alphabet, variable alphabet).  Basically, the stack can contain both
variables and terminal symbols, which represent the location in the
parse tree.


How it works:

For each input character scanned, a stack character is popped from the
top of the stack.  If the stack is empty, the program
ends.  Otherwise, that stack character is analyzed further.

If the popped character and the input character are identical, nothing
is pushed to the stack, and the scanner procedes.  This usually
happens when we are inside a variable already, and are fufilling its
terminals.  An example would be the end parenthesis ')' of some
expression.

If that popped character is a variable, then the input character must
be immediately derived by that variable.


Later:

Additional stack characters will be added, allowing for parse trees to
be more easily formed.  Output will be in XML, JSON, or some other
familar encoding.

 */

package main

import (
	"fmt"
)

var (
	stack   = "" // using string for simplicity
	exInput = "021300211"
)

var rulemap = map[rune]([]string){
	'$': []string{"AB"},
	'A': []string{"0A1", "2"},
	'B': []string{"1B", "3A"},
}

func init() {
	// initialize by pushing the starting production.
	pushString(rulemap['$'][0])
}

// pop() splits the stack into two seperate strings, with one of those
// strings containing a single character (the last character).  A type conversion needs to be done in order to return a type rune.
func pop() (rune, error) {
	if len(stack) <= 0 {
		return '\b', fmt.Errorf("stack empty.")
	}
	lastChar := rune(stack[len(stack)-1])
	stack = stack[:(len(stack) - 1)]
	return lastChar, nil
}

func push(r rune) {
	stack += string(r)
}

func pushString(s string) {
	for i := len(s) - 1; i >= 0; i-- {
		stack += string(s[i])
	}
}

func transition(r rune) error {
	v, err := pop()
	if err != nil {
		return err
	}

	// assuming the stack was able to pop a symbol, check to see
	// what variable the symbol is.  if the popped symbol is a
	// terminal symbol, then we can simply check for equality.
	list, exists := rulemap[v]
	if !exists {
		if v == r {
			return nil
		}
		return fmt.Errorf(
			"Invalid symbol. expected(%q), got:(%q)", v, r)
	}

	// grab the first character from each of the possibilities in
	// this production.  The list is an array of strings, and we
	// only want to look at the first character.
	for i := range list {

		// If we find a match, then the input character is
		// valid! Hurray!  It matches the variable, and we can
		// continue.
		if rune(list[i][0]) == r {
			if len(list[i]) > 1 {
				pushString(list[i][1:])
			}
			return nil
		}
	}

	// No matches found?! We aren't in the right context!
	return fmt.Errorf("variable(%q) wasn't expected symbol(%q).", v, r)
}

// prints the reverse of the stack; used for displaying output of the
// stack in a pushdown automata style.
func stackRev() string {
	s := make([]byte, len(stack))
	for i := len(stack) - 1; i >= 0; i-- {
		s[len(stack)-1-i] += (stack[i])
	}
	return string(s)
}

func main() {
	fmt.Printf("- %q ⊢ %s \n", '$', stackRev())
	for i, c := range exInput {
		err := transition(c)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%d %q ⊢ %s \n", i, exInput[i], stackRev())
	}
	if len(stack) == 0 {
		fmt.Println("string accepted!")
	} else {
		fmt.Println("Unexpected End: stack should be empty.")
	}
}
