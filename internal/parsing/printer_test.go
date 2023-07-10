package parsing

import (
	"gox/internal/scanning"
	"testing"
)

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
