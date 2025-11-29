#include "readonly.h"
#include "darray.h"
#include "table.h"



DECLARE_ARRAY(Value, constants)={};
Table strings={};



inline size_t internString(Table *strings, Token token, const char *src) {
    Value val;
    // Create a temporary stack object for lookup
    BsObjString lookupKey = { .value = (char*)src + token.start, .len = token.len };

    // Check if it exists
    BsObjString *found = (BsObjString *)Tgets(strings, &lookupKey, &val);

    if (found) {
        assert(val.type == TYPE_NUMBER);
        return val.value.num;
    } else {
        // Not found: Allocate and Add
        BsObjStringFromSource *newObj = (BsObjStringFromSource *)malloc(sizeof(BsObjStringFromSource));
        newObj->obj.type = OBJ_TYPE_STRING_SOURCE;
        newObj->value = src + token.start;
        newObj->len = token.len;

        // Wrap in Value for the constant table
        Value newVal;
        newVal.type = TYPE_OBJ;
        newVal.value.obj = (BsObj *)newObj;

        size_t outIndex = addConstant(newVal);

        // Store index in string table for future lookups
        Value indexVal = { .type = TYPE_NUMBER, .value.num = outIndex };
        Tset(strings, (BsObjString *)newObj, indexVal);
        return outIndex;
    }
}

size_t addConstant(Value c){
    size_t index=constants.len;
    append(constants,Value, c);
    return index;
};
