#include "disassembler.h"
#include "chunk.h"
#include "readonly.h"
#include "line.h"
#include "darray.h"
#include <stdio.h>
#include <stdint.h>

extern DECLARE_ARRAY(Value, constants);
extern DECLARE_ARRAY(u_int8_t, chunk);

static const char *opcodeNames[] = {
    "OP_CONSTANT_ZER0",
    "OP_CONSTANT",
    "OP_CONSTANT_LONG",
    "OP_ADD",
    "OP_SUB",
    "OP_MUL",
    "OP_DIV",
    "OP_NEGATE",
    "OP_EQUAL",
    "OP_EQUAL_NOT",
    "OP_GREATOR",
    "OP_LESS_NOT",
    "OP_LESS",
    "OP_GREATOR_NOT",
    "OP_GLOBALVAR_DEF",
    "OP_GLOBALVAR_GET",
    "OP_GLOBALVAR_ASSIGN",
    "OP_PRINT",
    "OP_RETURN"
};

static void printConstant(uint8_t constantIndex) {
    if (constantIndex >= constants.len) {
        printf("<invalid constant index %d>", constantIndex);
        return;
    }

    Value constant = constants.arr[constantIndex];
    switch (constant.type) {
        case TYPE_NUMBER:
            printf("%g", BS_NUMBER_TO_C_DOUBLE(constant));
            break;
        case TYPE_BOOL:
            printf("%s", BS_BOOLEAN_TO_C_BOOL(constant) ? "true" : "false");
            break;
        case TYPE_NIL:
            printf("nil");
            break;
        case TYPE_OBJ: {
            BsObj *obj = constant.value.obj;
            if (obj->type == OBJ_TYPE_STRING_SOURCE) {
                BsObjStringFromSource *str = (BsObjStringFromSource *)obj;
                printf("'%.*s'", str->len, str->value);
            } else if (obj->type == OBJ_TYPE_STRING_ALLOC) {
                BsObjStringFromAlloc *str = (BsObjStringFromAlloc *)obj;
                printf("'%s'", str->value);
            } else {
                printf("<unknown object>");
            }
            break;
        }
        default:
            printf("<unknown type>");
            break;
    }
}

static void printConstantLong(unsigned int constantIndex) {
    if (constantIndex >= constants.len) {
        printf("<invalid constant index %u>", constantIndex);
        return;
    }

    Value constant = constants.arr[constantIndex];
    switch (constant.type) {
        case TYPE_NUMBER:
            printf("%g", BS_NUMBER_TO_C_DOUBLE(constant));
            break;
        case TYPE_BOOL:
            printf("%s", BS_BOOLEAN_TO_C_BOOL(constant) ? "true" : "false");
            break;
        case TYPE_NIL:
            printf("nil");
            break;
        case TYPE_OBJ: {
            BsObj *obj = constant.value.obj;
            if (obj->type == OBJ_TYPE_STRING_SOURCE) {
                BsObjStringFromSource *str = (BsObjStringFromSource *)obj;
                printf("'%.*s'", str->len, str->value);
            } else if (obj->type == OBJ_TYPE_STRING_ALLOC) {
                BsObjStringFromAlloc *str = (BsObjStringFromAlloc *)obj;
                printf("'%s'", str->value);
            } else {
                printf("<unknown object>");
            }
            break;
        }
        default:
            printf("<unknown type>");
            break;
    }
}


//Handle Line Information
DisassemblerStatus disassembleInstruction(uint8_t *ip, uint8_t *start) {
    size_t offset = ip - start;
    int line = getLine(offset);

    BS_OP_CODES instruction = *ip;
    if (instruction > OP_RETURN || instruction < OP_CONSTANT_ZER0) {
        printf("UNKNOWN_OPCODE: %d\n", instruction);
        return DISASSEMBLER_ERROR;
    }

    // Print: Line | Offset | Opcode
    printf("%4d | %06zu | %-20s | ", line, offset, opcodeNames[instruction]);
    ip++;

    switch (instruction) {
        case OP_CONSTANT_ZER0:
            printf("constant[0] = ");
            printConstant(CONSTANT_ZERO_INDEX);
            break;

        case OP_CONSTANT: {
            uint8_t constantIndex = *ip++;
            printf("constant[%3d] = ", constantIndex);
            printConstant(constantIndex);
            break;
        }

        case OP_CONSTANT_LONG: {
            unsigned int constantIndex = 0;
            uint8_t highByte = *ip++;
            uint8_t midByte = *ip++;
            uint8_t lowByte = *ip++;
            constantIndex = constantIndex | lowByte;
            constantIndex = constantIndex | ((unsigned int)midByte << 8);
            constantIndex = constantIndex | ((unsigned int)highByte << 16);
            printf("constant[%6u] = ", constantIndex);
            printConstantLong(constantIndex);
            break;
        }

        case OP_GLOBALVAR_DEF:
        case OP_GLOBALVAR_GET:
        case OP_GLOBALVAR_ASSIGN: {
            uint8_t varConstIndex = *ip++;
            printf("constant[%3d] = ", varConstIndex);
            printConstant(varConstIndex);
            break;
        }

        case OP_ADD:
        case OP_SUB:
        case OP_MUL:
        case OP_DIV:
        case OP_NEGATE:
        case OP_EQUAL:
        case OP_EQUAL_NOT:
        case OP_GREATOR:
        case OP_LESS_NOT:
        case OP_LESS:
        case OP_GREATOR_NOT:
        case OP_PRINT:
        case OP_RETURN:
            // No additional details needed
            break;

        default:
            printf("<unhandled instruction>");
            return DISASSEMBLER_ERROR;
    }

    printf("\n");
    return DISASSEMBLER_OK;
}

DisassemblerStatus disassembleChunk(const char *filename) {

    printf("\n");
    printf("===== Disassembly of '%s' =====\n", filename);
    printf("LINE | OFFSET | INSTRUCTION          | DETAILS\n");
    printf("-----|--------|----------------------|----------------------------------\n");

    if (chunk.len == 0) {
        printf("<empty chunk>\n");
        return DISASSEMBLER_OK;
    }

    uint8_t *ip = chunk.arr;
    uint8_t *end = chunk.arr + chunk.len;

    while (ip < end) {
        DisassemblerStatus status = disassembleInstruction(ip, chunk.arr);
        if (status == DISASSEMBLER_ERROR) {
            return status;
        }

        // Advance to next instruction
        BS_OP_CODES instruction = *ip;
        ip++;

        switch (instruction) {
            case OP_CONSTANT_ZER0:
                // No operand
                break;

            case OP_CONSTANT:
            case OP_GLOBALVAR_DEF:
            case OP_GLOBALVAR_GET:
            case OP_GLOBALVAR_ASSIGN:
                // One byte operand
                ip++;
                break;

            case OP_CONSTANT_LONG:
                // Three byte operand
                ip += 3;
                break;

            case OP_ADD:
            case OP_SUB:
            case OP_MUL:
            case OP_DIV:
            case OP_NEGATE:
            case OP_EQUAL:
            case OP_EQUAL_NOT:
            case OP_GREATOR:
            case OP_LESS_NOT:
            case OP_LESS:
            case OP_GREATOR_NOT:
            case OP_PRINT:
            case OP_RETURN:
                // No operands
                break;

            default:
                printf("Error: Unknown opcode %d at offset %zu\n", instruction, ip - chunk.arr - 1);
                return DISASSEMBLER_ERROR;
        }
    }

    printf("===== End of Disassembly =====\n\n");
    return DISASSEMBLER_OK;
}
