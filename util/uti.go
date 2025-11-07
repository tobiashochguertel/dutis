package util

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

type Uti struct {
	Name       string
	Path       string
	Identifier string
}

var kMDItemCFBundleIdentifierPattern = regexp.MustCompile(`kMDItemCFBundleIdentifier\s+=\s+"(.+)"`)
var kMDItemContentTypePattern = regexp.MustCompile(`kMDItemContentType\s+=\s+"(.+)"`)

func ListUti(path string) map[string]Uti {
	files, err := os.ReadDir(path)
	r := make(map[string]Uti)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan Uti)
	wg := &sync.WaitGroup{}
	wg.Add(len(files))

	for _, file := range files {
		go func(file os.DirEntry, wg *sync.WaitGroup) {
			defer wg.Done()

			fp := path + "/" + file.Name()
			cmd := exec.Command("mdls", "-name", "kMDItemCFBundleIdentifier", fp)
			out, err := cmd.Output()
			if err != nil {
				log.Fatal(err)
			}
			match := kMDItemCFBundleIdentifierPattern.FindStringSubmatch(string(out))
			if len(match) > 0 {
				c <- Uti{file.Name(), fp, match[1]}
			}
		}(file, wg)
	}

	go func(group *sync.WaitGroup) {
		wg.Wait()
		close(c)
	}(wg)

	for v := range c {
		r[v.Name] = v
	}
	return r
}

func ListApplicationsUti() map[string]Uti {
	return ListUti("/Applications")
}

func SetDefaultApplication(uti string, suffix string) error {
	fmt.Println("Set default application for", suffix, "to", uti)
	cmd := exec.Command("duti", "-s", uti, suffix, "all")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("duti error: %w, output: %s", err, string(output))
	}
	return nil
}

func getFileContentType(path string) string {
	cmd := exec.Command("mdls", "-name", "kMDItemContentType", path)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	match := kMDItemContentTypePattern.FindStringSubmatch(string(out))
	return match[1]
}

func cleanApplicationPath(path string) string {
	// Remove file:// prefix
	path = strings.TrimPrefix(path, "file://")
	
	// URL decode (fixes %20, etc.)
	decoded, err := url.QueryUnescape(path)
	if err == nil {
		path = decoded
	}
	
	// Remove trailing slash
	path = strings.TrimSuffix(path, "/")
	
	// Extract just the app name from various system paths
	if strings.Contains(path, "/") {
		// Common prefixes to remove
		prefixes := []string{
			"/Applications/",
			"/System/Applications/",
			"/System/Volumes/Preboot/Cryptexes/App/System/Applications/",
			"/Setapp/",
		}
		
		cleaned := false
		for _, prefix := range prefixes {
			if strings.HasPrefix(path, prefix) {
				path = strings.TrimPrefix(path, prefix)
				cleaned = true
				break
			}
		}
		
		// For any remaining paths (user directories, .cache, etc.), extract just the app name
		if !cleaned && strings.Contains(path, "/") {
			parts := strings.Split(path, "/")
			for i := len(parts) - 1; i >= 0; i-- {
				if strings.HasSuffix(parts[i], ".app") {
					// Found the app, include parent folder if it's a utility folder
					if i > 0 && (parts[i-1] == "Utilities" || parts[i-1] == "TeX") {
						path = parts[i-1] + "/" + parts[i]
					} else {
						path = parts[i]
					}
					break
				}
			}
		}
	}
	
	return path
}

func LSCopyAllRoleHandlersForContentType(suf string) []string {
	// Check cache first
	if cached, ok := LoadRecommendedAppsCache(suf); ok {
		return cached
	}

	contentFile, _ := os.CreateTemp("/tmp", "dutis-content.*"+suf)
	defer func(name string) {
		os.Remove(name)
	}(contentFile.Name())
	contentFileContentType := getFileContentType(contentFile.Name())

	scriptFile, _ := os.CreateTemp("/tmp", "dutis-script.*.swift")
	defer func(name string) {
		os.Remove(name)
	}(scriptFile.Name())

	if _, err := scriptFile.Write([]byte(`
import CoreServices
import Foundation

let args = CommandLine.arguments
guard args.count > 1 else {
    print("Missing argument")
    exit(1)
}

let fileType = args[1]

guard let bundleIds = LSCopyAllRoleHandlersForContentType(fileType as CFString, LSRolesMask.all)  else {
    print("Failed to fetch bundle Ids for specified filetype")
    exit(1)
}

(bundleIds.takeRetainedValue() as NSArray)
    .compactMap { bundleId -> NSArray? in
        guard let retVal = LSCopyApplicationURLsForBundleIdentifier(bundleId as! CFString, nil) else { return nil }
        return retVal.takeRetainedValue() as NSArray
    }
    .flatMap { $0 }
    .forEach { print($0) }
`)); err != nil {
		return []string{}
	}

	cmd := exec.Command("swift", scriptFile.Name(), contentFileContentType)
	out, err := cmd.Output()
	if err != nil {
		return []string{}
	}
	
	applicationFullPathList := strings.Split(string(out), "\n")
	var cleanedList []string
	seen := make(map[string]bool)
	
	for _, path := range applicationFullPathList {
		if path == "" {
			continue
		}
		cleaned := cleanApplicationPath(path)
		if cleaned == "" {
			continue
		}
		// Deduplicate
		if !seen[cleaned] {
			seen[cleaned] = true
			cleanedList = append(cleanedList, cleaned)
		}
	}
	
	// Save to cache
	_ = SaveRecommendedAppsCache(suf, cleanedList)
	
	return cleanedList
}
