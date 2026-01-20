#include "vm.h"
#include "chunk.h"
#include "readonly.h"
#include "stack.h"
#include "table.h"
#include <assert.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>

#define READ_BYTE_CODE() (*frame.ip++)
#define EVALUATE_BIN_EXP(operator)                                             \
  do {                                                                         \
    Value b = pop();                                                           \
    Value a = pop();                                                           \
    push(C_DOUBLE_TO_BS_NUMBER(                                                \
        (BS_NUMBER_TO_C_DOUBLE(a) operator BS_NUMBER_TO_C_DOUBLE(b))));        \
  } while (0)

Table globals = {};

ProgramStatus run() {
  Tinit(&globals);
  frame.ip = frame.chunk.arr;
  for (;;) {
    switch (READ_BYTE_CODE()) {
    case OP_CONSTANT_ZER0: {
      Value value = frame.constants->arr[CONSTANT_ZERO_INDEX];
      push(value);
    } break;
    case OP_CONSTANT: {
      uint8_t index = READ_BYTE_CODE();
      Value value = frame.constants->arr[index];
      push(value);
    } break;
    case OP_CONSTANT_LONG: {
      unsigned int index = 0;
      uint8_t highByte = READ_BYTE_CODE();
      uint8_t midByte = READ_BYTE_CODE();
      uint8_t lowByte = READ_BYTE_CODE();
      index = index | lowByte;
      index = index | ((unsigned int)midByte << 8);
      index = index | ((unsigned int)highByte << 16);
      Value value = frame.constants->arr[index];
      push(value);
    } break;
    case OP_POP: {
      pop();
    } break;
    case OP_ADD: {
      Value b = pop();
      Value a = pop();
      if (a.type != b.type) {
        fputs("Add requires operands of the same the type", stderr);
        return ERROR;
      }
      switch (a.type) {
      case TYPE_NUMBER:
        push(C_DOUBLE_TO_BS_NUMBER(
            (BS_NUMBER_TO_C_DOUBLE(a) + BS_NUMBER_TO_C_DOUBLE(b))));
        break;
      case TYPE_OBJ: {
        // Operands a,b can be the same string type or different but they can be
        // used interchangly because they have common fields
        BsObjStringFromSource *BsObjStringA =
            (BsObjStringFromSource *)a.value.obj;
        BsObjStringFromSource *BsObjStringB =
            (BsObjStringFromSource *)b.value.obj;

        size_t ResultStrLen = BsObjStringA->len + BsObjStringB->len;

        BsObjStringFromAlloc *BsObjStringResult =
            (BsObjStringFromAlloc *)malloc(sizeof(BsObjStringFromAlloc) +
                                           ResultStrLen + 1);

        BsObjStringResult->value =
            (char *)(BsObjStringResult + (size_t)sizeof(BsObjStringFromAlloc));
        BsObjStringResult->len = ResultStrLen;

        BsObjStringResult->obj = (BsObj){.type = OBJ_TYPE_STRING_ALLOC};

        memcpy(BsObjStringResult->value, BsObjStringA->value,
               BsObjStringA->len);
        memcpy(BsObjStringResult->value + BsObjStringA->len,
               BsObjStringB->value, BsObjStringB->len);

        BsObjStringResult->value[ResultStrLen] = '\0';
        push(CREATE_BS_OBJ(BsObjStringResult));
      } break;
      default:
        fputs("Add requires operands to be either a number or a string type",
              stderr);
        return ERROR;
      }

    } break;
    case OP_SUB:
      EVALUATE_BIN_EXP(-);
      break;
    case OP_MUL:
      EVALUATE_BIN_EXP(*);
      break;
    case OP_DIV:
      EVALUATE_BIN_EXP(/);
      break;
    case OP_EQUAL: {
      Value b = pop();
      Value a = pop();
      if (a.type != b.type) {
        fputs("Equal comparison requires operands of the same the type",
              stderr);
        return ERROR;
      }
      switch (a.type) {
      case TYPE_NUMBER:
        push(C_BOOL_TO_BS_BOOLEAN(BS_NUMBER_TO_C_DOUBLE(a) ==
                                  BS_NUMBER_TO_C_DOUBLE(b)));
        break;
      case TYPE_BOOL:
        push(C_BOOL_TO_BS_BOOLEAN(BS_BOOLEAN_TO_C_BOOL(a) ==
                                  BS_BOOLEAN_TO_C_BOOL(b)));
        break;
      case TYPE_OBJ: {

        BsObjString *BsObjStringA = (BsObjString *)a.value.obj;
        BsObjString *BsObjStringB = (BsObjString *)b.value.obj;

        if (BsObjStringA->len != BsObjStringB->len) {
          push(C_BOOL_TO_BS_BOOLEAN(false));
        } else {
          push(
              C_BOOL_TO_BS_BOOLEAN(BsObjStringA->value == BsObjStringB->value));
        }
      } break;
      default:
        fputs("Add requires operands to be either a number, bool or a string "
              "type",
              stderr);
        return ERROR;
      }

    } break;
    case OP_EQUAL_NOT: {
      Value b = pop();
      Value a = pop();
      if (a.type != b.type) {
        fputs("Add requires operands of the same the type", stderr);
        return ERROR;
      }
      switch (b.type) {
      case TYPE_NUMBER:
        push(C_BOOL_TO_BS_BOOLEAN(
            !(BS_NUMBER_TO_C_DOUBLE(a) == BS_NUMBER_TO_C_DOUBLE(b))));
        break;
      default:
        fputs("Comparison operation requires operands to be a number", stderr);
        return ERROR;
      }
    } break;
    case OP_GREATOR: {
      Value b = pop();
      Value a = pop();
      if (a.type != b.type) {
        fputs("Add requires operands of the same the type", stderr);
        return ERROR;
      }
      switch (b.type) {
      case TYPE_NUMBER:
        push(C_BOOL_TO_BS_BOOLEAN(
            (BS_NUMBER_TO_C_DOUBLE(a) > BS_NUMBER_TO_C_DOUBLE(b))));
        break;
      default:
        fputs("Comparison operation requires operands to be a number", stderr);
        return ERROR;
      }
    } break;
    case OP_GREATOR_NOT: {
      Value b = pop();
      Value a = pop();
      if (a.type != b.type) {
        fputs("Add requires operands of the same the type", stderr);
        return ERROR;
      }
      switch (a.type) {
      case TYPE_NUMBER:
        push(C_BOOL_TO_BS_BOOLEAN(
            !(BS_NUMBER_TO_C_DOUBLE(a) > BS_NUMBER_TO_C_DOUBLE(b))));
        break;
      default:
        fputs("Comparison operation requires operands to be a number", stderr);
        return ERROR;
      }
    } break;
    case OP_LESS: {
      Value b = pop();
      Value a = pop();
      if (a.type != b.type) {
        fputs("Add requires operands of the same the type", stderr);
        return ERROR;
      }
      switch (b.type) {
      case TYPE_NUMBER:
        push(C_BOOL_TO_BS_BOOLEAN(BS_NUMBER_TO_C_DOUBLE(a) <
                                  BS_NUMBER_TO_C_DOUBLE(b)));
        break;
      default:
        fputs("Comparison operation requires operands to be a number", stderr);
        return ERROR;
      }
    } break;
    case OP_LESS_NOT:
      break;
    case OP_PRINT: {
      Value result = pop();
      switch (result.type) {
      case TYPE_OBJ: {
        switch (result.value.obj->type) {
        case OBJ_TYPE_STRING_SOURCE:
          printf("%.*s\n", ((BsObjStringFromSource *)result.value.obj)->len,
                 ((BsObjStringFromSource *)result.value.obj)->value);
          break;
        case OBJ_TYPE_STRING_ALLOC:
          printf("%s\n", ((BsObjStringFromAlloc *)result.value.obj)->value);
          break;
        default:
          break;
        }
      } break;
      case TYPE_BOOL:
        printf("%d\n", BS_BOOLEAN_TO_C_BOOL(result));
        break;
      case TYPE_NUMBER:
        printf("%f\n", BS_NUMBER_TO_C_DOUBLE(result));
        break;
      case TYPE_NIL:
        printf("null\n");
        break;
      }
    } break;
    // TODO: OP_GLOBALVAR_LONG_DEF
    case OP_GLOBALVAR_DEF: {
      uint8_t varConstIndex = READ_BYTE_CODE();
      Value var = frame.constants->arr[varConstIndex];
      assert(var.type == TYPE_OBJ);
      Value evalrhs = pop();
      Tset(&globals, (BsObjString *)var.value.obj, evalrhs);
    } break;
    case OP_GLOBALVAR_GET: {
      uint8_t *cacheHashIndexIp = frame.ip + 1;
      uint8_t varConstIndex = READ_BYTE_CODE();

      Value var = frame.constants->arr[varConstIndex];
      assert(var.type == TYPE_OBJ);

      uint16_t cacheHashIndex = 0;

      // We can support only 65536 globals at a time
      cacheHashIndex = cacheHashIndex | (uint8_t)READ_BYTE_CODE();
      cacheHashIndex = cacheHashIndex << 8;
      cacheHashIndex = cacheHashIndex | (uint8_t)READ_BYTE_CODE();

      if (globals.len && cacheHashIndex < globals.cap &&
          globals.arr[cacheHashIndex].key ==
              (BsObjString *)((BsObjStringFromSource *)var.value.obj)->value) {
        push(globals.arr[cacheHashIndex].value);
        break;
      }

      Value value={0};
      if (Tget(&globals, (BsObjString *)var.value.obj, &value,
               (uint32_t *)&cacheHashIndex)) {
        *cacheHashIndexIp++ = (cacheHashIndex >> 8) & 0xFF;
        *cacheHashIndexIp = cacheHashIndex & 0xFF;
        //Test the cacheIndex
        push(value);
        break;
      };
      fprintf(stderr, "%.*s is undefined.\n",
              ((BsObjStringFromSource *)var.value.obj)->len,
              ((BsObjStringFromSource *)var.value.obj)->value);
      return ERROR;
    } break;
    // TODO: OP_GLOBALVAR_LONG_ASSIGN
    case OP_GLOBALVAR_ASSIGN: {
      uint8_t varConstIndex = READ_BYTE_CODE();
      Value var = frame.constants->arr[varConstIndex];
      assert(var.type == TYPE_OBJ);
      Value evalrhs = pop();
      // TODO: check for undefined vars
      if (Tset(&globals, (BsObjString *)var.value.obj, evalrhs)) {
        fprintf(stderr, "%.*s is undefined.\n",
                ((BsObjStringFromSource *)var.value.obj)->len,
                ((BsObjStringFromSource *)var.value.obj)->value);
        return ERROR;
      };

    } break;
    case OP_LOCALVAR_GET: {
      uint8_t slot = *frame.ip++;
      push(getStackItem(slot));
    } break;
    case OP_LOCALVAR_ASSIGN: {
      uint8_t slot = *frame.ip++;
      updateStackItem(slot, getStackItem(slot));
    } break;
    case OP_RETURN:
      return SUCCESS;
    default:
      fputs("Invalid bytecode instruction\n", stderr);
      return ERROR;
    }
  }
  return SUCCESS;
};
