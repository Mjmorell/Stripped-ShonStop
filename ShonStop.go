package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/gen2brain/dlgs"

	"github.com/AlecAivazis/survey"

	"github.com/mjmorell/ShonStop/backup"
	"github.com/mjmorell/ShonStop/console"
)

func main() {
	defer func() {
		r := recover()
		console.Error(r.(error), 9, false)
	}()

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		fmt.Printf(console.Red(fmt.Sprintln("\ncaught sig:", sig)))
		console.Error(errors.New("You stopped the program"), 1, false)
	}()

	console.VersionCheck()

	bInfo := backup.NewBackup()

	bInfo.SetSudo()

	console.Head()

	// Process Setup

	// survey.AskOne(&survey.Input{Message: "What is your Username?"}, &bInfo.Technician)
	// if bInfo.Technician == "" {
	// 	console.Error(errors.New("Empty Username"), 11, false)
	// }

	// if bInfo.Unix {
	// bORr := 0
	// prompt := &survey.Select{
	// 	Message: "Backup or Restore?",
	// 	Options: []string{"Backup", "Restore"},
	// }
	// survey.AskOne(prompt, &bORr)
	// if bORrs == 1 { // RESTORE
	// restore(bInfo)
	// os.Exit(0)
	// }
	// }

	// Username setup

	fmt.Println()
	fmt.Println(console.HIRed("Cloud-Inclusive Mode should ONLY be used in very specific usecases."))
	fmt.Println(console.HIRed("THIS CAN (AND WILL) CAUSE DUPLICATES with cloud services"))
	survey.AskOne(&survey.Confirm{Message: "Should I proceed in a Cloud-Inclusive Manner? (Default:N)"}, &bInfo.Paranoid)

	// Identifier?
	survey.AskOne(&survey.Input{Message: "What is the Identifier for this backup?"}, &bInfo.Identifier)
	if bInfo.Identifier == "" {
		console.Error(errors.New("Empty Identifier"), 12, false)
	}
	// Log Setup
	bInfo.SetLog()
	bInfo.MainLog.Println("Identifier:", bInfo.Identifier)

	console.Head()
	// Info on disks:
	fmt.Println(console.Yellow("Getting information on disks..."))

	if bInfo.Unix {
		bInfo.AllDisks = backup.GetVolumes()
	} else {
		bInfo.AllDisks = backup.GetDrives()
	}
	bInfo.MainLog.Println()
	bInfo.MainLog.Println("Available Disks:")
	for _, each := range bInfo.AllDisks {
		bInfo.MainLog.Println(each.Name, "--", each.AbsPath)
		bInfo.MainLog.Println(" - Free:", console.PrintByte(each.Usage.Free(), true), "--", strconv.FormatFloat(100-float64(each.Usage.Usage()*100), 'f', 2, 32)+"%")
		bInfo.MainLog.Println(" - Used:", console.PrintByte(each.Usage.Used(), true), "--", strconv.FormatFloat(float64(each.Usage.Usage()*100), 'f', 2, 32)+"%")
		bInfo.MainLog.Println()
	}

	diskList := []string{}
	for _, each := range bInfo.AllDisks {
		diskList = append(diskList, fmt.Sprintf(each.AbsPath+" - Freespace: "+console.PrintByte(each.Usage.Free(), true)))
	}
	diskList = append(diskList, "Other")
	{
		fmt.Println()

		src := -1

		for src == -1 {
			prompt := &survey.Select{
				Message: "Please choose a " + console.Green("Source") + console.Bold(" to back up from:"),
				Options: diskList,
			}
			survey.AskOne(prompt, &src)
		}
		//source := QueryList(diskList, "Please choose a source to back up from:")

		if src == len(diskList)-1 {
			survey.AskOne(&survey.Confirm{Message: "Proceed in Folder Mode? (Default:N)"}, &bInfo.FolderMode)
			for tempBool := false; !tempBool && !bInfo.FolderMode; {
				survey.AskOne(&survey.Confirm{Message: console.Red("The manually chosen directly MUST contain\n    users that will be backed up\n    DO YOU UNDERSTAND?")}, &tempBool)
			}
			if bInfo.Unix {
				bInfo.Source.AbsPath = backup.SelectMacFolder("Select your source")
			} else {
				bInfo.Source.AbsPath, _, _ = dlgs.File("Choose Source Folder", "", true)
			}
			fi, err := os.Stat(bInfo.Source.AbsPath)
			if !fi.Mode().IsDir() || err != nil {
				console.Error(errors.New("you chose "+bInfo.Source.AbsPath+"\nNot a folder"), 044, false)
			}

			bInfo.Source = backup.GetInformation(bInfo.Source.AbsPath)

		} else {
			bInfo.Source = bInfo.AllDisks[src]
			bInfo.AllDisks = backup.RemoveDisk(bInfo.AllDisks, src)
			diskList = append(diskList[:src], diskList[src+1:]...)
		}
	}

	{
		fmt.Println()

		dst := -1
		for dst == -1 {
			prompt := &survey.Select{
				Message: "Please choose a " + console.Green("Destination") + console.Bold(" to back up to:"),
				Options: diskList,
			}
			survey.AskOne(prompt, &dst)
		}
		//destination := QueryList(diskList, "Please choose a destination to back up to:")

		if dst == len(diskList)-1 {
			for tempBool := false; !tempBool; {
				survey.AskOne(&survey.Confirm{Message: console.Red("The program will create a folder named: " + bInfo.Identifier + "\n    in the chosen folder. \n    DO YOU UNDERSTAND?")}, &tempBool)
			}
			if bInfo.Unix {
				bInfo.Destination.AbsPath = backup.SelectMacFolder("Select your Destination")
			} else {
				bInfo.Destination.AbsPath, _, _ = dlgs.File("Choose Destination Folder", "", true)
			}
			fi, err := os.Stat(bInfo.Destination.AbsPath)
			if !fi.Mode().IsDir() || err != nil {
				console.Error(errors.New("you chose "+bInfo.Destination.AbsPath+"\nNot a folder"), 045, false)
			}

			bInfo.Destination.Name = bInfo.Destination.AbsPath
		} else {
			bInfo.Destination = bInfo.AllDisks[dst]
			bInfo.AllDisks = backup.RemoveDisk(bInfo.AllDisks, dst)
			diskList = append(diskList[:dst], diskList[dst+1:]...)
		}
		fmt.Println()

	}

	bInfo.MainLog.Println("Source:", bInfo.Source.Name)
	bInfo.MainLog.Println("Destination:", bInfo.Destination.Name)
	bInfo.MainLog.Println()

	bInfo.GetUsers()

	if len(bInfo.AllUsers) < 1 {
		bInfo.ErrLog.Println("No users found!")
		console.Error(errors.New("found no users in the selected source "+bInfo.Source.AbsPath), 046, false)
	}

	bInfo.MainLog.Println("All Users:")

	//fmt.Println(console.Cyan("This is for information only. Select users in the prompted window."))
	for _, eachUser := range bInfo.AllUsers {
		bInfo.MainLog.Println(eachUser.Name)
		bInfo.MainLog.Println(" - Approx Size (/1024) :" + console.PrintByte(eachUser.Size, true))
		bInfo.MainLog.Println(" - Approx Size (/1000) :" + console.PrintShortByte(eachUser.Size, true))
		bInfo.MainLog.Println(" - Size in Bytes:" + console.PrintByte(eachUser.Size, false, "b"))
		//fmt.Println(console.HIGreen(eachUser.Name) + console.Cyan(" -- "+console.PrintByte(eachUser.Size, true)))
	}

	console.Head()

	if !bInfo.FolderMode && len(bInfo.AllUsers) > 1 {
		userString := []string{}
		for _, eachUser := range bInfo.AllUsers {
			userString = append(userString, eachUser.Name+" -- "+console.PrintByte(eachUser.Size, true))
		}
		usersToBackup := []int{}
		for len(usersToBackup) == 0 {
			console.Head()
			prompt := &survey.MultiSelect{
				Message: "What users do you want to back up?",
				Options: userString,
			}
			survey.AskOne(prompt, &usersToBackup)
		}
		for _, selected := range usersToBackup {
			bInfo.SelectedUsers = append(bInfo.SelectedUsers, bInfo.AllUsers[selected])
		}
	} else if !bInfo.FolderMode && len(bInfo.AllUsers) == 1 {
		check := false
		survey.AskOne(&survey.Confirm{Message: "Only one user was found. They are:\n  " +
			bInfo.AllUsers[0].Name + " -- " + console.PrintByte(bInfo.AllUsers[0].Size, true) + "\n  " +
			"Back them up?"}, &check)

		if !check {
			console.Error(errors.New("You selected to skip "+bInfo.AllUsers[0].Name), 001, false)
		}
		bInfo.SelectedUsers = bInfo.AllUsers
	} else if bInfo.FolderMode {
		bInfo.SelectedUsers = bInfo.AllUsers
	} else {
		err := errors.New("error occurred in setting up users for choice")
		console.Error(err, 211, false)
	}

	console.Head()
	bInfo.MainLog.Println("Chosen users to back up:")
	for _, each := range bInfo.SelectedUsers {
		bInfo.MainLog.Println("  ", each.Name+" "+console.PrintByte(each.Size, false, "b")+" bytes")
	}

	bInfo.GoBackup()

	console.Head()
	bInfo.VerifySizes()

	dir, _ := os.Getwd()

	backup.CopyFile(filepath.Join(dir, bInfo.Identifier+"_MAIN.log"), filepath.Join(bInfo.Destination.AbsPath, bInfo.Identifier+"_MAIN.log"))
	backup.CopyFile(filepath.Join(dir, bInfo.Identifier+"_ERROR.log"), filepath.Join(bInfo.Destination.AbsPath, bInfo.Identifier+"_ERROR.log"))
	backup.CopyFile(filepath.Join(dir, bInfo.Identifier+"_FILE.log"), filepath.Join(bInfo.Destination.AbsPath, bInfo.Identifier+"_FILE.log"))
	backup.CopyFile(filepath.Join(dir, bInfo.Identifier+"_END-REPORT.log"), filepath.Join(bInfo.Destination.AbsPath, bInfo.Identifier+"_END-REPORT.log"))

	fmt.Println(console.Cyan("Program has completed. Please verify the numbers listed."))
	fmt.Println(console.Magenta("  Note:"), "I have not touched the root. Please back that up manually.")
	console.ExitWait()

	os.Exit(0)
}
