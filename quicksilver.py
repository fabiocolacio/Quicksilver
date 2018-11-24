#!/bin/python3

import requests
import json
import base64
import getpass

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

def get_messages(peer, jwt):
    params = { 'peer': peer }
    headers = { 'Session': jwt }
    res = requests.get(host + '/get', params=params, headers=headers, verify=False)
    return json.loads(res.text)

if __name__ == '__main__':
    jwt = login()
    session = unwrap_jwt(jwt)
    print('')

    my_id = session['Uid']
    peer_id, peer_name  = lookup_user()

    messages = get_messages(peer_name, jwt)
    for message in messages:
        author = message['Username']
        time = message['Timestamp']
        msg = message['Message']

        print('[%s] %s: %s' % (time, author, msg))

    msg = input('\rChat with %s: ' % peer_name)
