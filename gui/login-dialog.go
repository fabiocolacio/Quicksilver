package gui

import(
    "errors"
    "github.com/gotk3/gotk3/gtk"
)

var(
    ErrNoLogin = errors.New("No login credentals provided.")
)

func LoginDialogRun(parent *gtk.Window) (user string, passwd string, err error) {
    dialog, err := gtk.DialogNew()
    if err != nil {
        return user, passwd, err
    }

    dialog.SetTransientFor(parent)
    dialog.SetModal(true)
    dialog.SetTitle("Login")
    dialog.SetDefaultSize(300, 300)

    if _, err = dialog.AddButton("Cancel", gtk.RESPONSE_CLOSE); err != nil {
        return user, passwd, err
    }
    if _, err = dialog.AddButton("Login", gtk.RESPONSE_OK); err != nil {
        return user, passwd, err
    }

    dialog.SetDefaultResponse(gtk.RESPONSE_OK)

    grid, err := gtk.GridNew()
    if err != nil {
        return user, passwd, err
    }

    content, err := dialog.GetContentArea()
    if err != nil {
        return user, passwd, err
    }
    content.PackStart(grid, true, true, 0)

    label, err := gtk.LabelNew("Username: ")
    if err != nil {
        return user, passwd, err
    }
    grid.Attach(label, 0, 0, 1, 1)

    userEntry, err := gtk.EntryNew()
    if err != nil {
        return user, passwd, err
    }
    grid.AttachNextTo(userEntry, label, gtk.POS_RIGHT, 1, 1)

    label, err = gtk.LabelNew("Password: ")
    if err != nil {
        return user, passwd, err
    }
    grid.Attach(label, 0, 1, 1, 1)

    passwdEntry, err := gtk.EntryNew()
    passwdEntry.SetVisibility(false)
    if err != nil {
        return user, passwd, err
    }
    grid.AttachNextTo(passwdEntry, label, gtk.POS_RIGHT, 1, 1)

    grid.ShowAll()

    status := dialog.Run()
    defer dialog.Destroy()
    if status != gtk.RESPONSE_OK {
        return user, passwd, ErrNoLogin
    }

    user, err = userEntry.GetText()
    if err != nil {
        return user, passwd, err
    }

    passwd, err = passwdEntry.GetText()
    if err != nil {
        return user, passwd, err
    }

    return user, passwd, err
}
