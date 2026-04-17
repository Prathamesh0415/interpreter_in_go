# Interpreter in Go — Complete Study Guide

---

## Table of Contents

1. [What Is This Project? (Overview)](#1-what-is-this-project-overview)
2. [Complete Flow — How It Works Step by Step](#2-complete-flow--how-it-works-step-by-step)
3. [Important Files — Purpose & Internals](#3-important-files--purpose--internals)

---

## 1. What Is This Project? (Overview)

### What it is

This project is a **custom programming language interpreter** written in Go. It implements a small but complete scripting language called **Monkey** (based on the famous book *"Writing An Interpreter In Go"* by Thorsten Ball).

### What does "interpreter" mean?

There are two common ways to run code you write:

| Approach | How it works | Examples |
|---|---|---|
| **Compiler** | Translates your code to machine code first, then runs it | C, C++, Go itself |
| **Interpreter** | Reads your code and executes it line by line, live | Python, Ruby, JavaScript (Node) |

This project is an **interpreter** — you type code, it runs immediately.

### What can the language do?

```
>> let x = 5;
>> let y = 10;
>> let add = fn(a, b) { return a + b; };
>> add(x, y)
15
>> if (x < y) { x } else { y }
5
>> let result = add(3, 4) * 2;
>> result
14
```

**Supported features:**
- Integer arithmetic: `+`, `-`, `*`, `/`
- Boolean values: `true`, `false`
- Comparison: `==`, `!=`, `<`, `>`
- Prefix operators: `-5`, `!true`
- Variable binding: `let x = 10;`
- If/else expressions
- Functions: `fn(x, y) { x + y }`
- Function calls: `add(1, 2)`
- Return statements inside functions
- Closures (functions that remember their outer scope)

### Real-world use cases of interpreters like this

- **Python, Ruby, PHP** all work similarly internally
- **Shell scripts** (bash) are interpreted
- **JavaScript engines** in browsers interpret JS
- **Database query languages** (like SQL) parse and evaluate expressions
- **Game scripting** — games use embedded interpreters for game logic
- **Template engines** — Django, Jinja2 parse and evaluate template syntax

---

## 2. Complete Flow — How It Works Step by Step

When you type `let x = 5 + 3;` and hit Enter, here is exactly what happens inside the program:

```
Input string
    │
    ▼
┌─────────────┐
│    LEXER    │  Breaks text into tokens: [LET] [x] [=] [5] [+] [3] [;]
└──────┬──────┘
       │  stream of tokens
       ▼
┌─────────────┐
│   PARSER    │  Builds a tree structure (AST) from tokens
└──────┬──────┘
       │  Abstract Syntax Tree
       ▼
┌─────────────┐
│  EVALUATOR  │  Walks the tree, computes values, returns result
└──────┬──────┘
       │  Object (result)
       ▼
┌─────────────┐
│    REPL     │  Prints the result back to you: 8
└─────────────┘
```

### Step-by-step walkthrough: `let x = 5 + 3;`

#### Step 1 — Lexer turns text into tokens

The **Lexer** reads the string character by character and produces a stream of **Tokens**.

```
Input:  "let x = 5 + 3;"

Output tokens:
  Token{Type: LET,       Literal: "let"}
  Token{Type: IDENT,     Literal: "x"}
  Token{Type: ASSIGN,    Literal: "="}
  Token{Type: INT,       Literal: "5"}
  Token{Type: PLUS,      Literal: "+"}
  Token{Type: INT,       Literal: "3"}
  Token{Type: SEMICOLON, Literal: ";"}
  Token{Type: EOF,       Literal: ""}
```

Think of it like splitting a sentence into individual words.

#### Step 2 — Parser builds the AST

The **Parser** takes those tokens and builds a **tree** that represents the structure/meaning of the code.

```
For "let x = 5 + 3;"

LetStatement
├── Name: Identifier("x")
└── Value: InfixExpression
             ├── Operator: "+"
             ├── Left:  IntegerLiteral(5)
             └── Right: IntegerLiteral(3)
```

This tree is called an **Abstract Syntax Tree (AST)**. "Abstract" means it captures the *meaning*, not the exact text.

#### Step 3 — Evaluator walks the AST

The **Evaluator** does a recursive walk of the tree:

```
Eval(LetStatement)
  → Eval(InfixExpression "+")
      → Eval(IntegerLiteral 5) → Integer{5}
      → Eval(IntegerLiteral 3) → Integer{3}
      → 5 + 3 = 8
  → env.Set("x", Integer{8})
  → stores x = 8 in the environment
```

#### Step 4 — REPL shows result

The REPL (Read-Eval-Print Loop) prints the result.

```
>> let x = 5 + 3;
8
```

---

### Flow for a function call: `add(2, 3)`

Assuming `let add = fn(a, b) { return a + b; };` was already defined:

```
Eval(CallExpression)
  → Eval("add") → looks up in env → returns Function object
  → Eval arguments: [2, 3] → [Integer{2}, Integer{3}]
  → applyFunction(Function, [2, 3])
      → create NEW enclosed environment
          (outer = current env, inner = {a:2, b:3})
      → Eval(body: "return a + b;")
          → Eval(ReturnStatement)
              → Eval(InfixExpression "+")
                  → Eval(a) → looks in inner env → Integer{2}
                  → Eval(b) → looks in inner env → Integer{3}
                  → 2 + 3 = Integer{5}
              → wrap in ReturnValue{Integer{5}}
      → unwrap ReturnValue → Integer{5}
  → return Integer{5}
```

---

## 3. Important Files — Purpose & Internals

### `token/token.go` — The Token Vocabulary

**Purpose:** Defines every possible "word" (token) in the language.

```go
type TokenType string   // just a string like "INT", "LET", "+"

type Token struct {
    Type    TokenType   // what kind of thing it is
    Literal string      // the actual text, e.g., "42" or "myVar"
}
```

**The keyword map:**
```go
var keyword = map[string]TokenType {
    "fn":     FUNCTION,
    "let":    LET,
    "true":   TRUE,
    "false":  FALSE,
    "if":     IF,
    "else":   ELSE,
    "return": RETURN,
}
```

`LookupIdent( )` checks if an identifier is a keyword or a user-defined variable name. When the lexer reads "let", it calls `LookupIdent("let")` and gets back `LET` token type, not `IDENT`.

---

### `lexer/lexer.go` — The Tokenizer

**Purpose:** Takes raw source code (a string) and produces tokens one at a time.

**The struct:**
```go
type Lexer struct {
    input        string  // full source code
    position     int     // current character index
    readPosition int     // next character index (one ahead)
    ch           byte    // current character being examined
}
```

**Two-pointer trick:** The lexer keeps TWO positions:
- `position` — the character you are currently looking at
- `readPosition` — the character you are about to look at next

This allows **peeking ahead** without consuming the character.

**Why peek?** For two-character tokens like `==` and `!=`:
```go
case '=':
    if l.peekChar() == '=' {   // peek without advancing
        ch := l.ch             // save '='
        l.readChar()           // now advance to second '='
        tok = Token{Type: EQ, Literal: "=="}
    } else {
        tok = newToken(ASSIGN, l.ch)   // just single '='
    }
```

**Key functions:**

| Function | What it does |
|---|---|
| `New(input)` | Creates lexer, calls `readChar()` to load first char |
| `NextToken()` | Returns the next token, advances position |
| `readChar()` | Moves to next character |
| `peekChar()` | Looks at next character WITHOUT moving |
| `readIdentifier()` | Reads a whole word like `let` or `myVar` |
| `readNumber()` | Reads a whole number like `123` |
| `skipWhitespace()` | Skips spaces, tabs, newlines |

---

### `ast/ast.go` — The Abstract Syntax Tree

**Purpose:** Defines all the node types that can appear in the syntax tree.

**The core interfaces:**
```go
type Node interface {
    TokenLiteral() string  // token literal for debugging
    String() string        // human-readable representation
}

type Statement interface {
    Node
    statementNode()   // marker method — just marks it as a statement
}

type Expression interface {
    Node
    expressionNode()  // marker method — marks it as an expression
}
```

**Statements vs Expressions — important distinction:**

| Type | What it is | Examples |
|---|---|---|
| **Statement** | Does something, produces no value | `let x = 5;`, `return 5;` |
| **Expression** | Produces a value | `5 + 3`, `add(2,3)`, `true` |

**All AST node types:**

| Node | Represents | Example |
|---|---|---|
| `Program` | Entire program | root node |
| `LetStatement` | Variable declaration | `let x = 5;` |
| `ReturnStatement` | Return from function | `return x + 1;` |
| `ExpressionStatement` | An expression used as a statement | `5 + 3;` |
| `Identifier` | A variable name | `x`, `add` |
| `IntegerLiteral` | A number | `42` |
| `Boolean` | true or false | `true` |
| `PrefixExpression` | Unary operator | `-5`, `!true` |
| `InfixExpression` | Binary operator | `5 + 3`, `x == y` |
| `IfExpression` | If/else block | `if (x) { ... }` |
| `BlockStatement` | A `{ ... }` block | `{ let x = 1; x }` |
| `FunctionLiteral` | Function definition | `fn(x) { x + 1 }` |
| `CallExpression` | Function call | `add(1, 2)` |

**Why does `IfExpression` extend `Expression` not `Statement`?**

Because in this language, `if` returns a value! You can write:
```
let result = if (x > 5) { 10 } else { 20 };
```
This is a key design decision — everything is an expression.

---

### `parser/parser.go` — The Pratt Parser

**Purpose:** Takes the token stream from the Lexer and builds an AST.

**The Parser struct:**
```go
type Parser struct {
    l              *lexer.Lexer
    currToken      token.Token   // token being examined now
    peekToken      token.Token   // next token (look-ahead)
    prefixParseFns map[token.TokenType]prefixParseFn
    infixParseFns  map[token.TokenType]infixParseFn
    errors         []string
}
```

**Why TWO tokens (curr + peek)?**

The parser always needs to look one token ahead to make decisions. Example: when you see `5`, is it `5;` (end of expression) or `5 + 3` (part of infix)? You need to peek at the next token to decide.

**Operator Precedence — the most important concept in the parser:**

```go
const (
    _ int = iota  // 0, unused
    LOWEST        // 1
    EQUALS        // 2  ==, !=
    LESSGREATER   // 3  <, >
    SUM           // 4  +, -
    PRODUCT       // 5  *, /
    PREFIX        // 6  -x, !x
    CALL          // 7  add(x)
)
```

// 1 + 2 * 3 * 4 --> (1 + ((2 * 3) * 4))

Higher number = binds tighter. This ensures `2 + 3 * 4` is parsed as `2 + (3 * 4)` not `(2 + 3) * 4`.

**Pratt Parsing — the core algorithm:**

This parser uses a technique called **Pratt Parsing** (also called Top-Down Operator Precedence).

The key idea: for each token type, register a **parsing function**:
- **prefix function** — how to parse this token at the START of an expression
- **infix function** — how to parse this token in the MIDDLE of an expression

```go
// Registration (in New()):
parser.registerPrefix(token.INT,    parser.parseIntegerLiteral)
parser.registerPrefix(token.IDENT,  parser.parseIdentifier)
parser.registerPrefix(token.BANG,   parser.parsePrefixExpression)
parser.registerPrefix(token.MINUS,  parser.parsePrefixExpression)

parser.registerInfix(token.PLUS,    parser.parseInfixExpression)
parser.registerInfix(token.ASTERIK, parser.parseInfixExpression)
parser.registerInfix(token.LPAREN,  parser.parseCallExpression)
```

**The `parseExpression` function — heart of the parser:**
*
```go
func (p *Parser) parseExpression(precedence int) ast.Expression {
    // 1. Find and call the prefix parse function for current token
    prefix := p.prefixParseFns[p.currToken.Type]
    leftExp := prefix()

    // 2. Loop: while next token has higher precedence, keep consuming
    for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
        infix := p.infixParseFns[p.peekToken.Type]
        p.nextToken()
        leftExp = infix(leftExp)  // leftExp becomes left side of infix
    }

    return leftExp
}
```

**Tracing `1 + 2 * 3` through `parseExpression(LOWEST)`:**

```
Call parseExpression(LOWEST=1)
  → prefix("1") → IntegerLiteral(1)
  → peekToken = "+", precedence("+")=SUM=4 > LOWEST=1, so enter loop
  → call infix(IntegerLiteral(1))  [parseInfixExpression]
      → left = IntegerLiteral(1), operator = "+"
      → call parseExpression(SUM=4)   ← recursive call with higher precedence
          → prefix("2") → IntegerLiteral(2)
          → peekToken = "*", precedence("*")=PRODUCT=5 > SUM=4, enter loop
          → call infix(IntegerLiteral(2))
              → left=2, op="*", right=parseExpression(PRODUCT=5)
                  → prefix("3") → IntegerLiteral(3)
                  → peekToken=";", 5 > LOWEST? No. Stop.
                  → return IntegerLiteral(3)
              → return InfixExpression(2 * 3)
          → return InfixExpression(2 * 3)
      → right = InfixExpression(2 * 3)
      → return InfixExpression(1 + (2 * 3))  ✓ correct!
```

**`expectPeek` — the validation helper:**
```go
func (p *Parser) expectPeek(t token.TokenType) bool {
    if p.peekTokenIs(t) {
        p.nextToken()   // only advance if it IS the expected token
        return true
    } else {
        p.peekError(t)  // add error message
        return false
    }
}
```
This is how the parser enforces grammar rules. For `let x = 5;`, the parser:
1. Sees `LET` → calls `parseLetStatement`
2. `expectPeek(IDENT)` → ensures `x` follows
3. `expectPeek(ASSIGN)` → ensures `=` follows
4. Then parses the value expression

---

### `object/object.go` — The Runtime Value System

**Purpose:** Defines what VALUES look like during execution. Every value in the language is an `Object`.

**The Object interface:**
```go
type Object interface {
    Type() ObjectType   // what kind of value: "INTEGER", "BOOLEAN", etc.
    Inspect() string    // how to print it
}
```

**All object types:**

| Object | Holds | Example value |
|---|---|---|
| `Integer` | `int64` | `42` |
| `Boolean` | `bool` | `true` |
| `Null` | nothing | `null` |
| `ReturnValue` | wraps another Object | used to bubble up `return` |
| `Error` | error message string | `"type mismatch: INTEGER + BOOLEAN"` |
| `Function` | params + body + env | `fn(x) { x + 1 }` |

**The Function object is special:**
```go
type Function struct {
    Parameters []*ast.Identifier  // the parameter names
    Body       *ast.BlockStatement // the function body (unevaluated AST)
    Env        *Environment        // THE CLOSURE — captures outer scope!
}
```

Notice the `Function` object stores the **AST node** (not evaluated code). When called, the body gets evaluated fresh each time. The `Env` field is what makes **closures** work.

**Why `ReturnValue` is a wrapper object?**

When you do `return 5` inside a function, the evaluator produces `ReturnValue{Integer{5}}`. This wrapper "bubbles up" through the evaluation stack. When `evalProgram` (top-level) sees it, it **unwraps** it:
```go
case *object.ReturnValue:
    return result.Value  // strip the wrapper, return just 5
```
Without this wrapper trick, a `return` inside a nested `if` inside a function would not work — the value would just be the last value, not a controlled return.

---

### `object/environment.go` — Variable Storage & Scoping

**Purpose:** Where variables live. Implements scoped variable lookup.

```go
type Environment struct {
    store map[string]Object  // variables in THIS scope
    outer *Environment       // parent scope (for closures)
}
```

**Two types of environments:**

```go
// Global/top-level environment
func NewEnvironment() *Environment {
    return &Environment{store: make(map[string]Object), outer: nil}
}

// Enclosed environment for function calls
func NewEnclosedEnvironment(outer *Environment) *Environment {
    env := NewEnvironment()
    env.outer = outer   // link to parent
    return env
}
```

**How variable lookup works (Get):**
```go
func (e *Environment) Get(name string) (Object, bool) {
    obj, ok := e.store[name]
    if !ok && e.outer != nil {
        obj, ok := e.outer.Get(name)  // look in OUTER scope
        return obj, ok
    }
    return obj, ok
}
```

This is **scope chain** — if the variable isn't in the current scope, look in the outer scope. This is how closures work:

```
let x = 10;                    // stored in global env
let getX = fn() { return x; } // fn body uses x
getX()                         // works! finds x in outer scope
```

When `getX()` is called:
1. New enclosed env is created (outer = global env)
2. `x` is not in the function's env
3. `Get` walks up to `outer` (global env)
4. Finds `x = 10` there
5. Returns `10`

---

### `evaluator/evaluator.go` — The Execution Engine

**Purpose:** Recursively walks the AST and produces values.

**The main `Eval` function — a giant type switch:**
```go
func Eval(node ast.Node, env *object.Environment) object.Object {
    switch node := node.(type) {
        case *ast.Program:          return evalProgram(node, env)
        case *ast.LetStatement:     // evaluate value, store in env
        case *ast.Identifier:       return evalIdentifier(node, env)
        case *ast.IntegerLiteral:   return &object.Integer{Value: node.Value}
        case *ast.Boolean:          return nativeBoolToBooleanObject(node.Value)
        case *ast.PrefixExpression: // eval right, apply prefix operator
        case *ast.InfixExpression:  // eval left + right, apply operator
        case *ast.BlockStatement:   return evalBlockStatement(node, env)
        case *ast.IfExpression:     return evalIfExpression(node, env)
        case *ast.ReturnStatement:  // eval value, wrap in ReturnValue{}
        case *ast.FunctionLiteral:  return &object.Function{...}
        case *ast.CallExpression:   // eval function, eval args, call it
    }
}
```

**Singleton booleans — optimization:**
```go
var (
    NULL  = &object.Null{}
    TRUE  = &object.Boolean{Value: true}
    FALSE = &object.Boolean{Value: false}
)
```

Only ONE instance of `true`, `false`, and `null` ever exists. So `true == true` works with pointer comparison (`left == right`) — no need to compare field values. This is efficient and safe.

**Error propagation:**
```go
func isError(obj object.Object) bool {
    return obj != nil && obj.Type() == object.ERROR_OBJ
}

// Used everywhere like:
right := Eval(node.Right, env)
if isError(right) { return right }  // short-circuit on error
```

Errors "bubble up" automatically. Every `Eval` call checks for errors before proceeding. This means an error anywhere in deep nested code will cleanly propagate to the top.

**`evalBlockStatement` vs `evalProgram` — subtle difference:**

```go
// evalProgram: UNWRAPS ReturnValue (we reached top level)
case *object.ReturnValue:
    return result.Value   // strip wrapper

// evalBlockStatement: does NOT unwrap (let it bubble up)
if result.Type() == object.RETURN_VALUE_OBJ {
    return result  // keep wrapper so it bubbles up further
}
```

This distinction is critical for nested functions:
```
let outer = fn() {
    let inner = fn() { return 99; };
    inner();   // returns ReturnValue{99}
    return 1;  // this should NOT run
};
```
If `evalBlockStatement` unwrapped, `inner()` would return `99` (unwrapped), then execution would continue and return `1`. By NOT unwrapping in blocks, `ReturnValue{99}` bubbles all the way out.

---

### `repl/repl.go` — The Interactive Shell

**Purpose:** The Read-Eval-Print Loop. Takes user input, runs it, shows output.

```go
func Start(in io.Reader, out io.Writer) {
    scanner := bufio.NewScanner(in)
    env := object.NewEnvironment()   // one env for whole session

    for {
        fmt.Printf(">> ")
        scanned := scanner.Scan()    // READ
        line := scanner.Text()

        l := lexer.New(line)
        p := parser.New(l)
        program := p.ParseProgram()  // PARSE
        
        ast.PrintAST(program, "")    // debug: print the AST tree
        
        evaluated := evaluator.Eval(program, env)  // EVAL
        io.WriteString(out, evaluated.Inspect())    // PRINT
    }
}
```

Key point: `env` is created **once** outside the loop. This means variables you define in one line persist to the next:
```
>> let x = 5;
>> x
5
```
Both lines share the same `env`.

---