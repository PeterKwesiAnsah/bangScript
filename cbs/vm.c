#include "chunk.h"
#include "compiler.h"
#include "vm.h"
#include "darray.h"
#include "readonly.h"
#include "stack.h"
#include <stdint.h>
#define READ_BYTE_CODE(currentInsPointer) (*currentInsPointer++)


extern DECLARE_ARRAY(Value, constants);
extern DECLARE_ARRAY(u_int8_t, chunk);
uint8_t *ip=NULL;
ProgramStatus run(){
    ip=chunk.arr;
    for(;;){
        switch(READ_BYTE_CODE(ip)){
            case OP_CONSTANT:{
                uint8_t index = READ_BYTE_CODE(ip);
                Value value= constants.arr[index];
                push(value);
            }
            case OP_CONSTANT_LONG:{
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
            case OP_ADD:
            case OP_SUB:
            case OP_MUL:
            case OP_DIV:
            case OP_RETURN:
             return SUCCESS;
            default:
             fputs("Invalid bytecode instruction",stderr);
             return ERROR;
        }
    }


    return SUCCESS;
};
