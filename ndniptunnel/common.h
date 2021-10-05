//
// @Author: Wei Guohua
// @Description: 
// @Version: 1.0.0
// @Date: 2021-02-03 09:13
// @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
//
#ifndef _COMMON_H_
#define _COMMON_H_

#include <stdlib.h>

#ifdef __cplusplus
extern "C" {
#endif

//void OnData(char*);
void GoOnData(char *buf, int size) ;
void GoOnInterest(char *name) ;
void GoOnTimeout();
void GoOnNack();

#ifdef __cplusplus
}
#endif

#endif 
