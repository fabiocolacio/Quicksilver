package gui

import(
    "github.com/gotk3/gotk3/gtk"
)

type UI struct {
    Window  *gtk.Window

    outText *gtk.TextView
    inText  *gtk.TextView
}

func UINew() (*UI, error) {
    win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
    if err != nil {
        return nil, err
    }
    win.SetTitle("Quicksilver")
    win.SetDefaultSize(500, 500)
    win.Connect("destroy", gtk.MainQuit)
    win.ShowAll()

    vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 8)
    if err != nil {
        return nil, err
    }
    win.Add(vbox)

    txtBuffer, err := gtk.TextBufferNew(nil)
    if err != nil {
        return nil, err
    }

    incoming, err := gtk.TextViewNewWithBuffer(txtBuffer)
    if err != nil {
        return nil, err
    }
    incoming.SetEditable(false)
    incoming.SetCursorVisible(false)
    incoming.SetPixelsBelowLines(5)
    incoming.SetWrapMode(gtk.WRAP_WORD)

    topScroll, err := gtk.ScrolledWindowNew(nil, nil)
    if err != nil {
        return nil, err
    }
    topScroll.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_ALWAYS)
    topScroll.SetShadowType(gtk.SHADOW_IN)
    topScroll.Add(incoming)
    vbox.PackStart(topScroll, true, true, 0)

    outgoing, err := gtk.TextViewNew()
    if err != nil {
        return nil, err
    }
    outgoing.SetSizeRequest(0, 20)
    outgoing.SetAcceptsTab(false)
    outgoing.SetWrapMode(gtk.WRAP_NONE)

    botScroll, err := gtk.ScrolledWindowNew(nil, nil)
    if err != nil {
        return nil, err
    }
    botScroll.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_NEVER)
    botScroll.SetShadowType(gtk.SHADOW_IN)
    botScroll.Add(outgoing)
    vbox.PackStart(botScroll, false, false, 0)

    ui := &UI{
        Window: win,
    }

    return ui, err
}
