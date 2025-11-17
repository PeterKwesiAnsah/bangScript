#ifndef TABLE_H
#define TABLE_H

#include "readonly.h"
#include <stdbool.h>

struct KVnode {
BsObjString *key;
Value value;
struct KVnode *next;
};

typedef struct KVnode Tnode;

DECLARE_ARRAY_TYPE(Tnode,Table);

bool Tset(Table,BsObjString *, Value);
bool Tget(Table,BsObjString *, Value *);
bool Tdelet(Table,BsObjString *);

#endif
