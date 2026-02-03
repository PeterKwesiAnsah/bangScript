#ifndef TABLE_H
#define TABLE_H

#include "readonly.h"
#include <stdbool.h>
#include <stdint.h>

#define LOAD_FACTOR_MAX 0.75
#define LOAD_FACTOR_MIN 0.1
#define INIT_TABLE_SIZE 256

#define TABLE_EXPAND(tbl_ptr)                                                  \
  do {                                                                         \
    Table _tnew;                                                               \
    if (Tgrow((tbl_ptr), &_tnew)) {                                            \
      Tcopy((tbl_ptr), &_tnew);                                                \
      *(tbl_ptr) = _tnew;                                                      \
    }                                                                          \
  } while (0)

struct KVnode {
  BsObjString *key;
  Value value;
};

typedef struct KVnode Tnode;

// len -> number of filled Tnodes
// capacity -> Total
DECLARE_ARRAY_TYPE(Tnode, Table);

void Tinit(Table *);

bool Tgrow(Table *, Table *);
void Tcopy(const Table *, Table *);

bool Tset(Table *, BsObjString *, Value);
bool Tget(Table *, BsObjString *, Value *, uint32_t *);
bool Tdelete(Table *, BsObjString *);

BsObjString *Tgets(Table *, BsObjString *, Value *);
bool Tsets(Table *Tinstance, BsObjString *key, Value value);

#endif
