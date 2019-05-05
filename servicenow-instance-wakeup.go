package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"os"
	"time"
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

	log.Println("Loading config file...")

	// Read config into struct if exists
	userDetails = readConfig()

	if userDetails == nil || len(userDetails.Username) == 0 || len(userDetails.Password) == 0 {
		log.Println("No username or password provided. Use the -username and -password flags to set the username or password. e.g. program -username user@email.tld or setup a config.json with the details")
		os.Exit(1)
	}

	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.DisableGPU,
	}

	log.Printf("Starting app with debug=%t and headless=%t", userDetails.Debug, userDetails.ChromeHeadless)

	// navigate to a page, wait for an element, click
	if !userDetails.Debug {
		log.SetOutput(ioutil.Discard)
	}

	if userDetails.ChromeHeadless {
		opts = append(opts, chromedp.Headless)
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	// run task list
	err = chromedp.Run(ctx, wakeUpInstance(userDetails.Username, userDetails.Password))

	if err != nil {
		log.Fatal(err)
	}
}

func wakeUpInstance(username string, password string) chromedp.Tasks {
	return chromedp.Tasks{
		// This is the url which you are redirected to
		// after opening an inactive instance https://developer.servicenow.com/app.do#!/instance?wu=true
		chromedp.Navigate(`https://developer.servicenow.com/ssologin.do?relayState=%2Fapp.do%23%21%2Finstance%3Fwu%3Dtrue`),
		chromedp.WaitVisible(`.logo`),
		chromedp.SendKeys(`#username`, username, chromedp.ByID),
		chromedp.SendKeys(`#password`, password, chromedp.ByID),
		chromedp.Click(`#submitButton`, chromedp.ByID),
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
