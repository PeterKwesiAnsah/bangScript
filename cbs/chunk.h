#ifndef CHUNK_H
#define CHUNK_H
enum {
    OP_CONSTANT_ZER0,
    OP_CONSTANT,
    OP_CONSTANT_LONG,
    OP_ADD  ,
    OP_SUB  ,
    OP_MUL  ,
    OP_DIV  ,
    OP_NEGATE,
    //for debugging purposes
    OP_PRINT,
    OP_RETURN
} opcodes;
#endif
