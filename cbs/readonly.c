#include "readonly.h"
#include "darray.h"
#include "arena.h"

DECLARE_ARRAY(value, constants);

size_t addConstant(value c){
    size_t index=constants.len;
    append(constants, c, sizeof(value));
    return index;
};
