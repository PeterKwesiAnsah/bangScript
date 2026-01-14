#include "compiler.h"
#include "parser.h"
#include "readonly.h"
#include "table.h"
#include "vm.h"
#include <stdbool.h>


extern Parser parser;
extern Table strings;


extern inline size_t internString(Table *, Token, const char *);


extern Frame frame;

// parse,complile statements
CompilerStatus declaration(const char *src) {
  switch (parser.current.tt) {
  // variable declaration
  case TOKEN_VAR: {
    advance();
    Token cur = parser.current;
    if (frame.compiler->scopeDepth > 0) {
      // define local variable as it is in scope
      frame.compiler->locals[frame.compiler->len].name = cur;
      frame.compiler->locals[frame.compiler->len].depth = -1;
    }
    // if(cur.tt!=TOKEN_IDENTIFIER){
    // compile error
    //}
    advance();
    // var foo=0;
    //  Handle multiple global var definations
    if (parser.current.tt == TOKEN_EQUAL) {
      advance();
      expression();
    } else {
      // var foo;
      if (parser.current.tt != TOKEN_SEMICOLON) {
        fprintf(stderr, "Expected a semi-colon but got %d\n",
                parser.current.tt);
        parser.hadError = true;
        return true;
      }
      WRITE_BYTECODE(frame.chunk, OP_CONSTANT, cur.line);
      WRITE_BYTECODE(frame.chunk, CONSTANT_NIL_INDEX, cur.line);
      advance();
    }

    if (frame.compiler->scopeDepth > 0) {
      frame.compiler->locals[frame.compiler->len++].depth = frame.compiler->scopeDepth;
      break;
    }
    // if local var, we skip this???
    size_t BsobjStringConstIndex = internString(&strings, cur, src);
    WRITE_BYTECODE(frame.chunk, OP_GLOBALVAR_DEF, cur.line);
    WRITE_BYTECODE(frame.chunk, BsobjStringConstIndex, cur.line);
  } break;
  case TOKEN_PRINT: {
    Token printToken = parser.current;
    advance();
    expression();
    WRITE_BYTECODE(frame.chunk, OP_PRINT, printToken.line);
  } break;
  case TOKEN_LEFT_BRACE: {
    const unsigned scope = frame.compiler->scopeDepth++;

    while (!declaration(src) && parser.current.tt != TOKEN_RIGHT_BRACE &&
           parser.current.tt != TOKEN_EOF) {
    }
    if (parser.current.tt != TOKEN_EOF || parser.hadError) {
      //
      //
    }
    // sync stack and locals
    for (unsigned i = frame.compiler->len; i >= 0 && frame.compiler->locals[i].depth != scope;
         i--) {
      WRITE_BYTECODE(frame.chunk, OP_POP, parser.current.line);
      frame.compiler->len--;
    }

  } break;
  case TOKEN_EOF: {
    WRITE_BYTECODE(frame.chunk, OP_RETURN, parser.current.line);
    return !parser.hadError;
  }
  default:
    expression();
    break;
  }

  return !parser.hadError;
}
CompilerStatus compile(const char *src) {

  Tinit(&strings);
  addConstant(C_DOUBLE_TO_BS_NUMBER(0),frame.constants);
  addConstant(C_BOOL_TO_BS_BOOLEAN(true),frame.constants);
  addConstant(C_BOOL_TO_BS_BOOLEAN(false),frame.constants);
  addConstant((Value){.type = TYPE_NIL, .value = {}},frame.constants);

  // set the ball rolling
  advance();
  while (!declaration(src) && parser.current.tt != TOKEN_EOF);
  return !parser.hadError;
}
