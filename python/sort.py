#!/usr/bin/python
# -*- coding: utf-8 -*-

def by_name(t):
	return t[1]

L = [('acb', 75), ('Adam', 92), ('Bart', 66), ('Lisa', 88)]
L2 = sorted(L, key=by_name,reverse=True)
print(L2)
