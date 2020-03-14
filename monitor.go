package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const timesMonitoring = 2
const delay = 5 * time.Second

func main() {

	fmt.Println("Website monitor")

	for {
		displayOptions()
		command := getCommandInput()
		fmt.Println("")

		switch command {
		case 1:
			monitor()
		case 2:
			displayLogs()
		case 0:
			fmt.Println("Exiting program")
			os.Exit(0)
		default:
			fmt.Println("Unknown command")
			os.Exit(-1)
		}
	}

}

func monitor() {
	fmt.Println("Monitoring websites...")
	sites := readWebsites()

	for i := 0; i < timesMonitoring; i++ {
		for _, site := range sites {
			testWebsite(site)
		}
		time.Sleep(delay)
		fmt.Println("")
	}

	fmt.Println("")
}

func testWebsite(site string) {
	resp, err := http.Get(site)
	if err != nil {
		fmt.Println("Unexpected error:", err)
	}
	if resp != nil && resp.StatusCode == 200 {
		fmt.Println(site, "is UP")
		registerLog(site, true)

	} else {
		fmt.Println(site, "is DOWN")
		registerLog(site, false)
	}
}

func readWebsites() []string {
	var websites []string
	file, err := os.Open("websites.txt")

	if err != nil {
		fmt.Println("Unexpected error:", err)
		return websites
	}

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		line = strings.TrimSpace(line)
		websites = append(websites, line)
	}

	file.Close()
	return websites
}

func displayOptions() {
	fmt.Println("1 - Monitor websites")
	fmt.Println("2 - Display logs")
	fmt.Println("0 - Exit")
}

func getCommandInput() int {
	var command int
	fmt.Scan(&command)
	return command
}

func registerLog(site string, status bool) {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Unexpected error:", err)
	}
	delimiter := " - "
	now := time.Now().Format("02/01/2006 15:04:05")
	file.WriteString(now + delimiter + site + delimiter + " online: " + strconv.FormatBool(status) + "\n")

	file.Close()
}

func displayLogs() {
	fmt.Println("Displaying logs..." + "\n")
	file, err := ioutil.ReadFile("log.txt")
	if err != nil {
		fmt.Println("Unexpected error:", err)
	}
	fmt.Println(string(file))
}
