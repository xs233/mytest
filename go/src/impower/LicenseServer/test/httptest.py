# -*- coding: utf-8 -*-

import json
import os
import unittest
import md5
import hashlib
import urllib
import time
import base64

import requests
#from Crypto.Cipher import AES

host_api = 'http://127.0.0.1:8080'

session = requests.session()

class AccountTest(unittest.TestCase):

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_register_username(self):
        print "\n----------test_register_username----------"
        password = '96E79218965EB72C92A549DD5A330112'
        data = '{"userName":"user_add","userPassword":"' + password + '"}'
        response = session.post(host_api + '/v1/api/users', data=data)
        print(response.text)
    def test_reset_password(self):
        print "\n----------test_reset_password----------"
        data = '{"userName":"user_add","oldPassword":"96E79218965EB72C92A549DD5A330112","newPassword":"111111"}'
        #print "req cookie:",session.cookies
        response = session.put(host_api + '/v1/api/password', data=data)
        #print "rep cookie:",session.cookies
        print(response.text)
    def test_login_username(self):
        print "\n----------test_login_username----------"
        password = '96E79218965EB72C92A549DD5A330112'
        data = '{"userName":"admin","userPassword":"' + password + '"}'
        response = session.post(host_api + '/v1/api/sessions', data=data)
        print(response.text)
    def test_logout_username(self):
        print "\n----------test_logout_username----------"
        #print "req cookie:",session.cookies
        response = session.delete(host_api + '/v1/api/sessions/1')
        #print "rep cookie:",session.cookies
        print(response.text)
    def test_del_user(self):
        print "\n----------test_del_user----------"
        response = session.delete(host_api + '/v1/api/users/100')
        print(response.text)
    def test_get_user_info(self):
        print "\n----------test_get_user_info----------"
        response = session.get(host_api + '/v1/api/users/100')
        print(response.text)
    def test_get_all_user_info(self):
        print "\n----------test_get_all_user_info----------"
        response = session.get(host_api + '/v1/api/users')
        print(response.text)
		
class LicenseTest(unittest.TestCase):
    def test_apply_license(self):
        print "\n----------test_apply_license----------"
        mac = "A1:B2:C3:D4:E5:00"
        mac_base64 = base64.b64encode(mac)
        response = session.get(host_api + '/v1/api/license?mac=' + mac_base64)
        print(response.text)

    def test_apply_p2p(self):
        print "\n----------test_apply_p2p----------"
        mac = "A1:B2:C3:D4:E5:00"
        mac_base64 = base64.b64encode(mac)
        response = session.get(host_api + '/v1/api/p2p?mac=' + mac_base64)
        print(response.text)


if __name__ == '__main__':
    unittest.main()
