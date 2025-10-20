#ifndef COMPILER_H
#define COMPILER_H
typedef enum {
    COMPILER_ERROR,
    COMPILER_OK,
} CompilerStatus;

CompilerStatus compile(const char *);
#endif
