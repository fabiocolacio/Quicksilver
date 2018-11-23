#!/bin/python3

import gi
gi.require_version('Gtk', '3.0')
from gi.repository import Gtk

from login import LoginWidget

def login_cb(jwt):
    print("Logged in with jwt:", jwt)

if __name__ == '__main__':
    win = Gtk.Window()
    win.connect("destroy", Gtk.main_quit)

    login = LoginWidget(login_cb)
    win.add(login)

    win.set_default_size(500, 400)
    win.show_all()
    Gtk.main()
