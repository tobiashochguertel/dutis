package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/tobiashochguertel/dutis/util"
	"os"
	"strings"
	"sync"
)

const (
	Version    = "v0.2.1-fork"
	Repository = "https://github.com/tobiashochguertel/dutis"
)

var (
	utiMap               map[string]util.Uti
	utiMapOnce           sync.Once
	consecutiveInterrupts = 0
)

func getUtiMap() map[string]util.Uti {
	utiMapOnce.Do(func() {
		if cached, ok := util.LoadUtiCache(); ok {
			fmt.Println("\033[2;37m(using cached application data)\033[0m")
			utiMap = cached
		} else {
			fmt.Println("\033[2;37m(scanning applications...)\033[0m")
			utiMap = util.ListApplicationsUti()
			_ = util.SaveUtiCache(utiMap)
		}
	})
	return utiMap
}

const YouSelectPrompt = "You selected "

func chooseUti() string {
	fmt.Println("Please input uti.(Tab for auto complement)")

	promptHandler := func(d prompt.Document) []prompt.Suggest {
		var p []prompt.Suggest
		for _, v := range getUtiMap() {
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
	p := prompt.New(
		func(s string) {
			consecutiveInterrupts = 0
		},
		completer,
		prompt.OptionPrefix(prefix),
		prompt.OptionAddASCIICodeBind(
			prompt.ASCIICodeBind{
				ASCIICode: []byte{3},
				Fn: func(buf *prompt.Buffer) {
					consecutiveInterrupts++
					if consecutiveInterrupts >= 2 {
						fmt.Println("\nExiting...")
						os.Exit(0)
					} else {
						buf.DeleteBeforeCursor(len(buf.Document().TextBeforeCursor()))
						fmt.Println("\nPress Ctrl+C again to exit")
					}
				},
			},
		),
	)

	return p.Input()
}

func printRecommend(suf string) {
	fmt.Printf("\n\033[1;35m%s Recommended Applications %s\033[0m\n", 
		strings.Repeat("─", 10), strings.Repeat("─", 10))
	
	recommendApplications := util.LSCopyAllRoleHandlersForContentType(suf)
	if len(recommendApplications) > 0 {
		fmt.Printf("\033[2;37mFound %d application(s) for %s files:\033[0m\n\n", 
			len(recommendApplications), suf)
		
		for i, app := range recommendApplications {
			if app == "" {
				continue
			}
			// Use different colors for variety
			color := "\033[0;36m" // Cyan
			if i%2 == 1 {
				color = "\033[0;34m" // Blue
			}
			fmt.Printf("  %s• %s\033[0m\n", color, app)
		}
	} else {
		fmt.Println("\033[2;33m  No recommended applications found\033[0m")
	}
	
	fmt.Printf("\033[1;35m%s\033[0m\n\n", strings.Repeat("─", 46))
}

func printVersion() {
	fmt.Printf("\033[1;36m%s (%s)\033[0m\n", Version, Repository)
	
	// Show cache status
	if _, ok := util.LoadUtiCache(); ok {
		fmt.Printf("\033[2;32m✓ Cache ready\033[0m\n")
	} else {
		fmt.Printf("\033[2;33m○ Will build cache on first use\033[0m\n")
	}
	fmt.Println()
}

func showHelp() {
	fmt.Println("Usage: dutis [command]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  (none)              Interactive mode to set file associations")
	fmt.Println("  apply               Apply all configured associations from config")
	fmt.Println("  list                List all configured associations")
	fmt.Println("  remove <suffix>     Remove association for a suffix")
	fmt.Println("  version, -v         Show version information")
	fmt.Println("  --refresh-cache     Refresh the application cache")
	fmt.Println("  help, --help, -h    Show this help message")
	fmt.Println()
	fmt.Println("Config file: ~/.dutis/config.yaml")
}

func handleCommands() bool {
	if len(os.Args) < 2 {
		return false
	}

	command := os.Args[1]

	switch command {
	case "--refresh-cache":
		fmt.Println("Refreshing application cache...")
		utiMap := util.ListApplicationsUti()
		if err := util.SaveUtiCache(utiMap); err != nil {
			fmt.Printf("Error saving cache: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Cache refreshed with %d applications\n", len(utiMap))
		return true

	case "apply":
		config, err := util.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}
		if err := config.ApplyAll(); err != nil {
			os.Exit(1)
		}
		return true

	case "list":
		config, err := util.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}
		associations := config.ListAssociations()
		if len(associations) == 0 {
			fmt.Println("No associations configured yet.")
			fmt.Println("Run 'dutis' to set file associations interactively.")
		} else {
			fmt.Printf("Configured associations (%d):\n\n", len(associations))
			fmt.Printf("%-15s %-30s %s\n", "SUFFIX", "APPLICATION", "BUNDLE ID")
			fmt.Println(strings.Repeat("-", 80))
			for _, assoc := range associations {
				fmt.Printf("%-15s %-30s %s\n", assoc.Suffix, assoc.Application, assoc.BundleID)
			}
		}
		return true

	case "remove":
		if len(os.Args) < 3 {
			fmt.Println("Error: suffix required")
			fmt.Println("Usage: dutis remove <suffix>")
			os.Exit(1)
		}
		suffix := os.Args[2]
		config, err := util.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}
		if _, ok := config.GetAssociation(suffix); !ok {
			fmt.Printf("No association found for suffix: %s\n", suffix)
			os.Exit(1)
		}
		if err := config.RemoveAssociation(suffix); err != nil {
			fmt.Printf("Error removing association: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Removed association for: %s\n", suffix)
		return true

	case "version", "--version", "-v":
		fmt.Printf("%s\n", Version)
		fmt.Printf("Repository: %s\n", Repository)
		return true

	case "help", "--help", "-h":
		showHelp()
		return true
	}

	return false
}

func main() {
	if handleCommands() {
		return
	}

	printVersion()
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
	if utiItem, ok := getUtiMap()[utiName]; ok {
		if err := util.SetDefaultApplication(utiItem.Identifier, suf); err != nil {
			fmt.Printf("Error setting default application: %v\n", err)
			return
		}
		
		// Save to config
		config, err := util.LoadConfig()
		if err != nil {
			fmt.Printf("Warning: Could not load config: %v\n", err)
		} else {
			if err := config.AddAssociation(suf, utiItem.Name, utiItem.Identifier); err != nil {
				fmt.Printf("Warning: Could not save to config: %v\n", err)
			} else {
				fmt.Printf("\033[2;32m✓ Saved to config (~/.dutis/config.yaml)\033[0m\n")
			}
		}
	} else {
		fmt.Printf("uti %s not found\n", utiName)
	}
}
