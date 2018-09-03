# PushDown Automata Generator

This is a first example of a program that automatically generates
parsers on the fly.  What "on the fly" means is that the program is
given an input string and a Context-Free Grammar, and the output will
be parsed and structured.

## Note: Project Status

pushed down to the bottom of the stack of priorities. 

## Purpose

Designing and Parsing strings should be simple.  Domain-specific
languages should be casual and easy to use, no matter what they are.
Heavy-weight or confusing tools are hard to get started with, and/or
hard to use and share.


## Goals

- simple and clear definitions.  Anybody who has seen context-free
  grammars, programming language specs, or encoding specs should know
  immediatelyp how to use the tool.

- helpful error messages.

- no bells and whistles.

- made practical for converting arbitrary encodings into structured
  XML, JSON, or another easy-to-use encoding.




## How it works:

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


## Current status

Currently, things are hard-coded just because it's a first attempt.

- input alphabet **Ʃ={0,1,2,3}**
- Variables **V={A,B}**.
- The stack alphabet equivalent to the union(input alphabet, variable
  alphabet). 
- '$' is the initial variable.


Context-Free Grammar:
~~~
$ →  AB
A →  0A1 ∣ 2
B →  1B ∣ 3A
~~~

Grammar (in code form):
~~~go
var rulemap = map[rune]([]string){
	'$': []string{"AB"},
	'A': []string{"0A1", "2"},
	'B': []string{"1B", "3A"},
}
~~~

Input String:
~~~
021300211
~~~

Internals of how it builds up the stack in order to parse the input string:
~~~
- '$' ⊢ AB 
0 '0' ⊢ A1B 
1 '2' ⊢ 1B 
2 '1' ⊢ B 
3 '3' ⊢ A 
4 '0' ⊢ A1 
5 '0' ⊢ A11 
6 '2' ⊢ 11 
7 '1' ⊢ 1 
8 '1' ⊢  
string accepted!
~~~


Output (**TODO**):
~~~xml
<todo>
	<todo2>
		<term> organized XML output goes here! </term>
	</todo2>
</todo>
~~~
