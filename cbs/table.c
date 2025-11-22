#include "table.h"
#include "darray.h"
#include <assert.h>
#include <stddef.h>
#include <stdio.h>


// returns true if Tcur is closer to LOAD_FACTOR_MAX, false otherwise
static inline bool Tgrow(Table Tcur,Table *Tnew){
    if(((double)(Tcur.len+1)/(Tcur.cap))>= LOAD_FACTOR_MAX){

        size_t cap=(size_t)(Tcur.len+1)/(LOAD_FACTOR_MIN);
        Table temp={.len=Tcur.len,.arr=NULL};

        grow(temp, Tnode, sizeof(Tnode), cap);
        temp.cap=cap;
        *Tnew=temp;
        return true;
    }
    return false;
}

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
