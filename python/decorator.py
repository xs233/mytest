#!/usr/bin/python
# -*- coding: utf-8 -*-
import functools

def log(*text):
    def decorator(func):
        @functools.wraps(func)
        def wrapper(*args, **kw):
	    if text:
	    	print('%s %s():' % (text[0], func.__name__))
	    else:
		print('call %s():' % func.__name__)
            print('begin call')
            func(*args, **kw)
	    print('end call')
	    return
        return wrapper
    return decorator

@log()
def f():
    print('f()')

f()
