#!/usr/bin/python
# -*- coding: utf-8 -*-

'my first class'

__author__ = 'xusheng'

class Student(object):

    def __init__(self, name, score):
        self.name = name
        self.score = score

    def print_score(self):
	self.age = 19
        print('%s: %s: %s' % (self.name, self.score,self.age))

class Screen(object):
	@property
	def width(self):
		return self._width
	@width.setter
	def width(self,value):
		self._width = value

	@property
	def height(self):
		return self._height
	@height.setter
	def height(self,value):
		self._height = value

	@property
	def resolution(self):
		return self._width * self._height

if __name__ == '__main__':
	'''
	aa = Student('Bob',66)
	aa.age = 18
	aa.print_score()
	'''

	bb = Screen()
	bb.width = 1024
	bb.height = 768
	print(bb.resolution)
