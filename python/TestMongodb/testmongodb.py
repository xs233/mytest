#!/usr/bin/python
# -*- coding: utf-8 -*-

from pymongo import MongoClient

if __name__ == '__main__':
    client = MongoClient('mongodb://10.46.215.19:27017/')
    tb = client.test.feature
    print(tb.find_one())
