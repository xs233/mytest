from elasticsearch import Elasticsearch
import json,time

bt = time.time()
print(bt)
es = Elasticsearch()
count = es.count(index="target", body={"query": {"match": {"device":"66"}}})
print(count)

if not count['count']:
    sys.exit(1)

scroll = es.search(index="target",scroll="3m",size=1000,body={"sort": ["_doc"],"query": {"match": {"device":"66"}}})
'''
with open('text.txt','ab+') as f:
    for item in scroll['hits']['hits']:
        json.dump(item, f)
        f.write('\n')
'''

while(1):
    res = es.scroll(scroll_id=scroll['_scroll_id'],scroll="3m")
    if not res['hits']['hits']:
        print("data over")
        break
'''
    with open('text.txt', 'ab+') as f:
        for item in res['hits']['hits']:
            json.dump(item, f)
            f.write('\n')
'''
print("spent time : %f" % (time.time()-bt))

