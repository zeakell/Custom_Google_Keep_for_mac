package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	appName  = "Google Keep by KH"
	bundleID = "com.user.googlekeep"
)

func main() {
	fmt.Println("🚀 Memulai proses build & packaging...")

	// 1. Kompilasi Binary Utama
	fmt.Println("⚙️  Mengompilasi binary...")
	cmd := exec.Command("go", "build", "-ldflags", "-s -w", "-o", appName, "main-dmg.go")

	// Tambahkan ini agar error dari compiler muncul di terminal Anda:
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ Gagal kompilasi:\n%s\n", string(output))
		return
	}

	// 2. Buat Struktur Folder .app
	appFolder := appName + ".app"
	contentsDir := filepath.Join(appFolder, "Contents")
	macOSDir := filepath.Join(contentsDir, "MacOS")
	resDir := filepath.Join(contentsDir, "Resources")

	os.RemoveAll(appFolder) // Bersihkan build lama
	os.MkdirAll(macOSDir, 0755)
	os.MkdirAll(resDir, 0755)

	// 3. Pindahkan Binary & Buat Info.plist
	os.Rename(appName, filepath.Join(macOSDir, appName))

	plist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>Google Keep by KH</string>
    <key>CFBundleIconFile</key>
    <string>icon.icns</string>
    <key>CFBundleIdentifier</key>
    <string>com.user.googlekeep</string>
    <key>CFBundleName</key>
    <string>Google Keep by KH</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.13</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>LSBackgroundOnly</key>
    <false/>
</dict>
</plist>`, appName, bundleID, appName)
	os.WriteFile(filepath.Join(contentsDir, "Info.plist"), []byte(plist), 0644)

	// 4. Proses Ikon (Sama seperti sebelumnya)
	if _, err := os.Stat("icon.png"); err == nil {
		fmt.Println("🎨 Merakit ikon .icns...")
		processIcon("icon.png", filepath.Join(resDir, "icon.icns"))
	}

	// 5. MEMBUAT INSTALLER DMG
	fmt.Println("📦 Membuat Installer DMG...")
	createDMG(appFolder)

	fmt.Println("✨ Selesai! Cek file 'Google Keep Installer.dmg' di folder Anda.")
}

func createDMG(appFolder string) {
	// Buat folder sementara untuk isi DMG
	stagingGDir := "dmg_staging"
	os.Mkdir(stagingGDir, 0755)
	defer os.RemoveAll(stagingGDir)

	// Copy aplikasi ke folder staging
	exec.Command("cp", "-R", appFolder, stagingGDir).Run()

	// Buat link ke folder /Applications agar user bisa drag-and-drop
	exec.Command("ln", "-s", "/Applications", filepath.Join(stagingGDir, "Applications")).Run()

	// Jalankan hdiutil (tool bawaan macOS) untuk membuat DMG
	dmgName := "Google Keep Installer by KH.dmg"
	os.Remove(dmgName) // hapus jika sudah ada

	cmd := exec.Command("hdiutil", "create",
		"-volname", "Google Keep Install",
		"-srcfolder", stagingGDir,
		"-ov", "-format", "UDZO",
		dmgName)

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Gagal membuat DMG: %v\n", err)
	}
}

func processIcon(src, dst string) {
	// ... (fungsi processIcon yang sama dengan sebelumnya)
	iconset := "tmp.iconset"
	os.Mkdir(iconset, 0755)
	defer os.RemoveAll(iconset)
	sizes := []int{16, 32, 128, 256, 512}
	for _, s := range sizes {
		out := filepath.Join(iconset, fmt.Sprintf("icon_%dx%d.png", s, s))
		exec.Command("sips", "-z", fmt.Sprint(s), fmt.Sprint(s), src, "--out", out).Run()
	}
	exec.Command("iconutil", "-c", "icns", iconset, "-o", dst).Run()
}
