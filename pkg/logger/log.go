package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	baseDir        = "logs"            // โฟลเดอร์หลักเก็บ log
	rotationWindow = 10 * time.Minute  // หมุนไฟล์ทุก 10 นาที
)

var (
	mutex       sync.Mutex
	currentID   string   // ตัวระบุช่วงเวลา 10 นาที เช่น "2025-08-09_05-10"
	currentFile *os.File // ไฟล์ที่กำลังเขียนอยู่
)

// สร้าง path โฟลเดอร์/ไฟล์ จากเวลาปัจจุบัน
// dir: logs/2025-08-09
// file: logs/2025-08-09/log_2025-08-09_05-10.log
func makePaths(now time.Time) (bucketID, dirPath, filePath string) {
	trunc := now.Truncate(rotationWindow)
	day := trunc.Format("2006-01-02")
	min := trunc.Format("15-04")
	bucketID = day + "_" + min

	dirPath = filepath.Join(baseDir, day)
	fileName := fmt.Sprintf("log_%s.log", bucketID)
	filePath = filepath.Join(dirPath, fileName)
	return
}

// เปิดไฟล์ให้ตรงช่วง 10 นาทีปัจจุบัน (ถ้าเปลี่ยนช่วง จะปิดของเก่าแล้วเปิดของใหม่)
func ensureFileLocked() error {
	now := time.Now()
	bucketID, dirPath, filePath := makePaths(now)

	if currentFile != nil && bucketID == currentID {
		return nil
	}
	if currentFile != nil {
		_ = currentFile.Close()
		currentFile = nil
	}
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	currentID = bucketID
	currentFile = f
	return nil
}

// เขียนบรรทัด log ตามรูปแบบ: "YYYY-MM-DD HH:MM:SS.mmm [module] [level] message"
func write(module, level, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	line := fmt.Sprintf("%s [%s] [%s] %s\n", timestamp, module, strings.ToLower(level), message)

	mutex.Lock()
	defer mutex.Unlock()

	if err := ensureFileLocked(); err != nil {
		// ถ้าเปิดไฟล์ไม่ได้ พิมพ์ลงคอนโซลกันหาย
		fmt.Printf("logger error: %v | %s", err, line)
		return
	}
	_, _ = currentFile.WriteString(line)
}

// ปิดไฟล์ตอนโปรเซสจะจบ (เรียกจาก main: defer logger.Close())
func Close() {
	mutex.Lock()
	defer mutex.Unlock()
	if currentFile != nil {
		_ = currentFile.Close()
		currentFile = nil
	}
}

// public helper (ใช้งานใน service/handlers)
func Infof(module, format string, a ...any)  { write(module, "info", fmt.Sprintf(format, a...)) }
func Warnf(module, format string, a ...any)  { write(module, "warn", fmt.Sprintf(format, a...)) }
func Errorf(module, format string, a ...any) { write(module, "error", fmt.Sprintf(format, a...)) }
