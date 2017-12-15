from pyspark import SparkConf, SparkContext

conf = SparkConf().setMaster("spark://iZ28ur81pw2Z:7077").setAppName("My App")
#conf = SparkConf().setMaster("local").setAppName("My App")
sc = SparkContext(conf = conf)


with open('/home/spark/python/text.txt',"rb") as f:
	poet_list = f.readlines()
lines = sc.parallelize(poet_list)
#lines = sc.textFile('/home/spark/python/text.txt')
#aalines = lines.filter(lambda x : "aa" in x).saveAsTextFile('/home/spark/python/res')
print(lines.count())
print(lines.take(10))
#lines.saveAsTextFile('/home/spark/python/res')

sc.stop()


'''
	file_context = file_object.read()
	lines = sc.parallelize(file_context)
	print(lines.first())	

finally:
	file_object.close()
	sc.stop()
'''
'''
#lines = sc.textFile("file:///home/spark/env/spark-2.1.0-bin-hadoop2.7/README1.md")
lines = sc.textFile("README.md")
print(lines.first())
sc.stop()
'''
