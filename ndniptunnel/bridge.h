//
// @Author: Wei Guohua
// @Description: 
// @Version: 1.0.0
// @Date: 2021-02-03 09:13
// @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
//
#ifndef _BRIDGE_H_
#define _BRIDGE_H_

#ifdef __cplusplus
extern "C" {
#endif

int sendInterest(char *name) ;
int sendData(char *buf, int size, char *name) ;
int start(char *buf, int size) ;


#ifdef __cplusplus
}
#endif

#endif 
