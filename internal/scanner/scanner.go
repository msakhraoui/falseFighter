// Package scanner hashes files and collects their metadata
package scanner

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"syscall"

	"github.com/msakhraoui/falseFighter/internal/record"
)

func Snapshot(path string) (record.Record, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return record.Record{}, fmt.Errorf("lstat %s: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return record.Record{}, fmt.Errorf("%s: not a regular file", path)
	}
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return record.Record{}, fmt.Errorf("%s: no stat_t available", path)
	}

	file, err := os.Open(path)
	if err != nil {
		return record.Record{}, fmt.Errorf("open %s: %w", path, err)
	}
	defer func() { _ = file.Close() }()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return record.Record{}, fmt.Errorf("hash %s: %w", path, err)
	}
	return record.Record{
		Path:   path,
		SHA256: hex.EncodeToString(hash.Sum(nil)),
		Mode:   uint32(info.Mode()),
		UID:    st.Uid,
		GID:    st.Gid,
		Size:   info.Size(),
		MTime:  info.ModTime().Unix(),
		Inode:  st.Ino,
	}, nil
}

func Walk(root string) ([]record.Record, error) {
	var records []record.Record
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if path == root {
				return err
			}
			slog.Warn("skip unreadable entry", "path", path, "err", err)
			return nil
		}

		if d.IsDir() {
			return nil
		}
		if !d.Type().IsRegular() {
			return nil
		}

		s, err := Snapshot(path)
		if err != nil {
			slog.Warn("skip file", "path", path, "err", err)
			return nil
		}
		records = append(records, s)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Walk %s: %w", root, err)
	}
	return records, nil
}
