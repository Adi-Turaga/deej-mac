package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"bufio"

	"github.com/andybrewer/mack"
	"go.bug.st/serial"
)

func convertToInt(volumeString []string) []int {

	// Takes string array of volumes and converts them into an array of ints

	var volumes []int
	for _, i := range volumeString {
		j, err := strconv.Atoi(i)
		if err != nil {
			fmt.Println("Unable to parse...moving to next string")
		}
		volumes = append(volumes, j)
	}
	return volumes
}

func checkChange(vol1 []int, vol2 []int) bool {

	// Checks if any of the volumes changed

	for i := 0; i < len(vol1); i++ {
		if int(math.Abs(float64(vol1[i]-vol2[i]))) > 10 {
			return true
		}
	}
	return false
}

func setVolume(volumes []int, apps []string) {

	// Sets volume of the applications

	var script []string
	for i := 0; i < len(apps); i++ {
		fmt.Println("Volume:", volumes[i])
		command := fmt.Sprintf("set vol of (a reference to (the first audio application whose bundleID is equal to \"%s\")) to %d", apps[i], int(volumes[i]/10))
		script = append(script, command)
	}
	if len(script) > 0 {
		fullScript := strings.Join(script, "\n")
		mack.Tell("Background Music", fullScript)
	}
}

func main() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found")
	}
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}
	mode := &serial.Mode{
		BaudRate: 9600,
	}
	port, err := serial.Open("/dev/tty.usbserial-1430", mode) // CHANGE SERIAL PORT
	if err != nil {
		log.Fatal(err)
	}

	// apps to be used (use 4 of choice) (how to format in README)
	apps := []string{
		"com.microsoft.VSCode",
		"com.google.Chrome",
		"com.apple.Music",
		"com.apple.Safari",
	}

	var prev_volumes []int

	scanner := bufio.NewScanner(port)

	for scanner.Scan() {

		fmt.Println(scanner.Text())

		volumeString := strings.Split(scanner.Text(), "|")
		volumes := convertToInt(volumeString)

		/*fmt.Println(prev_volumes, volumes)
		fmt.Println(checkChange(prev_volumes, volumes))*/

		if checkChange(prev_volumes, volumes) {
			setVolume(volumes, apps)
		}

		prev_volumes = volumes

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

}
