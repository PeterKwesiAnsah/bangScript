#include "readonly.h"
#include "darray.h"


DECLARE_ARRAY(Value, constants)={};



size_t addConstant(Value c){
    size_t index=constants.len;
    append(constants,Value, c);
    return index;
};
