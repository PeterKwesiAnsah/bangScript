#include "compiler.h"
#include "chunk.h"
#include "parser.h"
#include "readonly.h"
#include <stdbool.h>
#include "scanner.h"
#include "table.h"
#include <stdint.h>

#define IN_SCOPE_LOCALS_LIMIT (UINT8_MAX+1)


extern Parser parser;
extern Table strings;
extern Chunk chunk;


extern inline size_t internString(Table *, Token , const char *);

typedef struct {
    Token name;
    unsigned depth;
} Local;

typedef struct {
    Local locals[IN_SCOPE_LOCALS_LIMIT];
    unsigned len;
    unsigned scopeDepth;
} Compiler;

Compiler current;

CompilerStatus compile(const char *src){

    Tinit(&strings);
    addConstant(C_DOUBLE_TO_BS_NUMBER(0));
    addConstant(C_BOOL_TO_BS_BOOLEAN(true));
    addConstant(C_BOOL_TO_BS_BOOLEAN(false));
    addConstant((Value){.type=TYPE_NIL,.value={}});

    //set the ball rolling
    advance();
    for(;;){
        switch (parser.current.tt) {
            // variable declaration
            case TOKEN_VAR:
            {
                advance();
                Token cur=parser.current;
                if(current.scopeDepth>0){
                    //define local variable as it is in scope
                    current.locals[current.len].name=cur;
                    current.locals[current.len].depth=-1;
                }
                //if(cur.tt!=TOKEN_IDENTIFIER){
                    //compile error
                //}
                advance();
                //var foo=0;
                // Handle multiple global var definations
                if(parser.current.tt==TOKEN_EQUAL){
                    advance();
                    expression();
                }else {
                     //var foo;
                    if(parser.current.tt!=TOKEN_SEMICOLON){
                        fprintf(stderr,"Expected a semi-colon but got %d\n",parser.current.tt);
                        parser.hadError=true;
                        return true;
                    }
                    WRITE_BYTECODE(chunk, OP_CONSTANT, cur.line);
                    WRITE_BYTECODE(chunk, CONSTANT_NIL_INDEX, cur.line);
                    advance();
                }

                if(current.scopeDepth>0){
                    current.locals[current.len++].depth=current.scopeDepth;
                    break;
                }
                //if local var, we skip this???
                size_t BsobjStringConstIndex=internString(&strings, cur, src);
                WRITE_BYTECODE(chunk, OP_GLOBALVAR_DEF, cur.line);
                WRITE_BYTECODE(chunk, BsobjStringConstIndex, cur.line);
            }
            break;
            case TOKEN_PRINT:{
                Token printToken=parser.current;
                advance();
                expression();
                WRITE_BYTECODE(chunk, OP_PRINT, printToken.line);
            }
            break;
            case TOKEN_EOF:{
                WRITE_BYTECODE(chunk, OP_RETURN, parser.current.line);
              return !parser.hadError;
            }
            default:
                expression();
            break;
        }
    }
    return !parser.hadError;
}
