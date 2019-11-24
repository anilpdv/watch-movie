package main

import (
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

const key = "6l1fPOMWd42gn7A4"
const secret_key = "oa7fqt4ubocfydfvifr9zcwo4umsih"

func format(video_id string, ipaddress string) string {
	return fmt.Sprintf("https://videospider.in/getticket.php?key=%s&secret_key=%s&video_id=%s&ip=%s", key, secret_key, video_id, ipaddress)
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

func getUrl(u string, video_id string) string {
	resp, err := http.Get(u)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	ticket, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("https://videospider.stream/getvideo?key=%s&video_id=%s&ticket=%s", key, video_id, string(ticket))

}

func getIpaddress() string {

	ipUrl := "https://api.ipify.org?format=text"

	resp, err := http.Get(ipUrl)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(ip)
}

func main() {

	app := &cli.App{
		Name:  "watch",
		Usage: "fight the lonliness with fun!",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "id",
				Value: "",
				Usage: "imdb id, to generate movie link.",
			},
		},
		Action: func(c *cli.Context) error {
			name := "Nefertiti"
			if c.NArg() > 0 {
				name = c.Args().Get(0)
			}

			ip := getIpaddress()
			url := format(name, ip)
			fmt.Println(url)
			streamurl := getUrl(url, name)
			openbrowser(streamurl)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
