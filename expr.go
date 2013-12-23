package main

import (
	"fmt"
	"io/ioutil"
	"math"
)

//LexType
/*
const (
	OPERATOR = 1000 + iota
	VARIABLE
	CONSTANT
	KEYWORD
	TOKEN
)
type LexType int
*/
const (
	OPERATOR = "operator"
	VARIABLE = "variable"
	CONSTANT = "constant"
	KEYWORD  = "keyword"
	TOKEN    = "token"
)

type LexType string

type Lexical struct {
	Type   LexType
	String string
	Value  int
}

var lex []*Lexical = make([]*Lexical, 0, 20)

// tokens
const (
	COMMA = 300 + iota
	SEMICOLON
)

// keywords:
const (
	VAR = 200 + iota
	ECHO
)

// operators:
const (
	_ = iota
	EVAL
	ADDEVAL
	SUBEVAL
	MULEVAL
	DIVEVAL
	OR
	AND
	BITOR
	BITAND
	EQUAL
	NOTEQUAL
	RIGHTSHIFT
	LEFTSHIFT
	ADD
	SUB
	MUL
	DIV
	MOD
	NOT
	BITNOT
	INC
	DEC
	PARENL
	PARENR
	MAX
)

var Priority []int

type OP int

type Expr struct {
	Operation OP
	LeafData  string//key of idTable
	Left      *Expr
	Right     *Expr
}

type VarData struct {
	//0:invalid, 1:StrValue OK, 2:NumValue OK, 4+1:Only StrValue
	flags    int
	StrValue string
	NumValue float64
}
type ConstData VarData // alias for VarData

type VarType interface {
	ToString() string
	ToNumber() float64
}

var idTable map[string]VarType = make(map[string]VarType)

func (id *VarData) ToString() string {
	if id.flags == 0 {
		return ""
	}
	if (id.flags & 1) == 1 {
		return id.StrValue
	}
	if (id.flags & 2) == 2 {
		id.StrValue = fmt.Sprintf("%f", id.NumValue)
		id.flags |= 1
		return id.StrValue
	}
	fmt.Println("ERROR:: no such flat:", id.flags)
	return ""
}
func (id *ConstData) ToString() string {
	return (*VarData)(id).ToString()
}

func (id *VarData) ToNumber() float64 {
	if id.flags == 0 {
		return 0
	}
	if (id.flags & 2) == 2 {
		return id.NumValue
	}
	if (id.flags & 4) == 4 {
		fmt.Println("ERROR:: cannot convert string to number:", id.StrValue)
		return 0 // cannot convert string to number.
	}
	if (id.flags & 1) == 1 {
		var tmp float64 = 0
		n, err := fmt.Sscanf(id.StrValue, "%f", &tmp)
		if n > 0 && err == nil {
			id.NumValue = tmp
			id.flags |= 2
			return id.NumValue
		} else {
			fmt.Println("ERROR:: cannot convert string to number:", id.StrValue)
			id.flags |= 4
			return 0
		}
	}
	fmt.Println("ERROR:: no such flag:", id.flags)
	return 0
}
func (id *ConstData) ToNumber() float64 {
	return (*VarData)(id).ToNumber()
}

func parse(data []byte) {
	length := len(data)
	pa := 0
	for i := 0; i < length; i++ {
		switch data[i] {
		case ' ', '\t', '\n', '\r':
			if i-pa > 0 {
				parseElem(data, pa, i)
			}
			pa = i + 1
		case ';', '(', ')', ',':
			if i-pa > 0 {
				parseElem(data, pa, i)
				pa = i
			}
			parseElem(data, pa, i+1)
			pa = i + 1
		case '=':
			parseElem(data, pa, i+1)
			pa = i + 1
		}
	}
}

