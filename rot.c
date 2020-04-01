#include <stdio.h>
#include <string.h>
#include <stdlib.h>

int main(int argc, char **argv){
	
	if(argc < 3){
		printf("Usage: %s #rot string\n", argv[0]);
		return 1;
	}
	
	int num = atoi(argv[1]);

	char *output = malloc(strlen(argv[2]));
	char *out = output;
	char *letter = argv[2];

	int work = 0;
	int offset = 0;

	while(*letter){

		if((int)(*letter) > 64 && (int)(*letter) < 91){
			offset = 65;
		}

		else if((int)(*letter) > 96 && (int)(*letter) < 123){
			offset = 97;
		}			
		else{
			offset = 0;
			*out = *letter;
		}

		if(offset){
			work = (int)(*letter) - offset;
			*out = (char)(offset + ((26 + work + num) % 26));
		}


		out++;
		letter++;
	}
	*out = '\0';

	printf("%s\n", output);

	free(output);

	return 0;
}
