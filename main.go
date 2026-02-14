package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/sys/unix"
)

import flag "github.com/spf13/pflag"

const version = "0.5.0"

type FileInfo struct {
	path string
	ino  uint64
}

func hashFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close() //nolint:errcheck

	hash, err := blake2b.New256(nil)
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func dedupDirectory(directory string, dryRun bool, quiet bool, noCross bool) error {
	hashes := make(map[string]FileInfo)

	var rootStat unix.Stat_t
	if err := unix.Stat(directory, &rootStat); err != nil {
		return &os.PathError{Op: "stat", Path: directory, Err: err}
	}

	err := filepath.WalkDir(directory, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		var fileStat unix.Stat_t
		if err := unix.Lstat(path, &fileStat); err != nil {
			return &os.PathError{Op: "lstat", Path: path, Err: err}
		}
		if d.Type().IsDir() {
			// Crossing mount point
			if fileStat.Dev != rootStat.Dev {
				if !noCross {
					if err := dedupDirectory(path, dryRun, quiet, noCross); err != nil {
						return err
					}
				}
				return filepath.SkipDir
			}
			return nil
		} else if !d.Type().IsRegular() {
			return nil
		}

		fileHash, err := hashFile(path)
		if err != nil {
			return err
		}
		if destFile, exists := hashes[fileHash]; exists {
			// Same file?
			if destFile.ino == fileStat.Ino {
				return nil
			}
			if !dryRun {
				if err := os.Remove(path); err != nil {
					return err
				}
				if err := os.Link(destFile.path, path); err != nil {
					return err
				}
			}
			if !quiet {
				fmt.Printf("'%s' => '%s'\n", path, destFile.path)
			}
		} else {
			hashes[fileHash] = FileInfo{path: path, ino: fileStat.Ino}
		}
		return nil
	})
	return err
}

func main() {
	var opts struct {
		dryRun  bool
		noCross bool
		quiet   bool
		version bool
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] DIRECTORY...\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.BoolVarP(&opts.dryRun, "dry-run", "n", false, "dry run")
	flag.BoolVarP(&opts.quiet, "quiet", "q", false, "be quiet about it")
	flag.BoolVarP(&opts.noCross, "one-file-system", "x", false, "do not cross filesystems")
	flag.BoolVarP(&opts.version, "version", "", false, "print version and exit")
	flag.Parse()

	if opts.version {
		fmt.Printf("dedup v%s %v %s/%s\n", version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	log.SetFlags(0)
	log.SetPrefix("ERROR: ")

	for _, directory := range flag.Args() {
		if err := dedupDirectory(directory, opts.dryRun, opts.quiet, opts.noCross); err != nil {
			log.Fatal(err)
		}
	}
}