func parseElem(data []byte, begin int, end int) {
	str := string(data[begin:end])
	other := false
	//println("\t\t\t\t", str)
	var token *Lexical = &Lexical{String: str, Value: 0}
	switch str {
	case "var":
		//println("key word var")
		token.Type = KEYWORD
		token.Value = VAR
		lex = append(lex, token)
	case "echo":
		token.Type = KEYWORD
		token.Value = ECHO
		//println("key word echo")
		lex = append(lex, token)
	case ";":
		token.Type = TOKEN
		token.Value = SEMICOLON
		//println("semi-colon")
		lex = append(lex, token)
		// The end of a statement.
	case "+":
		//println("add")
		token.Type = OPERATOR
		token.Value = ADD
		lex = append(lex, token)
	case "-":
		//println("sub")
		token.Type = OPERATOR
		token.Value = SUB
		lex = append(lex, token)
	case "*":
		//println("mul")
		token.Type = OPERATOR
		token.Value = MUL
		lex = append(lex, token)
	case "/":
		//println("div")
		token.Type = OPERATOR
		token.Value = DIV
		lex = append(lex, token)
	case "%":
		//println("mod")
		token.Type = OPERATOR
		token.Value = MOD
		lex = append(lex, token)
	case "++":
		//println("inc")
		token.Type = OPERATOR
		token.Value = INC
		lex = append(lex, token)
	case "--":
		//println("dec")
		token.Type = OPERATOR
		token.Value = DEC
		lex = append(lex, token)
	case "=":
		//println("evaluate")
		token.Type = OPERATOR
		token.Value = EVAL
		lex = append(lex, token)
	case "(":
		//println("left paren")
		token.Type = OPERATOR
		token.Value = PARENL
		lex = append(lex, token)
	case ")":
		//println("right paren")
		token.Type = OPERATOR
		token.Value = PARENR
		lex = append(lex, token)
	default:
		other = true
	}
	if other {
		if str[0] >= '0' && str[0] <= '9' {
			//println("constant", str)
			token.Type = CONSTANT
			lex = append(lex, token)
		} else {
			//println("variable", str)
			token.Type = VARIABLE
			lex = append(lex, token)
		}
	}
}

func initer() {
	xxx := []int{
		0,
		EVAL,
		EVAL, /*ADDEVAL*/
		EVAL, /*SUBEVAL*/
		EVAL, /*MULEVAL*/
		EVAL, /*DIVEVAL*/
		OR,
		AND,
		BITOR,
		BITAND,
		EQUAL,
		EQUAL, /*NOTEQUAL*/
		RIGHTSHIFT,
		RIGHTSHIFT, /*LEFTSHIFT*/
		ADD,
		ADD, /*SUB*/
		MUL,
		MUL, /*DIV*/
		MUL, /*MOD*/
		NOT,
		BITNOT,
		INC,
		INC, /*DEC*/
		PARENL,
		PARENR,
		MAX,
	}
	Priority = make([]int, len(xxx))
	for i := 0; i < MAX; i++ {
		Priority[i] = xxx[i]
	}
}

func main() {
	initer()
	//fmt.Println("AND=", AND)
	//fmt.Println("MAX=", MAX)
	fileName := "/home/togry/go/tester"
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("error reading files:", err)
		return
	}
	fmt.Println("file:", fileName, " size:", len(data))
	parse(data)

	//for _, d := range(lex) {
	//	fmt.Println("Type:", d.Type, " Str:", d.String, " Val:", d.Value)
	//}

	generate()
}

// syntax: (to be extended...)
// STMT :: VARI = EXPR ; | "var" VARI ; | "echo" EXPR ;
// VARI :: [a-zA-Z][a-zA-Z0-9]*
// EXPR :: EXPR + EXPR | EXPR - EXPR
// EXPR :: EXPR * EXPR | EXPR / EXPR | EXPR % EXPR
// EXPR :: ( EXPR ) | CONS | VARI
// CONS :: [1-9][0-9]*(.[0-9]+)?

// state machine:
//       "var"      VARI        ;
//      /----->(10)----->(11)------>(12)
//     /                         ; ^
//    /       "echo"              /
// (0)-------------------->(100)-/
//   \                  /   /   ^
//    \ VARI        =  /   /+-*/%\
//     \---->(200)----/    \-----/

// generator generates the binary syntax tree of
// the lexical elements.
var expr *Expr

