---
layout: default
---

# Buggy
<img align="bottom" src="./assets/img/BuggySayingHello.png" alt="drawing" width="100">

***Programming language and interpreter***

## Content

1. [Content](#content)
2. [Introduction](#introduction)
3. [Interpreter](#interpreter)
   1. [Syntax](#syntax)
      * [Variable definition](#variable-definition)
      * [Function definition and calling](#function-definition-and-calling)
      * [If-else-statements](#if-else-statements)
      * [Array definition](#array-definition)
      * [Hash definition](#hash-definition)
   2. [Built-in functions](#built-in-functions)
4. [FAQ](#faq)
5. [Errors](#errors)

### Introduction
Buggy is a programming language with a self-written on Golang interpreter for it. Here is current list of supported features of Buggy:
* C-like syntax  
* variable bindings
* integers and booleans
* arithmetic expressions
* built-in functions
* first-class and higher-order functions
* closures
* a hash(dictionary) data structure
* an array data structure
* a string data structure

### Interpreter
#### Syntax
##### Variable definition
Integer
```
>>let foo = 50
>>let bar = -5
>>foo
50
>>bar
-5
```
---
Booleans
```
>>let foo = true
>>let bar = false
>>foo
true
>>bar
false
```
---
String
```
>>let foo = "1bc"
>>foo
1bc
```
---
##### Function definition and calling
Function with no arguments
```
>>let foo = fn() { 5 }
>>foo()
5
```
Function with one argument
```
>>let foo = fn(a) { a }
>>foo(2)
2
```
Function with several arguments
```
>>let foo = fn(a, b, c) { a + b + c }
>>foo(1, 2, 3)
6
```
Function with function as argument
```
>>let double = fn(x) { x * 2 }
>>let cover = fn(a, f) { f(a) }
>>cover(2, double)
4
```
Recursive call of function
```
>>let double = fn(x) { x * 2 }
>>let cover = fn(a, f) { f(a) }
>>cover(2, double)
4
```
---
##### If-else statements
General structure

```
if (expression){  
  FirstBlockOfStatements
} else {
  SecondBlockOfStatements
}
```
>`expression` should take value `true`(or any different from `null` and `false` value) to evaluate `FirstBlockOfStatements` and should take value `false` or `null` to evaluate `SecondBlockOfStatements`

Example
```
>>if (true) { 5 }
5
>>if (false) { 5 } else { 3 }
3
>>if (1) { 5 } else { 3 }
5
```
---
##### Array definition

```
>>let foo = [1, "a", fn() { 5 }]
>>foo
[1, a, fn() {
  5
}]
```

##### Hash definition
```
>>let foo = {1: 2, 2: "a", 3: fn() { 5 }}
>>foo
{1: 2: a, 3: fn() {
  5
}}
```


#### Built-in functions
help()
> Prints out help with list of built-in functions and link here.


```
>>help()
```
---
~~help(object.Builtin)~~. [^*]
> Prints out built-in function description.


```
>>help(object)
```
---
len(s)
> If `s` is `object.String`, returns length of string (number of letters in string).
> If `s` is `object.Array`, returns number of elements in array.


```
>>let foo = [1, 2, 3]
>>len(foo)
3
>>let bar = "Buggy"
>>len(bar)
5
```
---
first(array)
> Returns first element of `array`, if `array` is empty returns `NULL`

```
>>let foo = [1, 2]
>>first(foo)
1
>>first([])
null
```
---
last(array)
> Returns last element of `array`, if `array` is empty returns `NULL`

```
>>let foo = [1, 2]
>>last(foo)
2
>>last([])
null
```
---
rest(array)
> Returns new array object with all elements of previous `array` except first. If `array` is empty returns `NULL`

```
>>let foo = [1, 2]
>>rest(foo)
[2]
>>rest(rest(foo))
[]
>>rest(rest(rest(foo)))
null
```
---
push(array, newVal)
> Returns new array object with all elements of previous `array` and `newVal` added to the end.

```
>>let foo = [1, 2]
>>push(foo, 3)
[1, 2, 3]
>>push([], 3)
[3]
```
---
say(args)
> Prints string representation of `args` (one `arg` by the line) and returns `null`.
>
```
>>say("hello")
hello
null
>>let foo = fn(x) { x };
>>say(foo)
fn(x){
  x
}
null
```
### Examples

### Errors
<img src="./assets/img/DeadBuggy.png" alt="drawing" width="50"/>

#### Parser errors
```
>>let a = 9223372036854775808
...
could not parse 9223372036854775808 as Integer
```
---
```
>>let a = (5
...
expected peek type was = ) got = EOF instead
```
---
```
>>}
...
no prefix parse function for } found
```

#### Evaluator errors
```
>>let a = -"a"
ERROR: unknown operator: -STRING
```
---
```
>>let a = "ab" / 2
ERROR: unknown operator: STRING/INTEGER
```
---
```
>>let a = "ab" - "ab"
ERROR: unknown operator: STRING-STRING
```
---
```
>>let a = [1] * [2]
ERROR: unknown operator: ARRAY * ARRAY
```
---
```
>>let a = [1] * 2
ERROR: type mismatch: ARRAY * INTEGER
```
---
```
>>b   
ERROR: identifier not found: b
```
---
```
>>5()
ERROR: not a function: INTEGER
```
---
```
>>let a = {1:"a"}
>>a[a]
ERROR: unusable as hash key: HASH
```
---
```
>>let a = [1]
>>let b = {a:2}
ERROR: unusable as hash key: ARRAY
```
---
```
>>let a = [1]
>>a[{1:2}]
ERROR: index operator not supported: HASH
```
#### Built-in errors
There are only two types of errors for built-in functions:
* `wrong number of arguments. got = #, want <OR=#`
* `argument to BUILTIN not supported, got ARGUMENT.TYPE`


### FAQ

[^*]: Under development.
