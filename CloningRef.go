package main

import (
	"flag"

	"log"
	"os"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

func main() {
	flag.Parse()
	mountpoint := flag.Arg(0)
	rootdir := flag.Arg(1)

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

	err = fs.Serve(c, filesys{rootdir})
	if err != nil {
		log.Fatal(err)
	}

	// check if the mount process has an error to report
	<-c.Ready
	if err := c.MountError; err != nil {
		log.Fatal(err)
	}
}

type filesys struct {
	rootdir string
}

type dir struct {
	dirpath string
}
type file struct {
	filepath string
}

//filesys implements FS
func (f filesys) Root() (fs.Node, error) {
	return dir{f.rootdir}, nil
}

func (d dir) Attr(ctx context.Context, attr *fuse.Attr) error {
	a := syscall.Stat_t{}
	syscall.Lstat(d.dirpath, &a)
	attr.Inode = uint64(a.Ino)
	attr.Mode = os.ModeDir | os.FileMode(a.Mode)
	attr.Size = uint64(a.Size)
	return nil
}

func (d dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	direc, err := os.Open(d.dirpath)
	if err != nil {
		log.Fatal(err)
	}
	files, err := direc.Readdir(10)
	direc.Close()
	if err != nil {
		log.Fatal(err)
	}
	for _, filer := range files {
		if name == filer.Name() {
			if filer.IsDir() {
				return dir{d.dirpath + "/" + name}, nil
			} else {
				return file{d.dirpath + "/" + name}, nil
			}
		}
	}
	return nil, fuse.ENOENT

}

func (d dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	open, err := os.Open(d.dirpath)
	if err != nil {
		log.Fatal(err)
	}
	files, err := open.Readdir(10)
	open.Close()
	if err != nil {
		log.Fatal(err)
	}
	a := []fuse.Dirent{}
	for _, file := range files {
		info2 := syscall.Stat_t{}
		syscall.Lstat(d.dirpath+"/"+file.Name(), &info2)
		if file.IsDir() {
			a = append(a, fuse.Dirent{Inode: info2.Ino, Name: file.Name(), Type: fuse.DT_Dir})
		} else {
			a = append(a, fuse.Dirent{Inode: info2.Ino, Name: file.Name(), Type: fuse.DT_File})
		}

	}
	return a, nil

}

//file implements node and handle
func (f file) Attr(ctx context.Context, attr *fuse.Attr) error {
	attri := syscall.Stat_t{}
	syscall.Lstat(f.filepath, &attri)
	attr.Inode = uint64(attri.Ino)
	attr.Mode = os.FileMode(attri.Mode)
	attr.Size = uint64(attri.Size)
	return nil
}

func (f file) ReadAll(ctx context.Context) ([]byte, error) {
	file, err := os.Open(f.filepath)
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, 10000000000)
	_, err = file.Read(data)
	return data, err
}
