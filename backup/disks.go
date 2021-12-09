package backup

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mjmorell/ShonStop/console"

	"github.com/ricochet2200/go-disk-usage/du"
)

//Disk holds a lot of information for disk paths, such as "C:\\" or "/Volumes/Macintosh\ HD/"
type Disk struct {
	Usage *du.DiskUsage

	RootFiles   []RootFile
	RootFolders []RootFolder

	AbsPath string
	Name    string
}

//RemoveDisk removes a disk entry from a list
func RemoveDisk(slice []Disk, s int) []Disk {
	return append(slice[:s], slice[s+1:]...)
}

//GetVolumes gets a list of available mounted disks/partitions on a mac, in the /Volumes/ folder
func GetVolumes() (r []Disk) {
	files, err := ioutil.ReadDir("/Volumes/")
	if err != nil {
		console.Error(err, 41, false)
	}

	for _, f := range files {
		_, err := ioutil.ReadDir(filepath.Join("/", "Volumes", f.Name()))
		if err != nil {
			continue
		}
		r = append(r, GetInformation(filepath.Join("/", "Volumes", f.Name())))
		r[len(r)-1].Name = f.Name()
	}
	return
}

//GetDrives gets a list of available mounted disks on a Windows machine, in the [letter]:\ type
func GetDrives() (r []Disk) {
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		_, err := os.Open(string(drive) + ":\\")
		if err != nil {
			continue
		}

		r = append(r, GetInformation(string(drive)+":\\"))

	}
	return
}

//GetInformation gets info on disks on a windows machine and appends useful information
func GetInformation(src string) (drive Disk) {
	drive.Usage = du.NewDiskUsage(src)
	drive.AbsPath = src
	entries, _ := ioutil.ReadDir(src)

	for _, entry := range entries {
		if entry.Name() == "_backupFlag" {
			drive.AbsPath = filepath.Join(src, "backups")
			return
		}
		if entry.Name() == "Users" {
			drive.AbsPath = filepath.Join(src, "Users")
			return
		}
	}
	return
}
