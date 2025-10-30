#ifndef CHUNK_H
#define CHUNK_H
typedef enum {
    OP_CONSTANT_ZER0,
    OP_CONSTANT,
    OP_CONSTANT_LONG,
    OP_ADD  ,
    OP_SUB  ,
    OP_MUL  ,
    OP_DIV  ,
    OP_NEGATE,
    OP_EQUAL,
    OP_GREATOR,
    OP_LESS,
    //for debugging purposes
    OP_PRINT,
    OP_RETURN
} BS_OP_CODES;
#endif
