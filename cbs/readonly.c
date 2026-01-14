#include "readonly.h"
#include "darray.h"
#include "table.h"
#include "vm.h"



Table strings={};
extern Frame frame;


static uint32_t hashString(const char* key, int length) {
    uint32_t hash = 2166136261u;
    for (int i = 0; i < length; i++) {
        hash ^= (uint8_t)key[i];
        hash *= 16777619;
    }
    return hash;
}

 size_t internString(Table *strings, Token token, const char *src) {
    Value val;
    // Create a temporary stack object for lookup
    BsObjString lookupKey = { .value = (char *) src + token.start, .len = token.len, .hash=hashString(src + token.start, token.len) };

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
        newObj->hash=hashString((char *)src + token.start,token.len);

        // Wrap in Value for the constant table
        Value newVal;
        newVal.type = TYPE_OBJ;
        newVal.value.obj = (BsObj *)newObj;

        size_t outIndex = addConstant(newVal,frame.constants);

        // Store index in string table for future lookups
        Value indexVal = { .type = TYPE_NUMBER, .value.num = outIndex };

        TABLE_EXPAND(strings);

        Tsets(strings, (BsObjString *)newObj, indexVal);

        return outIndex;
    }
}

size_t addConstant(Value BsValue, Constants *constants){
    size_t index=constants->len;
    appendPtr(constants,Value, BsValue);
    return index;
};
