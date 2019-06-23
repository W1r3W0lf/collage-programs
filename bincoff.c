#include <stdio.h>
#include <stdlib.h>

int* pascal(int size){

	int* thisRow;

	if (size == 1){
		thisRow = malloc(sizeof(int));
		thisRow[0] = 1;
	}else{
	
		int* lastRow = pascal(size-1);
	
		thisRow = malloc(sizeof(int)*size);

		thisRow[0] = 1;
		thisRow[size-1] = 1;
		for(int i = 1 ; i < size-1 ; i++){
			thisRow[i] = lastRow[i-1] + lastRow[i];
		}

		free(lastRow);
	}

	return thisRow;
}

int main(int argc, char** argv){
	
	if(argc < 2){
		printf("Not enough arguments\n");
		return 1;
	}
	

	int size = atoi(argv[1])+1;

	int* pass = pascal(size);

	if (argc == 2){
		for(int i=0; i<size; i++){
			printf("%d,",pass[i]);
		}
		printf("\n");
	} else{
		printf("%d\n",pass[atoi(argv[2])]);
	}

	
	
	free(pass);

	return 0;
}
