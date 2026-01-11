package main

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/getlantern/systray"
	webview "github.com/webview/webview_go"
)

func init() {
	// Memastikan runtime Go selalu menggunakan thread OS yang sama untuk UI
	runtime.LockOSThread()
}

func main() {
	// 1. Jalankan Systray di goroutine terpisah
	go func() {
		systray.Run(onReady, onExit)
	}()

	// 2. Jalankan WebView di Thread Utama (Main Thread)
	// Ini krusial untuk macOS
	w := webview.New(false)
	defer w.Destroy()

	w.SetTitle("Google Keep by KH")
	w.SetSize(1024, 768, webview.HintNone)

	// Script untuk membuka link eksternal di browser default
	w.Bind("openExternal", func(url string) {
		exec.Command("open", url).Run()
	})

	w.Init(`
		document.addEventListener('click', function(e) {
			const target = e.target.closest('a');
			if (target && target.href && !target.href.includes('keep.google.com')) {
				e.preventDefault();
				window.openExternal(target.href);
			}
		}, true);
	`)

	w.Navigate("https://keep.google.com/")

	// Loop utama aplikasi berhenti di sini
	w.Run()
}

func onReady() {
	systray.SetTitle("Keep by KH")
	systray.SetTooltip("Keep Google - by KH")

	mQuit := systray.AddMenuItem("Keluar", "Tutup Aplikasi")
	go func() {
		<-mQuit.ClickedCh
		os.Exit(0)
	}()
}

func onExit() {}
