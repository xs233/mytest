#include <iostream>
#define CRYPTOPP_ENABLE_NAMESPACE_WEAK 1
#include <cryptopp/md5.h>
#include <cryptopp/hex.h>

int main()
{
	byte digest[ CryptoPP::Weak::MD5::DIGESTSIZE ];
	std::string message = "abcdefghijklmnopqrstuvwxyz";

	CryptoPP::Weak::MD5 hash;
	hash.CalculateDigest( digest, (const byte*)message.c_str(), message.length() );

	CryptoPP::HexEncoder encoder;
	std::string output;

	encoder.Attach( new CryptoPP::StringSink( output ) );
	encoder.Put( digest, sizeof(digest) );
	encoder.MessageEnd();

	std::cout << output << std::endl;
	return 0;
}
