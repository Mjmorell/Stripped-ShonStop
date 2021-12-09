package backup

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mjmorell/ShonStop/console"
	"github.com/schollz/progressbar"
)

//import "errors"

//User holds pertanant information for each user in a OS
type User struct {
	Size uint64

	AllFiles   []RootFile
	AllFolders []RootFolder

	AbsPath string
	Name    string
}

//GetUsers sets up the users available for a backup
func (b *Backup) GetUsers() {
	skippable := make(map[string]bool)

	if b.FolderMode {
		b.Paranoid = true

		b.AllUsers = append(b.AllUsers, User{Name: b.Source.AbsPath,
			AbsPath: b.Source.AbsPath})

		temp := strings.Split(b.AllUsers[0].AbsPath, b.FD)

		b.AllUsers[0].Name = temp[len(temp)-1]

		b.getUserInfo()
		return
	}

	for _, each := range []string{"public", "default", "admini~1"} {
		skippable[each] = true
	}

	entries, err := ioutil.ReadDir(b.Source.AbsPath)
	if err != nil {
		b.ErrLog.Println("GetUsers ERROR!")
		b.ErrLog.Println(err.Error())
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if !skippable[strings.ToLower(entry.Name())] {
			b.AllUsers = append(b.AllUsers, User{Name: entry.Name(),
				AbsPath: filepath.Join(b.Source.AbsPath, entry.Name())})
		}
	}

	b.getUserInfo()

}

//getUserInfo gets and sets up each individual user
func (b *Backup) getUserInfo() {
	for index := range b.AllUsers {
		fmt.Println(console.Yellow("Getting info on " + b.AllUsers[index].Name))
		entries, err := ioutil.ReadDir(b.AllUsers[index].AbsPath)
		b.MainLog.Println("     " + b.AllUsers[index].Name)
		for _, entry := range entries {
			if entry.IsDir() {
				b.MainLog.Println("    \\ " + entry.Name() + " - " + console.PrintByte(directorySize(filepath.Join(b.AllUsers[index].AbsPath, entry.Name())), true))
			} else {
				b.MainLog.Println("       " + entry.Name())
			}
		}
		b.MainLog.Println()
		if err != nil {
			console.Error(err, 43, false)
		}

		for _, entry := range entries {
			if regexableFile(entry.Name()) {
				continue
			}
			if b.Unix && macSkippables(entry.Name(), b.Paranoid) {
				b.FileLog.Println("Skipped " + entry.Name())
				//fmt.Println("Skipped " + entry.Name())
				continue
			}
			if !b.Unix && winSkippables(entry.Name(), b.Paranoid) {
				b.FileLog.Println("Skipped " + entry.Name())
				//fmt.Println("Skipped " + entry.Name())
				continue
			}
			if entry.IsDir() {
				b.AllUsers[index].AllFolders = append(b.AllUsers[index].AllFolders, RootFolder{AbsPath: filepath.Join(b.AllUsers[index].AbsPath, entry.Name()),
					Name: entry.Name()})
			} else {
				b.AllUsers[index].AllFiles = append(b.AllUsers[index].AllFiles, RootFile{AbsPath: filepath.Join(b.AllUsers[index].AbsPath, entry.Name()),
					Name: entry.Name(),
					Size: uint64(entry.Size())})
			}
		}

		for subIndex := range b.AllUsers[index].AllFolders {
			b.AllUsers[index].AllFolders[subIndex].Size = directorySize(b.AllUsers[index].AllFolders[subIndex].AbsPath)
		}

		b.AllUsers[index].Size = sumFolder(b.AllUsers[index].AllFolders) + sumFile(b.AllUsers[index].AllFiles)
	}
}

func sumFolder(list []RootFolder) (total uint64) {
	for _, each := range list {
		total += each.Size
	}
	return total
}

func sumFile(list []RootFile) (total uint64) {
	for _, each := range list {
		total += each.Size
	}
	return total
}

