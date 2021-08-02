package main

import (
	"coursera_task3/model"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

var dataPool = sync.Pool{
	New: func() interface{} {
		return new(model.User)
	},
}



// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(out, "found users:")
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	seenBrowsers := make(map[string]bool, 0)
	uniqueBrowsers := 0

	lines := strings.Split(string(fileContents), "\n")

	users := make([]string, 0)
	i := -1
	for _, line := range lines {
		user := dataPool.Get().(*model.User)
		err := user.UnmarshalJSON([]byte(line))
		if err != nil {
			panic(err)
		}
		dataPool.Put(user)
		browsers := user.Browsers

		var userAndroid = false
		var userMSIE = false
		for _, browser := range browsers {
			isAndroid := false
			isMSIE := false
			isBrowserSeen, ok := seenBrowsers[browser]
			if !ok {
				isAndroid = strings.Contains(browser, "Android")
				isMSIE = strings.Contains(browser, "MSIE")
				if isAndroid || isMSIE {
					seenBrowsers[browser] = true
					uniqueBrowsers++
				}
			} else {
				if isBrowserSeen {
					isAndroid = strings.Contains(browser, "Android")
					isMSIE = strings.Contains(browser, "MSIE")
				}
			}
			if !userAndroid {
				userAndroid = isAndroid
			}
			if !userMSIE {
				userMSIE = isMSIE
			}
		}
		i++
		if !userAndroid || !userMSIE {
			continue
		}
		email := strings.Split(user.Email, "@")
		users = append(users, fmt.Sprintf("[%d] %s <%s [at] %s>", i, user.Name, email[0], email[1]))
	}

	fmt.Fprintln(out, strings.Join(users, "\n"))
	fmt.Fprintln(out, "\nTotal unique browsers", uniqueBrowsers)
}
