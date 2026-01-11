package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	appName  = "Google Keep"
	bundleID = "com.user.googlekeep"
)

func main() {
	fmt.Println("🚀 Memulai build aplikasi dengan Go...")

	// 1. Kompilasi Binary Utama
	fmt.Println("⚙️  Mengompilasi binary...")
	cmd := exec.Command("go", "build", "-ldflags", "-s -w", "-o", appName, "main.go")
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Gagal kompilasi: %v\n", err)
		return
	}

	// 2. Buat Struktur Folder .app
	contentsDir := filepath.Join(appName+".app", "Contents")
	macOSDir := filepath.Join(contentsDir, "MacOS")
	resDir := filepath.Join(contentsDir, "Resources")

	os.MkdirAll(macOSDir, 0755)
	os.MkdirAll(resDir, 0755)

	// 3. Pindahkan Binary
	os.Rename(appName, filepath.Join(macOSDir, appName))

	// 4. Buat file Info.plist
	fmt.Println("📝 Membuat Info.plist...")
	plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>%s</string>
    <key>CFBundleIconFile</key>
    <string>icon.icns</string>
    <key>CFBundleIdentifier</key>
    <string>%s</string>
    <key>CFBundleName</key>
    <string>%s</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.13</string>
</dict>
</plist>`, appName, bundleID, appName)

	os.WriteFile(filepath.Join(contentsDir, "Info.plist"), []byte(plistContent), 0644)

	// 5. Proses Ikon (Memanggil tool sistem macOS via Go)
	if _, err := os.Stat("icon.png"); err == nil {
		fmt.Println("🎨 Memproses ikon...")
		processIcon("icon.png", filepath.Join(resDir, "icon.icns"))
	}

	fmt.Printf("✅ Selesai! Aplikasi %s.app telah dibuat.\n", appName)
}

func processIcon(src string, dst string) {
	iconset := "tmp.iconset"
	os.Mkdir(iconset, 0755)
	defer os.RemoveAll(iconset)

	sizes := []int{16, 32, 64, 128, 256, 512}
	for _, s := range sizes {
		// Normal size
		out := filepath.Join(iconset, fmt.Sprintf("icon_%dx%d.png", s, s))
		exec.Command("sips", "-z", fmt.Sprint(s), fmt.Sprint(s), src, "--out", out).Run()
		// @2x size
		out2x := filepath.Join(iconset, fmt.Sprintf("icon_%dx%d@2x.png", s, s))
		exec.Command("sips", "-z", fmt.Sprint(s*2), fmt.Sprint(s*2), src, "--out", out2x).Run()
	}
	exec.Command("iconutil", "-c", "icns", iconset, "-o", dst).Run()
}
