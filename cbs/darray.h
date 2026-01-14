#ifndef DYNAMIC_ARRAY
#define DYNAMIC_ARRAY
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>

#define DECLARE_ARRAY_TYPE(type, name)                                         \
  typedef struct {                                                             \
    size_t cap;                                                                \
    size_t len;                                                                \
    type *arr;                                                                 \
  } name;

#define DECLARE_ARRAY(type, name)                                              \
  struct {                                                                     \
    size_t cap, len;                                                           \
    type *arr;                                                                 \
  } name

#define DEFAULT_SLICE_CAP 256

#define growDarrPtr(arrptr, type, size, count)                                 \
  do {                                                                         \
    {                                                                          \
      type *ptr = (type *)realloc(arrptr->arr, count * size);                  \
      if (ptr != NULL) {                                                       \
        arrptr->cap = count;                                                   \
        arrptr->arr = ptr;                                                     \
      } else {                                                                 \
        fputs("Not enough memory", stderr);                                    \
        exit(1);                                                               \
      }                                                                        \
    }                                                                          \
  } while (0)

#define growDarr(array, type, size, count)                                     \
  do {                                                                         \
    {                                                                          \
      type *ptr = (type *)realloc(array.arr, count * size);                    \
      if (ptr != NULL) {                                                       \
        array.cap = count;                                                     \
        array.arr = ptr;                                                       \
      } else {                                                                 \
        fputs("Not enough memory", stderr);                                    \
        exit(1);                                                               \
      }                                                                        \
    }                                                                          \
  } while (0)

#define append(array, type, el)                                                \
  do {                                                                         \
    if (array.len >= array.cap) {                                              \
      size_t count = array.cap == 0 ? DEFAULT_SLICE_CAP : (2 * array.cap);     \
      growDarr(array, type, sizeof(type), count);                              \
      array.arr[array.len++] = (type)el;                                       \
    } else {                                                                   \
      array.arr[array.len++] = (type)el;                                       \
    }                                                                          \
  } while (0)

#define appendPtr(arrayPtr, type, el)                                          \
  do {                                                                         \
    if ((arrayPtr)->len >= (arrayPtr)->cap) {                                  \
      size_t count =                                                           \
          (arrayPtr)->cap == 0 ? DEFAULT_SLICE_CAP : (2 * (arrayPtr)->cap);    \
      growDarrPtr(arrayPtr, type, sizeof(type), count);                        \
      (arrayPtr)->arr[(arrayPtr)->len++] = (type)el;                           \
    } else {                                                                   \
      (arrayPtr)->arr[(arrayPtr)->len++] = (type)el;                           \
    }                                                                          \
  } while (0)
#endif
