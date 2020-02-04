package codescan

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getEnumBasicLitValue(t *testing.T) {

	verifyGetEnumBasicLitValue(t, ast.BasicLit{Kind: token.INT, Value: "0"}, int64(0))
	verifyGetEnumBasicLitValue(t, ast.BasicLit{Kind: token.INT, Value: "-1"}, int64(-1))
	verifyGetEnumBasicLitValue(t, ast.BasicLit{Kind: token.INT, Value: "42"}, int64(42))
	verifyGetEnumBasicLitValue(t, ast.BasicLit{Kind: token.INT, Value: ""}, nil)
	verifyGetEnumBasicLitValue(t, ast.BasicLit{Kind: token.INT, Value: "word"}, nil)

	verifyGetEnumBasicLitValue(t, ast.BasicLit{Kind: token.FLOAT, Value: "0"}, float64(0))
	verifyGetEnumBasicLitValue(t, ast.BasicLit{Kind: token.FLOAT, Value: "-1"}, float64(-1))
	verifyGetEnumBasicLitValue(t, ast.BasicLit{Kind: token.FLOAT, Value: "42"}, float64(42))
	verifyGetEnumBasicLitValue(t, ast.BasicLit{Kind: token.FLOAT, Value: "1.1234"}, float64(1.1234))
	verifyGetEnumBasicLitValue(t, ast.BasicLit{Kind: token.FLOAT, Value: "1.9876"}, float64(1.9876))
	verifyGetEnumBasicLitValue(t, ast.BasicLit{Kind: token.FLOAT, Value: ""}, nil)
	verifyGetEnumBasicLitValue(t, ast.BasicLit{Kind: token.FLOAT, Value: "word"}, nil)

	verifyGetEnumBasicLitValue(t, ast.BasicLit{Kind: token.STRING, Value: "Foo"}, "Foo")
	verifyGetEnumBasicLitValue(t, ast.BasicLit{Kind: token.STRING, Value: ""}, "")
	verifyGetEnumBasicLitValue(t, ast.BasicLit{Kind: token.STRING, Value: "0"}, "0")
	verifyGetEnumBasicLitValue(t, ast.BasicLit{Kind: token.STRING, Value: "1.1"}, "1.1")

}

func verifyGetEnumBasicLitValue(t *testing.T, basicLit ast.BasicLit, expected interface{}) {
	actual := getEnumBasicLitValue(&basicLit)

	assert.Equal(t, expected, actual)
}
