#!/usr/bin/python
# -*- coding: utf-8 -*-
def fib(nMax):
	n,a,b = 0,0,1
	while n < nMax:
		print('l-%d' %(b))
		yield b
		print('r-%d' %(b))
		a,b = b,a+b
		n=n+1
	return

g = fib(int(raw_input('please enter the fib num:')))
#for s in g:
#	print(s)
