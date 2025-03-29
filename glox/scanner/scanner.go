/*
	TODO: Scanner or lexer reads characters from a file stream,creating lexems in the process and producing tokens
   	source_code ----->Scanner--->Lexems--->Tokens
    1. Producing Lexems
    Use buffer lexems (heap allocated)-as we will be creating tokens.
    	-Keep track of the line number.
     	-Keep track of the offset(start and end) of a lexem.

    2.Identifying Lexems
    Use regex exp. to identify lexems as either
    // Single-character tokens.
    LEFT_PAREN, RIGHT_PAREN, LEFT_BRACE, RIGHT_BRACE,
    COMMA, DOT, MINUS, PLUS, SEMICOLON, SLASH, STAR,
    // One or two character tokens.
    BANG, BANG_EQUAL,
    EQUAL, EQUAL_EQUAL,
    GREATER, GREATER_EQUAL,
    LESS, LESS_EQUAL,
    // Literals.
    IDENTIFIER, STRING, NUMBER,
    // Keywords.
    AND, CLASS, ELSE, FALSE, FUN, FOR, IF, NIL, OR,
    PRINT, RETURN, SUPER, THIS, TRUE, VAR, WHILE,
    EOF

    3. Create Token
    Use Lexem, Lexem Type to create token
    Token is a struct with members
    	-Lexem
     	-Lexem/Token Type
      	-Line
    4. Terminate tokens
    	Finally add an EOF token to mark the end of the file stream.
*/

package scanner

type Tokentype int
type Token struct {
	//Represents the token type of the word
	ttype Tokentype

	// line represents the line number in which the token was emitted.
	line int

	lexem string
}

const (
	LEFT_PAREN Tokentype = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR
	// One or two character tokens.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL
	// Literals.
	IDENTIFIER
	STRING
	NUMBER
	// Keywords.
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

type sprop struct {
	// line represents the current line number in the source code.
	line int

	// start is the index where the current token starts.
	start int

	// current is the index of the character currently being scanned.
	current int
}

func (sp *sprop) ScanTokens(source []byte) []*Token {
	tokens := make([]*Token, 0, 5)
	return tokens
}
