from flask import Flask, request

import requests
import os
import time

app = Flask(__name__)

def log_time(method):
    def duration(*args, **kw):
        ts = time.time()
        result = method(*args, **kw)
        te = time.time()
        
        print('%r - response_time: %2.2f ms' % (method.__name__, (te - ts) * 1000))
        
        return result
    
    return duration


@app.route('/<path>')
@log_time
def get(path): # pragma: no cover
    return mimik(path, 
                 request.headers, 
                 os.environ.get("MIMIK_TYPE"), 
                 os.environ.get("MIMIK_DESTINATION"),
                 os.environ.get("MIMIK_SIMULATE_ERROR"))    

def do_history(past, present):    
    if past is None:      
        return present
    
    return '{} -> {}'.format(past, present)

def get_version():
    try:
        labels = open('/tmp/etc/pod_labels')

        for label in labels:
            values = label.split('=')

            if values[0] == 'version':
                return values[1]

        raise Exception("No version label")
    except:
        return "v1"

def mimik(path, headers, type, destination, error):
    if error == 'true':
        return 'Error', 503

    past = headers.get("mimik-history")
    present = '{} ({})'.format(path, get_version())

    history = do_history(past, present)
    
    if type != "edge":        
        response = requests.get(destination, headers={"mimik-history": history, "Authorization": headers.get("Authorization")})

        return response.content

    return history

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)  # pragma: no cover