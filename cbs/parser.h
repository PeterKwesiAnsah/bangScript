//will contain Patt parser implementation
#ifndef PARSER_H
#define PARSER_H
#include "readonly.h"
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

typedef void (*ParseFn)(bool isAssignExp);
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

void errorAt(Token *, const char *,const char *);

//production rules for expressions
static void number(bool);
static void string(bool);
static void grouping(bool);
static void unary(bool);
static void binary(bool);
static void boolean(bool);
static void identifier(bool);
static void nil(bool);
#endif
