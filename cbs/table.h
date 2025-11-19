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

//len -> can be used to keep track of filled Tnodes
//capacity -> Size of Array of Head Tnodes
DECLARE_ARRAY_TYPE(Tnode,Table);

bool Tset(Table *,BsObjString *, Value);
inline void Tinit(Table *);
bool Tget(Table *,BsObjString *, Value *);
bool Tdelete(Table *,BsObjString *);

#endif
