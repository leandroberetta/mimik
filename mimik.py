from flask import Flask, request

import requests
import os

app = Flask(__name__)

@app.route('/<path>')
def get(path): # pragma: no cover
    return mimik(path, 
                 request.headers, 
                 os.environ.get("MIMIK_TYPE"), 
                 os.environ.get("MIMIK_DESTINATION"))    

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

def mimik(path, headers, type, destination):
    past = headers.get("mimik-history")
    present = '{} ({})'.format(path, get_version())

    history = do_history(past, present)
    
    if type != "edge":        
        response = requests.get(destination, headers={"mimik-history": history})

        return response.content

    return history

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)  # pragma: no cover