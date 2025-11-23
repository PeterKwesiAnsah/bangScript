#include "readonly.h"
#include "darray.h"
#include "table.h"


DECLARE_ARRAY(Value, constants)={};
Table strings={};

size_t addConstant(Value c){
    size_t index=constants.len;
    append(constants,Value, c);
    return index;
};
