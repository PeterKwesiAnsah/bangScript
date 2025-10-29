//Number Constants and String literals
#ifndef READONLY_H
#define READONLY_H
#include "darray.h"
#include <stddef.h>
#define CONSTANT_LIMIT 256
#define CONSTANT_LONG_LIMIT 16777216
#define CONSTANT_ZERO_INDEX 0



typedef enum {
    TYPE_BOOL,
    TYPE_NUMBER,
    TYPE_OBJ
} BsType;

typedef enum {
    OBJ_TYPE_STRING_SOURCE,
    OBJ_TYPE_STRING_ALLOC
} BsObjType;

typedef struct {
    BsObjType type;
} BsObj;

struct BsValue {
    union {
        double num;
        BsObj *obj;
    } value;
    BsType type;
};

typedef struct BsValue Value;

typedef struct {
    BsObj obj;
    unsigned int len;
    const char *value;
} BsObjStringFromSource;

typedef struct {
    BsObj obj;
    unsigned int len;
    char *value;
    char payload[];
} BsObjStringFromAlloc;

#define C_DOUBLE_TO_BS_NUMBER(double) ((Value){.value={.num=double},.type=TYPE_NUMBER})
#define BS_NUMBER_TO_C_DOUBLE(number) (number.value.num)
#define CREATE_BS_OBJ(objPointer) ((Value){.value={.obj=(BsObj *)objPointer},.type=TYPE_OBJ})
size_t addConstant(Value);
#endif
