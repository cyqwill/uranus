#!/usr/bin/env python
"""
simply websocket client
talks to Uranus, this client
send a msg to Uranus, the Uranus will
send that msg to the target ID client.

Usage:

python3 alice.py ws://localhost:9000/v1/ws


"""
import time
import numpy as np
import argparse
import code
import sys
import threading
import ssl
import six
import websocket
import json
import requests
import pickle
import os


OPCODE_DATA = (websocket.ABNF.OPCODE_TEXT, websocket.ABNF.OPCODE_BINARY)

# set target address
target_addr = 'usrHtNZ3klmi2'
send_addr = ''


class VAction(argparse.Action):
    def __call__(self, parser, args, values, option_string=None):
        if values is None:
            values = "1"
        try:
            values = int(values)
        except ValueError:
            values = values.count("v") + 1
        setattr(args, self.dest, values)


def parse_args():
    parser = argparse.ArgumentParser(description="Uranus Python Client for Testing.")
    parser.add_argument('-u', "--url", metavar="ws_url", default="ws://localhost:9000/v1/ws",
                        help="websocket url. ex. ws://echo.websocket.org/")
    parser.add_argument("-v", "--verbose", default=0, nargs='?', action=VAction,
                        dest="verbose",
                        help="set verbose mode. If set to 1, show opcode. "
                             "If set to 2, enable to trace  websocket module")
    parser.add_argument("-n", "--nocert", action='store_true',
                        help="Ignore invalid SSL cert")
    parser.add_argument("-s", "--subprotocols", nargs='*',
                        help="Set subprotocols")
    parser.add_argument("-o", "--origin",
                        help="Set origin")
    parser.add_argument("--eof-wait", default=0, type=int,
                        help="wait time(second) after 'EOF' received.")
    parser.add_argument("--timings", action="store_true",
                        help="Print timings in seconds")
    parser.add_argument("--headers",
                        help="Set custom headers. Use ',' as separator")

    return parser.parse_args()


class InteractiveConsole(code.InteractiveConsole):
    def write(self, data):
        sys.stdout.write("\033[2K\033[E")
        # sys.stdout.write("\n")
        sys.stdout.write("\033[34m< " + data + "\033[39m")
        sys.stdout.write("\n> ")
        sys.stdout.flush()

    def read(self):
        content = input('> ')
        return send_msg(content, target_addr)


def main(token, addr):
    start_time = time.time()
    args = parse_args()
    if args.verbose > 1:
        websocket.enableTrace(True)

    ws = websocket.create_connection(args.url)
    console = InteractiveConsole()
    print("Press Ctrl+C to quit")

    def recv():
        try:
            frame = ws.recv_frame()
        except websocket.WebSocketException:
            return websocket.ABNF.OPCODE_CLOSE, None
        if not frame:
            raise websocket.WebSocketException("Not a valid frame %s" % frame)
        elif frame.opcode in OPCODE_DATA:
            return frame.opcode, frame.data
        elif frame.opcode == websocket.ABNF.OPCODE_CLOSE:
            ws.send_close()
            return frame.opcode, None
        elif frame.opcode == websocket.ABNF.OPCODE_PING:
            ws.pong(frame.data)
            return frame.opcode, frame.data
        return frame.opcode, frame.data

    def recv_ws():
        while True:
            opcode, data = recv()
            msg = None
            if six.PY3 and opcode == websocket.ABNF.OPCODE_TEXT and isinstance(data, bytes):
                data = str(data, "utf-8")
            if not args.verbose and opcode in OPCODE_DATA:
                msg = data
            elif args.verbose:
                msg = "%s: %s" % (websocket.ABNF.OPCODE_MAP.get(opcode), data)
            msg_dict = json.loads(msg)
            content = msg_dict['content']
            sender = msg_dict["sender"]
            if msg is not None:
                if args.timings:
                    console.write(str(time.time() - start_time) + ": " + content)
                else:
                    console.write(content)
                # send back a msg
                ws.send(send_msg(content + " 你说的是这个吧。我是个回音机器人", sender))
            if opcode == websocket.ABNF.OPCODE_CLOSE:
                break

    thread = threading.Thread(target=recv_ws)
    thread.daemon = True
    thread.start()

    # ws.send(hi_msg())
    # ws.send(bytes("{'hi': 'ffff'}", encoding='utf-8'))
    ws.send(hi_msg(token, addr))

    while True:
        try:
            message = console.read()
            # TODO: wrap this message to Uranus received.
            ws.send(message)
        except KeyboardInterrupt:
            print('Bye.')
            return
        except EOFError:
            time.sleep(args.eof_wait)
            return


def hi_msg(token, addr):
    # change this to one account token and user_addr
    msg = {
        "token": token,
        "user_addr": addr,
        "ua": "py/macos",
        "device": "mac",
        "location": "湖南长沙"
    }
    out_msg = {
        "msg_type": "hi",
        "payload": msg
    }
    msg_str = json.dumps(out_msg)
    b = bytes(msg_str, 'utf-8')
    return b


def send_msg(content, target):
    msg = {
        "target": target,
        "sender": send_addr,
        "content": content,
        "msg_type": 0,
    }
    out_msg = {
        "msg_type": "send",
        "payload": msg
    }
    msg_str = json.dumps(out_msg)
    return bytes(msg_str, encoding='utf-8')


def login():
    login_url = "http://localhost:9000/api/v1/users_login"
    acc = input("[login] input your user account (like lucasjin): ")
    pwd = input("[login] input your account password: ")
    data = {"user_acc": acc, "user_password": pwd}
    rp = requests.post(login_url, data=data)
    print(rp)
    if rp.ok:
        rp = rp.json()
        print(rp)
        if rp['status'] == 'success':
            return rp["data"]["token"], rp["data"]["user_addr"]
        else:
            print('login failed.')
            exit()
    else:
        print('server not response.')
        exit()


if __name__ == "__main__":
    try:
        f_name = 'bob.pkl'
        if os.path.exists(f_name):
            with open(f_name, 'rb') as f:
                a = pickle.load(f)
            token = a['token']
            addr = a['user_addr']
        else:
            token, addr = login()
            with open(f_name, 'wb') as f:
                pickle.dump({"token": token, "user_addr": addr}, f)
        print('my address: ', addr)
        send_addr = addr
        main(token, addr)
    except Exception as e:
        print(e)
