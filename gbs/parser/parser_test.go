package parser

import (
	"bangScript/gbs/scanner"
	"reflect"
	"testing"
)

// Mock token creation helper
func makeToken(tokenType scanner.Tokentype, lexeme string, line int) *scanner.Token {
	return &scanner.Token{
		Ttype: tokenType,
		Lexem: lexeme,
		Line:  line,
	}
}

// TODO: we parse source as []byte without EOF
func TestPrimaryExpression(t *testing.T) {
	tests := []struct {
		name     string
		tokens   Tokens
		expected Primary
	}{
		{
			name:     "Number literal",
			tokens:   Tokens{makeToken(scanner.NUMBER, "42", 1), makeToken(scanner.EOF, "", 1)},
			expected: Primary{Node: makeToken(scanner.NUMBER, "42", 1)},
		},
		{
			name:     "String literal",
			tokens:   Tokens{makeToken(scanner.STRING, "hello", 1), makeToken(scanner.EOF, "", 1)},
			expected: Primary{Node: makeToken(scanner.STRING, "hello", 1)},
		},
		{
			name:     "Boolean true",
			tokens:   Tokens{makeToken(scanner.TRUE, "true", 1), makeToken(scanner.EOF, "", 1)},
			expected: Primary{Node: makeToken(scanner.TRUE, "true", 1)},
		},
		{
			name:     "Boolean false",
			tokens:   Tokens{makeToken(scanner.FALSE, "false", 1), makeToken(scanner.EOF, "", 1)},
			expected: Primary{Node: makeToken(scanner.FALSE, "false", 1)},
		},
		{
			name:     "Nil literal",
			tokens:   Tokens{makeToken(scanner.NIL, "nil", 1), makeToken(scanner.EOF, "", 1)},
			expected: Primary{Node: makeToken(scanner.NIL, "nil", 1)},
		},
		{
			name:     "Identifier",
			tokens:   Tokens{makeToken(scanner.IDENTIFIER, "variable", 1), makeToken(scanner.EOF, "", 1)},
			expected: Primary{Node: makeToken(scanner.IDENTIFIER, "variable", 1)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			primResult, err := tt.tokens.expression()
			if err != nil {
				t.Fatalf("Expected Primary expression, got %s", err.Error())
			}
			if !reflect.DeepEqual(primResult, tt.expected) {
				t.Errorf("Expected %+v, got %+v", tt.expected, primResult)
			}
		})
	}
}

