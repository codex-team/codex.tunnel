#!/usr/bin/python3

import re
import time
import os

config = """server {{
    listen 80;
    server_name {};

    proxy_set_header Host $http_host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $remote_addr;
    proxy_set_header X-Forwarded-Proto $scheme;

    location / {{
        proxy_pass http://127.0.0.1:{}/;
    }}
}}
"""

host, port = input(), input()
if not port.isnumeric():
    print("Port should be numeric")
    exit(0)

port = int(port)
if port < 1024 or port > 65535:
    print("Port is invalid")
    exit(0)

if len(host) < 1 or len(host) > 24:
    print("Host length is invalid")
    exit(0)

if not re.match(r'^[A-Za-z][A-Za-z0-9]*$', host):
    print("Host should be alphanumeric and start with a letter")
    exit(0)

host = host + ".tun.ifmo.su"

with open("/etc/nginx/sites-enabled/{}.conf".format(host), "w") as w:
    w.write(config.format(host, port))

os.system("service nginx restart")
while True:
    time.sleep(5)
