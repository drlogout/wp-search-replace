package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

var sequence = []string{"http", "https"}
var reverseSequence = []string{"https", "http"}

func main() {
	args := os.Args[1:]

	if len(args) != 2 {
		fmt.Println(os.Args[0] + " <search> <replace>")
		os.Exit(1)
	}

	replace, err := parse(args[1])
	if err != nil {
		log.Fatal(err)
	}

	search, err := parse(args[0])
	if err != nil {
		log.Fatal(err)
	}

	//  remove paths with subdomain
	err = wpSearchReplace("http://www."+search.Host, "http://"+search.Host)
	if err != nil {
		log.Fatal(err)
	}
	err = wpSearchReplace("https://www."+search.Host, "https://"+search.Host)
	if err != nil {
		log.Fatal(err)
	}

	if replace.Scheme == "https" {
		sequence = reverseSequence
	}

	for _, scheme := range sequence {
		search.Scheme = scheme
		err = wpSearchReplace(search.String(), replace.String())
		if err != nil {
			log.Fatal(err)
		}
	}

}

func parse(s string) (*url.URL, error) {
	if !strings.HasPrefix(s, "http") {
		s = "http://" + s
	}

	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(u.Host, "www.") {
		split := strings.Split(u.String(), ".")
		u, err = url.Parse(u.Scheme + "://" + strings.Join(split[1:], "."))
		if err != nil {
			return nil, err
		}
	}

	return u, nil
}

func wpSearchReplace(s, r string) error {
	stdout, stderr := bytes.Buffer{}, bytes.Buffer{}

	cmd := exec.Command("wp", "search-replace", "--allow-root", "--all-tables", "--precise", s, r)
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()

	fmt.Println()
	fmt.Println("Search:", s, "Replace:", r)
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if err != nil {
		fmt.Println(errStr)
		return err
	}

	if outStr != "" {
		scanner := bufio.NewScanner(strings.NewReader(outStr))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "Success") {
				fmt.Println(line)
			}
		}
	}

	return nil
}