func TestUnaryExpression(t *testing.T) {
	minusT := makeToken(scanner.MINUS, "-", 1)
	numberT := makeToken(scanner.NUMBER, "42", 1)
	EOFT := makeToken(scanner.EOF, "", 1)
	idT := makeToken(scanner.IDENTIFIER, "isOpen", 1)
	bangT := makeToken(scanner.BANG, "!", 1)
	tests := []struct {
		name     string
		tokens   Tokens
		expected Unary
	}{
		{
			name: "Negation",
			tokens: Tokens{
				minusT,
				numberT,
				EOFT,
			},
			expected: Unary{
				operator: minusT,
				right:    Primary{Node: numberT},
			},
		},
		{
			name: "Logical Not",
			tokens: Tokens{
				bangT,
				idT,
				EOFT,
			},
			expected: Unary{
				operator: bangT,
				right:    Primary{Node: idT},
			},
		},
		{
			name: "Nested Unary",
			tokens: Tokens{
				bangT,
				bangT,
				idT,
				EOFT,
			},
			expected: Unary{
				operator: bangT,
				right: Unary{
					operator: bangT,
					right:    Primary{Node: idT},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UnaryResult, err := tt.tokens.expression()
			if err != nil {
				t.Fatalf("Expected Unary expression, got %T", err)
			}
			if !reflect.DeepEqual(UnaryResult, tt.expected) {
				t.Errorf("Expected %+v, got %+v", tt.expected, UnaryResult)
			}
		})
	}
}

func TestBinaryExpression(t *testing.T) {
	// Common tokens
	numberA := makeToken(scanner.NUMBER, "5", 1)
	numberB := makeToken(scanner.NUMBER, "10", 1)
	stringA := makeToken(scanner.STRING, "hello", 1)
	stringB := makeToken(scanner.STRING, "world", 1)
	idA := makeToken(scanner.IDENTIFIER, "a", 1)
	idB := makeToken(scanner.IDENTIFIER, "b", 1)
	/**
	trueT := makeToken(scanner.TRUE, "true", 1)
	falseT := makeToken(scanner.FALSE, "false", 1)
	*/

	EOFT := makeToken(scanner.EOF, "", 1)

	// Operators
	plus := makeToken(scanner.PLUS, "+", 1)
	minus := makeToken(scanner.MINUS, "-", 1)
	star := makeToken(scanner.STAR, "*", 1)
	slash := makeToken(scanner.SLASH, "/", 1)
	greater := makeToken(scanner.GREATER, ">", 1)
	greaterEqual := makeToken(scanner.GREATER_EQUAL, ">=", 1)
	less := makeToken(scanner.LESS, "<", 1)
	lessEqual := makeToken(scanner.LESS_EQUAL, "<=", 1)
	equalEqual := makeToken(scanner.EQUAL_EQUAL, "==", 1)
	bangEqual := makeToken(scanner.BANG_EQUAL, "!=", 1)
	/**
	and := makeToken(scanner.AND, "and", 1)
	or := makeToken(scanner.OR, "or", 1)
	*/

	tests := []struct {
		name     string
		tokens   Tokens
		expected Binary
	}{
		// Arithmetic operators
		{
			name:   "Addition",
			tokens: Tokens{numberA, plus, numberB, EOFT},
			expected: Binary{
				left:     Primary{Node: numberA},
				operator: plus,
				right:    Primary{Node: numberB},
			},
		},
		{
			name:   "Subtraction",
			tokens: Tokens{numberA, minus, numberB, EOFT},
			expected: Binary{
				left:     Primary{Node: numberA},
				operator: minus,
				right:    Primary{Node: numberB},
			},
		},
		{
			name:   "Multiplication",
			tokens: Tokens{numberA, star, numberB, EOFT},
			expected: Binary{
				left:     Primary{Node: numberA},
				operator: star,
				right:    Primary{Node: numberB},
			},
		},
		{
			name:   "Division",
			tokens: Tokens{numberA, slash, numberB, EOFT},
			expected: Binary{
				left:     Primary{Node: numberA},
				operator: slash,
				right:    Primary{Node: numberB},
			},
		},

		// String concatenation
		{
			name:   "String concatenation",
			tokens: Tokens{stringA, plus, stringB, EOFT},
			expected: Binary{
				left:     Primary{Node: stringA},
				operator: plus,
				right:    Primary{Node: stringB},
			},
		},

		// Comparison operators
		{
			name:   "Greater than",
			tokens: Tokens{numberA, greater, numberB, EOFT},
			expected: Binary{
				left:     Primary{Node: numberA},
				operator: greater,
				right:    Primary{Node: numberB},
			},
		},
		{
			name:   "Greater than or equal",
			tokens: Tokens{numberA, greaterEqual, numberB, EOFT},
			expected: Binary{
				left:     Primary{Node: numberA},
				operator: greaterEqual,
				right:    Primary{Node: numberB},
			},
		},
		{
			name:   "Less than",
			tokens: Tokens{numberA, less, numberB, EOFT},
			expected: Binary{
				left:     Primary{Node: numberA},
				operator: less,
				right:    Primary{Node: numberB},
			},
		},
		{
			name:   "Less than or equal",
			tokens: Tokens{numberA, lessEqual, numberB, EOFT},
			expected: Binary{
				left:     Primary{Node: numberA},
				operator: lessEqual,
				right:    Primary{Node: numberB},
			},
		},
		{
			name:   "Equal",
			tokens: Tokens{idA, equalEqual, idB, EOFT},
			expected: Binary{
				left:     Primary{Node: idA},
				operator: equalEqual,
				right:    Primary{Node: idB},
			},
		},
		{
			name:   "Not equal",
			tokens: Tokens{idA, bangEqual, idB, EOFT},
			expected: Binary{
				left:     Primary{Node: idA},
				operator: bangEqual,
				right:    Primary{Node: idB},
			},
		},
		// Logical operators
		/*
			{
				name:   "Logical AND",
				tokens: Tokens{trueT, and, falseT, EOFT},
				expected: Binary{
					left:     Primary{Node: trueT},
					operator: and,
					right:    Primary{Node: falseT},
				},
			},
			{
				name:   "Logical OR",
				tokens: Tokens{trueT, or, falseT, EOFT},
				expected: Binary{
					left:     Primary{Node: trueT},
					operator: or,
					right:    Primary{Node: falseT},
				},
			},
		*/

		// Complex expressions
		{
			name:   "Chained Binary expressions (1 + 2 * 3)",
			tokens: Tokens{numberA, plus, numberB, star, makeToken(scanner.NUMBER, "3", 1), EOFT},
			expected: Binary{
				left:     Primary{Node: numberA},
				operator: plus,
				right: Binary{
					left:     Primary{Node: numberB},
					operator: star,
					right:    Primary{Node: makeToken(scanner.NUMBER, "3", 1)},
				},
			},
		},
		// Mixed expressions
		{
			name:   "Binary with Unary right operand",
			tokens: Tokens{numberA, plus, minus, numberB, EOFT},
			expected: Binary{
				left:     Primary{Node: numberA},
				operator: plus,
				right: Unary{
					operator: minus,
					right:    Primary{Node: numberB},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.tokens.expression()
			if err != nil {
				t.Fatalf("Expected Binary expression, got error: %s", err.Error())
			}

			// Type check
			binResult, ok := result.(Binary)
			if !ok {
				t.Fatalf("Expected Binary expression, got %T", result)
			}

			if !reflect.DeepEqual(binResult, tt.expected) {
				t.Errorf("Expected %+v, got %+v", tt.expected, binResult)
			}
		})
	}
}

func TestParenthesizeExpression(t *testing.T) {
	// Tokens
	leftParen := makeToken(scanner.LEFT_PAREN, "(", 1)
	rightParen := makeToken(scanner.RIGHT_PAREN, ")", 1)
	numberT := makeToken(scanner.NUMBER, "42", 1)
	plus := makeToken(scanner.PLUS, "+", 1)
	number2T := makeToken(scanner.NUMBER, "10", 1)
	EOFT := makeToken(scanner.EOF, "", 1)

	tests := []struct {
		name     string
		tokens   Tokens
		expected Exp
	}{
		{
			name:     "Simple grouping",
			tokens:   Tokens{leftParen, numberT, rightParen, EOFT},
			expected: Primary{Node: numberT},
		},
		{
			name:   "Grouped Binary expression",
			tokens: Tokens{leftParen, numberT, plus, number2T, rightParen, EOFT},
			expected: Binary{
				left:     Primary{Node: numberT},
				operator: plus,
				right:    Primary{Node: number2T},
			},
		},
		{
			name:     "Nested grouping",
			tokens:   Tokens{leftParen, leftParen, numberT, rightParen, rightParen, EOFT},
			expected: Primary{Node: numberT},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.tokens.expression()
			if err != nil {
				t.Fatalf("Expected grouping expression, got error: %s", err.Error())
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

func TestOperatorPrecedence(t *testing.T) {
	// Tokens for numbers and operations
	num1 := makeToken(scanner.NUMBER, "1", 1)
	num2 := makeToken(scanner.NUMBER, "2", 1)
	num3 := makeToken(scanner.NUMBER, "3", 1)
	//num4 := makeToken(scanner.NUMBER, "4", 1)
	plus := makeToken(scanner.PLUS, "+", 1)
	minus := makeToken(scanner.MINUS, "-", 1)
	star := makeToken(scanner.STAR, "*", 1)
	slash := makeToken(scanner.SLASH, "/", 1)
	//greater := makeToken(scanner.GREATER, ">", 1)
	//less := makeToken(scanner.LESS, "<", 1)
	leftParen := makeToken(scanner.LEFT_PAREN, "(", 1)
	rightParen := makeToken(scanner.RIGHT_PAREN, ")", 1)
	EOFT := makeToken(scanner.EOF, "", 1)

	tests := []struct {
		name        string
		tokens      Tokens
		expected    Exp
		description string
	}{
		{
			name:   "Multiplication before addition",
			tokens: Tokens{num1, plus, num2, star, num3, EOFT},
			expected: Binary{
				left:     Primary{Node: num1},
				operator: plus,
				right: Binary{
					left:     Primary{Node: num2},
					operator: star,
					right:    Primary{Node: num3},
				},
			},
			description: "1 + 2 * 3 should evaluate as 1 + (2 * 3)",
		},
		{
			name:   "Division before subtraction",
			tokens: Tokens{num1, minus, num2, slash, num3, EOFT},
			expected: Binary{
				left:     Primary{Node: num1},
				operator: minus,
				right: Binary{
					left:     Primary{Node: num2},
					operator: slash,
					right:    Primary{Node: num3},
				},
			},
			description: "1 - 2 / 3 should evaluate as 1 - (2 / 3)",
		},
		/*
			 * 	{
					name:   "Comparison before logical AND",
					tokens: Tokens{num1, greater, num2, and, num3, less, num4, EOFT},
					expected: Binary{
						left: Binary{
							left:     Primary{Node: num1},
							operator: greater,
							right:    Primary{Node: num2},
						},
						operator: and,
						right: Binary{
							left:     Primary{Node: num3},
							operator: less,
							right:    Primary{Node: num4},
						},
					},
					description: "1 > 2 and 3 < 4 should evaluate as (1 > 2) and (3 < 4)",
				},
		*/
		{
			name:   "Parentheses override precedence",
			tokens: Tokens{leftParen, num1, plus, num2, rightParen, star, num3, EOFT},
			expected: Binary{
				left: Binary{
					left:     Primary{Node: num1},
					operator: plus,
					right:    Primary{Node: num2},
				},
				operator: star,
				right:    Primary{Node: num3},
			},
			description: "(1 + 2) * 3 should evaluate as (1 + 2) * 3",
		},

		{
			name:   "Nested Parentheses override precedence",
			tokens: Tokens{leftParen, leftParen, num1, plus, num2, rightParen, rightParen, star, num3, EOFT},
			expected: Binary{
				left: Binary{
					left:     Primary{Node: num1},
					operator: plus,
					right:    Primary{Node: num2},
				},
				operator: star,
				right:    Primary{Node: num3},
			},
			description: "(((1 + 2)) * 3 should evaluate as (1 + 2) * 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.tokens.expression()
			if err != nil {
				t.Fatalf("Expected expression, got error: %s", err.Error())
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("%s: Expected %+v, got %+v", tt.description, tt.expected, result)
			}
		})
	}
}
