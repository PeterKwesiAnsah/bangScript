#ifndef CHUNK_H
#define CHUNK_H
enum {
    OP_CONSTANT,
    OP_CONSTANT_LONG,
    OP_ADD  ,
    OP_SUB  ,
    OP_MUL  ,
    OP_DIV  ,
    OP_NEGATE,
    OP_RETURN
} opcodes;
#endif
