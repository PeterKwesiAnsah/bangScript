#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include "compiler.h"
#include "vm.h"


typedef enum {
    DISASSEMBLER_MODE,
    HELP_MODE,
    REPL_MODE,
    SCRIPT_MODE
} OperationMode;

const char *src;
OperationMode mode=SCRIPT_MODE;

int main(int argc,char *args[]){
    if (argc==1){
    mode=REPL_MODE;
    }else{
        //run in script mode
        // bangscript --help
        // bangscript -h
        // bangscript <filename> --disassembler
        // bangscript <filename> -d
        // bangscript <filename>
        const char * flags[]={[DISASSEMBLER_MODE]="diassembler",[HELP_MODE]="help"};

        for(int i=1;i<argc-1;i++){
            const char * arg=args[i];
            switch (*arg++) {
                case '-':{
                    //handle flag options
                    switch (*arg++) {
                        case 'h':{
                            //either help short hand or invalid command
                            // Str length needs to be 2
                            if(strlen(args[i])==2){
                                //help shorthand
                                mode=HELP_MODE;
                            }
                            fprintf(stderr, "Unknown flag: %s\n", arg);
                            exit(1);
                        }
                        break;
                        case 'd':{
                            //either disassembler short hand or invalid command
                            // Str length needs to be 2
                            // // Str length needs to be 2
                            //
                            if(strlen(args[i])==2){
                                //help shorthand
                                mode=DISASSEMBLER_MODE;
                            }
                            fprintf(stderr, "Unknown flag: %s\n", arg);
                            exit(1);

                        }
                        break;
                        case '-':{
                            //dissembler or help
                            for(int j=0; j < sizeof(flags) / sizeof *flags;j++){
                                unsigned int flagLen=strlen(flags[j]);
                                if(flagLen==strlen(arg) && !memcmp(flags[j], arg,flagLen)){
                                    mode=j;
                                }
                                continue;
                            }
                        }
                        break;
                        default:
                        //invalid command
                        break;
                    }
                }
                break;
                default:
                //must be file name
                break;
            }
        }

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
    }//else{
        //fprintf(stderr,"Usage: %s [path to script]\n",args[0]);
        //exit(1);
    //}
    return 0;
}
