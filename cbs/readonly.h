//Number Constants and String literals
#ifndef READONLY_H
#define READONLY_H
#include "darray.h"
#include <stdbool.h>
#include <stddef.h>
#define CONSTANT_LIMIT 256
#define CONSTANT_LONG_LIMIT 16777216

typedef enum {
    CONSTANT_ZERO_INDEX,
    CONSTANT_TRUE_BOOL_INDEX,
    CONSTANT_FALSE_BOOL_INDEX,
    CONSTANT_NIL_INDEX
    //0,true,false,nil
}  CONSTANT_LITERAL_INDEXES;




typedef enum {
    TYPE_NIL,
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
        bool boolean;
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

typedef struct {
    BsObj obj;
    unsigned int len;
    char *value;
} BsObjString;



#define C_DOUBLE_TO_BS_NUMBER(cdouble) ((Value){.value={.num=cdouble},.type=TYPE_NUMBER})
#define BS_NUMBER_TO_C_DOUBLE(bsnumber) (bsnumber.value.num)

#define CREATE_BS_OBJ(objPointer) ((Value){.value={.obj=(BsObj *)objPointer},.type=TYPE_OBJ})

#define C_BOOL_TO_BS_BOOLEAN(cbool) ((Value){.value={.boolean=cbool},.type=TYPE_BOOL})
#define BS_BOOLEAN_TO_C_BOOL(bsbool) (bsbool.value.boolean)
size_t addConstant(Value);
#endif
