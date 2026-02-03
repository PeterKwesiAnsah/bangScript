#ifndef VM_H
#define VM_H
#include "chunk.h"
#include "compiler.h"
#include "readonly.h"

typedef enum { SUCCESS, ERROR } ProgramStatus;

// TODO: create a frame struct that contains chunk struct,ip and

typedef struct {
  Chunk chunk;
  const char *src;
  uint8_t *ip;
  Constants *constants;
  Compiler *compiler;
} Frame;

extern Frame frame;

ProgramStatus run();

#endif
