//will contain Patt parser implementation
#ifndef PARSER_H
#define PARSER_H
#include "scanner.h"
#include "stdbool.h"

typedef enum {
PREC_NONE,
PREC_ASSIGNMENT, // =
PREC_OR, // or
PREC_AND, // and
PREC_EQUALITY, // == !=
PREC_COMPARISON, // < > <= >=
PREC_TERM, // + -
PREC_FACTOR, // * /
PREC_UNARY, // ! -
PREC_CALL, // . ()
PREC_PRIMARY
} Precedence;

typedef void (*ParseFn)();
typedef struct {
ParseFn prefix;
ParseFn infix;
Precedence precedence;
} ParseRule;

typedef struct {
    Token current,previous;
    bool hadError;
} Parser ;
void expression();

static void parsePrecedence(Precedence precedence);
void advance();

//production rules for expressions
static void number();
static void grouping();
static void unary();
static void binary();
#endif
