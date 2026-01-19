// temporal location for operands
#ifndef STACK_H
#define STACK_H
#include "darray.h"
#include "readonly.h"
#define MAX_STACK_SIZE 256
void push(Value);
Value pop();
Value getStackItem(size_t);
void updateStackItem(size_t,Value);
#endif
