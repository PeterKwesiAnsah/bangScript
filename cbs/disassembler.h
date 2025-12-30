#ifndef DISASSEMBLER_H
#define DISASSEMBLER_H

#include "chunk.h"

typedef enum {
    DISASSEMBLER_OK,
    DISASSEMBLER_ERROR
} DisassemblerStatus;

DisassemblerStatus disassembleChunk(const char *filename);
DisassemblerStatus disassembleInstruction(uint8_t *ip, uint8_t *start);

#endif