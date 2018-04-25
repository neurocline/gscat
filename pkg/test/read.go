// read.go
// iterate directories and read files

package test

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/karrick/godirwalk"
	"github.com/MichaelTJones/walk"
	"github.com/iafan/cwalk"
)

func ReadTest(basepath string) {
	fmt.Printf("------------\nReading\n------------\n")
	read_filepath(basepath)
	fmt.Println()
	read_tjones_walk(basepath)
	fmt.Println()
	read_iafan_cwalk(basepath)
	fmt.Println()
	read_karrick_godirwalk(basepath)
}

func read_ignore(p string) {
	f, err := os.Open(p)
	if err != nil {
		fmt.Printf("Error opening %s: %s\n", p, err)
		return
	}
	defer f.Close()
	buf := make([]byte, 1000000)
	var total int64
	var count int
	for err == nil {
		count, err = f.Read(buf)
		total += int64(count)
		if count == 0 {
			break
		}
	}
}

func read_filepath(basepath string) {
	defer print_elapsed(time.Now())
	fmt.Println("calling path/filepath.Walk")
	dot := 0
	count := 0
	dircount := 0
	filecount := 0
	var filebytes int64 = 0
	othercount := 0
	fmt.Printf("Searching in %s\n", basepath)
	filepath.Walk(basepath, func(root string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		count++
		switch mode := info.Mode(); {
		case mode.IsRegular():
			filecount++
			filebytes = filebytes + info.Size()
			read_ignore(root)
		case mode.IsDir():
			dircount++
		default:
			othercount++
		}

		dot++
		if dot == 1000 {
			fmt.Fprintf(os.Stderr, ".")
			dot = 0
		}
		return nil
	})
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Printf("total=%d dirs=%d files=%d bytes=%dMB other=%d\n", count, dircount, filecount, filebytes/1000000, othercount)
}

func read_tjones_walk(basepath string) {
	defer print_elapsed(time.Now())
	fmt.Println("calling MichaelTJones/walk.Walk")

	// These get parallel access, so need a mutex to update them
	var mu sync.Mutex
	dot := 0
	count := 0
	dircount := 0
	filecount := 0
	var filebytes int64 = 0
	othercount := 0

	fmt.Printf("Searching in %s\n", basepath)
	cwalk.Walk(basepath, func(root string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		count++
		switch mode := info.Mode(); {
		case mode.IsRegular():
			filecount++
			filebytes = filebytes + info.Size()
			mu.Unlock()
			read_ignore(root)
			mu.Lock()
		case mode.IsDir():
			dircount++
		default:
			othercount++
		}

		dot++
		if dot == 1000 {
			fmt.Fprintf(os.Stderr, ".")
			dot = 0
		}
		return nil
	})
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Printf("total=%d dirs=%d files=%d bytes=%dMB other=%d\n", count, dircount, filecount, filebytes/1000000, othercount)
}

func read_iafan_cwalk(basepath string) {
	defer print_elapsed(time.Now())
	fmt.Println("calling iafan/cwalk.Walk")

	// These get parallel access, so need a mutex to update them
	var mu sync.Mutex
	dot := 0
	count := 0
	dircount := 0
	filecount := 0
	var filebytes int64 = 0
	othercount := 0
	fmt.Printf("Searching in %s\n", basepath)
	walk.Walk(basepath, func(root string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		count++
		switch mode := info.Mode(); {
		case mode.IsRegular():
			filecount++
			filebytes = filebytes + info.Size()
			mu.Unlock()
			read_ignore(root)
			mu.Lock()
		case mode.IsDir():
			dircount++
		default:
			othercount++
		}

		dot++
		if dot == 1000 {
			fmt.Fprintf(os.Stderr, ".")
			dot = 0
		}
		return nil
	})
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Printf("total=%d dirs=%d files=%d bytes=%dMB other=%d\n", count, dircount, filecount, filebytes/1000000, othercount)
}

// Once you call Stat or Lstat the performance difference for godirwalk
// is gone, and the parallel walkers race far ahead, relatively speaking.

func read_karrick_godirwalk(basepath string) {
	defer print_elapsed(time.Now())
	fmt.Println("Calling github.com/karrick/godirwalk.Walk")
	dot := 0
	count := 0
	dircount := 0
	filecount := 0
	var filebytes int64 = 0
	othercount := 0
	fmt.Printf("Searching in %s\n", basepath)
	godirwalk.Walk(basepath, &godirwalk.Options{
		Unsorted: true,
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			count++
			switch mode := de.ModeType(); {
			case mode.IsRegular():
				filecount++
				fi, err := os.Lstat(osPathname)
				if err == nil {
					filebytes = filebytes + fi.Size()
				read_ignore(osPathname)
				}
			case mode.IsDir():
				dircount++
			default:
				othercount++
			}

			dot++
			if dot == 1000 {
				fmt.Fprintf(os.Stderr, ".")
				dot = 0
			}
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
	})
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Printf("total=%d dirs=%d files=%d bytes=%dMB other=%d\n", count, dircount, filecount, filebytes/1000000, othercount)
}
