package main

import (
	"context"
	"flag"
	"log"

	"encoding/json"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/client"
	"io/ioutil"
	"os"
)

type User struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	ChromeHeadless bool   `json:"headless"`
	Debug          bool   `json:"debug"`
}

func main() {
	var err error
	userDetails := &User{}

	userDetails.Username = *flag.String("username", "", "write the username/email with which you are loggin in to the developers account")
	userDetails.Password = *flag.String("password", "", "write the password with which you are loggin in to the developers account")
	userDetails.ChromeHeadless = *flag.Bool("headless", false, "bool, if we need headless mode with chrome or not, default:false")
	userDetails.Debug = *flag.Bool("debug", false, "bool, if you want debug output or not, default:false")
	flag.Parse()

	// Read config into struct if exists
	userDetails = readConfig()

	if userDetails == nil || len(userDetails.Username) == 0 || len(userDetails.Password) == 0 {
		log.Println("No username or password provided. Use the -username and -password flags to set the username or password. e.g. program -username user@email.tld or setup a config.json with the details")
		os.Exit(1)
	}

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	//c, err := getInstance(userDetails, ctxt)
	// create chrome instance
	var c *chromedp.CDP

	if !userDetails.Debug {
		log.SetOutput(ioutil.Discard)
	}

	if userDetails.ChromeHeadless {
		c, err = chromedp.New(ctxt, chromedp.WithTargets(client.New().WatchPageTargets(ctxt)), chromedp.WithLog(log.Printf))
	} else {
		c, err = chromedp.New(ctxt, chromedp.WithLog(log.Printf))
	}

	if err != nil {
		log.Fatal(err)
	}

	// run task list
	err = c.Run(ctxt, wakeUpInstance(userDetails.Username, userDetails.Password))
	if err != nil {
		log.Fatal(err)
	}

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

func wakeUpInstance(username string, password string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(`https://developer.servicenow.com/ssologin.do?relayState=%2Fapp.do%23!%2Fdashboard`),
		chromedp.WaitVisible(`#logo`),
		chromedp.SendKeys(`#username`, username, chromedp.ByID),
		chromedp.SendKeys(`#password`, password, chromedp.ByID),
		chromedp.Click(`#submitButton`, chromedp.ByID),
		chromedp.WaitVisible(`#dp-hdr-userinfo-link`, chromedp.ByID),
		chromedp.WaitVisible(`#dp-hdr-br-manage-link`, chromedp.ByID),
		chromedp.Click(`#dp-hdr-br-manage-link`, chromedp.ByID),
		chromedp.WaitVisible(`#dp-hdr-br-link-instance`, chromedp.ByID),
		chromedp.Click(`#dp-hdr-br-link-instance`, chromedp.ByID),
		chromedp.WaitVisible(`#instanceWakeUpBtn`, chromedp.ByID),
		chromedp.Click(`#instanceWakeUpBtn`, chromedp.ByID),
		chromedp.WaitNotVisible(`#dp-instance-hib-overlay`, chromedp.ByID),
	}
}

func readConfig() *User {
	jsonFile, err := os.Open("config.json")

	if err != nil {
		log.Panic(err)
		return nil
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	userInfo := User{}

	err = json.Unmarshal(byteValue, &userInfo)

	if err != nil {
		return nil
	}

	return &userInfo
}

//func getInstance(userDetails *User, ctxt context.Context) (*chromedp.CDP, error) {
//	return c, err
//}
