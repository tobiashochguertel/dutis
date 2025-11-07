package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/mrtkrcm/dutis/util"
	"os"
	"strings"
)

var utiMap = util.ListApplicationsUti()

const YouSelectPrompt = "You selected "

func chooseUti() string {
	fmt.Println("Please input uti.(Tab for auto complement)")

	promptHandler := func(d prompt.Document) []prompt.Suggest {
		var p []prompt.Suggest
		for _, v := range utiMap {
			p = append(p, prompt.Suggest{Text: v.Name, Description: "uti: " + v.Identifier})
		}
		return prompt.FilterHasPrefix(p, d.GetWordBeforeCursor(), true)
	}

	t := inputWithDoubleCtrlC("> ", promptHandler)
	if t == "" {
		return ""
	}
	fmt.Println(YouSelectPrompt + t)
	return t
}

func chooseSuffix() string {
	fmt.Println("Please input suffix.(Tab for auto complement)")
	t := inputWithDoubleCtrlC("> ", util.SuffixCompleter)
	if t == "" {
		return ""
	}
	fmt.Println(YouSelectPrompt + t)
	return t
}

func choosePreset() {
	fmt.Println("Please input preset.(Tab for auto complement)")
	t := inputWithDoubleCtrlC("> ", util.PresetCompleter)
	if t != "" {
		fmt.Println(YouSelectPrompt + t)
	}
}

func inputWithDoubleCtrlC(prefix string, completer prompt.Completer) string {
	consecutiveInterrupts := 0

	p := prompt.New(
		func(s string) {
			// Input executor - this is called when Enter is pressed
			consecutiveInterrupts = 0
		},
		completer,
		prompt.OptionPrefix(prefix),
		prompt.OptionAddASCIICodeBind(
			prompt.ASCIICodeBind{
				ASCIICode: []byte{3}, // Ctrl+C (ASCII code 3)
				Fn: func(buf *prompt.Buffer) {
					consecutiveInterrupts++
					if consecutiveInterrupts >= 2 {
						fmt.Println("\nAre you sure you want to exit? (yes/no)")
						confirm := prompt.Input("> ", func(d prompt.Document) []prompt.Suggest {
							s := []prompt.Suggest{
								{Text: "yes", Description: "Exit the application"},
								{Text: "no", Description: "Continue"},
							}
							return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
						})
						if confirm == "yes" {
							os.Exit(0)
						}
						consecutiveInterrupts = 0
					} else {
						// Clear the buffer and show message
						buf.DeleteBeforeCursor(len(buf.Document().TextBeforeCursor()))
						fmt.Println("\nPress Ctrl+C again to exit, or continue typing.")
					}
				},
			},
		),
	)

	return p.Input()
}

func printRecommend(suf string) {
	pmp := strings.Repeat("=", 10) + " Recommended applications " + strings.Repeat("=", 10)
	fmt.Println(pmp)
	if recommendApplications := util.LSCopyAllRoleHandlersForContentType(suf); len(recommendApplications) > 0 {
		for _, n := range util.LSCopyAllRoleHandlersForContentType(suf) {
			fmt.Println(n)
		}
	} else {
		fmt.Println("No recommend applications")
	}
	fmt.Println(strings.Repeat("=", len(pmp)))
}

func main() {
	util.InstallDeps()
	//fmt.Println("Please select mode by number.(Tab for auto complement)\n(1). change default application by suffix\n(2).
	// change default application by preset")
	//t := prompt.Input("> ", mainCompleter)
	//fmt.Println("You selected " + t)
	t := "1"
	var suf string
	switch t {
	case "1":
		suf = chooseSuffix()
	case "2":
		choosePreset()
	}

	if suf == "" {
		return
	}
	printRecommend(suf)

	utiName := chooseUti()
	if utiName == "" {
		return
	}
	if utiItem, ok := utiMap[utiName]; ok {
		util.SetDefaultApplication(utiItem.Identifier, suf)
	} else {
		fmt.Printf("uti %s not found\n", utiName)
	}
}
