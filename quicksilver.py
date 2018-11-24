#!/bin/python3

import requests
import json
import base64
import getpass
import time
import os

from requests.packages.urllib3.exceptions import InsecureRequestWarning
requests.packages.urllib3.disable_warnings(InsecureRequestWarning)

host = 'https://localhost:9090'

def login():
    user = input('Username: ')
    passwd = getpass.getpass('Password: ')

    creds = json.dumps({
        'Username': user,
        'Password': passwd
    })

    res = requests.get(host + '/login', data=creds, verify=False)

    if res.status_code == 200:
        return res.text

    print('\nLogin failed:', res.text)
    print('Please try again.', end='\n\n')

    return login()

def lookup_user():
    user = input('Who would you like to talk to? ')

    params = { 'user': user }
    res = requests.get(host + '/lookup', params=params, verify=False)

    if res.status_code == 200:
        print('')
        return res.text, user

    print('User %s does not exist.' % user)

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

def poll_messages(peer, jwt):
    timestamp = None
    while True:
        messages = get_messages(peer_name, timestamp, jwt)
        for message in messages:
            author = message['Username' ]
            timestamp = message['Timestamp']
            msg = message['Message']
            print('\n[%s] %s: %s' % (timestamp, author, msg))
        time.sleep(2)

if __name__ == '__main__':
    jwt = login()
    session = unwrap_jwt(jwt)
    print('')

    my_id = session['Uid']
    peer_id, peer_name  = lookup_user()

    if os.fork() == 0:
        poll_messages(peer_name, jwt)

    while True:
        print('\rChat with %s: ' % peer_name, flush=True, end='')
        msg = input()
        send_message(msg, peer_name, jwt)
