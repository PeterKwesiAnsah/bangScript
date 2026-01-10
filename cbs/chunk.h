#ifndef CHUNK_H
#define CHUNK_H
#include "line.h"
#include "darray.h"

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
    OP_EQUAL_NOT,
    OP_GREATOR,
    OP_LESS_NOT,
    OP_LESS,
    OP_GREATOR_NOT,
    OP_GLOBALVAR_DEF,
    OP_GLOBALVAR_GET,
    OP_GLOBALVAR_ASSIGN,
    OP_LOCALVAR_GET,
    OP_LOCAVAR_ASSIGN,
    OP_POP,
    //for debugging purposes
    OP_PRINT,
    OP_RETURN
} BS_OP_CODES;

#define DECLARE_CHUNK_TYPE(name) \
typedef struct {   \
    size_t cap;\
    size_t len;\
    u_int8_t *arr;\
    uint8_t *ip;\
} name;

DECLARE_CHUNK_TYPE(Chunk)

#define DECLARE_CHUNK(name) \
    DECLARE_ARRAY(u_int8_t, name); \
    uint8_t *ip;

#define WRITE_BYTECODE(chunk,byte,line) do{ \
append(chunk,u_int8_t,byte);\
}while(0)



#endif
