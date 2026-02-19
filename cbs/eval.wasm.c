#include "compiler.h"
#include "eval-wasm.h"
#include "vm.h"


Frame frame;

#ifdef __EMSCRIPTEN__
#include <emscripten.h>
EMSCRIPTEN_KEEPALIVE
#endif
int evalWebSrc(char *input){
    frame.compiler = (Compiler *)alloca(sizeof(Compiler));
    frame.compiler->len = 0;
    frame.compiler->scopeDepth = 0;
    Constants constants = {0};
    frame.constants = &constants;
    frame.src=input;
    CompilerStatus status = compile();
    if (status == COMPILER_ERROR)
      return COMPILER_ERROR;
    return run();
}
