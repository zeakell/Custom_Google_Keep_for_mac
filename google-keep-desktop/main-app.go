package main

import (
	"os/exec"

	webview "github.com/webview/webview_go"
)

func main() {
	// Debug true agar bisa klik kanan > inspect element jika ada error
	w := webview.New(false)
	defer w.Destroy()

	w.SetTitle("Keep by KH")
	w.SetSize(1000, 700, webview.HintNone)

	// Script untuk menangani klik link secara cerdas
	w.Init(`
		document.addEventListener('click', function(e) {
			var target = e.target.closest('a');
			if (target && target.href) {
				if (!target.href.includes('keep.google.com')) {
					e.preventDefault();
					window.external.invoke(target.href);
				}
			}
		}, true);
	`)

	// Handler untuk menerima instruksi dari JavaScript di atas
	w.Bind("invoke", func(url string) {
		// Perintah macOS untuk membuka browser default
		exec.Command("open", url).Run()
	})

	w.Navigate("https://keep.google.com/")
	w.Run()
}
