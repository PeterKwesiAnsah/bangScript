#ifndef VM_H
#define VM_H
#include "chunk.h"
#include "compiler.h"
#include "readonly.h"

typedef enum {
   SUCCESS,
   ERROR
} ProgramStatus;

//TODO: create a frame struct that contains chunk struct,ip and

typedef struct {
    Chunk chunk;
    const char *src;
    u_int8_t *ip;
    Value *constants;
    Compiler *compiler;
} Frame;

extern Frame frame;

ProgramStatus run();

#endif
