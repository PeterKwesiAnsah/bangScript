#include "table.h"
#include "darray.h"
#include <stdbool.h>
#include <stddef.h>
#include <stdio.h>
#include <string.h>

void Tinit(Table *Tinstance) {
  Tinstance->len = 0;
  Tinstance->arr = NULL;
  growDarrPtr(Tinstance, Tnode, sizeof(Tnode), INIT_TABLE_SIZE);
  memset(Tinstance->arr, 0, INIT_TABLE_SIZE);
}

void Tcopy(const Table *Told, Table *Tnew) {

  size_t cap = Told->cap;
  Tnode *Toldarr = Told->arr;
  Tnode *Tnewarr = Tnew->arr;
  size_t i;

  size_t limit = (cap < 3) ? 0 : cap - 2;

  for (i = 0; i < limit; i += 3) {

    if (Toldarr[i].key != NULL) {
      uint32_t index = Toldarr[i].key->hash % Tnew->cap;
      Tnode *node = &Tnewarr[index];
      while (node->key != NULL) {
        index = (index + 1) % Tnew->cap;
        node = &Tnewarr[index];
      }
      node->key = Toldarr[i].key;
      node->value = Toldarr[i].value;
    }

    if (Toldarr[i + 1].key != NULL) {
      uint32_t index = Toldarr[i + 1].key->hash % Tnew->cap;
      Tnode *node = &Tnewarr[index];
      while (node->key != NULL) {
        index = (index + 1) % Tnew->cap;
        node = &Tnewarr[index];
      }
      node->key = Toldarr[i + 1].key;
      node->value = Toldarr[i + 1].value;
    }

    if (Toldarr[i + 2].key != NULL) {
      uint32_t index = Toldarr[i + 2].key->hash % Tnew->cap;
      Tnode *node = &Tnewarr[index];
      while (node->key != NULL) {
        index = (index + 1) % Tnew->cap;
        node = &Tnewarr[index];
      }
      node->key = Toldarr[i + 2].key;
      node->value = Toldarr[i + 2].value;
    }
  }

  for (; i < cap; i++) {
    if (Toldarr[i].key == NULL)
      continue;

    uint32_t index = Toldarr[i].key->hash % Tnew->cap;
    Tnode *node = &Tnewarr[index];

    while (node->key != NULL) {
      index = (index + 1) % Tnew->cap;
      node = &Tnewarr[index];
    }
    node->key = Toldarr[i].key;
    node->value = Toldarr[i].value;
  }
  // TODO: free Told
}
// returns true if Tcur is closer to LOAD_FACTOR_MAX, false otherwise
bool Tgrow(Table *Tcur, Table *Tnew) {
  if (((double)(Tcur->len + 1) / (Tcur->cap)) >= LOAD_FACTOR_MAX) {

    size_t cap = (size_t)(Tcur->len + 1) / (LOAD_FACTOR_MIN);
    Table temp = {.len = Tcur->len, .arr = NULL};

    growDarr(temp, Tnode, sizeof(Tnode), cap);
    memset(temp.arr, 0, cap);
    temp.cap = cap;
    *Tnew = temp;

    return true;
  }
  return false;
}

// true if it inserts into an empty bucket, false if it updated one
bool Tset(Table *Tinstance, BsObjString *key, Value value) {
  size_t cap = Tinstance->cap;
  uint32_t index = key->hash % cap;
  Tnode *node = &Tinstance->arr[index];

  while (node->key != key && node->key != NULL) {
    index = (index + 1) % cap;
    node = &Tinstance->arr[index];
  }

  if (node->key == key) {
    // filled node
    node->value = value;
    return false;
  }

  // regular empty node
  node->key = key;
  node->value = value;
  Tinstance->len++;

  return true;
}

// true if entry was  found, false otherwise
bool Tget(Table *Tinstance, BsObjString *key, Value *value,
          uint32_t *foundIndex) {
  size_t cap = Tinstance->cap;
  uint32_t index = key->hash % cap;
  Tnode *node = &Tinstance->arr[index];

  while (node->key != key && node->key != NULL) {
    index = (index + 1) % cap;
    node = &Tinstance->arr[index];
  }

  if (node->key == NULL) {
    return false;
  }

  *value = node->value;
  *foundIndex = index;
  return true;
}

bool Tsets(Table *Tinstance, BsObjString *key, Value value) {
  size_t cap = Tinstance->cap;
  uint32_t index = key->hash % cap;
  Tnode *node = &Tinstance->arr[index];

  while (node->key != NULL) {
    index = (index + 1) % cap;
    node = &Tinstance->arr[index];
  }

  // regular empty node
  node->key = key;
  node->value = value;
  Tinstance->len++;

  return true;
}

BsObjString *Tgets(Table *Tinstance, BsObjString *key, Value *value) {
  size_t cap = Tinstance->cap;
  uint32_t index = key->hash % cap;
  Tnode *node = &Tinstance->arr[index];

  while (node->key != NULL &&
         !(key->len == node->key->len && key->hash == node->key->hash &&
           !memcmp(key->value, node->key->value, node->key->len))) {
    index = (index + 1) % cap;
    node = &Tinstance->arr[index];
  }

  if (node->key == NULL) {
    return NULL;
  }

  // All strings(from source) at compile time finds it way to the constant
  // table, as with any value in the source, at runtime the VM needs a way to
  // look it up
  *value = node->value;
  return node->key;
}

bool Tdelete(Table *Tinstance, BsObjString *key) {
  size_t cap = Tinstance->cap;
  uint32_t index = key->hash % cap;
  Tnode *node = &Tinstance->arr[index];

  while (node->key != NULL && node->key != key) {
    index = (index + 1) % cap;
    node = &Tinstance->arr[index];
  }

  if (node->key == NULL) {
    return false;
  }

  uint32_t peekIndex = index;
  Tnode *hole = node;

  while (Tinstance->arr[peekIndex = ((peekIndex + 1) % cap)].key != NULL) {
    hole = &Tinstance->arr[peekIndex];
    bool sitsRightfully = (hole->key->hash % cap) == peekIndex;
    // we skip nodes that are rightfully in their buckets
    if (sitsRightfully)
      continue;
    *node = *hole;
  }

  hole->key = NULL;
  Tinstance->len--;

  return true;
}
