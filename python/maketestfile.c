#include<stdio.h>

void makefile()
{
	FILE* fp = fopen("text.txt","wb+");
	if (!fp)
		perror("open file failed");
	
	int i;
	for (i = 0;i < 1000000;++i)
	{
		fputc((char)(rand()%128),fp);
	}
	fclose(fp);
}

int main()
{
	makefile();
	return 0;
}
