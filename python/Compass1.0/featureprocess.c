#include <stdlib.h>
#include <stdbool.h>

int IMP_FRS_FeatSim(void*,void*,unsigned int,float*);

bool str_to_double_array(char* str,double* dArray) {
	char* pEnd;
	dArray[0] = strtod(str,&pEnd);
	int index = 0;
	while(*pEnd) {
		if (index >= 127)
			return false;
		dArray[++index] = strtod(++pEnd,&pEnd);
	}
	return true;
}

float feature_process(char* srcStr,char* desStr) {
	double dSrc[128] = {0.0};
	double dDes[128] = {0.0};
	if (!str_to_double_array(srcStr,dSrc))
		return -1.0;	
	if (!str_to_double_array(desStr,dDes))
		return -1.0;
	float fRes = 0.0;
	IMP_FRS_FeatSim(dSrc,dDes,1,&fRes);
	return fRes;
}
