//will contain Patt parser implementation
#ifndef PARSER_H
#define PARSER_H
#include "scanner.h"
#include "stdbool.h"
typedef struct {
    Token current,previous;
    bool hadError;
} Parser ;
void expression();
void advance();

//production rules for expressions
static void number();
static void grouping();
static void unary();
static void binary();
#endif
