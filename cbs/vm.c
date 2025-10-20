#include "chunk.h"
#include "compiler.h"
#include "vm.h"
#include "darray.h"
#include "readonly.h"
#include "stack.h"
#include <stdint.h>
#include <stdio.h>
#define READ_BYTE_CODE(currentInsPointer) (*currentInsPointer++)
#define EVALUATE_BIN_EXP(operator) do {\
    Value b=pop();\
    Value a=pop();\
    push((a operator b));\
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
            case OP_ADD:
            EVALUATE_BIN_EXP(+);
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
            case OP_PRINT:{
            Value result= pop();
            printf("%f\n",result);
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
