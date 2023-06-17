package interpreter

import (
	"fmt"
	"strconv"

	stack "github.com/emirpasic/gods/stacks/arraystack"
)

const (
	STOP  = "STOP"
	ADD   = "ADD"
	PUSH  = "PUSH"
	SUB   = "SUB"
	MUL   = "MUL"
	DIV   = "DIV"
	EQ    = "EQ"
	LT    = "LT"
	GT    = "GT"
	AND   = "AND"
	OR    = "OR"
	JUMP  = "JUMP"
	JUMPI = "JUMPI"
)
const EXECUTION_LIMIT = 100

// code:= []int{PUSH,2,PUSH,3,ADD,STOP}
type Interpreter struct {
	programCounter int
	stack          stack.Stack
	code           []string
}

func (inter *Interpreter) jump() error {
	destination, _ := inter.stack.Pop()
	if destination.(int) < 0 || destination.(int) > len(inter.code) {
		return fmt.Errorf("Invalid destination : %d", destination)
	}
	inter.programCounter = destination.(int)
	fmt.Println(inter.programCounter, inter.code[inter.programCounter])
	return nil
}

func New() *Interpreter {
	inter := new(Interpreter)
	inter.code = make([]string, 0)
	inter.programCounter = 0
	inter.stack = *stack.New()
	return inter
}
func (inter *Interpreter) RunCode(code []string) (int, error) {
	val := 0
	inter.code = code
	for inter.programCounter < len(code) {
		if inter.programCounter > EXECUTION_LIMIT {
			return 0, fmt.Errorf("Error: Exceeded the execution limit")
		}
		opCode := code[inter.programCounter]
		switch opCode {

		case STOP:
			res, _ := inter.stack.Pop()
			fmt.Printf("result of computation is:%d\n", res)
			fmt.Println("execution complete")
			return res.(int), nil
		case PUSH:
			inter.programCounter++
			value, err := strconv.Atoi(code[inter.programCounter])
			CheckErr(err)
			inter.stack.Push(value)
			inter.programCounter++
		case ADD, SUB, MUL, DIV, LT, GT, EQ, AND, OR:
			inter.programCounter++
			val1, _ := inter.stack.Pop()
			val2, _ := inter.stack.Pop()
			fmt.Println(val1, val2)
			if opCode == ADD {
				val = val1.(int) + val2.(int)
			} else if opCode == SUB {
				val = val1.(int) - val2.(int)
			} else if opCode == MUL {
				val = val1.(int) * val2.(int)
			} else if opCode == DIV {
				val = val1.(int) / val2.(int)
			} else if opCode == LT {
				if val1.(int) < val2.(int) {
					val = 1
				} else {
					val = 0
				}
			} else if opCode == GT {
				if val1.(int) > val2.(int) {
					val = 1
				} else {
					val = 0
				}
			} else if opCode == EQ {
				if val1.(int) == val2.(int) {
					val = 1
				} else {
					val = 0
				}
			} else if opCode == AND || opCode == OR {
				var log1, log2, log bool
				if val1.(int) == 1 {
					log1 = true
				} else {
					log1 = false
				}
				if val2.(int) == 1 {
					log2 = true
				} else {
					log2 = false
				}
				if opCode == AND {
					log = log1 && log2
				} else {
					log = log1 || log2
				}
				if log == true {
					val = 1
				} else {
					val = 0
				}
			}

			inter.stack.Push(val)
		case JUMP:
			err := inter.jump()
			if err != nil {
				return 0, err
			}
		case JUMPI:
			condition, _ := inter.stack.Pop()
			if condition.(int) == 1 {
				err := inter.jump()
				if err != nil {
					return 0, err
				}
			}
		default:
			return val, fmt.Errorf("Invalid op code")
		}

	}
	return 0, nil
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
