#include "compiler.h"
#include "parser.h"
#include "readonly.h"


CompilerStatus compile(const char *src){
    extern Parser parser;
    addConstant(C_DOUBLE_TO_BS_NUMBER(0));
    //Include , true , false and nil
    //set the ball rolling
    advance();
    expression();
    //expect parser.current to be "semi colon" token
    return !parser.hadError;
}
