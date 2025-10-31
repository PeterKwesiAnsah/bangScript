#include "chunk.h"
#include "vm.h"
#include "darray.h"
#include "readonly.h"
#include "stack.h"
#include <stdint.h>
#include <stdio.h>
#include <string.h>
#define READ_BYTE_CODE(currentInsPointer) (*currentInsPointer++)
#define EVALUATE_BIN_EXP(operator) do {\
    Value b=pop();\
    Value a=pop();\
    push(C_DOUBLE_TO_BS_NUMBER((BS_NUMBER_TO_C_DOUBLE(a) operator BS_NUMBER_TO_C_DOUBLE(b))));\
    }\
    while(0)



extern DECLARE_ARRAY(Value, constants);
extern DECLARE_ARRAY(u_int8_t, chunk);

uint8_t *ip=NULL;

ProgramStatus run(){
    ip=chunk.arr;
    for(;;){
        switch(READ_BYTE_CODE(ip)){
            case OP_CONSTANT_ZER0:{
                Value value= constants.arr[CONSTANT_ZERO_INDEX];
                push(value);
            }
            break;
            case OP_CONSTANT:
            {
                uint8_t index = READ_BYTE_CODE(ip);
                Value value= constants.arr[index];
                push(value);
            }
            break;
            case OP_CONSTANT_LONG:
            {
                unsigned int index =0;
                uint8_t highByte=READ_BYTE_CODE(ip);
                uint8_t midByte=READ_BYTE_CODE(ip);
                uint8_t lowByte=READ_BYTE_CODE(ip);
                index= index | lowByte;
                index= index | ((unsigned int)midByte << 8);
                index= index | ((unsigned int)highByte << 16);
                Value value= constants.arr[index];
                push(value);
            }
            break;
            case OP_ADD:{
                Value b=pop();
                Value a=pop();
                if(a.type!=b.type){
                    fputs("Add requires operands of the same the type",stderr);
                    return ERROR;
                }
                switch (a.type) {
                    case TYPE_NUMBER:
                        push(C_DOUBLE_TO_BS_NUMBER((BS_NUMBER_TO_C_DOUBLE(a) + BS_NUMBER_TO_C_DOUBLE(b))));
                    break;
                    case TYPE_OBJ:
                    {
                         //switch ((a.value.obj)->type) {
                             //case OBJ_TYPE_STRING_SOURCE:
                             //{
                                   // Operands a,b can be the same string type or different but they can be used interchangly because they have common fields
                                    BsObjStringFromSource *BsObjStringA=(BsObjStringFromSource *)a.value.obj;
                                    BsObjStringFromSource *BsObjStringB=(BsObjStringFromSource *)b.value.obj;

                                    size_t ResultStrLen=BsObjStringA->len + BsObjStringB->len;

                                    BsObjStringFromAlloc *BsObjStringResult=(BsObjStringFromAlloc *)malloc(sizeof(BsObjStringFromAlloc)+ResultStrLen+1);

                                    BsObjStringResult->value=(char *)(BsObjStringResult+(size_t)sizeof(BsObjStringFromAlloc));
                                    BsObjStringResult->len=ResultStrLen;

                                    BsObjStringResult->obj=(BsObj){.type=OBJ_TYPE_STRING_ALLOC};

                                    memcpy(BsObjStringResult->value,BsObjStringA->value,BsObjStringA->len);
                                    memcpy(BsObjStringResult->value + BsObjStringA->len,BsObjStringB->value,BsObjStringB->len);

                                    BsObjStringResult->value[ResultStrLen]='\0';
                                    push(CREATE_BS_OBJ(BsObjStringResult));
                                    // }
                             //break;
                            // default:
                             //fputs("Add requires operands to be either a number or a string type",stderr);
                             //return ERROR;
                             //}
                    }
                    break;
                    default:
                    fputs("Add requires operands to be either a number or a string type",stderr);
                    return ERROR;
                }

            }
            break;
            case OP_SUB:
            EVALUATE_BIN_EXP(-);
            break;
            case OP_MUL:
            EVALUATE_BIN_EXP(*);
            break;
            case OP_DIV:
            EVALUATE_BIN_EXP(/);
            break;
            case OP_EQUAL:
                {
                    Value b=pop();
                    Value a=pop();
                    if(a.type!=b.type){
                        fputs("Equal comparison requires operands of the same the type",stderr);
                        return ERROR;
                    }
                    switch (a.type) {
                        case TYPE_NUMBER:
                        push(C_BOOL_TO_BS_BOOLEAN(BS_NUMBER_TO_C_DOUBLE(a) == BS_NUMBER_TO_C_DOUBLE(b)));
                        break;
                        case TYPE_BOOL:
                        push(C_BOOL_TO_BS_BOOLEAN(BS_BOOLEAN_TO_C_BOOL(a) == BS_BOOLEAN_TO_C_BOOL(b)));
                        break;
                        case TYPE_OBJ:{

                            BsObjString * BsObjStringA=(BsObjString *)a.value.obj;
                            BsObjString * BsObjStringB=(BsObjString *)b.value.obj;

                            if(BsObjStringA->len!=BsObjStringB->len){
                                push(C_BOOL_TO_BS_BOOLEAN(false));
                            }else{
                                push(C_BOOL_TO_BS_BOOLEAN(!memcmp(BsObjStringA->value, BsObjStringB->value, BsObjStringA->len)));
                            }
                        }
                        break;
                        default:
                        fputs("Add requires operands to be either a number, bool or a string type",stderr);
                        return ERROR;
                    }

                }
                break;
            case OP_EQUAL_NOT:
                    {
                        Value b=pop();
                        Value a=pop();
                        if(a.type!=b.type){
                            fputs("Add requires operands of the same the type",stderr);
                            return ERROR;
                        }
                        switch (b.type) {
                            case TYPE_NUMBER:
                                push(C_BOOL_TO_BS_BOOLEAN(!(BS_NUMBER_TO_C_DOUBLE(a) == BS_NUMBER_TO_C_DOUBLE(b))));
                            break;
                            default:
                                fputs("Comparison operation requires operands to be a number",stderr);
                                return ERROR;
                        }
                    }
            break;
            case OP_GREATOR:
                {
                    Value b=pop();
                    Value a=pop();
                    if(a.type!=b.type){
                        fputs("Add requires operands of the same the type",stderr);
                        return ERROR;
                    }
                    switch (b.type) {
                        case TYPE_NUMBER:
                            push(C_BOOL_TO_BS_BOOLEAN((BS_NUMBER_TO_C_DOUBLE(a) > BS_NUMBER_TO_C_DOUBLE(b))));
                        break;
                        default:
                            fputs("Comparison operation requires operands to be a number",stderr);
                            return ERROR;
                    }
                }
            break;
            case OP_GREATOR_NOT: {
                Value b=pop();
                Value a=pop();
                if(a.type!=b.type){
                    fputs("Add requires operands of the same the type",stderr);
                    return ERROR;
                }
                switch (a.type) {
                    case TYPE_NUMBER:
                       push(C_BOOL_TO_BS_BOOLEAN(!(BS_NUMBER_TO_C_DOUBLE(a) > BS_NUMBER_TO_C_DOUBLE(b))));
                    break;
                    default:
                    fputs("Comparison operation requires operands to be a number",stderr);
                    return ERROR;
                }
            }
            break;
            case OP_LESS:
                {
                    Value b=pop();
                    Value a=pop();
                    if(a.type!=b.type){
                        fputs("Add requires operands of the same the type",stderr);
                        return ERROR;
                    }
                    switch (b.type) {
                        case TYPE_NUMBER:
                            push(C_BOOL_TO_BS_BOOLEAN(BS_NUMBER_TO_C_DOUBLE(a) < BS_NUMBER_TO_C_DOUBLE(b)));
                        break;
                        default:
                            fputs("Comparison operation requires operands to be a number",stderr);
                            return ERROR;
                    }
                }
            break;
            case OP_LESS_NOT:
            break;
            case OP_PRINT:
            {
            Value result= pop();
            switch (result.type) {
                case TYPE_OBJ:{
                    switch (result.value.obj->type) {
                        case OBJ_TYPE_STRING_SOURCE:
                        printf("%.*s\n",((BsObjStringFromSource *)result.value.obj)->len,((BsObjStringFromSource *)result.value.obj)->value);
                        break;
                        case OBJ_TYPE_STRING_ALLOC:
                        printf("%s\n",((BsObjStringFromAlloc *)result.value.obj)->value);
                        break;
                        default:
                        break;
                    }
                }
                break;
                case TYPE_BOOL:
                    printf("%d\n", BS_BOOLEAN_TO_C_BOOL(result));
                break;
                case TYPE_NUMBER:
                 printf("%f\n",BS_NUMBER_TO_C_DOUBLE(result));
                break;
            }
            }
            break;
            case OP_RETURN:
             return SUCCESS;
            default:
             fputs("Invalid bytecode instruction",stderr);
             return ERROR;
        }
    }
    return SUCCESS;
};
