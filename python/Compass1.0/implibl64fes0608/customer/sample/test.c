#include "../include/imp_frs_api.h"
#include <stdio.h>

void main()
{
        double A1[128];
        double A2[128];
        float f = 0;
		int i;
        for ( i = 0; i < 128; i++)
        {
                A1[i] =1.1;
                A2[i] = 0.15;
        }
        A1[23] = 0.215;
        A1[23] = 0.3415;
        int res=IMP_FRS_FeatSim(A1, A2, 1, &f);
        printf("%d,%f\n", res, f);
        printf("%f", cosine_similarity(A1, A2, 128));
        return ;

}

