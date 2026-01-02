#include "line.h"
#include "darray.h"
#include <stddef.h>
#include <stdint.h>


struct {
    uint8_t *cur;
    unsigned curlineNum;
    uint8_t occur;
} AddLineInfoTracker;


struct {
    size_t prevLine;
    unsigned prevtotalLineOffsets;
} GetLineInfoTrackerFast;


//TODO: line start,line end

// lines[index] -> (offset - 1) , the last bytecode offset for line number(index)
DECLARE_ARRAY(uint8_t, lines);

//indexes of arr, ranges from  0 to n-1 , where in reality are line numbers 1 to n
void addLine(unsigned int line){

    if(line==AddLineInfoTracker.curlineNum){
        *AddLineInfoTracker.cur+=1;
        return;
    }
    AddLineInfoTracker.curlineNum=line;

    while(line!=lines.len){
        //fill in the holes
         append(lines,uint8_t, 0);
    }

    size_t curIdx=lines.len;
    append(lines,uint8_t, 1);
    AddLineInfoTracker.cur=&lines.arr[curIdx];
}

unsigned int getLine(size_t offset){

    size_t totalLineOffsets=lines.arr[1];
    size_t line=1;

    for(; offset+1 > totalLineOffsets && line+1 < lines.len;line++){
        totalLineOffsets+=lines.arr[line+1];
    }

    return line;
}

//getLineFast caches previous sums, applicable when we are getting line information sequentially
unsigned int getLineFast(size_t offset){
    size_t prevtotalLineOffsets=GetLineInfoTrackerFast.prevtotalLineOffsets;
    size_t prevLine=GetLineInfoTrackerFast.prevLine;
    for(; offset+1 > prevtotalLineOffsets && prevLine+1 < lines.len; prevLine++){
        prevtotalLineOffsets+=lines.arr[prevLine+1];
    }
    prevtotalLineOffsets=prevtotalLineOffsets;
    prevLine=prevLine;
    return prevLine;
}
