package backup

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar"
)

//RootFolder holds a lot of information for each root folder found in a user's directory
type RootFolder struct {
	Size uint64

	Percentage float32

	AbsPath string
	Name    string
}

func copyDirectoryWin(src, dst string, errorLog, fileLog *log.Logger, buf []byte, bar *progressbar.ProgressBar) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	// Get the stats of a folder - such as symlink and permissions
	si, err := os.Stat(src)
	if err != nil {
		fmt.Printf("  ERROR! " + err.Error() + "\n - Has occured. Please verify this error!")
		errorLog.Println(src + " - " + err.Error())
		return err
	}

	// Make a directory in destination with same permissions
	if err = os.MkdirAll(dst, si.Mode()); err != nil {
		errorLog.Println(src + " - ERROR: MkDir")
		return
	}

	// Read the source and get all subentries of the folder
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		errorLog.Println(src + " - ERROR: DNE or ReadDir")
		return
	}

	// Iterate over each entry/item in the folder
	for _, entry := range entries {

		// Create a temporary src/dst for each entry
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// If an entry is a folder, start this loop all over again
			if err = copyDirectoryWin(srcPath, dstPath, errorLog, fileLog, buf, bar); err != nil {
				errorLog.Println(srcPath + " - ERROR: copyDirectory")
				continue
			}
		} else {
			bar.Add(int(entry.Size()))
			if err = copyFileWin(srcPath, dstPath, errorLog, fileLog, buf); err != nil {
				errorLog.Println(srcPath + " - ERROR: copyFile")
				continue
			}
		}
	}
	return
}

func copyDirectoryMac(src, dst string, errorLog, fileLog *log.Logger, buf []byte, bar *progressbar.ProgressBar) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	/*
		// Get the stats of a folder - such as symlink and permissions
		si, err := os.Stat(src)
		if err != nil {
			fmt.Printf("  ERROR! " + err.Error() + "\n - Has occured. Please verify this error!")
			errorLog.Println(src + " - " + err.Error())
			return err
		}
	*/
	// Make a directory in destination with same permissions
	if err = os.MkdirAll(dst, 0777); err != nil {
		errorLog.Println(src + " - ERROR: MkDir")
		return
	}

	// Read the source and get all subentries of the folder
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		errorLog.Println(src + " - ERROR: DNE or ReadDir")
		return
	}

	// Iterate over each entry/item in the folder
	for _, entry := range entries {

		// Create a temporary src/dst for each entry
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// If an entry is a folder, start this loop all over again
			if err = copyDirectoryMac(srcPath, dstPath, errorLog, fileLog, buf, bar); err != nil {
				errorLog.Println(srcPath + " - ERROR: copyDirectory")
				continue
			}
		} else {
			// If it's not a link
			/*
				if entry.Mode()&os.ModeSymlink != 0 {
					continue
				}
			*/
			// Copy that file
			bar.Add(int(entry.Size()))
			if err = copyFileMac(srcPath, dstPath, errorLog, fileLog, buf); err != nil {
				errorLog.Println(srcPath + " - ERROR: copyFile")
				continue
			}
		}
	}
	return
}

// copyDirectory copys src folder to dst folder, creating dst if it doesn't exist
func copyDirectory(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			err = copyDirectory(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}
	return
}

// appData copys the appdata from src (user's root folder) to the dst user's folder
func appData(src, dst string, errorLog *log.Logger) {
	//Firefox Bookmarks
	err := copyFirefoxAD(src+"\\AppData\\Roaming\\Mozilla\\Firefox\\Profiles", dst+"\\Desktop\\FirefoxData")
	if err != nil {
		errorLog.Println("XXXXX ERROR - FIREFOX XXXXX")
		errorLog.Print(err)
	}

	//Chrome Bookmarks
	os.MkdirAll(dst+"\\AppData\\Local\\Google\\Chrome\\User Data\\Default", 0777)
	err = CopyFile(src+"\\AppData\\Local\\Google\\Chrome\\User Data\\Default\\Bookmarks", dst+"\\AppData\\Local\\Google\\Chrome\\User Data\\Default\\Bookmarks")
	if err != nil {
		errorLog.Println("XXXXX ERROR - CHROME XXXXX")
		errorLog.Print(err)
	}

	//Microsoft Signatures
	os.MkdirAll(dst+"\\AppData\\Roaming\\Microsoft", 0777)
	err = copyDirectory(src+"\\AppData\\Roaming\\Microsoft\\Signatures", dst+"\\AppData\\Roaming\\Microsoft\\Signatures")
	if err != nil {
		errorLog.Println("XXXXX ERROR - MICROSOFT SIGNATURES XXXXX")
		errorLog.Print(err)
	}

	//Microsoft Sticky Notes (PRE WINDOWS 8.1)
	err = copyDirectory(src+"\\AppData\\Local\\Packages\\Microsoft.MicrosoftStickyNotes_8wekyb3bd8bbwe\\LocalState", dst+"\\AppData\\Local\\Packages\\Microsoft.MicrosoftStickyNotes_8wekyb3bd8bbwe\\LocalState")
	if err != nil {
		errorLog.Println("XXXXX ERROR - MICROSOFT STICKY NOTES XXXXX")
		errorLog.Print(err)
	}

	err = copyDirectory(src+"\\AppData\\Roaming\\Microsoft\\Document Building Blocks", dst+"\\AppData\\Roaming\\Microsoft\\Document Building Blocks")
	if err != nil {
		errorLog.Println("XXXXX ERROR - MICROSOFT Document Building Blocks XXXXX")
		errorLog.Print(err)
	}
}

