#include "compiler.h"
#include "parser.h"
#include "readonly.h"
#include <stdbool.h>
#include "scanner.h"
#include "table.h"


CompilerStatus compile(const char *src){
    extern Parser parser;
    extern Table strings;
    Tinit(&strings);
    addConstant(C_DOUBLE_TO_BS_NUMBER(0));
    addConstant(C_BOOL_TO_BS_BOOLEAN(true));
    addConstant(C_BOOL_TO_BS_BOOLEAN(false));
    addConstant((Value){.type=TYPE_NIL,.value={}});

    //set the ball rolling
    advance();
    switch (parser.current.tt) {
        // variable declaration
        case TOKEN_VAR:
        {
            advance();
            Token cur=parser.current;
            if(cur.tt!=TOKEN_IDENTIFIER){
                //compile error
            }
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
                advance();
            }



        }
        break;


        default:
            expression();
        break;

    }


    return !parser.hadError;
}
