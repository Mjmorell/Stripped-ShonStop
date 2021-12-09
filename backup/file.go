package backup

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	tm "github.com/buger/goterm"
	"github.com/mjmorell/ShonStop/console"
)

//RootFile holds information for files on a root of a directory, mainly user root files
type RootFile struct {
	Size uint64

	Percentage float32

	AbsPath string
	Name    string
}

func copyFileWin(src, dst string, errorLog, fileLog *log.Logger, buf []byte) error {
	// Skip all the bad files
	if regexableFile(src) {
		errorLog.Println(src + " - REGEX")
		return nil
	}

	// Open the file for reading
	in, err := os.Open(src)
	if err != nil {
		errorLog.Println(src + " - ERROR: Open")
		return err
	}
	defer in.Close()

	temp, _ := in.Stat()
	if temp.Size() > 524288000 {
		fmt.Printf("\n" + console.Cyan("  > Single file sized at ", console.PrintByte(uint64(temp.Size()), true), ". Using Robocopy. May take a minute or 20"))
		filename := ""
		{
			filenameArr := strings.Split(src, "\\")
			filename = filenameArr[len(filenameArr)-1]
		}
		src := strings.TrimSuffix(src, "\\"+filename)
		dst := strings.TrimSuffix(dst, "\\"+filename)
		_, err := exec.Command("cmd", "/C", "Robocopy", src, dst, filename, "/W:0").Output()
		if err != nil {
			errorLog.Println(err.Error())
		}

		//fmt.Println(string(out))

		fileLog.Println(filepath.Join(src, filename))
		return nil
	} else if temp.Size() < 16384 {
		// Create a output file at destination
		out, err := os.Create(dst)
		if err != nil {
			errorLog.Println(src + " - ERROR: Overwrite")
			return err
		}
		defer func() {
			if e := out.Close(); e != nil {
				err = e
			}
		}()

		// Copy the file
		_, err = io.CopyBuffer(out, in, buf)
		if err != nil {
			errorLog.Println(src + " - ERROR: Copy")
			return err
		}

		// Get stats of file
		si, _ := os.Stat(src)
		checkOut, _ := os.Stat(dst)
		if si.Size() != checkOut.Size() {
			tempSize := si.Size() - checkOut.Size()
			errorLog.Printf("\n" + src + " size miss/matched by: " + console.PrintByte(uint64(tempSize), true))
		}

		err = os.Chmod(dst, si.Mode())
		if err != nil {
			errorLog.Println(src + " chmod err")
		}

		// Print the file only AFTER it has been copied without error
		fileLog.Println(src)
		return nil
	}
	// Create a output file at destination
	out, err := os.Create(dst)
	if err != nil {
		errorLog.Println(src + " - ERROR: Overwrite")
		return err
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	// Copy the file
	_, err = io.Copy(out, in)
	if err != nil {
		errorLog.Println(src + " - ERROR: Copy")
		return err
	}

	// Commit the write
	err = out.Sync()
	if err != nil {
		errorLog.Println(src + " - ERROR: Sync")
		return err
	}

	// Get stats of file
	si, _ := os.Stat(src)
	checkOut, _ := os.Stat(dst)
	if si.Size() != checkOut.Size() {
		tempSize := si.Size() - checkOut.Size()
		errorLog.Printf("\n" + src + " size miss/matched by: " + console.PrintByte(uint64(tempSize), true))
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		errorLog.Println(src + " chmod err")
	}

	if temp.Size() > 524288000 {
		fmt.Println("\r                                                                                                    ")
		tm.MoveCursorUp(2)
		tm.Flush()
	}
	// Print the file only AFTER it has been copied without error
	fileLog.Println(src)
	return nil

}

func copyFileMac(src, dst string, errorLog, fileLog *log.Logger, buf []byte) error {

	// Skip all the bad files
	if regexableFile(src) {
		errorLog.Println(src + " - REGEX")
		return nil
	}

	// Open the file for reading
	in, err := os.Open(src)
	if err != nil {
		errorLog.Println(src + " - ERROR: Open")
		return err
	}
	defer in.Close()

	temp, _ := in.Stat()
	if temp.Size() > 104857600 {
		fmt.Printf("\n" + console.Cyan("  > Single file sized at ", console.PrintByte(uint64(temp.Size()), true), ". May take a minute or 20"))
	}

	// Create a output file at destination
	out, err := os.Create(dst)
	if err != nil {
		errorLog.Println(src + " - ERROR: Overwrite")
		return err
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	// Copy the file
	_, err = io.Copy(out, in)
	if err != nil {
		errorLog.Println(src + " - ERROR: Copy")
		errorLog.Println(err.Error())
	}

	// Commit the write
	_ = out.Sync()

	// Get stats of file
	si, _ := os.Stat(src)
	checkOut, _ := os.Stat(dst)
	if si.Size() != checkOut.Size() {
		tempSize := si.Size() - checkOut.Size()
		errorLog.Printf("\n" + src + " size miss/matched by: " + console.PrintByte(uint64(tempSize), true))
	}

	err = os.Chmod(dst, 0777)
	if err != nil {
		errorLog.Println(src + " chmod err")
	}

	/*
		// Copy perms/stats of og file
		if err = os.Chmod(dst, si.Mode()); err != nil {
			errorLog.Println(src + " - ERROR: Chmod")
			return err
		}
	*/
	// Print the file only AFTER it has been copied without error
	fileLog.Println(src)
	if temp.Size() > 104857600 {
		fmt.Println("\r                                                                                                    ")
		tm.MoveCursorUp(2)
		tm.Flush()
	}
	return nil
}

func regexableFile(filename string) bool {
	filename = strings.ToLower(filename)

	/*if matched, _ := regexp.MatchString(`.*\.dat`, strings.ToLower(filename)); matched {
		return true

	} else */
	if matched, _ := regexp.MatchString(`shonstop.exe`, filename); matched {
		return true

	} else if matched, _ := regexp.MatchString(`shonstop.*\.exe`, filename); matched {
		return true

	} else if matched, _ := regexp.MatchString(`.*\.lnk`, filename); matched {
		return true

	} else if matched, _ := regexp.MatchString("ntuser.*", strings.ToLower(filename)); matched {
		return true

	} else if matched, _ := regexp.MatchString("desktop.ini", filename); matched {
		return true

	}
	return false
}

//CopyFile is the most basic copy function, copying a file from src to dst
func CopyFile(src, dst string) error {
	// Open the file for reading
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	// Copy the file
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	// Commit the write
	err = out.Sync()
	if err != nil {
		return err
	}

	return nil
}
