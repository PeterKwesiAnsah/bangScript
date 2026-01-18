#include "parser.h"
#include "chunk.h"
#include "readonly.h"
#include "scanner.h"
#include "table.h"
#include "setjmp.h"
#include "vm.h"
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

Parser parser = {{}, {}, false};

extern const char *scanerr;
// table of a set of unique(textually) BsObjString
extern Table strings;

extern Frame frame;

extern jmp_buf buf;



extern inline size_t internString(Table *, Token, const char *);

ParseRule rules[] = {
    [TOKEN_LEFT_PAREN] = {grouping, NULL, PREC_NONE},
    [TOKEN_RIGHT_PAREN] = {NULL, NULL, PREC_NONE},
    [TOKEN_LEFT_BRACE] = {NULL, NULL, PREC_NONE},
    [TOKEN_RIGHT_BRACE] = {NULL, NULL, PREC_NONE},
    [TOKEN_COMMA] = {NULL, NULL, PREC_NONE},
    [TOKEN_DOT] = {NULL, NULL, PREC_NONE},
    [TOKEN_MINUS] = {unary, binary, PREC_TERM},
    [TOKEN_PLUS] = {NULL, binary, PREC_TERM},
    [TOKEN_SEMICOLON] = {NULL, NULL, PREC_NONE},
    [TOKEN_SLASH] = {NULL, binary, PREC_FACTOR},
    [TOKEN_STAR] = {NULL, binary, PREC_FACTOR},
    [TOKEN_BANG] = {unary, NULL, PREC_NONE},
    [TOKEN_BANG_EQUAL] = {NULL, binary, PREC_EQUALITY},
    [TOKEN_EQUAL] = {NULL, NULL, PREC_NONE},
    [TOKEN_EQUAL_EQUAL] = {NULL, binary, PREC_EQUALITY},
    [TOKEN_GREATER] = {NULL, binary, PREC_COMPARISON},
    [TOKEN_GREATER_EQUAL] = {NULL, binary, PREC_COMPARISON},
    [TOKEN_LESS] = {NULL, binary, PREC_COMPARISON},
    [TOKEN_LESS_EQUAL] = {NULL, binary, PREC_COMPARISON},
    [TOKEN_IDENTIFIER] = {identifier, NULL, PREC_PRIMARY},
    [TOKEN_STRING] = {string, NULL, PREC_NONE},
    [TOKEN_NUMBER] = {number, NULL, PREC_NONE},
    [TOKEN_AND] = {NULL, NULL, PREC_NONE},
    [TOKEN_CLASS] = {NULL, NULL, PREC_NONE},
    [TOKEN_ELSE] = {NULL, NULL, PREC_NONE},
    [TOKEN_TRUE] = {boolean, NULL, PREC_NONE},
    [TOKEN_FALSE] = {boolean, NULL, PREC_NONE},
    [TOKEN_FOR] = {NULL, NULL, PREC_NONE},
    [TOKEN_FUN] = {NULL, NULL, PREC_NONE},
    [TOKEN_IF] = {NULL, NULL, PREC_NONE},
    [TOKEN_NIL] = {NULL, NULL, PREC_NONE},
    [TOKEN_OR] = {NULL, NULL, PREC_NONE},
    [TOKEN_PRINT] = {NULL, NULL, PREC_NONE},
    [TOKEN_RETURN] = {NULL, NULL, PREC_NONE},
    [TOKEN_SUPER] = {NULL, NULL, PREC_NONE},
    [TOKEN_THIS] = {NULL, NULL, PREC_NONE},
};

void advance() {
  parser.previous = parser.current;
  for (;;) {
    parser.current = scanTokens(frame.src);
    if (parser.current.tt != TOKEN_ERROR)
      break;
    errorAt(&parser.current, "", (char *)0);
  }
}
static void parsePrecedence(Precedence precedence) {
  advance();
  ParseFn prefixRule = rules[parser.previous.tt].prefix;
  if (prefixRule == NULL) {
    errorAt(&parser.previous, "Invalid Expression.", frame.src);

    return;
  }
  prefixRule(precedence == PREC_ASSIGNMENT);

  while (precedence <= rules[parser.current.tt].precedence) {
    advance();
    ParseFn infixRule = rules[parser.previous.tt].infix;
    infixRule(precedence == PREC_ASSIGNMENT);
  }
}