// copyFirefoxAD copys the firefox data from source (a user directory) to the destination user's Desktop
func copyFirefoxAD(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyFirefoxAD(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			if entry.Name() == "places.sqlite" {
				err = CopyFile(srcPath, dstPath)
				if err != nil {
					return
				}
			} else {
				continue
			}
		}
	}

	return
}

// directorySize returns the byte-size of the source. If len(isRoot) != 0, it will skip Appdata/Library
func directorySize(src string, isRoot ...bool) (size uint64) {
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return 0
	}
	for _, entry := range entries {
		if strings.ToLower(entry.Name()) == "appdata" && len(isRoot) > 0 {
			continue
		}
		if strings.ToLower(entry.Name()) == "library" && len(isRoot) > 0 {
			continue
		}
		if entry.IsDir() {
			size += directorySize(filepath.Join(src, entry.Name()))
		} else if regexableFile(entry.Name()) {
			continue
		} else {
			size += uint64(entry.Size())
		}
	}
	return
}

// library copys the library folder data from src (user's root folder) to the dst user's folder
func library(src, dst string, errorLog *log.Logger) {
	//	appDataMisc=("Library/Application Support/Firefox" "/Library/Application Support/Google/Chrome/Default/Bookmarks" "Library/StickiesDatabase") #Library items that require no special options

	os.MkdirAll(filepath.Join(dst, "Library/Application Support/Firefox"), 0777)
	err := copyDirectory(filepath.Join(src, "Library/Application Support/Firefox"), filepath.Join(dst, "Library/Application Support/Firefox"))
	if err != nil {
		errorLog.Println("Library Error: Firefox data at", filepath.Join(src, "Library/Application Support/Firefox"), "Empty")
		errorLog.Println(err.Error())
	}
	os.MkdirAll(filepath.Join(dst, "/Library/Application Support/Google/Chrome/Default"), 0777) // "Bookmarks is the file needed"
	err = CopyFile(filepath.Join(src, "/Library/Application Support/Google/Chrome/Default/Bookmarks"), filepath.Join(dst, "/Library/Application Support/Google/Chrome/Default/Bookmarks"))
	if err != nil {
		errorLog.Println("Library Error: Chrome file at", filepath.Join(src, "/Library/Application Support/Google/Chrome/Default/Bookmarks"), "Empty")
		errorLog.Println(err.Error())
	}

	os.MkdirAll(filepath.Join(dst, "/Library/StickiesDatabase"), 0777)
	err = copyDirectory(filepath.Join(src, "Library/StickiesDatabase"), filepath.Join(dst, "Library/StickiesDatabase"))
	if err != nil {
		errorLog.Println("Library Error: StickyDatabase data at", filepath.Join(src, "Library/StickiesDatabase"), "Empty")
		errorLog.Println(err.Error())
	}

	os.MkdirAll(filepath.Join(dst, "/Library/Safari"), 0777)
	err = CopyFile(filepath.Join(src, "/Library/Safari/Bookmarks.plist"), filepath.Join(dst, "/Library/Safari/Bookmarks.plist"))
	if err != nil {
		errorLog.Println("Library Error: Safari file at", filepath.Join(src, "/Library/Safari/Bookmarks.plist"), "Empty")
		errorLog.Println(err.Error())
	}
}
