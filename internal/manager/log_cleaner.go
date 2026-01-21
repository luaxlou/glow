package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// LogCleaner handles log file cleanup based on age and total size
type LogCleaner struct {
	LogDir     string
	MaxAgeDays int
	MaxTotalMB int
}

// NewLogCleaner creates a new log cleaner instance
func NewLogCleaner(logDir string, maxAgeDays int, maxTotalMB int) *LogCleaner {
	return &LogCleaner{
		LogDir:     logDir,
		MaxAgeDays: maxAgeDays,
		MaxTotalMB: maxTotalMB,
	}
}

// Run performs the cleanup operation
func (lc *LogCleaner) Run() error {
	if lc.LogDir == "" {
		return fmt.Errorf("log directory not specified")
	}

	// Check if log directory exists
	if _, err := os.Stat(lc.LogDir); os.IsNotExist(err) {
		// Nothing to clean if directory doesn't exist
		return nil
	}

	// Phase 1: Clean by age
	if lc.MaxAgeDays > 0 {
		if err := lc.cleanByAge(); err != nil {
			return fmt.Errorf("failed to clean logs by age: %w", err)
		}
	}

	// Phase 2: Clean by total size
	if lc.MaxTotalMB > 0 {
		if err := lc.cleanByTotalSize(); err != nil {
			return fmt.Errorf("failed to clean logs by total size: %w", err)
		}
	}

	return nil
}

// cleanByAge removes log files older than MaxAgeDays
func (lc *LogCleaner) cleanByAge() error {
	cutoffTime := time.Now().AddDate(0, 0, -lc.MaxAgeDays)

	return filepath.Walk(lc.LogDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process log files (including rotated logs *.log.N*)
		if !lc.isLogFile(path) {
			return nil
		}

		// Check if file is older than cutoff
		if info.ModTime().Before(cutoffTime) {
			fmt.Printf("Removing old log file: %s (age: %d days)\n", path, int(time.Since(info.ModTime()).Hours()/24))
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove %s: %w", path, err)
			}
		}

		return nil
	})
}

// cleanByTotalSize removes oldest log files if total size exceeds MaxTotalMB
func (lc *LogCleaner) cleanByTotalSize() error {
	var logFiles []logFileInfo

	// Collect all log files with their info
	err := filepath.Walk(lc.LogDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process log files
		if !lc.isLogFile(path) {
			return nil
		}

		logFiles = append(logFiles, logFileInfo{
			Path:    path,
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})

		return nil
	})

	if err != nil {
		return err
	}

	// Calculate total size
	var totalSize int64
	for _, lf := range logFiles {
		totalSize += lf.Size
	}

	maxSizeBytes := int64(lc.MaxTotalMB) * 1024 * 1024

	// If total size is within limit, no action needed
	if totalSize <= maxSizeBytes {
		return nil
	}

	// Sort by modification time (oldest first)
	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i].ModTime.Before(logFiles[j].ModTime)
	})

	// Remove oldest files until under limit
	for _, lf := range logFiles {
		if totalSize <= maxSizeBytes {
			break
		}

		fmt.Printf("Removing log file to reduce size: %s (size: %d bytes)\n", lf.Path, lf.Size)
		if err := os.Remove(lf.Path); err != nil {
			return fmt.Errorf("failed to remove %s: %w", lf.Path, err)
		}

		totalSize -= lf.Size
	}

	return nil
}

// isLogFile checks if a file is a log file
func (lc *LogCleaner) isLogFile(path string) bool {
	base := filepath.Base(path)

	// Match *.log, *.log.N*, *.log.N.gz, etc.
	return strings.HasSuffix(base, ".log") ||
		strings.Contains(base, ".log.") ||
		strings.HasSuffix(base, ".log.gz")
}

type logFileInfo struct {
	Path    string
	Size    int64
	ModTime time.Time
}

// StartPeriodicCleanup starts a goroutine that periodically cleans logs
func StartPeriodicCleanup(logDir string, maxAgeDays int, maxTotalMB int, interval time.Duration) {
	cleaner := NewLogCleaner(logDir, maxAgeDays, maxTotalMB)

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		// Run initial cleanup
		if err := cleaner.Run(); err != nil {
			fmt.Printf("Log cleanup error: %v\n", err)
		}

		for range ticker.C {
			if err := cleaner.Run(); err != nil {
				fmt.Printf("Log cleanup error: %v\n", err)
			}
		}
	}()
}
