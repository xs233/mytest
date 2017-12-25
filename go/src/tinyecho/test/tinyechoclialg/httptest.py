# -*- coding: utf-8 -*-

import json
import os
import unittest
import md5
import hashlib
import urllib
import time
import struct

import requests

host_api = 'http://localhost:9000'

MSG_HEADER_LEN = 32
MSG_CODE_ALARM = 1

session = requests.session()

class PushMessageTest(unittest.TestCase):

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_push_message_opa_alarm(self):
        print "\n----------test_push_message_opa_alarm----------"

        url = host_api + '/publish'

        message = dict()
        message['code'] = 'MsgCodeOPAAlarm'
        message['content'] = 'hello world'
        message['publisher'] = 'silver'
        message['time'] = datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        message_str = json.dumps(message)

        channel = 'silver_opa_alarm'
        body = 'channel:' + channel + '\r\n' + 'message:' + message_str

        print 'url: \n', url
        print 'body: \n', body

        response = session.post(url, data=body)
        print(response.text)

    def test_push_message_alarm_bin(self):
        print "\n----------test_push_message_alarm_bin----------"

        channel = 'channel=imp_alarm_opa&channel=imp_alarm_pfr'
        url = host_api + '/publish?' + channel

        txt_data_json = dict()
        txt_data_json['cameraID'] = 'cameraID1'
        txt_data_json['algID'] = 'opa'
        txt_data_json['timestamp'] = (long)(time.time()*1000)
        txt_data = json.dumps(txt_data_json)

        bin_data = "123456"
        file = open('F:/Data/picture/720p.jpg', 'rb')
        try:
            bin_data = file.read()
        finally:
            file.close()

        message = self.package(MSG_CODE_ALARM, txt_data, bin_data)        
  
        print 'url: \n', url
        #print 'body: \n', body

        response = session.post(url, data=message)
        print(response.text)

    def package(self, cmd_code, txt_data, bin_data):
        FLAG = "EVIL"
        LENGTH = len(txt_data) + len(bin_data) + MSG_HEADER_LEN
        CHECKSUM = 0
        VERSION = 0x0100
        COMMANDCODE = cmd_code
        ERRORCODE = 0
        TEXTDATALENGTH = len(txt_data)
        BINDATALENGTH = len(bin_data)

        header = struct.pack(">4siiiiiii", FLAG, LENGTH, CHECKSUM, VERSION, COMMANDCODE, ERRORCODE, TEXTDATALENGTH, BINDATALENGTH)

        data = header + txt_data + bin_data
        return data

if __name__ == '__main__':
    unittest.main()
