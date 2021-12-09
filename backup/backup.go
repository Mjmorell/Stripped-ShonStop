package backup

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/mjmorell/ShonStop/console"
)

//Backup holds lots of pertanant information for a backup process
type Backup struct {
	Unix       bool
	Paranoid   bool
	FolderMode bool

	FD string

	Identifier string

	Source      Disk
	Destination Disk
	AllDisks    []Disk

	SelectedUsers []User
	AllUsers      []User

	MainLog *log.Logger
	FileLog *log.Logger
	ErrLog  *log.Logger
	EndLog  *log.Logger

	buf []byte
}

//NewBackup returns a Backup, testing sudo, setting fd, and setting OS
func NewBackup() (temp Backup) {
	temp.setFD()
	temp.setOS()
	return
}

//setFD sets Folder Separater character
func (b *Backup) setFD() {
	if runtime.GOOS == "windows" {
		b.FD = "\\"
	} else {
		b.FD = "/"
	}

	b.buf = make([]byte, 256)
}

//setOS sets the runtime. Not really needed but mainly for readability
func (b *Backup) setOS() {
	if runtime.GOOS != "windows" {
		b.Unix = true
	} else {
		b.Unix = false
	}
}

//SetSudo verifies admin access on a mac
func (b *Backup) SetSudo() {
	if b.Unix {
		out, err := exec.Command("id", "-u").Output()
		fmt.Println(out)
		if err != nil {
			console.Error(err, 901, false)
		}

		i, err := strconv.Atoi(string(out[:len(out)-1]))
		if err != nil {
			console.Error(err, 902, false)
		}
		fmt.Println(i)
		if i != 0 {
			console.Head()
			fmt.Println("Please run me as sudo.")
			fmt.Println("You can type", console.Red("sudo !!"), "for quicker response")
			fmt.Println()
			fmt.Println()
			os.Exit(21)
		}
	}
}

//SetLog sets up the log files for the program
func (b *Backup) SetLog() {
	main, err := os.OpenFile(b.Identifier+"_MAIN.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		console.Error(err, 31, false)
	}

	errors, err := os.OpenFile(b.Identifier+"_ERROR.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		console.Error(err, 32, false)
	}

	file, err := os.OpenFile(b.Identifier+"_FILE.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		console.Error(err, 33, false)
	}

	end, err := os.OpenFile(b.Identifier+"_END-REPORT.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		console.Error(err, 33, false)
	}

	b.MainLog = log.New(main, "", log.Ltime)
	b.ErrLog = log.New(errors, "", log.Ltime)
	b.FileLog = log.New(file, "", log.Ltime|log.Lmicroseconds)
	b.EndLog = log.New(end, "", 0)

	b.MainLog.Println("---- NEW START ----")
	b.ErrLog.Println("---- NEW START ----")
	b.FileLog.Println("---- NEW START ----")
	b.EndLog.Println()
}
