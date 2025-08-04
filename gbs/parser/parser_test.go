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
		expected primary
	}{
		{
			name:     "Number literal",
			tokens:   Tokens{makeToken(scanner.NUMBER, "42", 1), makeToken(scanner.EOF, "", 1)},
			expected: primary{node: makeToken(scanner.NUMBER, "42", 1)},
		},
		{
			name:     "String literal",
			tokens:   Tokens{makeToken(scanner.STRING, "hello", 1), makeToken(scanner.EOF, "", 1)},
			expected: primary{node: makeToken(scanner.STRING, "hello", 1)},
		},
		{
			name:     "Boolean true",
			tokens:   Tokens{makeToken(scanner.TRUE, "true", 1), makeToken(scanner.EOF, "", 1)},
			expected: primary{node: makeToken(scanner.TRUE, "true", 1)},
		},
		{
			name:     "Boolean false",
			tokens:   Tokens{makeToken(scanner.FALSE, "false", 1), makeToken(scanner.EOF, "", 1)},
			expected: primary{node: makeToken(scanner.FALSE, "false", 1)},
		},
		{
			name:     "Nil literal",
			tokens:   Tokens{makeToken(scanner.NIL, "nil", 1), makeToken(scanner.EOF, "", 1)},
			expected: primary{node: makeToken(scanner.NIL, "nil", 1)},
		},
		{
			name:     "Identifier",
			tokens:   Tokens{makeToken(scanner.IDENTIFIER, "variable", 1), makeToken(scanner.EOF, "", 1)},
			expected: primary{node: makeToken(scanner.IDENTIFIER, "variable", 1)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			primResult, err := tt.tokens.expression()
			if err != nil {
				t.Fatalf("Expected primary expression, got %s", err.Error())
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
		expected unary
	}{
		{
			name: "Negation",
			tokens: Tokens{
				minusT,
				numberT,
				EOFT,
			},
			expected: unary{
				operator: minusT,
				right:    primary{node: numberT},
			},
		},
		{
			name: "Logical Not",
			tokens: Tokens{
				bangT,
				idT,
				EOFT,
			},
			expected: unary{
				operator: bangT,
				right:    primary{node: idT},
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
			expected: unary{
				operator: bangT,
				right: unary{
					operator: bangT,
					right:    primary{node: idT},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unaryResult, err := tt.tokens.expression()
			if err != nil {
				t.Fatalf("Expected unary expression, got %T", err)
			}
			if !reflect.DeepEqual(unaryResult, tt.expected) {
				t.Errorf("Expected %+v, got %+v", tt.expected, unaryResult)
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
		expected binary
	}{
		// Arithmetic operators
		{
			name:   "Addition",
			tokens: Tokens{numberA, plus, numberB, EOFT},
			expected: binary{
				left:     primary{node: numberA},
				operator: plus,
				right:    primary{node: numberB},
			},
		},
		{
			name:   "Subtraction",
			tokens: Tokens{numberA, minus, numberB, EOFT},
			expected: binary{
				left:     primary{node: numberA},
				operator: minus,
				right:    primary{node: numberB},
			},
		},
		{
			name:   "Multiplication",
			tokens: Tokens{numberA, star, numberB, EOFT},
			expected: binary{
				left:     primary{node: numberA},
				operator: star,
				right:    primary{node: numberB},
			},
		},
		{
			name:   "Division",
			tokens: Tokens{numberA, slash, numberB, EOFT},
			expected: binary{
				left:     primary{node: numberA},
				operator: slash,
				right:    primary{node: numberB},
			},
		},

		// String concatenation
		{
			name:   "String concatenation",
			tokens: Tokens{stringA, plus, stringB, EOFT},
			expected: binary{
				left:     primary{node: stringA},
				operator: plus,
				right:    primary{node: stringB},
			},
		},

		// Comparison operators
		{
			name:   "Greater than",
			tokens: Tokens{numberA, greater, numberB, EOFT},
			expected: binary{
				left:     primary{node: numberA},
				operator: greater,
				right:    primary{node: numberB},
			},
		},
		{
			name:   "Greater than or equal",
			tokens: Tokens{numberA, greaterEqual, numberB, EOFT},
			expected: binary{
				left:     primary{node: numberA},
				operator: greaterEqual,
				right:    primary{node: numberB},
			},
		},
		{
			name:   "Less than",
			tokens: Tokens{numberA, less, numberB, EOFT},
			expected: binary{
				left:     primary{node: numberA},
				operator: less,
				right:    primary{node: numberB},
			},
		},
		{
			name:   "Less than or equal",
			tokens: Tokens{numberA, lessEqual, numberB, EOFT},
			expected: binary{
				left:     primary{node: numberA},
				operator: lessEqual,
				right:    primary{node: numberB},
			},
		},
		{
			name:   "Equal",
			tokens: Tokens{idA, equalEqual, idB, EOFT},
			expected: binary{
				left:     primary{node: idA},
				operator: equalEqual,
				right:    primary{node: idB},
			},
		},
		{
			name:   "Not equal",
			tokens: Tokens{idA, bangEqual, idB, EOFT},
			expected: binary{
				left:     primary{node: idA},
				operator: bangEqual,
				right:    primary{node: idB},
			},
		},
		// Logical operators
		/*
			{
				name:   "Logical AND",
				tokens: Tokens{trueT, and, falseT, EOFT},
				expected: binary{
					left:     primary{node: trueT},
					operator: and,
					right:    primary{node: falseT},
				},
			},
			{
				name:   "Logical OR",
				tokens: Tokens{trueT, or, falseT, EOFT},
				expected: binary{
					left:     primary{node: trueT},
					operator: or,
					right:    primary{node: falseT},
				},
			},
		*/

		// Complex expressions
		{
			name:   "Chained binary expressions (1 + 2 * 3)",
			tokens: Tokens{numberA, plus, numberB, star, makeToken(scanner.NUMBER, "3", 1), EOFT},
			expected: binary{
				left:     primary{node: numberA},
				operator: plus,
				right: binary{
					left:     primary{node: numberB},
					operator: star,
					right:    primary{node: makeToken(scanner.NUMBER, "3", 1)},
				},
			},
		},
		// Mixed expressions
		{
			name:   "Binary with Unary right operand",
			tokens: Tokens{numberA, plus, minus, numberB, EOFT},
			expected: binary{
				left:     primary{node: numberA},
				operator: plus,
				right: unary{
					operator: minus,
					right:    primary{node: numberB},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.tokens.expression()
			if err != nil {
				t.Fatalf("Expected binary expression, got error: %s", err.Error())
			}

			// Type check
			binResult, ok := result.(binary)
			if !ok {
				t.Fatalf("Expected binary expression, got %T", result)
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
			expected: primary{node: numberT},
		},
		{
			name:   "Grouped binary expression",
			tokens: Tokens{leftParen, numberT, plus, number2T, rightParen, EOFT},
			expected: binary{
				left:     primary{node: numberT},
				operator: plus,
				right:    primary{node: number2T},
			},
		},
		{
			name:     "Nested grouping",
			tokens:   Tokens{leftParen, leftParen, numberT, rightParen, rightParen, EOFT},
			expected: primary{node: numberT},
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
			expected: binary{
				left:     primary{node: num1},
				operator: plus,
				right: binary{
					left:     primary{node: num2},
					operator: star,
					right:    primary{node: num3},
				},
			},
			description: "1 + 2 * 3 should evaluate as 1 + (2 * 3)",
		},
		{
			name:   "Division before subtraction",
			tokens: Tokens{num1, minus, num2, slash, num3, EOFT},
			expected: binary{
				left:     primary{node: num1},
				operator: minus,
				right: binary{
					left:     primary{node: num2},
					operator: slash,
					right:    primary{node: num3},
				},
			},
			description: "1 - 2 / 3 should evaluate as 1 - (2 / 3)",
		},
		/*
			 * 	{
					name:   "Comparison before logical AND",
					tokens: Tokens{num1, greater, num2, and, num3, less, num4, EOFT},
					expected: binary{
						left: binary{
							left:     primary{node: num1},
							operator: greater,
							right:    primary{node: num2},
						},
						operator: and,
						right: binary{
							left:     primary{node: num3},
							operator: less,
							right:    primary{node: num4},
						},
					},
					description: "1 > 2 and 3 < 4 should evaluate as (1 > 2) and (3 < 4)",
				},
		*/
		{
			name:   "Parentheses override precedence",
			tokens: Tokens{leftParen, num1, plus, num2, rightParen, star, num3, EOFT},
			expected: binary{
				left: binary{
					left:     primary{node: num1},
					operator: plus,
					right:    primary{node: num2},
				},
				operator: star,
				right:    primary{node: num3},
			},
			description: "(1 + 2) * 3 should evaluate as (1 + 2) * 3",
		},

		{
			name:   "Nested Parentheses override precedence",
			tokens: Tokens{leftParen, leftParen, num1, plus, num2, rightParen, rightParen, star, num3, EOFT},
			expected: binary{
				left: binary{
					left:     primary{node: num1},
					operator: plus,
					right:    primary{node: num2},
				},
				operator: star,
				right:    primary{node: num3},
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
