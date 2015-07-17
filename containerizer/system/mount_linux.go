package system

import (
	"fmt"
	"os"
	"syscall"
)

type Mount struct {
	Type       MountType
	SourcePath string
	TargetPath string
	Flags      int
	Data       string
	Mode       MountMode
}

type MountMode int

const (
	MountCreateDir MountMode = iota
	MountCreateFile
	MountModeNone
)

type MountType string

const (
	Tmpfs  MountType = "tmpfs"
	Proc             = "proc"
	Devpts           = "devpts"
	Bind             = "bind"
)

func (m Mount) Mount() error {
	if m.Mode == MountCreateDir {
		if err := os.MkdirAll(m.TargetPath, 0700); err != nil {
			return fmt.Errorf("system: create mount point directory %s: %s", m.TargetPath, err)
		}
	} else if m.Mode == MountCreateFile {
		if _, err := os.OpenFile(m.TargetPath, os.O_CREATE|os.O_RDONLY, 0700); err != nil {
			return fmt.Errorf("system: create mount point file %s: %s", m.TargetPath, err)
		}
	}

	sourcePath := m.SourcePath
	if sourcePath == "" {
		sourcePath = string(m.Type)
	}

	if err := syscall.Mount(sourcePath, m.TargetPath, string(m.Type), uintptr(m.Flags), m.Data); err != nil {
		return fmt.Errorf("system: mount %s on %s: %s", m.Type, m.TargetPath, err)
	}

	return nil
}
