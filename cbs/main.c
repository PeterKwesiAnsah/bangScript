#include <stdio.h>
#include <stdlib.h>
#include "compiler.h"
#include "vm.h"

const char *src;
int main(int argc,char *args[]){
    if (argc==1){
    //run in REPL
    }else if(argc==2){
        //run in script mode
        FILE *fp=fopen(args[1], "r");
        if (fp==NULL){
            fprintf(stderr, "Failed to open File.Perhaps path to file is incorrect");
            exit(1);
        }
        fseek(fp, 0L, SEEK_END);
        size_t fileSize = ftell(fp);
        rewind(fp);
        char *buffer=(char *)malloc(fileSize + 1);
        size_t bytesRead = fread(buffer, sizeof(char), fileSize, fp);
        buffer[bytesRead] = '\0';
        src=buffer;
        fclose(fp);

        CompilerStatus status=compile(src);
        if(status== COMPILER_ERROR) return status;

        return run();
    }else{
        fprintf(stderr,"Usage: %s [path to script]\n",args[0]);
        exit(1);
    }
    return 0;
}
