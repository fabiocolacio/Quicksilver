package gui

import(
    "github.com/gotk3/gotk3/gtk"
    "github.com/gotk3/gotk3/gdk"
    "fmt"
    "log"
)

type UI struct {
    Window  *gtk.Window
    Callback func(msg string)

    buffer  *gtk.TextBuffer
    outText *gtk.TextView
    inText  *gtk.TextView
}

func UINew() (*UI, error) {
    ui := new(UI)

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
    outgoing.Connect("key-press-event", keyPressed, ui)

    botScroll, err := gtk.ScrolledWindowNew(nil, nil)
    if err != nil {
        return nil, err
    }
    botScroll.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_NEVER)
    botScroll.SetShadowType(gtk.SHADOW_IN)
    botScroll.Add(outgoing)
    vbox.PackStart(botScroll, false, false, 0)

    win.ShowAll()

    ui.Window = win
    ui.buffer = txtBuffer
    ui.outText = outgoing
    ui.inText = incoming

    return ui, err
}

func (ui *UI) ShowMessage(text, author, time string) {
    iter := ui.buffer.GetEndIter()
    txt := fmt.Sprintf("[%s] %s: %s\n", time, author, text)
    ui.buffer.Insert(iter, txt)
}

func keyPressed(textView *gtk.TextView, event *gdk.Event, ui *UI) bool {
    keyEvent := gdk.EventKeyNewFromEvent(event)
    if keyEvent.KeyVal() == gdk.KEY_Return {
        buf, err := textView.GetBuffer()
        if err != nil {
            log.Fatal(err)
        }

        start := buf.GetStartIter()
        end := buf.GetEndIter()
        txt, err := buf.GetText(start, end, false)
        if err != nil {
            log.Fatal(err)
        }

        if ui.Callback != nil {
            ui.Callback(txt)
        }

        buf.SetText("")

        return true
    }
    return false
}
