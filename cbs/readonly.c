#include "readonly.h"
#include "darray.h"
#include "arena.h"

DECLARE_ARRAY(Value, constants);



size_t addConstant(Value c){
    size_t index=constants.len;
    append(constants, c, sizeof(Value));
    return index;
};
