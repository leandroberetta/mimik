from flask import Flask, request

import requests
import os
import time

app = Flask(__name__)

@app.route('/<path>')
def get(path): # pragma: no cover
    return mimik(path, 
                 request.headers, 
                 os.environ.get('MIMIK_TYPE'), 
                 os.environ.get('MIMIK_DESTINATION'),
                 os.environ.get('MIMIK_SIMULATE_ERROR'))    

def chain_response(past, present):    
    if past is None:      
        return present
    
    return '{} -> {}'.format(past, present)

def get_version():
    version = 'v1'

    try:
        with open('/tmp/etc/pod_labels') as labels:
            for label in labels:
                values = label.split('=')

                if values[0] == 'version':
                    version = values[1]
    except: 
        return version                

    return version

def prepare_outbound_headers(inbound_headers):
    outbound_headers = {}

    outbound_headers['Authorization'] = inbound_headers.get('Authorization')

    outbound_headers['x-request-id'] = inbound_headers.get('x-request-id')
    outbound_headers['x-b3-traceid'] = inbound_headers.get('x-b3-traceid')      
    outbound_headers['x-b3-spanid'] = inbound_headers.get('x-b3-spanid')
    outbound_headers['x-b3-parentspanid'] = inbound_headers.get('x-b3-parentspanid')
    outbound_headers['x-b3-sampled'] = inbound_headers.get('x-b3-sampled')
    outbound_headers['x-b3-flags'] = inbound_headers.get('x-b3-flags')

    return outbound_headers

def mimik(path, headers, type, destination, error):
    if error == 'true':
        return 'Error', 503
    elif type == 'passthrough':
        present = requests.get(destination, headers=prepare_outbound_headers(headers)).text
        past = '{} ({})'.format(path, get_version())
        return chain_response(past, present)
    else:
        return '{} ({})'.format(path, get_version())

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)  # pragma: no cover