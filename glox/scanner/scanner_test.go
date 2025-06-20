package scanner

import (
	"testing"
)

// TestEmptySource ensures the scanner returns only EOF token for empty source
func TestEmptySource(t *testing.T) {
	source := ""
	tokens, _ := ScanTokens([]byte(source))

	if len(tokens) != 1 {
		t.Fatalf("Expected 1 token (EOF), got %d", len(tokens))
	}

	if tokens[0].Ttype != EOF {
		t.Errorf("Expected EOF token, got %v", tokens[0].Ttype)
	}
}

// TestSingleCharacterTokens tests recognition of all single-character tokens
func TestSingleCharacterTokens(t *testing.T) {
	source := "(){},.-+;*/"
	tokens, err := ScanTokens([]byte(source))

	if err != nil {
		t.Fatalf("Did not expect any errors but, got %s", err.Error())
	}

	// +1 for EOF token
	expectedCount := len(source) + 1
	if len(tokens) != expectedCount {
		t.Fatalf("Expected %d tokens, got %d", expectedCount, len(tokens))
	}

	expectedTypes := []Tokentype{
		LEFT_PAREN, RIGHT_PAREN, LEFT_BRACE, RIGHT_BRACE, COMMA,
		DOT, MINUS, PLUS, SEMICOLON, STAR, SLASH, EOF,
	}

	for i, expected := range expectedTypes {
		if i < len(tokens) && tokens[i].Ttype != expected {
			t.Errorf("Expected token[%d] to be %v, got %v", i, expected, tokens[i].Ttype)
		}
	}
}

// TestOneOrTwoCharacterTokens tests single and double character tokens
func TestOneOrTwoCharacterTokens(t *testing.T) {
	testCases := []struct {
		source   string
		expected []Tokentype
	}{
		{"!", []Tokentype{BANG, EOF}},
		{"!=", []Tokentype{BANG_EQUAL, EOF}},
		{"=", []Tokentype{EQUAL, EOF}},
		{"==", []Tokentype{EQUAL_EQUAL, EOF}},
		{"<", []Tokentype{LESS, EOF}},
		{"<=", []Tokentype{LESS_EQUAL, EOF}},
		{">", []Tokentype{GREATER, EOF}},
		{">=", []Tokentype{GREATER_EQUAL, EOF}},
	}

	for _, tc := range testCases {
		t.Run(tc.source, func(t *testing.T) {

			tokens, _ := ScanTokens([]byte(tc.source))

			if len(tokens) != len(tc.expected) {
				t.Fatalf("Expected %d tokens, got %d (%d-%d-%d)", len(tc.expected), len(tokens), tokens[0].Ttype, tokens[1].Ttype, tokens[2].Ttype)
			}

			for i, expected := range tc.expected {
				if tokens[i].Ttype != expected {
					t.Errorf("Expected token[%d] to be %v, got %v", i, expected, tokens[i].Ttype)
				}
			}
		})
	}
}

// TestLineComments tests handling of line comments
func TestLineComments(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		expected []Tokentype
	}{
		{
			"Single line comment",
			"// This is a comment\n42",
			[]Tokentype{NUMBER, EOF},
		},
		{
			"Comment at end of file",
			"42 // This is a comment",
			[]Tokentype{NUMBER, EOF},
		},
		{
			"Comment between tokens",
			"42 // Comment\n \"string\"",
			[]Tokentype{NUMBER, STRING, EOF},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//scanner := ScanTokens(tc.source)
			tokens, _ := ScanTokens([]byte(tc.source))

			if len(tokens) != len(tc.expected) {
				t.Fatalf("Expected %d tokens, got %d", len(tc.expected), len(tokens))
			}

			for i, expected := range tc.expected {
				if tokens[i].Ttype != expected {
					t.Errorf("Expected token[%d] to be %v, got %v", i, expected, tokens[i].Ttype)
				}
			}
		})
	}
}

