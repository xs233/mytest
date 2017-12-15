#!/usr/bin/python
# -*- coding: utf-8 -*-

'write config to redis'

__author__ = 'xusheng'

import redis
r = redis.StrictRedis(host='localhost', port=6379, db=0)
r.hmset('search_engine_config',{'spark_submit_path':'/home/spark/env/spark-2.1.1-bin-hadoop2.7/bin/spark-submit', \
				'spark_master':'spark://iZ28ur81pw2Z:7077', \
				'spark_jars':'/home/spark/env/spark-2.1.1-bin-hadoop2.7/jars/elasticsearch-spark-20_2.11-5.4.3.jar', \
				'spark_files_path':'/home/spark/python/Compass1.0/libfeatureprocess.so', \
				'spark_name':'MyApp', \
				'spark_interface_path':'/home/spark/python/Compass1.0/sparkinterface.py', \
				'es_user':'elastic', \
				'es_pass':'changeme', \
				'es_resource':'target/person', \
				'es_nodes':'10.46.215.19', \
				'threshold':0.7}) 
