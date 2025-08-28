#### BangScript
A JavaScript-like (not-JavaScript) scripting language built from scratch for fun and learning purposes. It’s written in Go, with plans to port the interpreter to C for performance, reduced runtime overhead and lower-level control.

#### Features
- Dynamic Type System
- Functions/Closures
- Automatic Memory Management
- Block c-style comments
- Multi-line string literals

##### Getting Started
###### Requirements
- Go 1.21+
- Git
###### To run the REPL
```bash
git clone https://github.com/peterkwesiansah/bangscript.git
cd bangScript/gbs
go run main.go
```
###### or a script
```bash
git clone https://github.com/peterkwesiansah/bangscript.git
cd bangScript/gbs
go build -o bs
./bs examples/hello.bs
```

#####  Running a Real Script
You can find this script in `examples/counter.bs`
```javascript
fun makeCounter() {
  var count = 0;
  fun inc() {
    count = count + 1;
    print count;
  }
  return inc;
}
var counter = makeCounter();

counter(); // → 1
counter(); // → 2
counter(); // → 3
```

#### What's Next
- [x] Lexer/Scanner
- [x] Parser
- [x] Resolver
- [x] Tree-walk Interpreter
- [x] Interactive REPL
- [x] Finalizing Closures
- [ ] Bytecode VM

#### Credits
This Language is heavily inspired by [lox](https://craftinginterpreters.com/)
