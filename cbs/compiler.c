#include "compiler.h"
#include "chunk.h"
#include "parser.h"
#include "readonly.h"
#include <stdbool.h>
#include "scanner.h"
#include "table.h"


extern Parser parser;
extern Table strings;
extern Chunk chunk;

extern inline size_t internString(Table *, Token , const char *);

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
                //if(cur.tt!=TOKEN_IDENTIFIER){
                    //compile error
                //}
                advance();
                //var foo=0;
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
                    WRITE_BYTECODE(chunk, OP_CONSTANT, 0);
                    WRITE_BYTECODE(chunk, CONSTANT_NIL_INDEX, 0);
                    advance();
                }

                size_t BsobjStringConstIndex=internString(&strings, cur, src);
                WRITE_BYTECODE(chunk, OP_GLOBALVAR_DEF, cur.line);
                WRITE_BYTECODE(chunk, BsobjStringConstIndex, 0);
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
                WRITE_BYTECODE(chunk, OP_RETURN, 0);
              return !parser.hadError;
            }
            default:
                expression();
            break;
        }
    }



    return !parser.hadError;
}
