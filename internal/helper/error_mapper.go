package helper

import (
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
)

// Kamus pesan error global untuk seluruh aplikasi
var globalConstraintMessages = map[string]string{
	"email":    "Alamat email sudah terdaftar",
	"username": "Username sudah digunakan oleh orang lain",
	"phone":    "Nomor telepon sudah terdaftar",
	"code":     "Kode ini sudah ada di sistem",
}

// TranslateError adalah pintu utama untuk memproses semua jenis error
func TranslateError(err error) any {
	if err == nil {
		return nil
	}

	// 1. Cek apakah ini error dari Validator (Struct Tags)
	if ve, ok := err.(validator.ValidationErrors); ok {
		return FormatValidationError(ve) // Balikin map[string]string
	}

	// 2. Cek apakah ini error dari Database (MySQL)
	if _, ok := err.(*mysql.MySQLError); ok {
		return ParseDatabaseError(err).Error() // Balikin string pesan ramah
	}

	// 3. Jika error umum lainnya (errors.New dari business logic)
	return err.Error()
}

func ParseDatabaseError(err error) error {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		switch mysqlErr.Number {
		case 1062: // Duplicate Entry
			return formatDuplicateError(mysqlErr.Message)
		case 1451, 1452: // Foreign Key
			return errors.New("Data tidak bisa diproses karena masih berhubungan dengan data lain")
		default:
			slog.Error("DATABASE_CRITICAL_ERROR", "msg", mysqlErr.Message)
			return errors.New("Terjadi kesalahan")
		}
	}
	return err
}

func formatDuplicateError(msg string) error {
	// 1. Regex untuk mencari nama constraint
	re := regexp.MustCompile(`key '(.+?)'`)
	match := re.FindStringSubmatch(msg)
	
	if len(match) > 1 {
		keyName := strings.ToLower(match[1])

		// 2. Cek kamus global kita
		for keyword, userMessage := range globalConstraintMessages {
			if strings.Contains(keyName, keyword) {
				return errors.New(userMessage)
			}
		}
		
		// 3. DEFAULT 1: Jika key ketemu tapi tidak ada di kamus (misal: 'uq_users_ktp')
		// Jangan kembalikan 'msg' aslinya! Kembalikan pesan umum.
		return errors.New("Data tersebut sudah terdaftar di sistem")
	}
	
	// 4. DEFAULT 2: Jika regex bahkan tidak menemukan nama key-nya
	return errors.New("Terjadi duplikasi data")
}

func FormatValidationError(ve validator.ValidationErrors) map[string]string {
	result := make(map[string]string)

	for _, f := range ve {
		field := f.Field()
		tag := f.Tag()
		param := f.Param()

		switch tag {
		case "required":
			result[field] = "Bagian ini wajib diisi"
		case "email":
			result[field] = "Format email tidak valid"
		case "min":
			result[field] = fmt.Sprintf("Minimal harus %s karakter", param)
		case "max":
			result[field] = fmt.Sprintf("Maksimal %s karakter saja", param)
		case "uuid":
			result[field] = "Format ID tidak valid"
		default:
			// Jika ada tag lain yang belum terdaftar tapi punya param
			if param != "" {
				result[field] = fmt.Sprintf("Gagal validasi %s: %s", tag, param)
			} else {
				result[field] = "Data tidak valid"
			}
		}
	}
	return result
}