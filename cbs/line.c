#include "line.h"
#include "darray.h"
#include <assert.h>
#include <stddef.h>
#include <stdint.h>

unsigned int curlineNum=0;
uint8_t occur=0;
uint8_t *cur;

unsigned int prevTotal=0;
size_t prevSumIndex=0;

//TODO: line start,line end

// lines[index] -> (offset - 1) , the last bytecode offset for line number(index)
DECLARE_ARRAY(uint8_t, lines);



//indexes of arr, ranges from  0 to n-1 , where in reality are line numbers 1 to n
void addLine(unsigned int line){

    assert(line !=0);
    if(line==curlineNum){
        *cur+=1;
        return;
    }
    curlineNum=line;

    while(line!=lines.len){
        //fill in the holes
         append(lines,uint8_t, 0);
    }

    size_t curIdx=lines.len;
    append(lines,uint8_t, 1);
    cur=&lines.arr[curIdx];
}

unsigned int getLine(size_t offset){

    size_t sum=lines.arr[1];
    size_t i=1;

    for(; offset+1 > sum && i < lines.len;i++){
        sum+=lines.arr[i];
    }

    return i;
}

//getLineFast caches previous sums, applicable when we are getting line information sequentially
unsigned int getLineFast(size_t offset){
    size_t sum=prevTotal;
    size_t i;
    for(i=prevSumIndex; offset>sum;i++){
        sum+=lines.arr[i];
    }
    prevTotal=sum;
    prevSumIndex=i;
    return i+1;
}
