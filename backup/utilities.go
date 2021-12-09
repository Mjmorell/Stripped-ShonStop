package backup

import (
	"github.com/mjmorell/ShonStop/console"
	"os/exec"
	"regexp"
	"strings"
)

func macSkippables(test string, paranoid bool) bool {
	test = strings.ToLower(test)
	//and .file
	if matched, _ := regexp.MatchString(`^\..*`, strings.ToLower(test)); matched {
		return true

	} else if matched, _ := regexp.MatchString(`^dropbox.*`, strings.ToLower(test)); matched {
		if paranoid {
			return false
		}
		return true

	} else if matched, _ := regexp.MatchString("^onedrive.*", strings.ToLower(test)); matched {
		if paranoid {
			return false
		}
		return true

	} else if matched, _ := regexp.MatchString("^google drive.*", strings.ToLower(test)); matched {
		if paranoid {
			return false
		}
		return true
	} else if matched, _ := regexp.MatchString(`^trash$`, strings.ToLower(test)); matched {
		return true
	} else if matched, _ := regexp.MatchString(`^library$`, strings.ToLower(test)); matched {
		return true
	}
	return false
}

func winSkippables(test string, paranoid bool) bool {
	test = strings.ToLower(test)
	//and .file
	if matched, _ := regexp.MatchString(`^\..*`, strings.ToLower(test)); matched {
		return true

	} else if matched, _ := regexp.MatchString(`^dropbox.*`, strings.ToLower(test)); matched {
		if paranoid {
			return false
		}
		return true

	} else if matched, _ := regexp.MatchString("^onedrive.*", strings.ToLower(test)); matched {
		if paranoid {
			return false
		}
		return true

	} else if matched, _ := regexp.MatchString("^google drive.*", strings.ToLower(test)); matched {
		if paranoid {
			return false
		}
		return true
	} else if matched, _ := regexp.MatchString(`^appdata$`, strings.ToLower(test)); matched {
		return true
	}
	return false
}

//SelectMacFolder prompts user to select a folder with a GUI
func SelectMacFolder(printout string) string {
	o, err := exec.Command("osascript", "-e", `(choose folder with prompt "Choose `+printout+`")`).Output()
	if err != nil {
		console.Error(err, 201, false)
	}

	out := strings.TrimSpace(string(o))

	out = strings.Split(out, "alias ")[1]
	//out = strings.Replace(out, " ", "\\ ", -1)
	out = "/Volumes/" + strings.Replace(out, ":", "/", -1)

	out = strings.TrimSuffix(out, "/")

	return out
}
