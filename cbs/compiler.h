#ifndef COMPILER_H
#define COMPILER_H
typedef enum {
    COMPILER_OK,
    COMPILER_ERROR
} CompilerStatus;

CompilerStatus compile(const char *);
#endif
