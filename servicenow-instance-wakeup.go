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
	var configFile string
	userDetails := &User{}

	flag.StringVar(&userDetails.Username, "username", "", "write the username/email with which you are logging in to the developers account")
	flag.StringVar(&userDetails.Password, "password", "", "write the password with which you are logging in to the developers account")
	flag.BoolVar(&userDetails.ChromeHeadless, "headless", false, "bool, if we need headless mode with chrome or not, default:false")
	flag.BoolVar(&userDetails.Debug, "debug", false, "bool, if you want debug output or not, default:false")
	flag.StringVar(&configFile, "config", "", "Provide the config file name, it can be a relative path or a full path, e.g. /home/user/servicenow-config.json or just simply 'config.json'")
	flag.Parse()

	// Read config into struct if exists
	if configFile != "" {
		log.Println("Your flags will be ignored and replaced by the values in the config file you specified...")
		log.Printf("Loading config file under the path [%s]", configFile)
		userDetails = readConfig(configFile)
	}

	if userDetails == nil || len(userDetails.Username) == 0 || len(userDetails.Password) == 0 {
		log.Println("No username or password provided. Use the -username and -password flags to set the username or password. e.g. program -username user@email.tld or setup a config.json with the details")
		os.Exit(1)
	}

	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.DisableGPU,
	}

	log.Printf("Starting the app with debug=%t/headless=%t/account=%s", userDetails.Debug, userDetails.ChromeHeadless, userDetails.Username)

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
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
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
		chromedp.Click(`#usernameSubmitButton`, chromedp.ByID),
		chromedp.WaitVisible(`#password`),
		chromedp.SendKeys(`#password`, password, chromedp.ByID),
		chromedp.Click(`#submitButton`, chromedp.ByID),
		chromedp.WaitVisible(`#instanceWakeUpBtn`, chromedp.ByID),
		chromedp.Click(`#instanceWakeUpBtn`, chromedp.ByID),
		chromedp.WaitNotVisible(`#dp-instance-hib-overlay`, chromedp.ByID),
	}
}

// Read the config file if required and load the json to the struct
func readConfig(config string) *User {
	// Load the specified config file from the path provided
	jsonFile, err := os.Open(config)

	if err != nil {
		log.Fatal(err)
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
