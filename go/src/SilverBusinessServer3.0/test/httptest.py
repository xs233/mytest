# -*- coding: utf-8 -*-

import json
import os
import unittest
import md5
import hashlib
import urllib
import time

import requests
#from Crypto.Cipher import AES

#host = 'http://localhost:8080'
#host = 'http://192.168.0.34:8080'
#host = 'http://115.28.143.67:8092'

#host_api = 'http://testapi.impowerinside.com'
#host_api = 'http://localhost:8080'
#host_data = 'http://data.impowerinside.com'
#host_data = 'http://localhost:8080'
#host_data = 'http://192.168.0.34:8080'
#host_openapi = 'http://openapi.impowerinside.com'
host_api = 'http://192.168.0.200:8081'

session = requests.session()
#测试账号接口的正确性
class AccountTest(unittest.TestCase):

    def setUp(self):
        pass

    def tearDown(self):
        pass

    #测试用户的登录
    def test_login_username(self):
        print "\n----------test_login_username|用户的登录----------"
        password = '96E79218965EB72C92A549DD5A330112'
        data = '{"userName":"admin","userPassword":"' + password + '"}'
        response = session.post(host_api + '/v1/api/sessions', data=data)
        print(response.text)

    #测试修改密码
    def test_revise_password(self): 
        print "\n----------test_revise_password|修改密码----------"
        data = '{"userName":"user_add","oldPassword":"96E79218965EB72C92A549DD5A330112","newPassword":"111111"}'
        #print "req cookie:",session.cookies
        response = session.put(host_api + '/v1/api/password/modify', data=data)
        #print "rep cookie:",session.cookies
        print(response.text)

    #测试重置密码
    def test_reset_password(self): 
        print "\n----------test_reset_password|重置密码----------"
        data = '{"userName":"user_add","oldPassword":"96E79218965EB72C92A549DD5A330112","newPassword":"111111"}'
        #print "req cookie:",session.cookies
        response = session.put(host_api + '/v1/api/password/reset', data=data)
        #print "rep cookie:",session.cookies
        print(response.text)
    
    #测试用户注销
    def test_logout_username(self):
        print "\n----------test_logout_username|用户注销----------"
        #print "req cookie:",session.cookies
        response = session.delete(host_api + '/v1/api/sessions/:sid')
        #print "rep cookie:",session.cookies
        print(response.text)

    #测试删除用户
    def test_del_user(self):
        print "\n----------test_del_user|删除用户----------"
        response = session.delete(host_api + '/v1/api/users/:uid')
        print(response.text)

    #测试查询用户个人信息
    def test_get_user_info(self):
        print "\n----------test_get_user_info|查询用户个人信息----------"
        response = session.get(host_api + '/v1/api/users/:uid')
        print(response.text)

    #查询所有用户信息
    def test_get_all_user_info(self):
        print "\n----------test_get_all_user_info|查询所有用户信息----------"
        response = session.get(host_api + '/v1/api/users')
        print(response.text)

    #测试添加用户
    def test_add_username(self):
        print "\n----------test_add_username|添加用户----------"
        password = '96E79218965EB72C92A549DD5A330112'
        data = '{"userName":"admin","userPassword":"' + password + '"}'
        response = session.post(host_api + '/v1/api/users', data=data)
        print(response.text)

