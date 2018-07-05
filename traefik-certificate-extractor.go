package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/fsnotify/fsnotify"
	"path/filepath"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	acmeJSONFile := flag.String("acmejson", "", "path of the acme.json-file")
	targetDir := flag.String("target", "", "directory where the certificates should be extracted to")
	watch := flag.Bool("watch", false, "should the extractor-tool keep watching the acme.json-file and rewrite the certificates")
	flag.Parse()

	if *acmeJSONFile == "" || *targetDir == "" {
		fmt.Print("you must specify -acmejson and -target\n")
		os.Exit(1)
	}

	extractCertsFromAcmeJSON(*acmeJSONFile, *targetDir)

	if *watch {
		watchAndExtractCertsFromAcmeJSON(*acmeJSONFile, *targetDir)
	}
}

func watchAndExtractCertsFromAcmeJSON(acmeJSONFile string, targetDir string) {
	watcher, err := fsnotify.NewWatcher()
	check(err)

	defer watcher.Close()

	acmeJSONFileAbs, err := filepath.Abs(acmeJSONFile)
	check(err)

	updateExtract := make(chan bool)

	go func() {
		var timer *time.Timer

		for {
			select {
			case event := <-watcher.Events:
				changedFileAbs, err := filepath.Abs(event.Name)
				check(err)

				if acmeJSONFileAbs != changedFileAbs {
					continue
				}

				if (event.Op&fsnotify.Write == fsnotify.Write) || (event.Op&fsnotify.Create == fsnotify.Create) {
					if timer != nil {
						timer.Stop()
					}
					timer = time.AfterFunc(time.Second*1, func() {
						updateExtract <- true
					})
				}

			case err := <-watcher.Errors:
				fmt.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(path.Dir(acmeJSONFile))
	check(err)

	for range updateExtract {
		fmt.Print("-- detected acme.json changes, gegenerating certificates\n")
		extractCertsFromAcmeJSON(acmeJSONFile, targetDir)
	}
}

func extractCertsFromAcmeJSON(acmeJSONFile string, targetDir string) {
	certificates := unmarshalAcmeJSON(acmeJSONFile)

	for _, cert := range certificates.Certificates {
		var err error

		fmt.Printf("%s\n", formatDomainName(cert.Domain))

		certTargetDir := path.Join(targetDir, cert.Domain.Main)

		err = os.MkdirAll(certTargetDir, 0700)
		check(err)

		extractCert(cert, certTargetDir)

		for _, san := range cert.Domain.SANs {
			sanSymlinkName := path.Join(targetDir, san)
			os.Symlink(cert.Domain.Main, sanSymlinkName)
		}
	}

	fmt.Print("--- done\n")
}

func extractCert(certificate *Certificate, targetDir string) {
	ioutil.WriteFile(path.Join(targetDir, "fullchain"), certificate.Certificate, 0600)
	ioutil.WriteFile(path.Join(targetDir, "privkey"), certificate.Key, 0600)
	ioutil.WriteFile(path.Join(targetDir, "all"), append(certificate.Key, certificate.Certificate...), 0600)
}

func formatDomainName(domain Domain) string {
	sans := ""
	if len(domain.SANs) > 0 {
		sans = " (" + strings.Join(domain.SANs, ", ") + ")"
	}

	return domain.Main + sans
}

func unmarshalAcmeJSON(acmeJSONFile string) Certificates {
	data, err := ioutil.ReadFile(acmeJSONFile)
	check(err)

	var certificates Certificates
	err = json.Unmarshal(data, &certificates)
	check(err)

	return certificates
}
