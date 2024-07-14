package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

func generateGPGKey(name, comment, email, passphrase string) (string, error) {
	// Simplified key generation using openpgp package (for demonstration purposes)
	// This does not create a complete keypair suitable for real use
	var entity *openpgp.Entity
	var err error

	if comment == "" {
		entity, err = openpgp.NewEntity(name, "", email, nil)
	} else {
		entity, err = openpgp.NewEntity(name, comment, email, nil)
	}
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	w, err := armor.Encode(buf, openpgp.PrivateKeyType, nil)
	if err != nil {
		return "", err
	}
	defer w.Close()

	if err := entity.SerializePrivate(w, nil); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func main() {
	a := app.New()
	w := a.NewWindow("GPG Key Generator")

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Enter your name")

	commentEntry := widget.NewEntry()
	commentEntry.SetPlaceHolder("Enter an optional comment")

	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Enter your email")

	passEntry := widget.NewPasswordEntry()
	passEntry.SetPlaceHolder("Enter your passphrase")

	result := widget.NewMultiLineEntry()
	result.SetPlaceHolder("Generated Key will appear here...")

	generateButton := widget.NewButton("Generate", func() {
		name := nameEntry.Text
		comment := commentEntry.Text
		email := emailEntry.Text
		pass := passEntry.Text

		if name == "" || email == "" || pass == "" {
			result.SetText("Name, Email, and Passphrase fields are required!")
			return
		}

		key, err := generateGPGKey(name, comment, email, pass)
		if err != nil {
			result.SetText(fmt.Sprintf("Error: %v", err))
			return
		}

		result.SetText(key)
	})

	content := container.NewVBox(
		widget.NewLabel("Name"),
		nameEntry,
		widget.NewLabel("Comment (Optional)"),
		commentEntry,
		widget.NewLabel("Email"),
		emailEntry,
		widget.NewLabel("Passphrase"),
		passEntry,
		generateButton,
		result,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 600))
	w.ShowAndRun()
}
