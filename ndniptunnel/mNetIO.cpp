//
// @Author: Wei Guohua
// @Description: 
// @Version: 1.0.0
// @Date: 2021-02-03 09:13
// @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
//
#include <iostream>
#include "mNetIO.h"
#include "common.h"

using namespace std;

//void GoOnData(char *buf, int size) ;
//void GoOnInterest(char *buf, int size) ;

MNetIO::MNetIO (){
}
MNetIO::~MNetIO (){

}

int MNetIO::sendInterest(char *name) {
	static uint64_t interestSeq = 0;
	string myName = string(name) + "/" + to_string(interestSeq++);
	Name interestName(myName) ;
//	interestName.appendSequenceNumber(interestSeq++);
	Interest interest(interestName);
	interest.setInterestLifetime(1_s);
//	interest.setMustBeFresh(true);
//	interest.setApplicationParameters((const uint8_t*)buf, size);
//	mKeyChain.sign(interest) ;
	mFace.expressInterest(interest,
			bind(&MNetIO::onData, this,  _1, _2),
			bind(&MNetIO::onNack, this, _1, _2),
			bind(&MNetIO::onTimeout, this, _1));
//	mFace.removeAllPendingInterests();
//	std::cout << "send interest : " << interestName.toUri() << std::endl;

	return 0 ;
}

int MNetIO::sendData(char *buf, int size, char *name) {
    Name dataName(name);
//    std::cout << "send data : " << dataName.toUri() << std::endl;
    // Create Data packet
    auto data = make_shared<Data>(dataName);
    data->setFreshnessPeriod(0_s);
    if (size > 0) {
        data->setContent(reinterpret_cast<const uint8_t *>(buf), size);
    }

    string dsha = "id:/localhost/identity/digest-sha256";
    ndn::security::SigningInfo si(dsha);
    // Sign Data packet with default identity
    mKeyChain.sign(*data, si);
    mFace.put(*data);
    return 0;
}

void MNetIO::onData(const Interest& interest, const Data& data){
//	std::cout << "onData : " << std::endl;
    char buf[9000];
    if(data.getContent().value_size() <= 0) return ;
    int payloadSz = data.getContent().value_size();
    memcpy(buf, data.getContent().value(), payloadSz);
    string tmp(buf, payloadSz) ;
    GoOnData(buf, payloadSz);
}

void MNetIO::onNack(const Interest& interest, const lp::Nack& nack){
//	std::cout << "onNack : " << nack.getReason()  << std::endl;
    GoOnNack();
}


void MNetIO::onTimeout(const Interest& interest){
//	std::cout << "onTimeout : " << interest  << std::endl;
    GoOnTimeout();
}

int MNetIO::start(char *prefix) {
	mFace.setInterestFilter(prefix,
			bind(&MNetIO::onInterest, this, _1, _2),
			RegisterPrefixSuccessCallback(),
			bind(&MNetIO::onRegisterFailed, this, _1, _2));
	mFace.processEvents(ndn::time::milliseconds::zero(), true) ;	
	return -1 ;
}

void MNetIO::onInterest(const InterestFilter& filter, const Interest& interest){
//	char buf[9000];
//	if(!interest.hasApplicationParameters() ||
//			interest.getApplicationParameters().value_size() <= 0) return ;
//	int payloadSz = interest.getApplicationParameters().value_size();
//	memcpy(buf, interest.getApplicationParameters().value(), payloadSz);
//	string tmp(buf, payloadSz) ;
	//onInterest(buf, payloadSz) ;
	string nameStr = interest.getName().toUri();
//	memcpy(buf, nameStr.c_str(), nameStr.size());
    char buf[9000];
    memcpy(buf, nameStr.c_str(), nameStr.size());
    buf[nameStr.size()] = '\0';
//    std::cout << "before oninterest" << std::endl;
	GoOnInterest(buf);
//    std::cout << "after oninterest" << std::endl;

//	// Create Data packet
//    auto data = make_shared<Data>(interest.getName());
//    data->setFreshnessPeriod(0_s);

//    string dsha = "id:/localhost/identity/digest-sha256";
//	ndn::security::SigningInfo si(dsha);
//    // Sign Data packet with default identity
//    mKeyChain.sign(*data, si);
//    // m_keyChain.sign(*data, signingByIdentity(<identityName>));
//    // m_keyChain.sign(*data, signingByKey(<keyName>));
//    // m_keyChain.sign(*data, signingByCertificate(<certName>));
//    // m_keyChain.sign(*data, signingWithSha256());
//    // Return Data packet to the requester
////    std::cout << "<< D: " << *data << std::endl;
//    mFace.put(*data);
}
void MNetIO::onRegisterFailed(const Name& prefix, const std::string& reason) {
	std::cerr << "ERROR: Failed to register prefix \""
		<< prefix << "\" in local hub's daemon (" << reason << ")"
		<< std::endl;
	mFace.shutdown();
}
