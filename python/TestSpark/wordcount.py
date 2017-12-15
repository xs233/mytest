#!/usr/bin/python
# -*- coding: utf-8 -*-

'word count'

__author__ = 'xusheng'

from pyspark import SparkConf, SparkContext

conf = SparkConf().setMaster("spark://iZ28ur81pw2Z:7077").setAppName("My App")
#conf = SparkConf().setMaster("local").setAppName("My App")
sc = SparkContext(conf = conf)

with open('/home/spark/env/spark-2.1.0-bin-hadoop2.7/README.md',"rb") as f:
	poet_list = f.readlines()
wordrdd = sc.parallelize(poet_list)
wordrdd1 = wordrdd.flatMap(lambda x:x.split()).map(lambda x:(x,1)).reduceByKey(lambda a,b:a+b).sortBy(lambda x:x[1],False)
print(wordrdd1.take(10))
sc.stop()
