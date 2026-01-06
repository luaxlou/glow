package manager

import (
	"fmt"
	"os"
	"sync"
)

type LogRotator struct {
	mu         sync.Mutex
	filename   string
	maxSize    int64
	maxBackups int
	file       *os.File
	size       int64
}

func NewLogRotator(filename string, maxSize int64, maxBackups int) (*LogRotator, error) {
	l := &LogRotator{
		filename:   filename,
		maxSize:    maxSize,
		maxBackups: maxBackups,
	}
	if err := l.open(); err != nil {
		return nil, err
	}
	return l, nil
}

func (l *LogRotator) open() error {
	info, err := os.Stat(l.filename)
	if err == nil {
		l.size = info.Size()
	} else if !os.IsNotExist(err) {
		return err
	}

	f, err := os.OpenFile(l.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	l.file = f
	return nil
}

func (l *LogRotator) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	writeLen := int64(len(p))
	if l.size+writeLen > l.maxSize {
		if err := l.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = l.file.Write(p)
	l.size += int64(n)
	return n, err
}

func (l *LogRotator) rotate() error {
	if l.file != nil {
		l.file.Close()
	}

	for i := l.maxBackups - 1; i >= 1; i-- {
		oldName := fmt.Sprintf("%s.%d", l.filename, i)
		newName := fmt.Sprintf("%s.%d", l.filename, i+1)

		if _, err := os.Stat(oldName); err == nil {
			os.Rename(oldName, newName)
		}
	}

	// Rename current to .1
	os.Rename(l.filename, fmt.Sprintf("%s.1", l.filename))

	l.size = 0
	return l.open()
}

func (l *LogRotator) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}
