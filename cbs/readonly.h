//Number Constants and String literals
#ifndef READONLY_H
#define READONLY_H
#include "darray.h"
#include <stddef.h>
#define CONSTANT_LIMIT 256
#define CONSTANT_LONG_LIMIT 16777216
typedef double Value;


size_t addConstant(Value);
#endif
