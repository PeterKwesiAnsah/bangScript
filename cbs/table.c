#include "table.h"
#include <assert.h>
#include <stdio.h>



Table table={};

bool Tset(Table *Tinstance,BsObjString *key, Value value){
    size_t cap=Tinstance->cap;
    u_int32_t index=key->hash % cap;
    Tnode *node=&Tinstance->arr[index];

    //TODO: check load factor, else this while loop never ends
    while(node->key!=NULL || node->key!=key){
        index=(index+1) % cap;
        node=&Tinstance->arr[index];
    }

    bool isEmpty=node->key==NULL;
    if(node->key==key){
        //filled node
        node->value=value;
        return false;
    }
     //regular empty node or head node
    node->key=key;
    node->value=value;
    Tinstance->len++;

    return isEmpty;
}