void expression() {
  parsePrecedence(PREC_ASSIGNMENT);
  if (parser.current.tt != TOKEN_SEMICOLON) {
    errorAt(&parser.previous, "Expected a semi-colon after token.", frame.src);
    return;
  }
  advance();
}
void grouping(bool isAssignExp) {
  parsePrecedence(PREC_ASSIGNMENT);
  if (parser.current.tt != TOKEN_RIGHT_PAREN) {
    errorAt(&parser.previous, "Expected a right parenthesi after token.",
            frame.src);
    return;
  }
  advance();
}

static void binary(bool isAssignExp) {
  TokenType operatorType = parser.previous.tt;
  unsigned int line = parser.previous.line;
  parsePrecedence((Precedence)(rules[parser.previous.tt].precedence + 1));
  switch (operatorType) {
  case TOKEN_PLUS:
    WRITE_BYTECODE(frame.chunk, OP_ADD, line);
    break;
  case TOKEN_MINUS:
    WRITE_BYTECODE(frame.chunk, OP_SUB, line);
    break;
  case TOKEN_STAR:
    WRITE_BYTECODE(frame.chunk, OP_MUL, line);
    break;
  case TOKEN_SLASH:
    WRITE_BYTECODE(frame.chunk, OP_DIV, line);
    break;
  case TOKEN_EQUAL_EQUAL:
    WRITE_BYTECODE(frame.chunk, OP_EQUAL, line);
    break;
  case TOKEN_BANG_EQUAL:
    WRITE_BYTECODE(frame.chunk, OP_EQUAL_NOT, line);
    break;
  case TOKEN_GREATER:
    WRITE_BYTECODE(frame.chunk, OP_GREATOR, line);
    break;
  case TOKEN_GREATER_EQUAL:
    WRITE_BYTECODE(frame.chunk, OP_LESS_NOT, line);
    break;
  case TOKEN_LESS:
    WRITE_BYTECODE(frame.chunk, OP_LESS, line);
    break;
  case TOKEN_LESS_EQUAL:
    WRITE_BYTECODE(frame.chunk, OP_GREATOR_NOT, line);
    break;
  default:
    return; // Unreachable.
  }
}

void unary(bool isAssignExp) {
  Token token = parser.previous;

  switch (token.tt) {
  case TOKEN_MINUS: {
    WRITE_BYTECODE(frame.chunk, OP_CONSTANT_ZER0, token.line);
    parsePrecedence(PREC_UNARY);
    WRITE_BYTECODE(frame.chunk, OP_SUB, token.line);
  } break;
  default:
    return;
  }
}

static void number(bool isAssignExp) {
  Token numToken = parser.previous;
  size_t constantIndex = addConstant(
      C_DOUBLE_TO_BS_NUMBER(atof(frame.src + numToken.start)), frame.constants);
  if (constantIndex >= CONSTANT_LIMIT) {
    // Write opcode
    WRITE_BYTECODE(frame.chunk, OP_CONSTANT_LONG, numToken.line);
    // Write operand as 3 bytes
    WRITE_BYTECODE(frame.chunk, constantIndex & 0xFF, numToken.line);
    WRITE_BYTECODE(frame.chunk, (constantIndex >> 8) & 0xFF, numToken.line);
    WRITE_BYTECODE(frame.chunk, (constantIndex >> 16) & 0xFF, numToken.line);
    return;
  }
  // write opCode
  WRITE_BYTECODE(frame.chunk, OP_CONSTANT, numToken.line);
  // write operand Index
  WRITE_BYTECODE(frame.chunk, constantIndex, numToken.line);
}
static void string(bool isAssignExp) {

  Token strToken = parser.previous;
  // Value value;

  size_t stringLiteralIndex = internString(&strings, strToken, frame.src);

  if (stringLiteralIndex >= CONSTANT_LIMIT) {
    // Write opcode
    WRITE_BYTECODE(frame.chunk, OP_CONSTANT_LONG, strToken.line);
    // Write operand as 3 bytes
    WRITE_BYTECODE(frame.chunk, stringLiteralIndex & 0xFF, strToken.line);
    WRITE_BYTECODE(frame.chunk, (stringLiteralIndex >> 8) & 0xFF,
                   strToken.line);
    WRITE_BYTECODE(frame.chunk, (stringLiteralIndex >> 16) & 0xFF,
                   strToken.line);
    return;
  }

  // write opCode
  WRITE_BYTECODE(frame.chunk, OP_CONSTANT, strToken.line);
  // write operand Index
  WRITE_BYTECODE(frame.chunk, stringLiteralIndex, strToken.line);
}

