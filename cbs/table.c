#include "table.h"
#include <assert.h>
#include <stdio.h>



//Table table={};



//true if it inserts into an empty bucket, false if it updated one
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
    //regular empty node
    node->key=key;
    node->value=value;
    Tinstance->len++;

    return isEmpty;
}

//true if entry was  found, false otherwise
bool Tget(Table *Tinstance,BsObjString *key, Value *value){
    size_t cap=Tinstance->cap;
    u_int32_t index=key->hash % cap;
    Tnode *node=&Tinstance->arr[index];

    while(node->key!=NULL || node->key!=key){
        index=(index+1) % cap;
        node=&Tinstance->arr[index];
    }

    if(node->key==NULL){
        return false;
    }

    *value=node->value;
    return true;
}
