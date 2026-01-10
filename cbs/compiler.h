#ifndef COMPILER_H
#define COMPILER_H
#include "chunk.h"
#include "scanner.h"
#include <stdint.h>

#define IN_SCOPE_LOCALS_LIMIT (UINT8_MAX + 1)

typedef enum {
    COMPILER_ERROR,
    COMPILER_OK,
} CompilerStatus;

typedef struct {
  Token name;
  unsigned depth;
} Local;

typedef struct {
  Local locals[IN_SCOPE_LOCALS_LIMIT];
  unsigned len;
  unsigned scopeDepth;
} Compiler;

CompilerStatus compile(const char *);
CompilerStatus declaration(const char *);
#endif
