#ifndef TABLE_H
#define TABLE_H

#include "readonly.h"
#include <stdbool.h>

#define LOAD_FACTOR_MAX 0.75
#define LOAD_FACTOR_MIN 0.1
#define INIT_TABLE_SIZE 256

struct KVnode {
BsObjString *key;
Value value;
};

typedef struct KVnode Tnode;

//len -> number of filled Tnodes
//capacity -> Total
DECLARE_ARRAY_TYPE(Tnode,Table);

void Tinit(Table *);

static bool Tgrow(Table, Table *);
void Tcopy(Table *,Table *);


bool Tset(Table *,BsObjString *, Value);
bool Tget(Table *,BsObjString *, Value *);
bool Tdelete(Table *,BsObjString *);

BsObjString *Tgets(Table *,BsObjString *,Value *);

#endif
