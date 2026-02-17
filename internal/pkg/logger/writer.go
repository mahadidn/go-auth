package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

/*
	Menulis log harian
*/

// DailyWriter akan diberi kemampuan mengecek tanggal dan write log-nya
type DailyWriter struct {
	LogDir 		string
	currentFile *os.File
	currentDate string
	mu			sync.Mutex // Mencegah bentrok saat banyak error terjadi bersamaan
}

// method Write ini buat DailyWriter dianggap sebagai io.Writer oleh Go
func (dw *DailyWriter) Write(p []byte) (n int, err error) {
	// locking dan unlocking agar tidak terjadi race condition
	dw.mu.Lock()
	defer dw.mu.Unlock()

	// buat tanggal format waktu sekarang
	now := time.Now().Format("2006-01-02")

	// cek apakah perlu ganti file (karna masuk hari baru) atau belum punya file terbuka
	if now != dw.currentDate {

		// pakai if condition buat ngecek, karna aplikasi yg baru nyala currentFile masih kosong, biar terhindar dari panic
		// kalau dia tidak nil, berarti ada datanya
		if dw.currentFile != nil {
			dw.currentFile.Close() // tutup file hari kemarin
		}

		// set-up file log baru
		fileName := fmt.Sprintf("%s/%s.log", dw.LogDir, now)
		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return 0, err
		}

		dw.currentFile = file
		dw.currentDate = now
	}

	return dw.currentFile.Write(p)
	
}
