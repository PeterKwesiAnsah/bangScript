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


Parser parser={{},{},false};

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
    [TOKEN_BANG]          = {unary,    NULL,  PREC_NONE},
    [TOKEN_BANG_EQUAL]    = {NULL,     binary, PREC_EQUALITY},
    [TOKEN_EQUAL]         = {NULL,     NULL,   PREC_NONE},
    [TOKEN_EQUAL_EQUAL]   = {NULL,     binary, PREC_EQUALITY},
    [TOKEN_GREATER]       = {NULL,     binary, PREC_COMPARISON},
    [TOKEN_GREATER_EQUAL] = {NULL,     binary,   PREC_COMPARISON},
    [TOKEN_LESS]          = {NULL,     binary, PREC_COMPARISON},
    [TOKEN_LESS_EQUAL]    = {NULL,     binary,   PREC_COMPARISON},
    [TOKEN_IDENTIFIER]    = {NULL,     NULL,   PREC_NONE},
    [TOKEN_STRING]        = {string,   NULL, PREC_NONE},
    [TOKEN_NUMBER]        = {number,   NULL,   PREC_NONE},
    [TOKEN_AND]           = {NULL,     NULL,   PREC_NONE},
    [TOKEN_CLASS]         = {NULL,     NULL,   PREC_NONE},
    [TOKEN_ELSE]          = {NULL,     NULL,   PREC_NONE},
    [TOKEN_TRUE]          = {boolean,  NULL,   PREC_NONE},
    [TOKEN_FALSE]         = {boolean,  NULL,   PREC_NONE},
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
        parser.hadError=true;
        return;
    }
    prefixRule();

while (precedence <= rules[parser.current.tt].precedence) {
    advance();
    ParseFn infixRule = rules[parser.previous.tt].infix;
    infixRule();
    }
}

void expression(){
    parsePrecedence(PREC_ASSIGNMENT);
    if(parser.current.tt!=TOKEN_SEMICOLON){
        fprintf(stderr,"Expected a semi-colon but got %d\n",parser.current.tt);
        parser.hadError=true;
        return;
    }
    advance();
    WRITE_BYTECODE(chunk,OP_PRINT, 0);
    WRITE_BYTECODE(chunk,OP_RETURN, 0);
}
void grouping(){
    parsePrecedence(PREC_ASSIGNMENT);
    if(parser.current.tt!=TOKEN_RIGHT_PAREN){
        fprintf(stderr,"Expected Right Paren but got %d\n",parser.current.tt);
        parser.hadError=true;
        return;
    }
    advance();
}

static void binary() {
    TokenType operatorType = parser.previous.tt;
    unsigned int line=parser.previous.line;
    parsePrecedence((Precedence)(rules[parser.previous.tt].precedence + 1));
    switch (operatorType) {
        case TOKEN_PLUS:
            WRITE_BYTECODE(chunk, OP_ADD,line);
            break;
        case TOKEN_MINUS:
          WRITE_BYTECODE(chunk, OP_SUB,line);
            break;
        case TOKEN_STAR:
            WRITE_BYTECODE(chunk, OP_MUL,line);
            break;
        case TOKEN_SLASH:
          WRITE_BYTECODE(chunk, OP_DIV,line);
            break;
        case TOKEN_EQUAL_EQUAL:
            WRITE_BYTECODE(chunk, OP_EQUAL,line);
            break;
        case TOKEN_BANG_EQUAL:
            WRITE_BYTECODE(chunk, OP_EQUAL_NOT,line);
            break;
        case TOKEN_GREATER:
            WRITE_BYTECODE(chunk, OP_GREATOR,line);
            break;
        case TOKEN_GREATER_EQUAL:
            WRITE_BYTECODE(chunk, OP_LESS_NOT,line);
            break;
        case TOKEN_LESS:
            WRITE_BYTECODE(chunk, OP_LESS,line);
            break;
        case TOKEN_LESS_EQUAL:
            WRITE_BYTECODE(chunk, OP_GREATOR_NOT,line);
            break;

        default:
            return; // Unreachable.
    }
}




void unary (){
    Token token=parser.previous;

    switch (token.tt) {
        case TOKEN_MINUS:
        {
            WRITE_BYTECODE(chunk, OP_CONSTANT_ZER0,0 );
            parsePrecedence(PREC_UNARY);
            WRITE_BYTECODE(chunk, OP_SUB,token.line);
        }
        break;
        default:
            return;
    }
}

static void number(){
    Token numToken=parser.previous;
    size_t constantIndex=addConstant(C_DOUBLE_TO_BS_NUMBER(atof(src+numToken.start)));
    if(constantIndex >= CONSTANT_LIMIT){
        // Write opcode
        WRITE_BYTECODE(chunk, OP_CONSTANT_LONG, numToken.line);
        // Write operand as 3 bytes
        WRITE_BYTECODE(chunk, constantIndex & 0xFF, numToken.line);
        WRITE_BYTECODE(chunk, (constantIndex >> 8) & 0xFF, numToken.line);
        WRITE_BYTECODE(chunk, (constantIndex >> 16) & 0xFF, numToken.line);
        return;
    }
    //write opCode
    WRITE_BYTECODE(chunk, OP_CONSTANT, numToken.line);
    //write operand Index
    WRITE_BYTECODE(chunk, constantIndex, numToken.line);
}
static void string(){
    Token strToken=parser.previous;
    Value BsObjvalue;
    BsObjvalue.type=TYPE_OBJ;
    BsObjStringFromSource *objString=(BsObjStringFromSource *)malloc(sizeof(BsObjStringFromSource));
    objString->obj=(BsObj){.type=OBJ_TYPE_STRING_SOURCE};
    objString->value=src+strToken.start;
    objString->len=strToken.len;
    BsObjvalue.value.obj=(BsObj *) objString;

    size_t stringLiteralIndex=addConstant(BsObjvalue);

    if(stringLiteralIndex >= CONSTANT_LIMIT){
        // Write opcode
        WRITE_BYTECODE(chunk, OP_CONSTANT_LONG, strToken.line);
        // Write operand as 3 bytes
        WRITE_BYTECODE(chunk, stringLiteralIndex & 0xFF, strToken.line);
        WRITE_BYTECODE(chunk, (stringLiteralIndex >> 8) & 0xFF, strToken.line);
        WRITE_BYTECODE(chunk, (stringLiteralIndex >> 16) & 0xFF, strToken.line);
        return;
    }

    //write opCode
    WRITE_BYTECODE(chunk, OP_CONSTANT, strToken.line);
    //write operand Index
    WRITE_BYTECODE(chunk, stringLiteralIndex, strToken.line);
}

static void boolean(){
    Token boolToken=parser.previous;
    //write opCode
    WRITE_BYTECODE(chunk, OP_CONSTANT, boolToken.line);
    //write operand Index
    WRITE_BYTECODE(chunk, (boolToken.tt==TOKEN_TRUE ? CONSTANT_TRUE_BOOL_INDEX:CONSTANT_FALSE_BOOL_INDEX),boolToken.line);
}