// TestBlockComments tests handling of block comments
func TestBlockComments(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		expected []Tokentype
	}{
		{
			"Simple block comment",
			"/* This is a block comment */42",
			[]Tokentype{NUMBER, EOF},
		},
		{
			"Empty block comment",
			"/**/42",
			[]Tokentype{NUMBER, EOF},
		},
		{
			"Multiline block comment",
			"/* This is a\nmultiline\nblock comment\n*/42",
			[]Tokentype{NUMBER, EOF},
		},
		{
			"Block comment at end of file",
			"42/* This is a comment */",
			[]Tokentype{NUMBER, EOF},
		},
		{
			"Block comment between tokens",
			"42/* Comment */ \"string\"",
			[]Tokentype{NUMBER, STRING, EOF},
		},
		{
			"Nested-looking block comment",
			"/* outer /* inner */ 42*/42",
			[]Tokentype{NUMBER, EOF},
		},
		{
			"Block comment with asterisks",
			"/* ** * ** */42",
			[]Tokentype{NUMBER, EOF},
		},
		{
			"Block comment with slashes",
			"/* // /// */42",
			[]Tokentype{NUMBER, EOF},
		},
		{
			"Mixed line and block comments",
			"// Line comment\n/* Block comment */\n42",
			[]Tokentype{NUMBER, EOF},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, _ := ScanTokens([]byte(tc.source))

			if len(tokens) != len(tc.expected) {
				t.Fatalf("Expected %d tokens, got %d", len(tc.expected), len(tokens))
			}

			for i, expected := range tc.expected {
				if tokens[i].Ttype != expected {
					t.Errorf("Expected token[%d] to be %v, got %v", i, expected, tokens[i].Ttype)
				}
			}
		})
	}
}

// TestNestedBlockComments tests handling of nested block comments
func TestNestedBlockComments(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		expected []Tokentype
	}{
		{
			"Simple nested block comment",
			"/* outer /* inner */ outer */42",
			[]Tokentype{NUMBER, EOF},
		},
		{
			"Multiple levels of nesting",
			"/* level1 /* level2 /* level3 */ level2 */ level1 */42",
			[]Tokentype{NUMBER, EOF},
		},
		{
			"Mixed nesting with line breaks",
			"/* level1\n/* level2 */\nlevel1 */42",
			[]Tokentype{NUMBER, EOF},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, _ := ScanTokens([]byte(tc.source))

			if len(tokens) != len(tc.expected) {
				t.Fatalf("Expected %d tokens, got %d", len(tc.expected), len(tokens))
			}

			for i, expected := range tc.expected {
				if tokens[i].Ttype != expected {
					t.Errorf("Expected token[%d] to be %v, got %v", i, expected, tokens[i].Ttype)
				}
			}
		})
	}
}

// TestIncompleteBlockComments tests error handling for unterminated block comments
func TestIncompleteBlockComments(t *testing.T) {
	testCases := []struct {
		name      string
		source    string
		expectErr bool
	}{
		{
			"Unterminated block comment",
			"/* This comment never ends",
			true,
		},
		{
			"Unterminated nested block comment",
			"/* outer /* inner */",
			true,
		},
		{
			"Almost complete block comment",
			"/* This comment almost ends *",
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ScanTokens([]byte(tc.source))

			hasErrors := err != nil
			if tc.expectErr != hasErrors {
				if tc.expectErr {
					t.Error("Expected scanner to report errors, but none were reported")
				} else {
					t.Errorf("Expected no errors, but got: %v", err)
				}
			}
		})
	}
}

// TestWhitespace tests that whitespace is properly ignored
func TestWhitespace(t *testing.T) {
	source := " \r\n\t42  \n  53"

	tokens, _ := ScanTokens([]byte(source))

	expectedTypes := []Tokentype{NUMBER, NUMBER, EOF}
	expectedLexemes := []string{"42", "53", ""}

	if len(tokens) != len(expectedTypes) {
		t.Fatalf("Expected %d tokens, got %d", len(expectedTypes), len(tokens))
	}

	for i, expected := range expectedTypes {
		if tokens[i].Ttype != expected {
			t.Errorf("Expected token[%d] to be %v, got %v", i, expected, tokens[i].Ttype)
		}
		if tokens[i].Lexem != expectedLexemes[i] {
			t.Errorf("Expected lexeme[%d] to be %q, got %q", i, expectedLexemes[i], tokens[i].Lexem)
		}
	}
}

// TestStrings tests string literals
func TestStrings(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		expected []Tokentype
		lexemes  []string
	}{
		{
			"Simple string",
			"\"hello world\"",
			[]Tokentype{STRING, EOF},
			[]string{"hello world", ""},
		},
		{
			"Empty string",
			"\"\"",
			[]Tokentype{STRING, EOF},
			[]string{"", ""},
		},
		{
			"Multiple strings",
			"\"one\" \"two\" \"three\"",
			[]Tokentype{STRING, STRING, STRING, EOF},
			[]string{"one", "two", "three", ""},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			tokens, _ := ScanTokens([]byte(tc.source))

			if len(tokens) != len(tc.expected) {
				t.Fatalf("Expected %d tokens, got %d", len(tc.expected), len(tokens))
			}

			for i, token := range tokens {
				if token.Ttype != tc.expected[i] {
					t.Errorf("Expected token[%d] to be %v, got %v", i, tc.expected[i], token.Ttype)
				}
				if token.Lexem != tc.lexemes[i] {
					t.Errorf("Expected lexeme[%d] to be %q, got %q", i, tc.lexemes[i], token.Lexem)
				}
			}
		})
	}
}