static void boolean(bool isAssignExp) {
  Token boolToken = parser.previous;
  // write opCode
  WRITE_BYTECODE(frame.chunk, OP_CONSTANT, boolToken.line);
  // write operand Index
  WRITE_BYTECODE(frame.chunk,
                 (boolToken.tt == TOKEN_TRUE ? CONSTANT_TRUE_BOOL_INDEX
                                             : CONSTANT_FALSE_BOOL_INDEX),
                 boolToken.line);
}

static void identifier(bool isAssignExp) {
  Token identifierToken = parser.previous;
  bool assignment = isAssignExp && parser.current.tt == TOKEN_EQUAL;

  uint8_t OP_CODE_GET = 0;
  uint8_t OP_CODE_SET = 0;

  int OP_CODE_OPERAND_INDEX = frame.compiler->len;

  if (frame.compiler->scopeDepth == 0)
    goto ParseCompileGlobals;

  for (; OP_CODE_OPERAND_INDEX >= 0; OP_CODE_OPERAND_INDEX--) {
    if (identifierToken.len ==
            frame.compiler->locals[OP_CODE_OPERAND_INDEX].name.len &&
        frame.compiler->locals[OP_CODE_OPERAND_INDEX].depth != -1 &&
        !memcmp(frame.src + identifierToken.start,
                frame.src +
                    frame.compiler->locals[OP_CODE_OPERAND_INDEX].name.start,
                identifierToken.len)) {
      break;
    }
  }

  if (OP_CODE_OPERAND_INDEX == -1) {
  ParseCompileGlobals:
    OP_CODE_OPERAND_INDEX = internString(&strings, identifierToken, frame.src);
    if (assignment) {
      advance();
      expression();
      // TODO: Handle long indexes
      WRITE_BYTECODE(frame.chunk, OP_GLOBALVAR_ASSIGN, identifierToken.line);
      WRITE_BYTECODE(frame.chunk, OP_CODE_OPERAND_INDEX, identifierToken.line);
    } else {
      WRITE_BYTECODE(frame.chunk, OP_GLOBALVAR_GET, identifierToken.line);
      WRITE_BYTECODE(frame.chunk, OP_CODE_OPERAND_INDEX, identifierToken.line);
      // cache hash index
      WRITE_BYTECODE(frame.chunk, 0, identifierToken.line);
      WRITE_BYTECODE(frame.chunk, 0, identifierToken.line);
    }
    return;
  } else if (assignment) {
    const unsigned localDepth =
        frame.compiler->locals[OP_CODE_OPERAND_INDEX].depth;
    frame.compiler->locals[OP_CODE_OPERAND_INDEX].depth = -1;
    advance();
    expression();
    frame.compiler->locals[OP_CODE_OPERAND_INDEX].depth = localDepth;
    // TODO: Handle long indexes
    WRITE_BYTECODE(frame.chunk, OP_LOCALVAR_ASSIGN, identifierToken.line);
    WRITE_BYTECODE(frame.chunk, OP_CODE_OPERAND_INDEX, identifierToken.line);
  } else {
    WRITE_BYTECODE(frame.chunk, OP_LOCALVAR_GET, identifierToken.line);
    WRITE_BYTECODE(frame.chunk, OP_CODE_OPERAND_INDEX, identifierToken.line);
  }
}

static void nil(bool isAssignExp) {
  Token nilToken = parser.previous;
  WRITE_BYTECODE(frame.chunk, OP_CONSTANT, nilToken.line);
  // write operand Index
  WRITE_BYTECODE(frame.chunk, CONSTANT_NIL_INDEX, nilToken.line);
}

void errorAt(Token *token, const char *message, const char *src) {
  fprintf(stderr, "[line %d] Error", token->line);
  if (token->tt == TOKEN_EOF) {
    fprintf(stderr, " at end");
  } else if (token->tt == TOKEN_ERROR) {
    // Nothing.
  } else {
    fprintf(stderr, " at '%.*s'", token->len, (char *)src + token->start);
  }
  fprintf(stderr, ": %s\n", message);
  parser.hadError = true;
  longjmp(buf,1);
}
