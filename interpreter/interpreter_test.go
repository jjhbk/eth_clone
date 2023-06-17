package interpreter_test

import (
	"jjeth/interpreter"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInterpreter(t *testing.T) {
	code := []string{"PUSH", "2", "PUSH", "3", "ADD", "STOP"}
	inter := interpreter.New()
	res, err := inter.RunCode(code)
	require.NoError(t, err)
	require.Equal(t, res, 5)

	code = []string{"PUSH", "2", "PUSH", "3", "SUB", "STOP"}
	inter = interpreter.New()
	res, err = inter.RunCode(code)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	code = []string{"PUSH", "2", "PUSH", "3", "MUL", "STOP"}
	inter = interpreter.New()
	res, err = inter.RunCode(code)
	require.NoError(t, err)
	require.Equal(t, res, 6)

	code = []string{"PUSH", "3", "PUSH", "18", "DIV", "STOP"}
	inter = interpreter.New()
	res, err = inter.RunCode(code)
	require.NoError(t, err)
	require.Equal(t, res, 6)

	code = []string{"PUSH", "3", "PUSH", "18", "LT", "STOP"}
	inter = interpreter.New()
	res, err = inter.RunCode(code)
	require.NoError(t, err)
	require.Equal(t, res, 0)

	code = []string{"PUSH", "3", "PUSH", "18", "GT", "STOP"}
	inter = interpreter.New()
	res, err = inter.RunCode(code)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	code = []string{"PUSH", "3", "PUSH", "3", "EQ", "STOP"}
	inter = interpreter.New()
	res, err = inter.RunCode(code)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	code = []string{"PUSH", "1", "PUSH", "0", "AND", "STOP"}
	inter = interpreter.New()
	res, err = inter.RunCode(code)
	require.NoError(t, err)
	require.Equal(t, res, 0)

	code = []string{"PUSH", "1", "PUSH", "0", "OR", "STOP"}
	inter = interpreter.New()
	res, err = inter.RunCode(code)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	code = []string{"PUSH", "6", "JUMP", "PUSH", "0", "JUMP", "PUSH", "1999996", "STOP"}
	inter = interpreter.New()
	res, err = inter.RunCode(code)
	require.NoError(t, err)
	require.Equal(t, res, 1999996)

	code = []string{"PUSH", "8", "PUSH", "1", "JUMPI", "PUSH", "0", "JUMP", "PUSH", "1999996", "STOP"}
	inter = interpreter.New()
	res, err = inter.RunCode(code)
	require.NoError(t, err)
	require.Equal(t, res, 1999996)
}
