import struct
import socket
import sys
import json
import md5

MSG_HEADER_LEN = 32

MSG_CODE_HEARTBEAT = 1
MSG_CODE_BIND_DEVICE = 2
MSG_CODE_UNBIND_DEVICE = 3
MSG_CODE_DEVICE_AUTH = 4
MSG_CODE_MODIFY_PASSWORD = 5
MSG_CODE_SET_ALG_PARA = 6
MSG_CODE_QUERY_ALG_PARA = 7
MSG_CODE_DOWNLOAD_IMAGE = 8
MSG_CODE_GET_RECORD_FILES = 9
MSG_CODE_DEVICE_UPDATE = 10

DEVICE_ID = "bbda7b6ebc1a4280b161764649a67bcb"
ACCOUNT = "admin"
PASSWORD = "123456"
NEWPASS = "123456"

def recv_all(sk, length):
    data = ""
    remain_len = length
    
    while remain_len > 0:
        once_len = remain_len
        if once_len > 1024:
            once_len = 1024
        once_data = sk.recv(once_len)
        remain_len -= len(once_data)
        data += once_data
    return data

def request(cmd, text_data, bin_data):
    sk = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server = ('192.168.0.235', 9001)
    #server = ('127.0.0.1', 9001)
    sk.connect(server)
    print "connect success:", server
    
    FLAG = "EVIL"
    LENGTH = len(text_data) + len(bin_data) + MSG_HEADER_LEN
    CHECKSUM = 0
    VERSION = 0x0100
    COMMANDCODE = cmd
    ERRORCODE = 0
    TEXTDATALENGTH = len(text_data)
    BINDATALENGTH = len(bin_data)
    TEXTDATA = text_data
    BINDATA = bin_data
    
    header = struct.pack(">4siiiiiii", FLAG, LENGTH, CHECKSUM,
                        VERSION, COMMANDCODE, ERRORCODE, TEXTDATALENGTH, BINDATALENGTH)
    sk.send(header)
    sk.send(TEXTDATA)
    sk.send(BINDATA)
    
    recv_head = sk.recv(MSG_HEADER_LEN)
    headers = struct.unpack(">4siiiiiii", recv_head)
    print "recv head:"
    print headers
    err_code = headers[5]
    if err_code != 0:
        print "failed! error code:", err_code
    
    txt_len = headers[6]
    print "txt_len:", txt_len
    if(txt_len > 0):
        txt = sk.recv(txt_len)
        print "recv txt:"
        print txt
    
    bin_len = headers[7]
    if(bin_len > 0):
        bin = recv_all(sk, bin_len)
        print "recv bin len:", len(bin)
    
        fp = open('1.jpg', 'wb')
        fp.write(bin)
        fp.close( )
    
    sk.close()
 
def test_search():
    m1 = md5.new()   
    m1.update("defender")   
    sendData = m1.hexdigest()
    print "search data:", sendData
    sendData = "975ECB719692FA2BC7255B0C2DD2F3A4"
    
    sk = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sk.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)
    sk.bind(("192.168.0.33", 0))
    sk.sendto(sendData, ("255.255.255.255", 9000))
    print "response:", sk.recv(32)
 

def test_bind():
   text = {"devID":DEVICE_ID, "devAcc":ACCOUNT,"devPass":PASSWORD,"userID":""}
   request(MSG_CODE_BIND_DEVICE, json.dumps(text), "")
   
def test_unbind():
   text = {"devID":DEVICE_ID, "devAcc":ACCOUNT,"userID":"id"}
   request(MSG_CODE_UNBIND_DEVICE, json.dumps(text), "")
   
def test_auth():
   text = {"devID":DEVICE_ID, "devAcc":ACCOUNT,"devPass":PASSWORD}
   request(MSG_CODE_DEVICE_AUTH, json.dumps(text), "")
   
def test_modify_pass():
   text = {"devID":DEVICE_ID, "devAcc":ACCOUNT,"devOldPass":PASSWORD,"devNewPass":NEWPASS}
   request(MSG_CODE_MODIFY_PASSWORD, json.dumps(text), "")

def test_alg_setpara():
   text = {"algID":"opa", "algPara": {"detectArea": "0.368000,0.320000,0.616000,0.784593"}}
   request(MSG_CODE_SET_ALG_PARA, json.dumps(text), "")

def test_alg_querypara():
   text = {"algID":"opa"}
   request(MSG_CODE_QUERY_ALG_PARA, json.dumps(text), "")
   
def test_download_image():
   text = {"imageURL": "2017-07-04-10:17:58.jpg"}
   request(MSG_CODE_DOWNLOAD_IMAGE, json.dumps(text), "")

def test_query_record_url():
   text = {"beginTime": 12233, "endTime":76664334}
   request(MSG_CODE_GET_RECORD_FILES, json.dumps(text), "")
   

if __name__ == "__main__":
    #test_search()
    
    #print "test_bind----------------------"
    #test_bind()
    
    #print "test_auth----------------------"
    #test_auth()
 	   
    #print "test_modify_pass----------------------"
    #test_modify_pass()
    
    #print "test_alg_querypara----------------------"
    #test_alg_querypara()
    
    #print "test_alg_setpara----------------------"
    #test_alg_setpara()
    
    print "test_download_image---------------------"
    test_download_image()
   
    #print "test_query_record_url-------------------"
    #test_query_record_url()
    
    #print "test_unbind----------------------"
    #test_unbind()
    
                    


