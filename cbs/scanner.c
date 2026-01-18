#include "scanner.h"
#include <ctype.h>
#include <string.h>

struct {
    unsigned int line;
    size_t start;
    signed long cur;
} state={1,0,0};

const char *scanerr;

Token scanTokens(const char *src){


    char c=0;
    //skip white spaces
    while ((c=src[state.cur],isspace(c))) {
        state.cur++;
        if(c=='\n')state.line++;
    }

     state.start=state.cur++;

    switch (c) {
        // Single-character tokens
        case '+':{
            return (Token){state.start,1,state.line,TOKEN_PLUS};
        }
        case '-':{
            return (Token){state.start,1,state.line,TOKEN_MINUS};
        }
        case '*':{
            return (Token){state.start,1,state.line,TOKEN_STAR};
        }
        case '(':{
            return (Token){state.start,1,state.line,TOKEN_LEFT_PAREN};
        }
        case ')':{
            return (Token){state.start,1,state.line,TOKEN_RIGHT_PAREN};
        }
        case '{':{
            return (Token){state.start,1,state.line,TOKEN_LEFT_BRACE};
        }
        case '}':{
            return (Token){state.start,1,state.line,TOKEN_RIGHT_BRACE};
        }
        case ';':{
            return (Token){state.start,1,state.line,TOKEN_SEMICOLON};
        }
        case ',':{
            return (Token){state.start,1,state.line,TOKEN_COMMA};
        }
        case '.':{
            return (Token){state.start,1,state.line,TOKEN_DOT};
        }
        //Slash or Comments
        case '/':{
            if (src[state.cur]== '/'){
                //consume till the end of line
                state.cur++;
                while ((c=src[state.cur],c!='\n' && !isEOF(c))){
                    state.cur++;
                }
                // consume the newline if present
                if (c == '\n') {
                    state.cur++;
                    state.line++;
                }
               return scanTokens(src);
            }
            return (Token){state.start,1,state.line,TOKEN_SLASH};
        }
        // One or two character tokens
        case '!':{
              if ((c=src[state.cur],c=='=')){
                  state.cur++;
                  return (Token){state.start,2,state.line,TOKEN_BANG_EQUAL};
              }
            return (Token){state.start,1,state.line,TOKEN_BANG};
        }
        case '<':{
            if ((c=src[state.cur],c=='=')){
                state.cur++;
                return (Token){state.start,2,state.line,TOKEN_LESS_EQUAL};
            }
            return (Token){state.start,1,state.line,TOKEN_LESS};
        }
        case '>':{
            if ((c=src[state.cur],c=='=')){
                state.cur++;
                return (Token){state.start,2,state.line,TOKEN_GREATER_EQUAL};
            }
            return (Token){state.start,1,state.line,TOKEN_GREATER};
        }
        case '=':{
            if ((c=src[state.cur],c=='=')){
                state.cur++;
                return (Token){state.start,2,state.line,TOKEN_EQUAL_EQUAL};
            }
            return (Token){state.start,1,state.line,TOKEN_EQUAL};
        }
        case '"':{
          while ((c=src[state.cur],!isEOF(c) && c!='"')){
              if(c=='\n'){
                  state.line++;
              }
                state.cur++;
          }
           if(isEOF(c)){
                scanerr="Unterminated String\n";
                return (Token){(size_t)scanerr,0,state.line,TOKEN_ERROR};
           }
            state.cur++;//consume closing "
           return (Token){state.start+1,(state.cur-(state.start+2)),state.line,TOKEN_STRING};
        };
        case '\0':{
            return (Token){state.start,0,state.line,TOKEN_EOF};
        }
        default:{
            if (isalpha(c) || c=='_'){
                while ((c=src[state.cur],isalnum(c) || c=='_')){
                    state.cur++;
                }
                size_t length=state.cur-state.start;
                char buf[length];
                memcpy(buf,src+state.start,length);

                switch (buf[0]) {
                    case 'a': {
                        if(length == 3 && memcmp((buf+1), "nd", 2) == 0) {
                            return (Token){state.start, length, state.line, TOKEN_AND};
                        }
                        break;
                    }
                    case 'c': {
                        if(length == 5 && memcmp((buf+1), "lass", 4) == 0) {
                            return (Token){state.start, length, state.line, TOKEN_CLASS};
                        }
                        break;
                    }
                    case 'e': {
                        if(length == 4 && memcmp((buf+1), "lse", 3) == 0) {
                            return (Token){state.start, length, state.line, TOKEN_ELSE};
                        }
                        break;
                    }
                    case 'f':
                        if (length > 1) {
                            switch (buf[1]) {
                                case 'a': {
                                    if(length == 5 && memcmp((buf+2), "lse", 3) == 0) {
                                        return (Token){state.start, length, state.line, TOKEN_FALSE};
                                    }
                                    break;
                                }
                                case 'o': {
                                    if(length == 3 && memcmp((buf+2), "r", 1) == 0) {
                                        return (Token){state.start, length, state.line, TOKEN_FOR};
                                    }
                                    break;
                                }
                                case 'u': {
                                    if(length == 3 && memcmp((buf+2), "n", 1) == 0) {
                                        return (Token){state.start, length, state.line, TOKEN_FUN};
                                    }
                                    break;
                                }
                            }
                        }
                        break;
                    case 'i': {
                        if(length == 2 && memcmp((buf+1), "f", 1) == 0) {
                            return (Token){state.start, length, state.line, TOKEN_IF};
                        }
                        break;
                    }
                    case 'n': {
                        if(length == 3 && memcmp((buf+1), "il", 2) == 0) {
                            return (Token){state.start, length, state.line, TOKEN_NIL};
                        }
                        break;
                    }
                    case 'o': {
                        if(length == 2 && memcmp((buf+1), "r", 1) == 0) {
                            return (Token){state.start, length, state.line, TOKEN_OR};
                        }
                        break;
                    }
                    case 'p': {
                        if(length == 5 && memcmp((buf+1), "rint", 4) == 0) {
                            return (Token){state.start, length, state.line, TOKEN_PRINT};
                        }
                        break;
                    }
                    case 'r': {
                        if(length == 6 && memcmp((buf+1), "eturn", 5) == 0) {
                            return (Token){state.start, length, state.line, TOKEN_RETURN};
                        }
                        break;
                    }
                    case 's': {
                        if(length == 5 && memcmp((buf+1), "uper", 4) == 0) {
                            return (Token){state.start, length, state.line, TOKEN_SUPER};
                        }
                        break;
                    }
                    case 't':
                        if (length > 1) {
                            switch (buf[1]) {
                                case 'h': {
                                    if(length == 4 && memcmp((buf+2), "is", 2) == 0) {
                                        return (Token){state.start, length, state.line, TOKEN_THIS};
                                    }
                                    break;
                                }
                                case 'r': {
                                    if(length == 4 && memcmp((buf+2), "ue", 2) == 0) {
                                        return (Token){state.start, length, state.line, TOKEN_TRUE};
                                    }
                                    break;
                                }
                            }
                        }
                        break;
                    case 'v': {
                        if(length == 3 && memcmp((buf+1), "ar", 2) == 0) {
                            return (Token){state.start, length, state.line, TOKEN_VAR};
                        }
                        break;
                    }
                    case 'w': {
                        if(length == 5 && memcmp((buf+1), "hile", 4) == 0) {
                            return (Token){state.start, length, state.line, TOKEN_WHILE};
                        }
                        break;
                    }
                    default:
                        break;
                }
                return (Token){state.start,length,state.line,TOKEN_IDENTIFIER};
            }else if (isdigit(c)){
                char metDot=0;
                scan_number:
                    while ((c=src[state.cur],isdigit(c))){
                        state.cur++;
                    }
                if(c=='.'){
                    if(metDot){
                        scanerr="Malformed number literal\n";
                        state.cur++;
                        return (Token){(size_t)scanerr,0,state.line,TOKEN_ERROR};
                    }
                     metDot=1;
                     state.cur++;
                    goto scan_number;
                }
                size_t length=state.cur-state.start;
                return (Token){state.start,length,state.line,TOKEN_NUMBER};
            }else{
                scanerr="Unexpected character\n";
                return (Token){(size_t)scanerr,0,state.line,TOKEN_ERROR};
            }
        }
    }
    return (Token){state.start,0,state.line,TOKEN_EOF};
}