func generate() {
	state := 0
	var last *Expr
	for _, elem := range lex {
		//println("---------------state=", state)
		switch elem.Type {
		case OPERATOR:
			switch elem.Value {
			case EVAL: //state should be 200
				if state != 200 {
					fmt.Println("ERROR invalid Expression")
				}
				state = 100
				tmp := expr
				expr = &Expr{Operation: EVAL, Left: tmp}
				last = expr
			case ADD, SUB:
				fallthrough
			case MUL, DIV, MOD:
				last = addExpression(nil, expr, elem)
			case PARENL, PARENR:
			}
		case VARIABLE:
			switch state {
			case 10: //variable declaration
				//add variable to variableTable
				if idTable[elem.String] == nil {
					idTable[elem.String] = new(VarData)
				} else {
					// what should I do?
				}
				state = 11
			case 0: //variable evaluation
				expr = &Expr{LeafData: elem.String}
				state = 200
			case 100: //expression
				//state = 100
				if expr == nil {
					expr = &Expr{LeafData: elem.String}
				} else {
					if last == nil {
						fmt.Println("ERROR: Wrong expression", elem.String)
						break
					}
					if last.Left == nil {
						last.Left = &Expr{LeafData: elem.String}
						//last = nil// just for binaries
					} else if last.Right == nil {
						last.Right = &Expr{LeafData: elem.String}
						//last = nil// just for binaries
					} else {
						fmt.Println("ERROR: wrong expression:", elem.String)
					}
				}
			}
		case CONSTANT: //constant expression
			// add to constant table, is it useful?
			if idTable[elem.String] == nil {
				idTable[elem.String] = &ConstData{flags: 1, StrValue: elem.String}
			} else {
				// ignore duplicate constants.
			}
			// add to expression
			if expr == nil {
				expr = &Expr{LeafData: elem.String}
			} else {
				if last.Left == nil {
					last.Left = &Expr{LeafData: elem.String}
					//last = nil// just for binaries
				} else if last.Right == nil {
					last.Right = &Expr{LeafData: elem.String}
					//last = nil// just for binaries
				} else {
					fmt.Println("ERROR: wrong Expression:", elem.String)
				}
			}
		case KEYWORD:
			//"var" or "echo"
			switch elem.Value {
			case VAR:
				state = 10
				expr = nil
				last = nil
			case ECHO:
				state = 100
				expr = nil //initial empty expression
			default:
				fmt.Println("ERROR: Don't know keyword:", elem.String)
			}
		case TOKEN:
			switch elem.String {
			case ";":
				//asume the state is 11 or 1xx
				state = 12
				//TODO do something.
				//output the generated syntax tree and digest it.
				//reset to default state
				val := evaluate(expr)
				fmt.Println("==============", val)
				state = 0
				expr = nil
				last = nil
			default:
				fmt.Println("ERROR: DON'T KNOW TOKEN:", elem.String)
			}
		}
	}
	//println("---------------state=", state)
	for i, v := range idTable {
		fmt.Println("id:", i, " value=", v.ToString())
	}
}
func addExpression(parent, root *Expr, elem *Lexical) *Expr {
	if root.Operation == 0 {
		//add here.
		tmp := root
		root = &Expr{Operation: OP(elem.Value), Left: tmp}
		last := root
		if parent == nil {
			expr = root
		} else if parent.Left == tmp {
			parent.Left = root
		} else if parent.Right == tmp {
			parent.Right = root
		}
		return last
	} else if Priority[elem.Value] < Priority[root.Operation] {
		tmp := root
		root = &Expr{Operation: OP(elem.Value), Left: tmp}
		last := root
		if parent == nil {
			expr = root
		} else if parent.Left == tmp {
			parent.Left = root
		} else if parent.Right == tmp {
			parent.Right = root
		}
		return last
	} else {
		if root.Left == nil {
			fmt.Println("ERROR: invalid expression:", elem.String)
			return nil
		}
		return addExpression(root, root.Right, elem)
	}
}

func evaluate(expr *Expr) float64 {
	if expr == nil {
		fmt.Println("EMPTY EXPRESSION")
		return math.NaN()
	}
	if expr.Operation == 0 {
		return idTable[expr.LeafData].ToNumber()
	}
	l := evaluate(expr.Left)
	r := evaluate(expr.Right)
	switch expr.Operation {
	case ADD:
		return l + r
	case SUB:
		return l - r
	case MUL:
		return l * r
	case DIV:
		return l / r
	case MOD:
		return float64(int(l) % int(r))
	case EVAL:
		id, err := idTable[expr.Left.LeafData].(*VarData)
		if err {
			id.flags = 2
			id.NumValue = r
		} else {
			fmt.Println("err: invalid type:", expr.Left.LeafData)
		}
		return r
	}
	return math.NaN()
}
