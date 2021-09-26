//
// @Author: Wei Guohua
// @Description: 
// @Version: 1.0.0
// @Date: 2021-02-03 09:13
// @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
//
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include "bridge.h"
#include "mNetIO.h"

static MNetIO mNetIO ;

int sendInterest(char *name){
	return mNetIO.sendInterest(name) ;
}

int sendData(char *buf, int size, char *name) {
    return mNetIO.sendData(buf, size, name);
}

int start(char *buf, int size){
	char prefix[2000] ;
	memcpy(prefix, buf, size) ;
	prefix[size] = 0;
	return mNetIO.start(buf) ;
}
