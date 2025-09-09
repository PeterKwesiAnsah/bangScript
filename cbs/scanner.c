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

const char *scanerr="";

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
            if (((c=fgetc(src)),state.cur++,state.fpos++,c=='/')){
                while ((c=fgetc(src),state.cur++,state.fpos++,c!='\n')){}
                //restore File Position to '\n'
                fseek(src, state.cur--, SEEK_SET);
                state.fpos--;
               return scanTokens(src);
            }
            //restore File Position to '\'+1
            fseek(src, state.cur--, SEEK_SET);
            state.fpos--;
            return (Token){state.start,1,state.line,TOKEN_SLASH};
        }
        // One or two character tokens
        case '!':{
              if (((c=fgetc(src)),state.cur++,state.fpos++,c=='=')){
                  return (Token){state.cur,2,state.line,TOKEN_BANG_EQUAL};
              }
              //restore File Position to '!'+1
              fseek(src, state.cur--, SEEK_SET);
              state.fpos--;
            return (Token){state.cur,1,state.line,TOKEN_BANG};
        }
        case '<':{
            if (((c=fgetc(src)),state.cur++,state.fpos++,c=='=')){
                return (Token){state.cur,2,state.line,TOKEN_LESS_EQUAL};
            }
            //restore File Position to '<'+1
            fseek(src, state.cur--, SEEK_SET);
            state.fpos--;
            return (Token){state.cur,1,state.line,TOKEN_LESS};
        }
        case '>':{
            if (((c=fgetc(src)),state.cur++,state.fpos++,c=='=')){
                return (Token){state.cur,2,state.line,TOKEN_GREATER_EQUAL};
            }
            //restore File Position to '>' + 1
            fseek(src, state.cur--, SEEK_SET);
            state.fpos--;
            return (Token){state.cur,1,state.line,TOKEN_GREATER};
        }
        case '=':{
            if (((c=fgetc(src)),state.cur++,state.fpos++,c=='=')){
                return (Token){state.cur,2,state.line,TOKEN_EQUAL_EQUAL};
            }
            //restore File Position to '=' + 1
            fseek(src, state.cur--, SEEK_SET);
            state.fpos--;
            return (Token){state.cur,1,state.line,TOKEN_EQUAL};
        }
        case '"':{
          while ((c=fgetc(src),state.cur++,state.fpos++,!isEOF(c) && c!='"')){
              if(c=='\n'){
                  state.line++;
              }
          }
           if(isEOF(c)){
                scanerr="Unterminated String\n";
                return (Token){state.start+1,0,state.line,TOKEN_ERROR};
           }
           return (Token){state.start+1,(state.cur-(state.start+1)),state.line,TOKEN_STRING};
        };
        case '\0':{
            return (Token){state.start,0,state.line,TOKEN_EOF};
        }
        default:{
            if (isalpha(c) || c=='_'){
            }else if (isdigit(c)){
                char metDot=0;
                scan_number:
                    while ((c=fgetc(src),state.cur++,state.fpos++,isdigit(c))){}
                if(c=='.'){
                    if(metDot){
                        scanerr="Malformed number literal\n";
                        return (Token){state.start,0,state.line,TOKEN_ERROR};
                    }
                     metDot=1;
                    goto scan_number;
                }
                //restore File Position to the last valid character + 1
                fseek(src, state.cur--, SEEK_SET);
                state.fpos--;
                size_t length=state.fpos-state.start;
                return (Token){state.start,length,state.line,TOKEN_NUMBER};
            }else{
                scanerr="Unexpected character\n";
                return (Token){state.start,0,state.line,TOKEN_ERROR};
            }
        }

    }
    return (Token){state.start,0,state.line,TOKEN_EOF};
}
