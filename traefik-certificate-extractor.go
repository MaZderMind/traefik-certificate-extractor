package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"encoding/json"
	"path"

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
	acme_json_file := flag.String("acmejson", "", "path of the acme.json-file")
	target_dir := flag.String("target", "", "directory where the certificates should be extracted to")
	watch := flag.Bool("watch", false, "should the extractor-tool keep watching the acme.json-file and rewrite the certificates")
	flag.Parse()

	if (*acme_json_file == "" || *target_dir == "") {
		fmt.Print("you must specify -acmejson and -target\n")
		os.Exit(1)
	}

	extract_certs_from_acme_json(*acme_json_file, *target_dir)

	if (*watch) {
		watch_and_extract_certs_from_acme_json(*acme_json_file, *target_dir)
	}
}

func watch_and_extract_certs_from_acme_json(acme_json_file string, target_dir string) {
	watcher, err := fsnotify.NewWatcher()
	check(err)

	defer watcher.Close()

	acme_json_file_abs, err := filepath.Abs(acme_json_file)
	check(err)

	update_extract := make(chan bool)

	go func() {
		var timer *time.Timer

		for {
			select {
			case event := <-watcher.Events:
				changed_file_abs, err := filepath.Abs(event.Name)
				check(err)

				if (acme_json_file_abs != changed_file_abs) {
					continue
				}

				if (event.Op & fsnotify.Write == fsnotify.Write) || (event.Op & fsnotify.Create == fsnotify.Create) {
					if timer != nil {
						timer.Stop()
					}
					timer = time.AfterFunc(time.Second * 1, func() {
						update_extract <- true
					})
				}

			case err := <-watcher.Errors:
				fmt.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(path.Dir(acme_json_file))
	check(err)

	for range update_extract {
		fmt.Print("-- detected acme.json changes, gegenerating certificates\n")
		extract_certs_from_acme_json(acme_json_file, target_dir)
	}
}

func extract_certs_from_acme_json(acme_json_file string, target_dir string) {
	account := unmarshal_acme_json(acme_json_file)

	for _, cert := range account.DomainsCertificate.Certs {
		var err error;

		fmt.Printf("%s\n", format_domain_name(cert.Domains))

		cert_target_dir := path.Join(target_dir, cert.Domains.Main)

		err = os.MkdirAll(cert_target_dir, 0700);
		check(err)

		extract_cert(cert.Certificate, cert_target_dir)

		for _, san := range cert.Domains.SANs {
			san_symlink_name := path.Join(target_dir, san)
			os.Symlink(cert.Domains.Main, san_symlink_name)
		}
	}

	fmt.Print("--- done\n")
}

func extract_cert(certificate *Certificate, target_dir string) {
	ioutil.WriteFile(path.Join(target_dir, "fullchain"), certificate.Certificate, 0600)
	ioutil.WriteFile(path.Join(target_dir, "privkey"), certificate.PrivateKey, 0600)
	ioutil.WriteFile(path.Join(target_dir, "all"), append(certificate.PrivateKey, certificate.Certificate...), 0600)
	ioutil.WriteFile(path.Join(target_dir, "url"), []byte(certificate.CertURL), 0600)
}

func format_domain_name(domain Domain) string {
	sans := ""
	if (len(domain.SANs) > 0) {
		sans = " (" + strings.Join(domain.SANs, ", ") + ")"
	}

	return domain.Main + sans;
}

func unmarshal_acme_json(acmejsonfile string) Account {
	data, err := ioutil.ReadFile(acmejsonfile)
	check(err)

	var account Account
	err = json.Unmarshal(data, &account)
	check(err)

	return account
}
