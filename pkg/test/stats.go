// stats.go
// iterate directories and gather stats

package test

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/karrick/godirwalk"
	"github.com/MichaelTJones/walk"
	"github.com/iafan/cwalk"
)

func StatsTest(basepath string) {
	Walk_filepath(basepath)
	fmt.Println()
	Walk_tjones_walk(basepath)
	fmt.Println()
	Walk_iafan_cwalk(basepath)
	fmt.Println()
	Walk_karrick_godirwalk(basepath)
}

func print_elapsed(startTime time.Time) {
	elapsed := time.Now().Sub(startTime)
	fmt.Printf("Elapsed: %.2f\n", elapsed.Seconds())
}

func Walk_filepath(basepath string) {
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
		case mode.IsDir():
			dircount++
		default:
			othercount++
		}

		dot++
		if dot == 10000 {
			fmt.Fprintf(os.Stderr, ".")
			dot = 0
		}
		return nil
	})
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Printf("total=%d dirs=%d files=%d bytes=%dMB other=%d\n", count, dircount, filecount, filebytes/1000000, othercount)
}

func Walk_tjones_walk(basepath string) {
	defer print_elapsed(time.Now())
	fmt.Println("calling MichaelTJones/walk.Walk")
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
		count++
		switch mode := info.Mode(); {
		case mode.IsRegular():
			filecount++
			filebytes = filebytes + info.Size()
		case mode.IsDir():
			dircount++
		default:
			othercount++
		}

		dot++
		if dot == 10000 {
			fmt.Fprintf(os.Stderr, ".")
			dot = 0
		}
		return nil
	})
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Printf("total=%d dirs=%d files=%d bytes=%dMB other=%d\n", count, dircount, filecount, filebytes/1000000, othercount)
}

func Walk_iafan_cwalk(basepath string) {
	defer print_elapsed(time.Now())
	fmt.Println("calling iafan/cwalk.Walk")
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
		count++
		switch mode := info.Mode(); {
		case mode.IsRegular():
			filecount++
			filebytes = filebytes + info.Size()
		case mode.IsDir():
			dircount++
		default:
			othercount++
		}

		dot++
		if dot == 10000 {
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

func Walk_karrick_godirwalk(basepath string) {
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
				}
			case mode.IsDir():
				dircount++
			default:
				othercount++
			}

			dot++
			if dot == 10000 {
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
