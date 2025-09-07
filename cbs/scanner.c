#include "scanner.h"
#include <ctype.h>
#include <stdio.h>
#include <stdlib.h>

struct {
    size_t fpos;
    unsigned int line;
    size_t start;
    signed long cur;
} state={0,0,0,-1};

Token scanTokens(FILE *src){
    state.start=state.fpos;

    char c=0;
    //skip white spaces
    while ((c=fgetc(src),state.cur++,state.fpos++,isspace(c))) {
        if(c=='\n')state.line++;;
    }
    switch (c) {
        // Single-character tokens
        case '+':{
            return (Token){state.cur,1,state.line,TOKEN_PLUS};
        }
        case '-':{
            return (Token){state.cur,1,state.line,TOKEN_MINUS};
        }
        case '*':{
            return (Token){state.cur,1,state.line,TOKEN_STAR};
        }
        case '(':{
            return (Token){state.cur,1,state.line,TOKEN_LEFT_PAREN};
        }
        case ')':{
            return (Token){state.cur,1,state.line,TOKEN_RIGHT_PAREN};
        }
        case '{':{
            return (Token){state.cur,1,state.line,TOKEN_LEFT_BRACE};
        }
        case '}':{
            return (Token){state.cur,1,state.line,TOKEN_RIGHT_BRACE};
        }
        case ';':{
            return (Token){state.cur,1,state.line,TOKEN_SEMICOLON};
        }
        case ',':{
            return (Token){state.cur,1,state.line,TOKEN_COMMA};
        }
        case '.':{
            return (Token){state.cur,1,state.line,TOKEN_DOT};
        }
        //Slash or Comments
        case '/':{
            if ((c=fgetc(src)),state.cur++,state.fpos++,!isEOF(c) && c=='/'){
                while ((c=fgetc(src),state.cur++,state.fpos++,c!='\n')){}
                fseek(src, state.cur--, SEEK_SET);
                state.fpos--;
               return scanTokens(src);
            }
            if (isEOF(c)){
                return (Token){state.start,0,state.line,TOKEN_EOF};
            }
            fseek(src, state.cur--, SEEK_SET);
            state.fpos--;
            return (Token){state.start,1,state.line,TOKEN_SLASH};
        }
        // One or two character tokens
        case '!':{
            return (Token){state.cur,1,state.line,TOKEN_BANG};
        }
        case '<':{
            return (Token){state.cur,1,state.line,TOKEN_LESS};
        }
        case '>':{
            return (Token){state.cur,1,state.line,TOKEN_GREATER};
        }
        case '=':{

        }

        case '\0':{
            return (Token){state.start,0,state.line,TOKEN_EOF};
        }

    }
    return (Token){state.start,0,state.line,TOKEN_EOF};
}
