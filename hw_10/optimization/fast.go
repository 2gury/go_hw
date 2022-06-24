package main

import (
	"bufio"
	"bytes"
	"fmt"
	users "go_hw/hw_9/optimization/models"
	"io"
	"log"
	"os"
	"strings"
)

type User struct {
	Browsers []string
	Company  string
	Country  string
	Email    string
	Job      string
	Name     string
	Phone    string
}

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	s := bufio.NewScanner(file)

	seenBrowsers := map[string]interface{}{}
	uniqueBrowsers := 0
	foundUsers := ""

	tmpUser := users.User{}
	i := 0
	for s.Scan() {
		err := tmpUser.UnmarshalJSON(s.Bytes())
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		for _, browser := range tmpUser.Browsers {
			if strings.Contains(browser, "Android") {
				isAndroid = true
				if _, ok := seenBrowsers[browser]; !ok {
					seenBrowsers[browser] = ""
					uniqueBrowsers++
				}
			}
			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				if _, ok := seenBrowsers[browser]; !ok {
					seenBrowsers[browser] = ""
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			i++
			continue
		}

		tmpUser.Email = strings.ReplaceAll(tmpUser.Email, "@", " [at] ")
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, tmpUser.Name, tmpUser.Email)
		i++
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}

func main() {
	out := new(bytes.Buffer)
	FastSearch(out)
	fastResult := out.String()
	log.Println(fastResult)
}
