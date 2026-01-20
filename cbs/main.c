#include "compiler.h"
#include "darray.h"
#include "disassembler.h"
#include "readonly.h"
#include "vm.h"
#include <alloca.h>
#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef enum {
  DISASSEMBLER_MODE,
  HELP_MODE,
  REPL_MODE,
  SCRIPT_MODE
} OperationMode;

Frame frame;

OperationMode mode = SCRIPT_MODE;

const char *readSourceFintoBuffer(const char *filename) {
  assert(filename);

  FILE *fp = fopen(filename, "r");
  if (fp == NULL) {
    perror("Failed to open file");
    exit(1);
  }
  fseek(fp, 0L, SEEK_END);
  size_t fileSize = ftell(fp);
  rewind(fp);
  char *buffer = (char *)malloc(fileSize + 1);
  size_t bytesRead = fread(buffer, sizeof(char), fileSize, fp);
  buffer[bytesRead] = '\0';
  frame.src = buffer;
  fclose(fp);
  return buffer;
}

int main(int argc, char *args[]) {
  const char *filename = NULL;
  frame.compiler = (Compiler *)alloca(sizeof(Compiler));
  frame.compiler->len = 0;
  frame.compiler->scopeDepth = 0;
  Constants constants = {0};
  frame.constants = &constants;

  if (argc == 1) {
    mode = REPL_MODE;
  } else {
    // bangscript --help
    // bangscript -h
    // bangscript <filename> --disassembler
    // bangscript <filename> -d
    // bangscript <filename>
    const char *flags[] = {
        [DISASSEMBLER_MODE] = "diassembler", [HELP_MODE] = "help"};

    for (int i = 1; i < argc; i++) {
      const char *arg = args[i];
      switch (*arg) {
      case '-': {
        arg++;
        // handle flag options
        switch (*arg++) {
        case 'h': {
          // either help short hand or invalid command
          //  Str length needs to be 2
          if (strlen(args[i]) == 2) {
            // help shorthand
            mode = HELP_MODE;
            break;
          }
          fprintf(stderr, "Unknown flag: %s\n", arg);
          exit(1);
        } break;
        case 'd': {
          // either disassembler short hand or invalid command
          //  Str length needs to be 2
          //  // Str length needs to be 2
          //
          if (strlen(args[i]) == 2) {
            // help shorthand
            mode = DISASSEMBLER_MODE;
            break;
          }
          fprintf(stderr, "Unknown flag: %s\n", arg);
          exit(1);
        } break;
        case '-': {
          for (int j = 0; j < sizeof(flags) / sizeof *flags; j++) {
            unsigned int flagLen = strlen(flags[j]);
            if (flagLen == strlen(arg) && !memcmp(flags[j], arg, flagLen)) {
              mode = j;
            }
            continue;
          }
        } break;
        default:
          // invalid command
          fprintf(stderr, "Unknown flag: %s\n", arg);
          exit(1);
          break;
        }
      } break;
      default:
        filename = arg;
        break;
      }
    }
  }

  switch (mode) {
  case REPL_MODE: {
    printf("Running bangscript in REPL mode\n");
    DECLARE_ARRAY_TYPE(char, Input)
    size_t scopeDepth = 0;
    size_t start = 0;
    Input src = {0};

  Loop:
    for (;;) {
      printf(">>> ");
    // read src from stdin
    Read:
      while (1) {
        int ch = getchar();
        if (ch == EOF) {
          return SUCCESS;
        } else if (ch == '\n') {
          // We don't break out easily
          //  if scopedepth is 0 and the buffer src looks well like a completed
          //  statement then we breakout else we keep on asking for input
          if (src.cap && src.len) {
            src.arr[src.len] = '\0';
          }
          if (scopeDepth > 0) {
            //TODO: support copy and paste
            printf("... ");
            continue;
          }
          break;
        }
        if (ch == '{') {
          scopeDepth++;
        } else if (ch == '}') {
          scopeDepth--;
        }
        append(src, char, ch);
      }
    Evaluate:
      if (src.len == 0)
        continue;
      if (!strcmp("q", src.arr) || !strcmp("quit", src.arr))
        return SUCCESS;
      frame.src = src.arr;
      CompilerStatus status = compile();
      if (status == COMPILER_ERROR)
        continue;
      run();
      src.len = 0;
    }
  } break;
  // TODO:
  case HELP_MODE:
    break;
  case DISASSEMBLER_MODE: {
    frame.src = readSourceFintoBuffer(filename);
    CompilerStatus status = compile();
    if (status == COMPILER_ERROR)
      return status;
    return disassembleChunk(filename);
  } break;
  case SCRIPT_MODE: {
    frame.src = readSourceFintoBuffer(filename);
    CompilerStatus status = compile();
    if (status == COMPILER_ERROR)
      return status;
    return run();
  } break;
  default:
    break;
  }
  return 0;
}
