#include "compiler.h"
#include "scanner.h"
#include "parser.h"

Parser parser;
CompilerStatus compile(const char *src){
    //set the ball rolling
    advance();
    expression();
    //expect parser.current to be "semi colon" token
    return !parser.hadError;
}
