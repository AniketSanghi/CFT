package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s MOUNTPOINT\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		os.Exit(2)
	}
	mountpoint := flag.Arg(0)

	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName("Linter"),
		fuse.Subtype("helloworld"),
		fuse.LocalVolume(),
		fuse.VolumeName("CFT"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	err = fs.Serve(c, Filesys{})
	if err != nil {
		log.Fatal(err)
	}

	// check if the mount process has an error to report
	<-c.Ready
	if err := c.MountError; err != nil {
		log.Fatal(err)
	}
}

type Filesys struct{}

func (Filesys) Root() (fs.Node, error) {
	return Directory{}, nil
}

type Directory struct{}

func (Directory) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 1
	a.Mode = os.ModeDir | 0555
	return nil
}

func (Directory) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if name == "hey" {
		return file{}, nil
	}
	return nil, fuse.ENOENT
}

var dir = []fuse.Dirent{
	{Inode: 2, Name: "hey", Type: fuse.DT_File},
}

func (Directory) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	return dir, nil
}

type file struct{}

func (file) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 2
	a.Mode = 0777
	a.Size = uint64(len("Hey there! This is Aniket Sanghi \n"))
	return nil
}

func (file) ReadAll(ctx context.Context) ([]byte, error) {
	return []byte("Hey there! This is Aniket Sanghi \n"), nil
}
