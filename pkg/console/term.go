// term.go
// Note - make multiple versions of this

// +build windows

package console

// Code taken from https://github.com/Azure/go-ansiterm/blob/master/winterm/api.go

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	kernel32DLL = syscall.NewLazyDLL("kernel32.dll")

	getConsoleScreenBufferInfoProc = kernel32DLL.NewProc("GetConsoleScreenBufferInfo")
)

type CONSOLE_SCREEN_BUFFER_INFO struct {
	Size              COORD
	CursorPosition    COORD
	Attributes        uint16
	Window            SMALL_RECT
	MaximumWindowSize COORD
}

type COORD struct {
	X int16
	Y int16
}

type SMALL_RECT struct {
	Left   int16
	Top    int16
	Right  int16
	Bottom int16
}

// GetConsoleScreenBufferInfo retrieves information about the specified console screen buffer.
// See http://msdn.microsoft.com/en-us/library/windows/desktop/ms683171(v=vs.85).aspx.
func GetConsoleScreenBufferInfo(handle uintptr) (*CONSOLE_SCREEN_BUFFER_INFO, error) {
	handle = getStdHandle(syscall.STD_OUTPUT_HANDLE)
	info := CONSOLE_SCREEN_BUFFER_INFO{}
	err := checkError(getConsoleScreenBufferInfoProc.Call(handle, uintptr(unsafe.Pointer(&info)), 0))
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// checkError evaluates the results of a Windows API call and returns the error if it failed.
func checkError(r1, r2 uintptr, err error) error {
	// Windows APIs return non-zero to indicate success
	if r1 != 0 {
		return nil
	}

	// Return the error if provided, otherwise default to EINVAL
	if err != nil {
		return err
	}
	return syscall.EINVAL
}

func getStdHandle(stdhandle int) uintptr {
	handle, err := syscall.GetStdHandle(stdhandle)
	if err != nil {
		panic(fmt.Errorf("could not get standard io handle %d", stdhandle))
	}
	return uintptr(handle)
}

func (info CONSOLE_SCREEN_BUFFER_INFO) String() string {
	return fmt.Sprintf("Size(%v) Cursor(%v) Window(%v) Max(%v)", info.Size, info.CursorPosition, info.Window, info.MaximumWindowSize)
}
