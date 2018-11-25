#!/bin/python3

import requests
import json
import base64
import time
import os
import sys
import gui
import gi

gi.require_version('Gtk', '3.0')
from gi.repository import Gtk

from requests.packages.urllib3.exceptions import InsecureRequestWarning
requests.packages.urllib3.disable_warnings(InsecureRequestWarning)

host = 'https://localhost:9090'

def login(user, passwd):
    creds = json.dumps({
        'Username': user,
        'Password': passwd
    })

    res = requests.get(host + '/login', data=creds, verify=False)

    if res.status_code == 200:
        return res.text

    return False

def lookup_user(user):
    params = { 'user': user }
    res = requests.get(host + '/lookup', params=params, verify=False)

    if res.status_code == 200:
        print('')
        return res.text

    return 0

def unwrap_jwt(jwt):
    token = jwt.split('.')
    return json.loads(base64.urlsafe_b64decode(token[1]))

def send_message(msg, to, jwt):
    params = { 'to': to }
    headers = { 'Session': jwt }

    res = requests.post(host + '/send', data=msg, params=params, headers=headers, verify=False)

def get_messages(peer, since, jwt):
    params = { 'peer': peer }
    headers = { 'Session': jwt }

    if since != None:
        params['since'] = since

    res = requests.get(host + '/get', params=params, headers=headers, verify=False)

    return json.loads(res.text)

# def poll_messages(peer_name, jwt, stdscrn):
#         timestamp = None
#         while True:
#             messages = get_messages(peer_name, timestamp, jwt)
#
#             if timestamp:
#                 stdscrn.addstr(0, 0, timestamp)
#                 stdscrn.refresh()
#
#             if messages:
#                 h, w = stdscrn.getmaxyx()
#
#                 for i in reversed(range(len(messages))):
#                     # stdscrn.move(0, 0)
#                     # stdscrn.deleteln()
#
#                     message = messages[i]
#                     author = message['Username' ]
#                     timestamp = message['Timestamp']
#                     msg = message['Message']
#                     display = '[%s] %s: %s' % (timestamp, author, msg)
#
#                     # stdscrn.addstr(h - 1 - i, 0, display)
#
#                 stdscrn.refresh()
#
#             time.sleep(2)
#
# def main(stdscrn):
#     curses.echo()
#
#     jwt = False
#     while not jwt:
#         stdscrn.addstr(0, 0, 'Please Login')
#         stdscrn.addstr(1, 0, 'Username: ')
#         stdscrn.addstr(2, 0, 'Password: ')
#         stdscrn.refresh()
#
#         user = stdscrn.getstr(1, len('Username: ')).decode('utf-8')
#         passwd = stdscrn.getstr(2, len('Password: ')).decode('utf-8')
#
#         jwt = login(user, passwd)
#
#         stdscrn.erase()
#         stdscrn.refresh()
#
#     session = unwrap_jwt(jwt)
#     my_id = session['Uid']
#
#     peer_id = False
#     peer_name = False
#     while not peer_id:
#         stdscrn.addstr(0, 0, 'Who would you like to chat with? ')
#
#         peer_name = stdscrn.getstr(0, len('Who would you like to chat with? ')).decode('utf-8')
#
#         peer_id = lookup_user(peer_name)
#
#         stdscrn.erase()
#         stdscrn.refresh()
#
#     win = stdscrn.subwin(curses.LINES - 1, curses.COLS, 0, 0)
#     msgfetch = threading.Thread(
#         group = None,
#         target = poll_messages,
#         name = 'msgfetch',
#         args = (peer_name, jwt, win))
#     msgfetch.start()
#
#     while True:
#         prompt = 'Chat with %s: ' % peer_name
#         stdscrn.addstr(curses.LINES - 1, 0, prompt)
#         stdscrn.refresh()
#
#         msg = stdscrn.getstr(curses.LINES - 1, len(prompt)).decode('utf-8')
#         stdscrn.clrtoeol()
#         stdscrn.refresh()
#
#         send_message(msg, peer_name, jwt)

if __name__ == '__main__':
    jwt = False
    session = False
    my_id = False
    my_name = False
    peer_id = False
    peer_name = False

    win = gui.ChatWindow()
    win.connect('destroy', Gtk.main_quit)

    login_dialog = gui.LoginDialog(win)
    status = login_dialog.run()
    if status == Gtk.ResponseType.OK:
        user = login_dialog.user_entry.get_text()
        passwd = login_dialog.passwd_entry.get_text()

        jwt = login(user, passwd)
        session = unwrap_jwt(jwt)
        my_id = session['Uid']
        my_name = user

        login_dialog.destroy()
    else:
        sys.exit(0)

    peer_dialog = gui.PeerDialog(win)
    status = peer_dialog.run()
    if status == Gtk.ResponseType.OK:
        user = peer_dialog.user_entry.get_text()

        peer_id = lookup_user(user)
        my_name = user

        peer_dialog.destroy()
    else:
        sys.exit(0)

    Gtk.main()
