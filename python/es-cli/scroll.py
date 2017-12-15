#!/usr/bin/python
# -*- coding: utf-8 -*-

from elasticsearch import Elasticsearch
import json

def get_scroll_id():
    es = Elasticsearch()
    count = es.count(index="target", body={"query": {"match": {"device":"66"}}})

    if not count['count']:
        return 0

    scroll = es.search(index="target",scroll="3m",size=10,body={"sort": ["_doc"],"query": {"match": {"device":"66"}}})
        return scroll['hits']['hits']


'''
while(1):
    res = es.scroll(scroll_id=scroll['_scroll_id'],scroll="3m")
    if not res['hits']['hits']:
        print("data over")
        break
'''
'''
    with open('text.txt', 'ab+') as f:
        for item in res['hits']['hits']:
            json.dump(item, f)
            f.write('\n')
'''

