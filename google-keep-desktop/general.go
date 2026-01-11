package main

import (
	webview "github.com/webview/webview_go"
)

func main() {
	// Menentukan ukuran jendela (lebar, tinggi, hint)
	w := webview.New(false)
	defer w.Destroy()

	w.SetTitle("Google Keep Desktop")
	w.SetSize(1024, 768, webview.HintNone)

	// Navigasi langsung ke URL Google Keep
	w.Navigate("https://keep.google.com/")

	// Menjalankan loop aplikasi
	w.Run()
}
