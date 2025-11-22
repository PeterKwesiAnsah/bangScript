#ifndef TABLE_H
#define TABLE_H

#include "readonly.h"
#include <stdbool.h>

#define LOAD_FACTOR_MAX 0.75
#define LOAD_FACTOR_MIN 0.1

struct KVnode {
BsObjString *key;
Value value;
};

typedef struct KVnode Tnode;

//len -> can be used to keep track of filled Tnodes
//capacity -> Size of Array of Head Tnodes
DECLARE_ARRAY_TYPE(Tnode,Table);

inline void Tinit(Table *);

static bool Tgrow(Table, Table *);
void Tcopy(Table *,Table *);


bool Tset(Table *,BsObjString *, Value);
bool Tget(Table *,BsObjString *, Value *);
bool Tdelete(Table *,BsObjString *);

#endif
