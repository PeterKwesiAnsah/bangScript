#include "compiler.h"
#include "scanner.h"
#include "parser.h"
#include "readonly.h"

Parser parser;
CompilerStatus compile(const char *src){
    addConstant(C_DOUBLE_TO_BS_NUMBER(0));
    //Include , true , false and nil
    //set the ball rolling
    advance();
    expression();
    //expect parser.current to be "semi colon" token
    return !parser.hadError;
}
