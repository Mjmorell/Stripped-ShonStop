package console

import "github.com/fatih/color"

var (
	//Red and all others are colors
	Red func(...interface{}) string
	//Green and all others are colors
	Green func(...interface{}) string
	//Yellow and all others are colors
	Yellow func(...interface{}) string
	//Blue and all others are colors
	Blue func(...interface{}) string
	//Magenta and all others are colors
	Magenta func(...interface{}) string
	//Cyan and all others are colors
	Cyan func(...interface{}) string
	//HIRed  and all others are colors
	HIRed func(...interface{}) string
	//HIGreen and all others are colors
	HIGreen func(...interface{}) string
	//HIYellow and all others are colors
	HIYellow func(...interface{}) string
	//HIBlue and all others are colors
	HIBlue func(...interface{}) string
	//HIMagenta and all others are colors
	HIMagenta func(...interface{}) string
	//HICyan and all others are colors
	HICyan func(...interface{}) string
	//Bold and all others are colors
	Bold func(...interface{}) string
	//Italic and all others are colors
	Italic func(...interface{}) string
	//Underline and all others are colors
	Underline func(...interface{}) string
)

func init() {
	Red = color.New(color.FgRed).SprintFunc()
	Green = color.New(color.FgGreen).SprintFunc()
	Yellow = color.New(color.FgYellow).SprintFunc()
	Blue = color.New(color.FgBlue).SprintFunc()
	Magenta = color.New(color.FgMagenta).SprintFunc()
	Cyan = color.New(color.FgCyan).SprintFunc()

	HIRed = color.New(color.FgHiRed).SprintFunc()
	HIGreen = color.New(color.FgHiGreen).SprintFunc()
	HIYellow = color.New(color.FgHiYellow).SprintFunc()
	HIBlue = color.New(color.FgHiBlue).SprintFunc()
	HIMagenta = color.New(color.FgHiMagenta).SprintFunc()
	HICyan = color.New(color.FgHiCyan).SprintFunc()

	Bold = color.New(color.Bold).SprintFunc()
	Italic = color.New(color.Italic).SprintFunc()
	Underline = color.New(color.Underline).SprintFunc()
}
