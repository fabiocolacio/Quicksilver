import gi

gi.require_version('Gtk', '3.0')
from gi.repository import Gtk

class LoginDialog(Gtk.Dialog):
    def __init__(self, parent):
        Gtk.Dialog.__init__(self)

        self.set_modal(True)
        self.set_transient_for(parent)

        self.set_title('Login')
        self.set_default_size(300, 300)

        self.add_buttons(
            'Cancel', Gtk.ResponseType.CLOSE,
            'Login', Gtk.ResponseType.OK)
        self.set_default_response(Gtk.ResponseType.OK)

        self.grid = Gtk.Grid.new()
        self.get_content_area().pack_start(self.grid, True, True, 0)

        label = Gtk.Label.new('Username: ')
        self.user_entry = Gtk.Entry.new()
        self.grid.attach(label, 0, 0, 1, 1)
        self.grid.attach_next_to(self.user_entry, label, Gtk.PositionType.RIGHT, 1, 1)

        label = Gtk.Label.new('Password: ')
        self.passwd_entry = Gtk.Entry.new()
        self.passwd_entry.set_visibility(False)
        self.grid.attach(label, 0, 1, 1, 1)
        self.grid.attach_next_to(self.passwd_entry, label, Gtk.PositionType.RIGHT, 1, 1)

        self.show_all()

class PeerDialog(Gtk.Dialog):
    def __init__(self, parent):
        Gtk.Dialog.__init__(self)

        self.set_modal(True)
        self.set_transient_for(parent)

        self.set_title('Select a peer')
        self.set_default_size(300, 300)

        self.add_buttons(
            'Cancel', Gtk.ResponseType.CLOSE,
            'Ok', Gtk.ResponseType.OK)
        self.set_default_response(Gtk.ResponseType.OK)

        self.grid = Gtk.Grid.new()
        self.get_content_area().pack_start(self.grid, True, True, 0)

        label = Gtk.Label.new('Peer: ')
        self.user_entry = Gtk.Entry.new()
        self.grid.attach(label, 0, 0, 1, 1)
        self.grid.attach_next_to(self.user_entry, label, Gtk.PositionType.RIGHT, 1, 1)

        self.show_all()

class ChatWindow(Gtk.Window):
    def __init__(self):
        Gtk.Window.__init__(self)

        self.vbox = Gtk.Box(orientation=Gtk.Orientation.VERTICAL, spacing=8)

        self.buf = Gtk.TextBuffer.new()
        self.in_text = Gtk.TextView.new_with_buffer(self.buf)
        self.in_text.set_editable(False)
        self.in_text.set_cursor_visible(False)
        self.in_text.set_pixels_below_lines(10)
        self.in_text.set_wrap_mode(Gtk.WrapMode.WORD)

        self.top_scroll = Gtk.ScrolledWindow.new(None, None)
        self.top_scroll.set_policy(Gtk.PolicyType.NEVER, Gtk.PolicyType.ALWAYS)
        self.top_scroll.set_shadow_type(Gtk.ShadowType.IN)
        self.top_scroll.add(self.in_text)
        self.vbox.pack_start(self.top_scroll, True, True, 0)

        self.out_text = Gtk.TextView.new()
        self.out_text.set_size_request(0, 20)
        self.out_text.set_accepts_tab(False)
        self.out_text.set_wrap_mode(Gtk.WrapMode.NONE)

        self.bot_scroll = Gtk.ScrolledWindow.new(None, None)
        self.bot_scroll.set_policy(Gtk.PolicyType.NEVER, Gtk.PolicyType.NEVER)
        self.bot_scroll.set_shadow_type(Gtk.ShadowType.IN)
        self.bot_scroll.add(self.out_text)
        self.vbox.pack_start(self.bot_scroll, False, False, 0)

        self.add(self.vbox)
        self.set_title('Quicksilver')
        self.set_default_size(600, 500)
        self.show_all()

    def show_message(self, text, author, time):
        iter = self.buf.get_end_iter()
        txt = '[%s] %s: %s\n' % (time, author, text)

        self.buf.insert(iter, txt)