//GoBackup is the actual backup function
func (b *Backup) GoBackup() {

	b.Destination.AbsPath = filepath.Join(b.Destination.AbsPath, b.Identifier)
	if _, err := os.Stat(b.Destination.AbsPath); !os.IsNotExist(err) {
		for i := 1; i <= 100; i++ {
			if i == 100 {
				console.Error(errors.New("Error in creating destination directory "+b.Identifier+".\nToo many folders under same Identifier.\nPlease check "+b.Destination.AbsPath), 53, false)
			}
			if _, err := os.Stat(b.Destination.AbsPath + "_(" + strconv.Itoa(i) + ")"); os.IsNotExist(err) {
				b.Destination.AbsPath = b.Destination.AbsPath + "_(" + strconv.Itoa(i) + ")"
				break
			}
		}
	}

	err := os.Mkdir(b.Destination.AbsPath, 0777)
	if err != nil {
		b.ErrLog.Println(err.Error())
		b.MainLog.Println(err.Error())
		console.Error(errors.New(err.Error()+"\nError in creating destination directory "+b.Identifier+".\nPlease check "+b.Destination.AbsPath), 54, false)
	}

	for _, eachUser := range b.SelectedUsers {
		usrPath := filepath.Join(b.Destination.AbsPath, eachUser.Name)
		fmt.Println()
		b.MainLog.Println()
		b.FileLog.Println()

		fmt.Println(console.Cyan("BACKING UP "+strings.ToUpper(eachUser.Name)), "-", console.Green(console.PrintByte(eachUser.Size, true)))
		b.MainLog.Println("Backing up " + eachUser.Name)
		b.FileLog.Println("Backing up " + eachUser.Name)

		os.Mkdir(usrPath, 0777)
		for _, eachFolder := range eachUser.AllFolders {
			fmt.Println(console.Cyan(">"), console.Green(eachFolder.Name), "-", console.Yellow(console.PrintByte(eachFolder.Size, true)))
			folderPath := filepath.Join(usrPath, eachFolder.Name)
			err := os.Mkdir(folderPath, 0777)
			if err != nil {
				console.Error(err, 55, false)
			}
			bar := progressbar.New(int(eachFolder.Size))
			bar.RenderBlank()
			if b.Unix {
				copyDirectoryMac(filepath.Join(eachUser.AbsPath, eachFolder.Name), folderPath, b.ErrLog, b.FileLog, b.buf, bar)
				fmt.Printf("\r                                                                                                              ")
				fmt.Println("\r 100% |████████████████████████████████████████| ", console.Green("Completed                                 "))
			} else {
				copyDirectoryWin(filepath.Join(eachUser.AbsPath, eachFolder.Name), folderPath, b.ErrLog, b.FileLog, b.buf, bar)
				fmt.Printf("\r                                                                                                              ")
				fmt.Println("\r 100% |████████████████████████████████████████| ", console.Green("Completed                                 "))
			}
		}

		for _, eachFile := range eachUser.AllFiles {
			filePath := filepath.Join(usrPath, eachFile.Name)
			//fmt.Println("  " + b.FD + eachFile.Name + " to " + filePath)

			if b.Unix {
				copyFileMac(filepath.Join(eachUser.AbsPath, eachFile.Name), filePath, b.ErrLog, b.FileLog, b.buf)
			} else {
				copyFileWin(filepath.Join(eachUser.AbsPath, eachFile.Name), filePath, b.ErrLog, b.FileLog, b.buf)
			}
		}

		if !b.Unix {
			fmt.Println(console.Cyan(">"), console.Green("AppData"))
			appData(eachUser.AbsPath, usrPath, b.ErrLog)
		} else {
			fmt.Println(console.Cyan(">"), console.Green("Library"))
			library(eachUser.AbsPath, usrPath, b.ErrLog)
		}
	}
}

//VerifySizes checks and shows the technician the final sizes of the selected users/folder
func (b *Backup) VerifySizes() {
	fmt.Println(console.Red("This can take quite a while to generate.\nPlease wait till the prompt shows to press Enter."))
	fmt.Println(console.Cyan("Note: /1024 is for size estimation on a PC. /1000 is for Mac."))
	fmt.Println()
	fmt.Println(console.Green("Final Size Reports:"))
	b.MainLog.Println("Final Sizing Reports:")
	for _, each := range b.SelectedUsers {

		fmt.Println(console.Green(" Loading " + each.Name))
		destination := filepath.Join(b.Destination.AbsPath, each.Name)
		tempSize := directorySize(destination, true)
		b.EndLog.Println("Note: /1024 is for size estimation on a PC. /1000 is for Mac.")

		b.EndLog.Println(each.Name)
		b.EndLog.Println("     Origin:", console.PrintByte(each.Size, false, "b"))
		b.EndLog.Println("Destination:", console.PrintByte(tempSize, false, "b"))
		b.EndLog.Println("     Origin (/1024):", console.PrintByte(each.Size, true))
		b.EndLog.Println("Destination (/1024):", console.PrintByte(tempSize, true))
		b.EndLog.Println("     Origin (/1000):", console.PrintShortByte(each.Size, true))
		b.EndLog.Println("Destination (/1000):", console.PrintShortByte(tempSize, true))

		fmt.Println(console.Cyan("                 Origin:"), console.PrintByte(each.Size, false, "b"))
		fmt.Println(console.Yellow("            Destination:"), console.PrintByte(tempSize, false, "b"))
		fmt.Println(console.Cyan("         Origin (/1024):"), console.PrintByte(each.Size, true))
		fmt.Println(console.Yellow("    Destination (/1024):"), console.PrintByte(tempSize, true))
		fmt.Println(console.Cyan("         Origin (/1000):"), console.PrintShortByte(each.Size, true))
		fmt.Println(console.Yellow("    Destination (/1000):"), console.PrintShortByte(tempSize, true))
		fmt.Println()

	}
}
