package parsing

import (
	"fmt"
	"gox/internal/scanning"
	"testing"
)

type AstPrinter struct {
}

func (r *AstPrinter) Print(expr Expr) (string, error) {
	return expr.accept(r).(string), nil
}

func (r *AstPrinter) visitForLiteral(expr *Literal) any {
	if expr.value == nil {
		return "nil"
	}
	return fmt.Sprint(expr.value)
}

func (r *AstPrinter) visitForUnary(expr *Unary) any {
	return r.parenthesize(expr.operator.Lexeme, expr.right)
}

func (r *AstPrinter) visitForBinary(expr *Binary) any {
	return r.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (r *AstPrinter) visitForGrouping(expr *Grouping) any {
	return r.parenthesize("group", expr.expression)
}

func (r *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	res := "(" + name
	for _, expr := range exprs {
		accept := expr.accept(r)
		res += " " + accept.(string)
	}
	res += ")"
	return res
}

func TestAstPrinter_Print(t *testing.T) {
	type args struct {
		expr Expr
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "basic",
		args: args{expr: &Binary{
			left: &Unary{
				operator: scanning.Token{
					TokenType: scanning.MINUS,
					Lexeme:    "-",
					Literal:   nil,
					Line:      1,
				},
				right: &Literal{value: 123},
			},
			operator: scanning.Token{
				TokenType: scanning.STAR,
				Lexeme:    "*",
				Literal:   nil,
				Line:      1,
			},
			right: &Grouping{expression: &Literal{value: 45.67}},
		}},
		want: "(* (- 123) (group 45.67))",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AstPrinter{}
			ast, err := r.Print(tt.args.expr)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			if ast != tt.want {
				t.Errorf("want: %s, got: %s", tt.want, ast)
				t.FailNow()
			}
			t.Log(ast)
		})
	}
}
