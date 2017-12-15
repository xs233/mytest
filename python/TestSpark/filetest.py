#!/usr/bin/python
# -*- coding: utf-8 -*-

'spark python test example'

__author__ = 'xusheng'

from pyspark import SparkConf, SparkContext
import sys,ctypes,json,hashlib,time
from operator import add

if __name__ == "__main__":
    conf = SparkConf()
    sc = SparkContext(conf=conf)
    rdd = sc.textFile("person66.json")
    
    '''
    def list_data(x):
            print x
            x = [x['timestamp'],x['person']['feature']['face']]
            return x
    '''

    '''
    counts = rdd.flatMap(lambda x: x.split(' ')) \
                  .map(lambda x: (x, 1)) \
                  .reduceByKey(add)
    output = counts.collect()
    #for (word, count) in output:
    #    print("%s: %i" % (word, count))
    
    '''
    def feature_calc(x):
        index = 100000
        n = 0
        while index:
            n += index
            index -= 1
        return True
    
    res = rdd.filter(feature_calc)
    print '-------------------------NumPart = ',res.getNumPartitions()
    #res.persist()
    print "before count = ",time.time()
    print('aft_count = %d' %res.count())
    print "after count = ",time.time()
    res.top(100)
    print "after top = ",time.time()
    
    sc.stop()
