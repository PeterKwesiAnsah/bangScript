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

DECLARE_ARRAY(uint8_t, lines);

//indexes of arr, ranges from  0 to n-1 , where in reality are line numbers 1 to n
void add(unsigned int line){
    if(line==curlineNum){
        *cur+=1;
        return;
    }
    curlineNum=line;
    while(line!=lines.len+1){
        //fill in the holes
         append(lines, 0, sizeof(uint8_t));
    }
    cur=&lines.arr[lines.len];
    append(lines, 1, sizeof(uint8_t));
}

unsigned int get(size_t offset){
    size_t sum=prevTotal;
    size_t i;
    for(i=prevSumIndex; offset>sum;i++){
        sum+=lines.arr[i];
    }
    prevTotal=sum;
    prevSumIndex=i;
    return i+1;
}
