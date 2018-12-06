import re

import logging
from flask import Flask, request
from flask import Response
from flask import abort

app = Flask(__name__)

@app.route('/register', methods=['POST'])
def register():
    try:
        data = request.json
        if not 'key' in data:
            return Response("Error: key is not provided", 400)
        if not 'password' in data:
            return Response("Error: password is not provided", 400)
        if not re.match(r'^[a-zA-Z0-9+/]+[=]*$', data['key']):
            return Response("Error: key is not a valid base64", 400)
        if not re.match(r'^[a-zA-Z0-9+/]+[=]*$', data['password']):
            return Response("Error: password is not a valid base64", 400)
        if len(data['password']) > 2048:
            return Response("Error: password is too long", 400)
        if valid_password != data['password']:
            return Response("Error: invalid password", 400)

        with open("/home/codex/.ssh/authorized_keys", "a+") as w:
            w.write("command=\"/home/codex/entry\",no-x11-forwarding,no-pty ssh-rsa {}\n".format(data['key']))

        return Response("OK", 200)

    except Exception as e:
        logging.debug(e)
        abort(400)


if __name__ == '__main__':
    valid_password = open("/home/config").read()
    app.run(debug=False, port=1339)
