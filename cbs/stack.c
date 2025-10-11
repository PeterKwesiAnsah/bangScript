#include "stack.h"
#include "readonly.h"
#include <stdio.h>
#include <stdlib.h>

Value stack[MAX_STACK_SIZE];

Value *top=stack;


void push(Value item){
    if(top+1==stack+MAX_STACK_SIZE){
        fprintf(stderr, "Maximum stack size exceeded");
        exit(EXIT_FAILURE);
    }
   *top=item;
   top++;
}

Value pop(){
    top--;
    return *top;
}
