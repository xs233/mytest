#!/usr/bin/python
# -*- coding: utf-8 -*-

'spark python interface'

__author__ = 'xusheng'

from pyspark import SparkConf, SparkContext
import sys,ctypes,json,redis,hashlib,time

class JsonParse(object):
	def parse(self,dict_json):
		if not (dict_json.get('searchRule') and dict_json.get('searchRule').get('feature') and dict_json.get('searchRule').get('feature').get('value')):
			return 0

		query_list = []
		if dict_json.get('beginTime') and dict_json.get('endTime'):
			time_dict = {"range":{"timestamp":{"gte":dict_json['beginTime'],"lte":dict_json['endTime']}}}
			query_list.append(time_dict)

		if dict_json.get('deviceUUIDList'):
			device_dict = {"terms":{"device":dict_json['deviceUUIDList'].split(',')}}
			query_list.append(device_dict)
		
		search_face = dict_json.get('searchRule').get('face')
		if search_face:
			if search_face.get('gender'):
				gender_dict = {"terms":{"person.face.gender":search_face['gender']}}
				query_list.append(gender_dict)
			if search_face.get('age'):
				age_dict = {"terms":{"person.face.age":search_face['age']}}
				query_list.append(age_dict)
			if search_face.get('hairStyle'):
				hair_style_dict = {"terms":{"person.face.hair_style":search_face['hairStyle']}}
				query_list.append(hair_style_dict)
			if search_face.get('isSpectacled'):
				is_spectacled_dict = {"terms":{"person.face.is_spectacled":search_face['isSpectacled']}}
				query_list.append(is_spectacled_dict)
		
		search_body = dict_json.get('searchRule').get('body')
		if search_body:
			if search_body.get('upperBodyColor'):
				upper_color_dict = {"terms":{"person.body.upper_body_color":search_body['upperBodyColor']}}
				query_list.append(upper_color_dict)
			if search_body.get('lowerBodyColor'):
				lower_color_dict = {"terms":{"person.body.lower_body_color":search_body['lowerBodyColor']}}
				query_list.append(lower_color_dict)
			if search_body.get('fullBodyColor'):
				full_color_dict = {"terms":{"person.body.full_body_color":search_body['fullBodyColor']}}
				query_list.append(full_color_dict)
		
		self.feature = dict_json['searchRule']['feature']['value']
		es_query_dict = {"query":{"bool":{"must":query_list}}}
		return json.dumps(es_query_dict)

if __name__ == "__main__":
	print "beginTime = ",time.time()
	args = sys.argv
	if len(args) != 3:
		sys.exit('args error!')

	conf = SparkConf()
	sc = SparkContext(conf = conf)

	try:
		es_config = json.loads(args[1])
		jparse = JsonParse()
		es_query = jparse.parse(json.loads(args[2]))
		if not es_query:
			sc.stop()
			sys.exit(1)
		print es_query
		print "parseTime = ",time.time()
		conf = {"es.net.http.auth.user":es_config['es_user'],
				"es.net.http.auth.pass":es_config['es_pass'],
				"es.resource" : es_config['es_resource'],
				"es.nodes":es_config['es_nodes'],
				"es.query":'%s' % es_query}
		rdd = sc.newAPIHadoopRDD("org.elasticsearch.hadoop.mr.EsInputFormat",
					 "org.apache.hadoop.io.NullWritable",
					 "org.elasticsearch.hadoop.mr.LinkedMapWritable",
					  conf=conf)
		#print('count = %d' %rdd.count())
		print "rddTime = ",time.time()
		
		def list_data(x):
			x = [x[0],x[1]['person']['feature']['face']]
			return x
		
		jparse.feature = jparse.feature.replace(' ','')
		bro_des_feat = sc.broadcast(jparse.feature)
		threshold = es_config['threshold']
		def feature_calc(x):
			src_feat = x[1]['person']['feature']['face']
			lib_handle = ctypes.CDLL('./libfeatureprocess.so')
			func = lib_handle.feature_process
			func.argtypes=[ctypes.c_char_p,ctypes.c_char_p]
			func.restype = ctypes.c_float
			res = func(src_feat,bro_des_feat.value)
			if float(res) > float(threshold):
				x[1]['person'] = res
				return True
			else:
				return False

		#res = rdd.map(list_data).filter(feature_calc)
		res = rdd.filter(feature_calc)
		print "mapfilterTime = ",time.time()
		print res.getNumPartitions()
		data = res.collect()
		#print res.top(1)
		print "collectTime = ",time.time()
		'''
		res.persist()
		if  res.take(1):
			print "take1Time = ",time.time()
			#print('aft_count = %d' %res.count())
			#print "countTime = ",time.time()
			data = res.collect()
			print "collectTime = ",time.time()
			r = redis.StrictRedis(host='localhost', port=6379, db=0)
			m = hashlib.md5()
			margs = "/v1/search/person" + args[2]
			m.update(margs)
			md5v = m.hexdigest()
			for i in data:
				if not r.zadd(md5v,i[1],i[0]):
					sc.stop()
					sys.exit(2)
		else:
			print('rdd is empty!')
		'''
		
	except Exception, e:
		print e
		sc.stop()
		sys.exit(3)
	sc.stop()
