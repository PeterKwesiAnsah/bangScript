#include "stack.h"
#include "readonly.h"
#include <stdio.h>
#include <stdlib.h>

value stack[MAX_STACK_SIZE];

value *top=stack;


void push(value item){
    if(top+1==stack+MAX_STACK_SIZE){
        fprintf(stderr, "Maximum stack size exceeded");
        exit(EXIT_FAILURE);
    }
   *top=item;
   top++;
}

value pop(){
    top--;
    return *top;
}
