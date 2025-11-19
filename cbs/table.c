#include "table.h"
#include <assert.h>
#include <stdio.h>



Table table={};

bool Tset(Table *Tinstance,BsObjString *key, Value value){
    size_t cap=Tinstance->cap;
    u_int32_t index=key->hash % cap;
    Tnode *head=&Tinstance->arr[index];
    Tnode *node=head;
    //TODO: check load factor, else this while loop never ends
    while(node->key!=NULL || node->key!=key){
        node=&Tinstance->arr[(++index)%cap];
    }

    Tnode *next=head->next;
    bool empty=node->key==NULL;

    if(empty && head!=node){
        //regular empty node
        head->next=node;
        node->next=next;

    }else if(node->key==key){
        //filled node
        node->value=value;
        return false;
    }
     //regular empty node or head node
    node->key=key;
    node->value=value;
    Tinstance->len++;

    return empty;
}
