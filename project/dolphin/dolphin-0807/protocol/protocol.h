#ifndef __PROTOCOL_H__
#define __PROTOCOL_H__

// driver ---> executor
struct AssignWorkHeader {
    char flag[8];
    int cores;
    int info_len;
    char reserved[16];
    char md5[32];
};

#endif
