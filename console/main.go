package console

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey"
	"github.com/tcnksm/go-latest"
)

//Version is version...
var Version string
var Changelog string

func init() {
	Version = "0.1.12"
}

func VersionCheck() {
	githubTag := &latest.GithubTag{
		Owner:      "mjmorell",
		Repository: "ShonStop",
	}

	res, err := latest.Check(githubTag, Version)
	if err != nil {
		return
	}
	if res.Outdated {
		Head()
		fmt.Println("Currently on Version", Version, "- There is an update to ", res.Current)
		fmt.Println("You should upgrade.")
		fmt.Println("Can be downloaded here:")
		fmt.Println(Cyan("https://github.com/Mjmorell/ShonStop/releases/latest"))
		Wait()
	}
}

//Clear runs a command to clear the console
func Clear() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

//Error prints out an error and a description for said error. This will not close/panic
func Error(err error, code int, safe bool) {
	fmt.Println(Bold(HIRed("ERROR: "+strconv.Itoa(code)) + " -- " + Red(err.Error())))
	if safe {
		var check bool
		survey.AskOne(&survey.Confirm{Message: HIMagenta(err.Error()) + ": This error is not necessarily crash-worthy and I can recover.\n" + Green("?") + " Do you want to continue? (Default: N)"}, &check)
		if !check {
			if runtime.GOOS == "windows" {
				ExitWait()
			}
			os.Exit(code)
		} else {
			return
		}
	} else {
		if runtime.GOOS == "windows" {
			fmt.Println(Yellow(" -- " + err.Error() + "\n\n  This is not a recoverable error. Please pass this information on to the developers:"))
			debug.PrintStack()
			ExitWait()
			os.Exit(code)
		} else {
			fmt.Println(Yellow(" -- " + err.Error() + "\n\n  This is not a recoverable error. Please pass this information on to the developers:"))
			debug.PrintStack()
			os.Exit(code)
		}
	}
}

//Flush flushes the console of all input
func Flush() {
	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < reader.Buffered(); i++ {
		reader.ReadByte()
	}
}

//Center centers the passed string in a console line.
func Center(str string, consoleWidth, stringWidth int) {
	if consoleWidth == 0 {
		consoleWidth = 100
	}
	if stringWidth == 0 {
		fmt.Printf(fmt.Sprintf("%[1]*s", -consoleWidth, fmt.Sprintf("%[1]*s", (consoleWidth+len(str))/2, str)))
		fmt.Printf("\n")
		return
	}
	for i := 0; i < (consoleWidth-stringWidth)/2; i++ {
		fmt.Printf(" ")
	}
	fmt.Printf(str + "\n")
}

//Wait pauses the program and waits for a user input
func Wait() {
	Flush()
	fmt.Println(Green("\n> ") + "Press " + Cyan("[Enter]") + " to continue")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

//ExitWait waits for user to press enter before exiting program after an error or at end of program
func ExitWait() {
	Flush()
	fmt.Println(Green("\n> ") + "Press " + Cyan("[Enter]") + " to close the program")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

//Head prints out a nicer header formatted well
func Head() {
	Flush()
	Clear()
	fmt.Printf(HIGreen("\n> ") + HICyan("ShonStop") + " - " + Green("v"+Version) + " - " + HIGreen("KB0016926"))
	fmt.Printf(HIGreen("\n> ") + "Developed by " + HICyan("mjmorell") + " with " + HIMagenta("<3") + " for LU under " + Green("GNU AGPLv3\n\n"))
}

//PrintByte prints out nicely a passed filesize:
//  ogSize - size of item
//  attach - should we attach the "kb"
//  del - if len(del) == 1, whatever is passed ("b","kb"/"k",...,"tb"/"t"), it will size it to that divisor
func PrintByte(ogSize uint64, attach bool, del ...string) string {
	size := float64(ogSize)
	if len(del) == 1 {
		switch strings.ToLower(del[0]) {
		case "b":
			return strconv.FormatFloat(size, 'f', 0, 32)
		case "kb", "k":
			return strconv.FormatFloat(size/1024, 'f', 2, 32)
		case "mb", "m":
			return strconv.FormatFloat(size/1024/1024, 'f', 2, 32)
		case "gb", "g":
			return strconv.FormatFloat(size/1024/1024/1024, 'f', 2, 32)
		case "tb", "t":
			return strconv.FormatFloat(size/1024/1024/1024/1024, 'f', 2, 32)
		default:
			return strconv.FormatFloat(size, 'f', 0, 32)
		}
	}

	if size < 1024 {
		return strconv.FormatFloat(size, 'f', 0, 32) + " Bytes"
	}

	for _, each := range []string{" B", " KiB", " MiB", " GiB", " TiB", " PiB"} {
		if size < 1024 {
			if attach {
				return strconv.FormatFloat(size, 'f', 2, 32) + each
			}
			return strconv.FormatFloat(size, 'f', 2, 32)
		}
		size /= 1024
	}
	return ""
}

func PrintShortByte(ogSize uint64, attach bool, del ...string) string {
	size := float64(ogSize)
	if len(del) == 1 {
		switch strings.ToLower(del[0]) {
		case "b":
			return strconv.FormatFloat(size, 'f', 0, 32)
		case "kb", "k":
			return strconv.FormatFloat(size/1000, 'f', 2, 32)
		case "mb", "m":
			return strconv.FormatFloat(size/1000/1000, 'f', 2, 32)
		case "gb", "g":
			return strconv.FormatFloat(size/1000/1000/1000, 'f', 2, 32)
		case "tb", "t":
			return strconv.FormatFloat(size/1000/1000/1000/1000, 'f', 2, 32)
		default:
			return strconv.FormatFloat(size, 'f', 0, 32)
		}
	}

	if size < 1000 {
		return strconv.FormatFloat(size, 'f', 0, 32) + " Bytes"
	}

	for _, each := range []string{" B", " KB", " MB", " GB", " TB", " PB"} {
		if size < 1000 {
			if attach {
				return strconv.FormatFloat(size, 'f', 2, 32) + each
			}
			return strconv.FormatFloat(size, 'f', 2, 32)
		}
		size /= 1000
	}
	return ""
}