#测试设备的接口正确性
class GroupTest(unittest.TestCase):

	#测试查询所有设备列表
    def test_query_all_device(self):
        print "\n----------test_query_all_device|查询所有设备列表----------"
        response = session.get(host_api + '/v1/api/devices')
        print(response.text)

    #测试修改设备名
    def test_change_device_name(self):
        print "\n----------test_change_device_name|修改设备名----------"
        response = session.get(host_api + '/v1/api/devices/:did')
        print(response.text)

    #测试用户查询设备组列表
    def test_query_all_group(self):
        print "\n----------test_query_all_group|用户查询设备组列表----------"
        response = session.get(host_api + '/v1/api/groups')
        print(response.text)

    #测试创建设备组列表
    def test_create_group(self):
	print "\n----------test_create_group|创建设备组列表----------"
	data = '{"groupName":"group1"}'
        response = session.post(host_api + '/v1/api/groups', data=data)
        print(response.text)

    #测试用户查询设备组详情
    def test_query_group(self):
	print "\n----------test_query_group|用户查询设备组详情----------"
        response = session.get(host_api + '/v1/api/groups/:gid')
        print(response.text)

    #测试修改设备组名
    def test_change_group_name(self):
	print "\n----------test_change_group_name|修改设备组名----------"
        response = session.get(host_api + '/v1/api/groups/:gid')
        print(response.text)

    #测试修改设备组名
    def test_modify_group(self):
	print "\n----------test_modify_group|修改设备组名----------"
	data = '{"groupName":"group_modify"}'
        response = session.put(host_api + '/v1/api/groups/:gid', data=data)
        print(response.text)

    #测试删除设备组
    def test_delete_group(self):
	print "\n----------test_delete_group|删除设备组----------"
        response = session.delete(host_api + '/v1/api/groups/:gid')
        print(response.text)

    #测试查询未分组设备
    def test_query_nongroup(self):
	print "\n----------test_query_nongroup|查询未分组设备----------"
        response = session.get(host_api + '/v1/api/nongroup/devices')
        print(response.text)

    #向设备组中添加设备
    def test_add_group_cameras(self):
	print "\n----------test_add_group_cameras|向设备组中添加设备----------"
	data = '{"deviceIDList":[2]}'        
	response = session.post(host_api + '/v1/api/groups/:gid/devices',data=data)
        print(response.text)

    #从设备组中删除设备列表
    def test_remove_group_cameras(self):
	print "\n----------test_remove_group_cameras|从设备组中删除设备列表----------"
	data = '{"deviceIDList":[1]}'        
	response = session.delete(host_api + '/api/groups/:gid/devices',data=data)
        print(response.text)

    #更新设备组设备列表
    def test_modify_group_cameras(self):
	print "\n----------test_modify_group_cameras|更新设备组设备列表----------"
	data = '{"addDeviceList":[1], "deleteDeviceList":[]}'        
	response = session.put(host_api + '/v1/api/groups/:gid/devices',data=data)
        print(response.text)

    #分页查询未分组的设备
    def test_query_nongroup_bypage(self):
	print "\n----------test_query_nongroup_bypage|分页查询未分组的设备----------"      
	paras = {
            'keyword':'',
            'offset':'0',
            'count':'10'
        }
        paras = urllib.urlencode(paras)
        response = session.get(host_api + '/v1/api/nongroup/page/devices'+ '?' + paras)
        print(response.text)

        
class FavoriteTest(unittest.TestCase):

    def test_seach_person(self):
	print "\n----------test_seach_person----------"      
	paras = {
            'offset':'0',
            'count':'10'
        }
        paras = urllib.urlencode(paras)
	data = '{"deviceUUIDList":"jjjj", "beginTime":123,"endTime":123,"searchRule":{"face":{"gender":"man","age":"10","hairStyle":"1","isSpectacled":"1"},"body":{"upperBodyColor":"man","lowerBodyColor":"10","fullBodyColor":"1","bodilyForm":"1","height":"1"}, "feature":{"type":1,"value":"10"}}}' 
        response = session.post(host_api + '/v1/api/search/person'+ '?' + paras,data=data)
        print(response.text)

    def test_save_seach_person_rule(self):
	print "\n----------test_save_seach_person_rule----------"      
	data = '{"searchRule":{"face":{"gender":"man","age":"10","hairStyle":"1","isSpectacled":"1"},"body":{"upperBodyColor":"man","lowerBodyColor":"10","fullBodyColor":"1","bodilyForm":"1","height":"1"}, "feature":{"type":1,"value":"10","imageURL":"10","selectedRect":"10"}}}' 
        response = session.post(host_api + '/v1/api/search/person/rules',data=data)
        print(response.text)

    def test_query_seach_person_rule(self):
	print "\n----------test_query_seach_person_rule----------"      
        response = session.get(host_api + '/v1/api/search/person/rules/15')
        print(response.text)

    def test_query_favorite_device(self):
	print "\n----------test_query_favorite_device----------"      
        response = session.get(host_api + '/v1/api/favorite/devices')
        print(response.text)

    def test_query_favorite_by_page(self):
	print "\n----------test_query_favorite_by_page----------"      
	paras = {
            'deviceID':'1,2',
            'beginTime':0,
	       'endTime':100,
            'offset':0,
	       'count':10
        }
        paras = urllib.urlencode(paras)
        response = session.get(host_api + '/v1/api/favorites'+ '?' + paras)
        print(response.text)

    def test_add_favorite(self):
	print "\n----------test_add_favorite----------"      
	data = '{"deviceID":3, "imageURL":"123","imageTime":123,"searchRuleID":10}' 
        response = session.post(host_api + '/v1/api/favorites',data=data)
        print(response.text)
    def test_cancel_favorite(self):
	print "\n----------test_cancel_favorite----------"      
        response = session.delete(host_api + '/v1/api/favorites/7')
        print(response.text)
if __name__ == '__main__':
    unittest.main()
