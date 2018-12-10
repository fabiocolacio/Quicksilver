package gui

import(
    "errors"
    "github.com/gotk3/gotk3/gtk"
)

var(
    ErrNoPeer = errors.New("No peer specified.")
)

func PeerDialogRun(parent *gtk.Window) (peer string, err error) {
    dialog, err := gtk.DialogNew()
    if err != nil {
        return
    }

    dialog.SetTransientFor(parent)
    dialog.SetModal(true)
    dialog.SetTitle("Select Peer")
    dialog.SetDefaultSize(300, 300)

    if _, err = dialog.AddButton("Cancel", gtk.RESPONSE_CLOSE); err != nil {
        return
    }
    if _, err = dialog.AddButton("Ok", gtk.RESPONSE_OK); err != nil {
        return
    }

    dialog.SetDefaultResponse(gtk.RESPONSE_OK)

    grid, err := gtk.GridNew()
    if err != nil {
        return
    }

    content, err := dialog.GetContentArea()
    if err != nil {
        return
    }
    content.PackStart(grid, true, true, 0)

    label, err := gtk.LabelNew("Peer: ")
    if err != nil {
        return
    }
    grid.Attach(label, 0, 0, 1, 1)

    peerEntry, err := gtk.EntryNew()
    if err != nil {
        return
    }
    grid.AttachNextTo(peerEntry, label, gtk.POS_RIGHT, 1, 1)

    grid.ShowAll()

    status := dialog.Run()
    defer dialog.Destroy()
    if status != gtk.RESPONSE_OK {
        return peer, ErrNoPeer
    }

    peer, err = peerEntry.GetText()
    if err != nil {
        return
    }

    return
}
