#!/usr/bin/env python3

from pathlib import Path
import os
from os.path import join
import subprocess as sp
import shlex
import yaml
import json
import time

print('I hope you started me from the root of the project')
print(f'cwd is {os.getcwd()}')
debug = Path('.debug')
debug.mkdir(exist_ok=True)

ngrok_conf = Path(join('dev', 'ngrok.yaml'))
godrunk_conf = Path('godrunk.yaml')

with open(ngrok_conf) as f:
    ng = yaml.safe_load(f)

if ng['log_format'] != 'json':
    print('ngrok log is not in json format')
    exit(1)

ng_log = ng['log']
os.truncate(ng_log, 0)

# start ngrok in background
cmd = shlex.split(f'ngrok start --config {ngrok_conf} godrunk')
ngrok = sp.Popen(cmd)

with open(join(debug, 'ngrok.pid'), 'w') as f:
    f.write(str(ngrok.pid))

print('trying to read logs of ngrok')
def parse_log_record(line):
    j = json.loads(line)

    found, address = False, None
    if 'obj' in j and j['obj'] == 'tunnels' and 'url' in j and j['url'].startswith('https'):
        address = j['url']
        found = True

    return found, address

with open(ng_log) as f:
    i = 0
    while True:
        line = f.readline()
        if not line:
            time.sleep(0.1)
        else:
            found, address = parse_log_record(line)
            if found:
                break
        i += 1
        if i == 1000:
            print('I have been waiting for so long...')
            exit(1)

with open(godrunk_conf) as f:
    y = yaml.safe_load(f)

y['port'] = ng['tunnels']['godrunk']['addr']
y['webhookAddress'] = address

with open(godrunk_conf, 'w') as f:
    yaml.dump(y, f)

print('Address and port have been succefully updated in godrunk.yaml')

ngrok.wait()