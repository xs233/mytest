#!/usr/bin/python
# -*- coding: utf-8 -*-

'Data process function'

__author__ = 'xusheng'

import ctypes

def calc(x):
	lib_handle = ctypes.CDLL('./libfeatureProcess.so')
	func = lib_handle.feature_process
	res = func('aaa',3,'bbb',3)
	print(res)
	if res > 50:
		return True
	else:
		return False

'''
def change_type(x):
       	#x[1] = {'device',x[1]['device'] * 2}
	return (x[0],{'value' : x[1]['device'] * 2})
'''
