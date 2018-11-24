import requests
import json
import gi
gi.require_version('Gtk', '3.0')
from gi.repository import Gtk

class LoginWidget(Gtk.Box):
    def __init__(self, cb):
        Gtk.Box.__init__(self, orientation=Gtk.Orientation.VERTICAL, spacing=8)

        self.login_cb = cb

        grid = Gtk.Grid.new()
        grid.set_row_spacing(6)
        grid.set_column_spacing(12)
        self.pack_start(grid, True, True, 0)

        label = Gtk.Label.new()
        label.set_markup('<span font_desc="14.0" weight="bold">Login</span>')
        label.set_xalign(0.0)
        grid.attach(label, 0, 0, 2, 1)

        label = Gtk.Label.new('Username: ')
        label.set_xalign(1.0)
        self.user_entry = Gtk.Entry.new()
        grid.attach(label, 0, 1, 1, 1)
        grid.attach(self.user_entry, 1, 1, 1, 1)

        label = Gtk.Label.new('Password: ')
        label.set_xalign(1.0)
        self.pass_entry = Gtk.Entry.new()
        self.pass_entry.set_visibility(False)
        grid.attach(label, 0, 2, 1, 1)
        grid.attach(self.pass_entry, 1, 2, 1, 1)

        button = Gtk.Button.new_with_label('Login')
        button.connect("clicked", self._try_login)
        grid.attach(button, 0, 3, 2, 1)

    def _try_login(self, btn):
        creds = json.dumps({
            "Username": self.user_entry.get_text(),
            "Password": self.pass_entry.get_text()
        })

        res = requests.get("https://localhost:9090/login", data=creds, verify=False)

        if res.status_code == 200:
            jwt = res.text
            self.login_cb(jwt)