// TestNumbers tests number literals
func TestNumbers(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		expected []Tokentype
		lexemes  []string
	}{
		{
			"Integer",
			"123",
			[]Tokentype{NUMBER, EOF},
			[]string{"123", ""},
		},
		{
			"Float",
			"123.456",
			[]Tokentype{NUMBER, EOF},
			[]string{"123.456", ""},
		},
		{
			"Multiple numbers",
			"42 3.14",
			[]Tokentype{NUMBER, NUMBER, EOF},
			[]string{"42", "3.14", ""},
		},
		{
			"Leading zero",
			"0.5",
			[]Tokentype{NUMBER, EOF},
			[]string{"0.5", ""},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, _ := ScanTokens([]byte(tc.source))

			if len(tokens) != len(tc.expected) {
				t.Fatalf("Expected %d tokens, got %d", len(tc.expected), len(tokens))
			}

			for i, token := range tokens {
				if token.Ttype != tc.expected[i] {
					t.Errorf("Expected token[%d] to be %v, got %v", i, tc.expected[i], token.Ttype)
				}
				if token.Lexem != tc.lexemes[i] {
					t.Errorf("Expected lexeme[%d] to be %q, got %q", i, tc.lexemes[i], token.Lexem)
				}
			}
		})
	}
}

// TestIdentifiers tests keywords and identifiers
func TestIdentifiers(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		expected []Tokentype
		lexemes  []string
	}{
		{
			"Keywords",
			"and class else false for fun if nil or print return super this true var while",
			[]Tokentype{
				AND, CLASS, ELSE, FALSE, FOR, FUN, IF, NIL, OR, PRINT,
				RETURN, SUPER, THIS, TRUE, VAR, WHILE, EOF,
			},
			[]string{
				"and", "class", "else", "false", "for", "fun", "if", "nil", "or", "print",
				"return", "super", "this", "true", "var", "while", "",
			},
		},
		{
			"Identifiers",
			"foo bar baz",
			[]Tokentype{IDENTIFIER, IDENTIFIER, IDENTIFIER, EOF},
			[]string{"foo", "bar", "baz", ""},
		},
		{
			"Mixed",
			"var foo = true; if bar",
			[]Tokentype{VAR, IDENTIFIER, EQUAL, TRUE, SEMICOLON, IF, IDENTIFIER, EOF},
			[]string{"var", "foo", "", "true", "", "if", "bar", ""},
		},
		{
			"Identifiers with numbers",
			"count123 x1 y2",
			[]Tokentype{IDENTIFIER, IDENTIFIER, IDENTIFIER, EOF},
			[]string{"count123", "x1", "y2", ""},
		},
		{
			"Identifiers with underscores",
			"_foo bar_baz _",
			[]Tokentype{IDENTIFIER, IDENTIFIER, IDENTIFIER, EOF},
			[]string{"_foo", "bar_baz", "_", ""},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, _ := ScanTokens([]byte(tc.source))

			if len(tokens) != len(tc.expected) {
				t.Fatalf("Expected %d tokens, got %d", len(tc.expected), len(tokens))
			}

			for i, token := range tokens {
				if token.Ttype != tc.expected[i] {
					t.Errorf("Expected token[%d] to be %v, got %v", i, tc.expected[i], token.Ttype)
				}
				if token.Lexem != tc.lexemes[i] {
					t.Errorf("Expected lexeme[%d] to be %q, got %q", i, tc.lexemes[i], token.Lexem)
				}
			}
		})
	}
}

// TestLoxProgramWithComments tests a Lox program with various types of comments
func TestLoxProgramWithComments(t *testing.T) {
	source := `
// This is a line comment
class Test {
  /* Block comment for initialization method */
  init() {
    this.value = 42; /* Inline block comment */
  }

  /* Multi-line
     block comment for
     getValue method */
  getValue() {
    return this.value; // Return value
  }

  /* Empty block */
  /**/
  /***/
}

/* Block comment before var declaration */ var test = Test();
print test.getValue(); /* Should print 42 */
`
	tokens, _ := ScanTokens([]byte(source))

	expectedTypes := []Tokentype{
		CLASS, IDENTIFIER, LEFT_BRACE,
		IDENTIFIER, LEFT_PAREN, RIGHT_PAREN, LEFT_BRACE,
		THIS, DOT, IDENTIFIER, EQUAL, NUMBER, SEMICOLON,
		RIGHT_BRACE,
		IDENTIFIER, LEFT_PAREN, RIGHT_PAREN, LEFT_BRACE,
		RETURN, THIS, DOT, IDENTIFIER, SEMICOLON,
		RIGHT_BRACE,
		RIGHT_BRACE,
		VAR, IDENTIFIER, EQUAL, IDENTIFIER, LEFT_PAREN, RIGHT_PAREN, SEMICOLON,
		PRINT, IDENTIFIER, DOT, IDENTIFIER, LEFT_PAREN, RIGHT_PAREN, SEMICOLON,
		EOF,
	}

	if len(tokens) != len(expectedTypes) {
		t.Fatalf("Expected %d tokens, got %d", len(expectedTypes), len(tokens))
	}

	for i, expected := range expectedTypes {
		if tokens[i].Ttype != expected {
			t.Errorf("Expected token[%d] to be %v, got %v", i, expected, tokens[i].Ttype)
		}
	}
}

