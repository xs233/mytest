#include <stdlib.h>
#include <time.h>

int feature_process(char* src_feature,int src_len,char* des_feature,int des_len)
{
	srand((unsigned)time( NULL ));  
	//randomize(); 
	return rand()%100 + 1;
}
