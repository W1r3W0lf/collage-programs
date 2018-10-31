#include <stdio.h>
#include <stdbool.h>
#include <math.h>
#include <stdlib.h>

unsigned long elrPolly(int x, int c){
  unsigned long n = x;
  n = n*n+n+c;
  return n;
}

bool isPrime(unsigned long number){
  double srt = sqrt(number)+1;
  for(int i=2; i <= (int)srt;i++){
    if(number%i==0)
      return false;
  }
  return true;
}

int sweetness(unsigned long* test){
  int numPrimes = 0;
  for(int i=0; i<1000; i++){
    if(isPrime(test[i])){
      numPrimes++;
    }
  }
  return numPrimes;
}

int fastSweetness(unsigned long* test){
  int total = 0;
  for(int i=0 ; i<1000; i++){
    total += test[i];
  }
  return total;
}

//#define MAXNUM 4000000


#define MAXNUM  100000

int findSweetConst(){
  int constint = 1;
  int testSweet = 0;
  int bestSweet = 0;
  long testSpace[1000];

  for(int c = 1 ; c < 1000 ; c+=2){
    for(int x = 0 ; x < 1000 ; x++){
      testSpace[x] = elrPolly(x,c);
    }
    testSweet = sweetness(testSpace);
    if(testSweet > bestSweet){
      bestSweet = testSweet;
      constint = c;
    }
  }

  return constint;
}

int findOffset(unsigned long* land ){
  unsigned long* finger;
  int testSweet;
  int bestSweet = 0;
  int myBestOffset = 0;

  for(int offset = 0 ; offset < (MAXNUM-1000) ; offset++){
    finger = &(land[offset]);
    testSweet = fastSweetness(finger);
    if(testSweet > bestSweet){
      myBestOffset = offset;
      bestSweet = testSweet;
      if(testSweet == 1000){
        for(int x = 0; x<1000; x++){
          printf("%lu\n",land[offset+x]);
        }
        return offset;
      }
    }

    if(offset%1000 == 0){
      printf("%d ~ %f : %d - %d\n",offset,(offset/(float)MAXNUM),bestSweet,testSweet);
    }

  }
  return myBestOffset;
}

void setPrimes(unsigned long* numbers){

  for(int x=0; x<MAXNUM; x++){

    if(x%10000==0){
      printf("%d : %f\n",x,(x/(float)MAXNUM));
    }

    if(isPrime(numbers[x])){
      numbers[x]=1;
    }else{
      numbers[x]=0;
    }
  }
}

int main(){

  unsigned long* space;
  space = malloc(MAXNUM*sizeof(unsigned long));

  printf("Finding Sweet Const\n");
  int sweetConst = findSweetConst();

  printf("Sweet C:%d\n",sweetConst);

  printf("calculating elrPollys\n");
  for(int offset = 0 ; offset < MAXNUM ; offset++){
    space[offset] = elrPolly(offset,sweetConst);
  }

  printf("setting primes\n");
  setPrimes(space);

  printf("Finding Offset\n");
  int bestOffset = findOffset(space);

  free(space);

  printf("\n***\nBest Const:%d\nBest Offset:%d\n***\n",sweetConst, bestOffset);
  return 0;
}