// TestComplexLoxProgram tests scanning a more complex Lox program with multiple features
func TestComplexLoxProgram(t *testing.T) {
	source := `
/*
 * Fibonacci function implementation
 * Calculates the nth Fibonacci number recursively
 */
fun fibonacci(n) {
  if (n <= 1) return n;
  return fibonacci(n - 2) + fibonacci(n - 1);
}

// Print first 10 Fibonacci numbers
for (var i = 0; i < 10; i = i + 1) {
  print fibonacci(i);
}

/* Counter class implementation */
class Counter {
  init() {
    this.count = 0; /* Initialize counter */
  }

  /* Increment method */
  increment() {
    this.count = this.count + 1; // Add 1 to counter
    return this.count;
  }
}

var counter = Counter();
print "Count: " + counter.increment(); /* Should print "Count: 1" */
`

	tokens, err := ScanTokens([]byte(source))

	// We won't test each individual token, but verify that the scanner
	// produces a reasonable number of tokens and no errors
	if len(tokens) < 50 {
		t.Errorf("Expected at least 50 tokens, got %d", len(tokens))
	}

	if err != nil {
		t.Errorf("Expected no errors, but got: %v", err)
	}

	// Verify last token is EOF
	lastToken := tokens[len(tokens)-1]
	if lastToken.Ttype != EOF {
		t.Errorf("Expected last token to be EOF, got %v", lastToken.Ttype)
	}
}

// TestLineTracking verifies that the scanner correctly tracks line numbers
func TestLineTracking(t *testing.T) {
	source := `line 1
line 2
/* Block comment
   spanning
   multiple lines */
line 6`
	tokens, _ := ScanTokens([]byte(source))

	// Should have 6 tokens: "line" "1" "line" "2" "line" "6" EOF
	expectedLines := []int{1, 1, 2, 2, 6, 6, 6}

	if len(tokens) != 7 {
		t.Fatalf("Expected 7 tokens, got %d", len(tokens))
	}

	for i, token := range tokens {
		if token.Line != expectedLines[i] {
			t.Errorf("Expected token[%d] to have line %d, got %d", i, expectedLines[i], token.Line)
		}
	}
}

// TestMultilineString tests that the scanner handles multiline strings correctly
func TestMultilineString(t *testing.T) {
	source := "\"This is a\nmultiline\nstring\""
	//scanner := ScanTokens(source)
	tokens, _ := ScanTokens([]byte(source))

	if len(tokens) != 2 { // STRING + EOF
		t.Fatalf("Expected 2 tokens, got %d", len(tokens))
	}

	if tokens[0].Ttype != STRING {
		t.Errorf("Expected STRING token, got %v", tokens[0].Ttype)
	}

	if tokens[0].Line != 3 {
		t.Errorf("Expected string token to end on line 3, got %d", tokens[0].Line)
	}
}

// TestCommentsAndLineTracking tests that the scanner correctly handles line tracking with comments
func TestCommentsAndLineTracking(t *testing.T) {
	source := `line 1
// Comment line
line 3
/* Block comment line 4
   Block comment line 5
   Block comment line 6 */
line 7`

	tokens, _ := ScanTokens([]byte(source))

	// Should have 6 tokens: "line" "1" "line" "3" "line" "7" EOF
	expectedLines := []int{1, 1, 3, 3, 7, 7, 7}
	expectedLexemes := []string{"line", "1", "line", "3", "line", "7", ""}

	if len(tokens) != len(expectedLexemes) {
		t.Fatalf("Expected 7 tokens, got %d", len(tokens))
	}

	for i, token := range tokens {
		if token.Line != expectedLines[i] {
			t.Errorf("Expected token[%d] to have line %d, got %d", i, expectedLines[i], token.Line)
		}
		if token.Lexem != expectedLexemes[i] {
			t.Errorf("Expected lexeme[%d] to be %q, got %q", i, expectedLexemes[i], token.Lexem)
		}
	}
}
