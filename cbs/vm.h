#ifndef VM_H
#define VM_H
typedef enum {
    OP_CONSTANT,       // Load a constant from the constant pool
    OP_NIL,            // Push 'nil'
    OP_TRUE,           // Push 'true'
    OP_FALSE,          // Push 'false'

    OP_POP,            // Pop the top of the stack

    OP_GET_LOCAL,      // Read a local variable
    OP_SET_LOCAL,      // Write a local variable

    OP_GET_GLOBAL,     // Read a global variable
    OP_DEFINE_GLOBAL,  // Define a new global variable
    OP_SET_GLOBAL,     // Write a global variable

    OP_GET_UPVALUE,    // Read a closed-over variable
    OP_SET_UPVALUE,    // Write a closed-over variable

    OP_GET_PROPERTY,   // Read an object property
    OP_SET_PROPERTY,   // Write an object property

    OP_GET_SUPER,      // Read a superclass method

    OP_EQUAL,          // ==
    OP_GREATER,        // >
    OP_LESS,           // <

    OP_ADD,            // +
    OP_SUBTRACT,       // -
    OP_MULTIPLY,       // *
    OP_DIVIDE,         // /

    OP_NOT,            // Logical not
    OP_NEGATE,         // Arithmetic negate (-x)

    OP_PRINT,          // Print top of stack

    OP_JUMP,           // Unconditional jump
    OP_JUMP_IF_FALSE,  // Jump if false
    OP_LOOP,           // Backward jump (for loops)

    OP_CALL,           // Function call
    OP_INVOKE,         // Method call
    OP_SUPER_INVOKE,   // Superclass method call

    OP_CLOSURE,        // Create a closure
    OP_CLOSE_UPVALUE,  // Close an upvalue
    OP_RETURN,         // Return from a function

    OP_CLASS,          // Define a class
    OP_INHERIT,        // Handle class inheritance
    OP_METHOD          // Define a method
} OpCode;

#endif
