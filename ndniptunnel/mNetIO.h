//
// @Author: Wei Guohua
// @Description: 
// @Version: 1.0.0
// @Date: 2021-02-03 09:13
// @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
//
#ifndef _MNETIO_H_
#define _MNETIO_H_
#include <ndn-cxx/face.hpp>
#include <iostream>
#include <deque>
using namespace std ;
using namespace ndn ;

class MNetIO
{
public:
	MNetIO ();
	~MNetIO ();
	int sendInterest(char* name) ;
    int sendData(char *buf, int size, char* name) ;
	int start(char *prefix) ;

private:
	void onInterest(const InterestFilter& filter, const Interest& interest);
	void onRegisterFailed(const Name& prefix, const std::string& reason);

	void onData(const Interest& interest, const Data& data);
	void onNack(const Interest& interest, const lp::Nack& nack);
	void onTimeout(const Interest& interest);
private:
	Face mFace;
	KeyChain mKeyChain;
	/* data */
};

#endif 
