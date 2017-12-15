from pyspark import SparkConf, SparkContext

conf = SparkConf().setMaster("spark://iZ28ur81pw2Z:7077").setAppName("My App")
#conf = SparkConf().setMaster("local").setAppName("My App")
sc = SparkContext(conf = conf)

conf = {"es.net.http.auth.user":"elastic",
"es.net.http.auth.pass":"changeme",
"es.resource" : "target/vericle",
"es.nodes":"139.129.218.165",
"es.query":'{"query":{"range":{"device":{"gt":1659916,"lt":1660000}}}}'}

def f(x):
	x[1]['type'] = 6
	return x

rdd = sc.newAPIHadoopRDD("org.elasticsearch.hadoop.mr.EsInputFormat",
                             "org.apache.hadoop.io.NullWritable",
                             "org.elasticsearch.hadoop.mr.LinkedMapWritable",
                             conf=conf)

#print(rdd.count())
print(rdd.map(f).first())
#rdd.saveAsTextFile("/home/spark/python/res")
sc.stop()
