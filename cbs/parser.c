#include "parser.h"
#include "chunk.h"
#include "darray.h"
#include "scanner.h"
#include <stdbool.h>
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include "line.h"
#include "readonly.h"
#define WRITE_BYTECODE(chunk,byte,line) do{ \
append(chunk, (u_int8_t)byte, sizeof(u_int8_t));\
addLine(line);\
}while(0)


extern const char *src;
extern const char *scanerr;
extern DECLARE_ARRAY(u_int8_t, chunk);




ParseRule rules[] = {
    [TOKEN_LEFT_PAREN]    = {grouping, NULL,   PREC_NONE},
    [TOKEN_RIGHT_PAREN]   = {NULL,     NULL,   PREC_NONE},
    [TOKEN_LEFT_BRACE]    = {NULL,     NULL,   PREC_NONE},
    [TOKEN_RIGHT_BRACE]   = {NULL,     NULL,   PREC_NONE},
    [TOKEN_COMMA]         = {NULL,     NULL,   PREC_NONE},
    [TOKEN_DOT]           = {NULL,     NULL,   PREC_NONE},
    [TOKEN_MINUS]         = {unary,    binary, PREC_TERM},
    [TOKEN_PLUS]          = {NULL,     binary, PREC_TERM},
    [TOKEN_SEMICOLON]     = {NULL,     NULL,   PREC_NONE},
    [TOKEN_SLASH]         = {NULL,     binary, PREC_FACTOR},
    [TOKEN_STAR]          = {NULL,     binary, PREC_FACTOR},
    [TOKEN_BANG]          = {NULL,     NULL,   PREC_NONE},
    [TOKEN_BANG_EQUAL]    = {NULL,     NULL,   PREC_NONE},
    [TOKEN_EQUAL]         = {NULL,     NULL,   PREC_NONE},
    [TOKEN_EQUAL_EQUAL]   = {NULL,     NULL,   PREC_NONE},
    [TOKEN_GREATER]       = {NULL,     NULL,   PREC_NONE},
    [TOKEN_GREATER_EQUAL] = {NULL,     NULL,   PREC_NONE},
    [TOKEN_LESS]          = {NULL,     NULL,   PREC_NONE},
    [TOKEN_LESS_EQUAL]    = {NULL,     NULL,   PREC_NONE},
    [TOKEN_IDENTIFIER]    = {NULL,     NULL,   PREC_NONE},
    [TOKEN_STRING]        = {NULL,     NULL,   PREC_NONE},
    [TOKEN_NUMBER]        = {number,   NULL,   PREC_NONE},
    [TOKEN_AND]           = {NULL,     NULL,   PREC_NONE},
    [TOKEN_CLASS]         = {NULL,     NULL,   PREC_NONE},
    [TOKEN_ELSE]          = {NULL,     NULL,   PREC_NONE},
    [TOKEN_FALSE]         = {NULL,     NULL,   PREC_NONE},
    [TOKEN_FOR]           = {NULL,     NULL,   PREC_NONE},
    [TOKEN_FUN]           = {NULL,     NULL,   PREC_NONE},
    [TOKEN_IF]            = {NULL,     NULL,   PREC_NONE},
    [TOKEN_NIL]           = {NULL,     NULL,   PREC_NONE},
    [TOKEN_OR]            = {NULL,     NULL,   PREC_NONE},
    [TOKEN_PRINT]         = {NULL,     NULL,   PREC_NONE},
    [TOKEN_RETURN]        = {NULL,     NULL,   PREC_NONE},
    [TOKEN_SUPER]         = {NULL,     NULL,   PREC_NONE},
    [TOKEN_THIS]          = {NULL,     NULL,   PREC_NONE},
};


Parser parser={{},{},false};
void advance(){
    parser.previous=parser.current;
    for(;;){
        parser.current=scanTokens(src);
        if(parser.current.tt!=TOKEN_ERROR) break;
        fprintf(stderr, "%s at line %d", scanerr,parser.current.line);
        parser.hadError=true;
    }
}
static void parsePrecedence(Precedence precedence) {
    advance();
    ParseFn prefixRule = rules[parser.previous.tt].prefix;
    if (prefixRule == NULL) {
        fputs("Invalid Expression",stderr);
        return;
    }
prefixRule();

while (precedence <= rules[parser.previous.tt].precedence) {
    advance();
    ParseFn infixRule = rules[parser.previous.tt].infix;
    infixRule();
    }
}

void expression(){
    parsePrecedence(PREC_ASSIGNMENT);
}
void grouping(){
    expression();
    if(parser.current.tt!=TOKEN_RIGHT_PAREN){
        fprintf(stderr,"Expected Right Paren but got %d",parser.current.tt);
        return;
    }
    advance();
}

void unary (){
    Token token=parser.previous;
    parsePrecedence(PREC_UNARY);
    switch (token.tt) {
        case TOKEN_MINUS:
        {
            WRITE_BYTECODE(chunk, OP_CONSTANT,0 );
            WRITE_BYTECODE(chunk, 0,0 );
            WRITE_BYTECODE(chunk, OP_SUB,0);
        }
        break;
        default:
            return;
    }
}

static void number(){
    Token numToken=parser.previous;
    size_t constantIndex=addConstant(atof(src+numToken.start));
    if(constantIndex >= CONSTANT_LIMIT){
        // Write opcode
        WRITE_BYTECODE(chunk, OP_CONSTANT_LONG, numToken.line);
        // Write operand as 3 bytes
        WRITE_BYTECODE(chunk, constantIndex & 0xFF, numToken.line);
        WRITE_BYTECODE(chunk, (constantIndex >> 8) & 0xFF, numToken.line);
        WRITE_BYTECODE(chunk, (constantIndex >> 16) & 0xFF, numToken.line);
    }
    //write opCode
    WRITE_BYTECODE(chunk, OP_CONSTANT, numToken.line);
    //write operand Index
    WRITE_BYTECODE(chunk, constantIndex, numToken.line);
}
