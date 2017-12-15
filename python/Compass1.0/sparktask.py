#!/usr/bin/python
# -*- coding: utf-8 -*-

'add spark task'

__author__ = 'xusheng'

import subprocess,json,redis

class SparkTask(object):
	__spark_config = {'spark_submit_path':'','spark_master':'','spark_jars':'','spark_files_path':'','spark_name':'','spark_interface_path':''}
	__es_config = {'es_user':'','es_pass':'','es_resource':'','es_nodes':'','threshold':''}
	def spark_init(self):
		r = redis.StrictRedis(host='localhost', port=6379, db=0)
		for i in self.__spark_config:
			value = r.hget('search_engine_config',i)
			if value:
				self.__spark_config[i] = value
			else:
			 	return False
		for i in self.__es_config:
			value = r.hget('search_engine_config',i)
			if value:
				self.__es_config[i] = value
			else:
			 	return False
		return True
	def search_person(self,query):
		print query
		return subprocess.call("%s --master %s --jars %s --files %s --name %s  --driver-cores 2 \
					%s '%s' '%s'" % (self.__spark_config['spark_submit_path'],
							 self.__spark_config['spark_master'],
							 self.__spark_config['spark_jars'],
							 self.__spark_config['spark_files_path'],
							 self.__spark_config['spark_name'],
							 self.__spark_config['spark_interface_path'],
							 json.dumps(self.__es_config),query),shell=True)
