package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type User struct {
	Username       string        `json:"username"`
	Password       string        `json:"password"`
	ChromeHeadless bool          `json:"headless"`
	Debug          bool          `json:"debug"`
	Timeout        time.Duration `json:"timeout"`
}

func main() {
	var err error
	var configFile string
	var timeout time.Duration
	userDetails := &User{}

	flag.StringVar(&userDetails.Username, "username", "", "write the username/email with which you are logging in to the developers account")
	flag.StringVar(&userDetails.Password, "password", "", "write the password with which you are logging in to the developers account")
	flag.BoolVar(&userDetails.ChromeHeadless, "headless", false, "bool, if we need headless mode with chrome or not, default:false")
	flag.BoolVar(&userDetails.Debug, "debug", false, "bool, if you want debug output or not, default:false")
	flag.StringVar(&configFile, "config", "", "Provide the config file name, it can be a relative path or a full path, e.g. /home/user/servicenow-config.json or just simply 'config.json'")
	flag.DurationVar(&timeout, "timeout", 60, "Set the timeout after which the app should exit. This is a number in seconds, default:60")
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

	err = wakeUpInstance(ctx, userDetails.Username, userDetails.Password, timeout)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func wakeUpInstance(ctx context.Context, username string, password string, timeout time.Duration) error {
	var cancel func()
	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, timeout*time.Second)
	defer cancel()

	initialURL := "https://developer.servicenow.com/ssologin.do?relayState=%2Fdev.do%23%21%2Fhome"

	fmt.Printf("Navigating to the webpage: %s\n", initialURL)
	// first navigate to the sso login page
	if err := chromedp.Run(ctx, chromedp.Navigate(initialURL)); err != nil {
		return fmt.Errorf("could not navigate to the SSO login page: %v", err)
	} else {
		fmt.Printf("Successfully navigated to the webpage...\n")
	}

	fmt.Printf("Searching for the .logo element...\n")
	if err := chromedp.Run(ctx, chromedp.WaitVisible(`.logo`)); err != nil {
		return fmt.Errorf("could not detect .logo element: %v", err)
	} else {
		fmt.Printf("Found .logo element\n")
	}

	fmt.Printf("Filling out the username field...\n")
	if err := chromedp.Run(ctx, chromedp.SendKeys(`#username`, username, chromedp.ByID)); err != nil {
		return fmt.Errorf("could not fill out the username: %v", err)
	} else {
		fmt.Printf("Filled username field with %s\n", username)
	}

	fmt.Printf("Clicking the next button...\n")
	if err := chromedp.Run(ctx, chromedp.Click(`#usernameSubmitButton`, chromedp.ByID)); err != nil {
		return fmt.Errorf("could not click the next button: %v", err)
	} else {
		fmt.Printf("Clicked Next button\n")
	}

	fmt.Printf("Searching for the password field...\n")
	if err := chromedp.Run(ctx, chromedp.WaitVisible(`#password`)); err != nil {
		return fmt.Errorf("could not detect password element: %v", err)
	} else {
		fmt.Printf("Found password field\n")
	}

	fmt.Printf("Filling out the password field...\n")
	if err := chromedp.Run(ctx, chromedp.SendKeys(`#password`, password, chromedp.ByID)); err != nil {
		return fmt.Errorf("could not fill out the password: %v", err)
	} else {
		fmt.Printf("Filled password field with your password ******\n")
	}

	fmt.Printf("Clicking the submit button...\n")
	if err := chromedp.Run(ctx, chromedp.Click(`#submitButton`, chromedp.ByID)); err != nil {
		return fmt.Errorf("could not click submit button: %v", err)
	} else {
		fmt.Printf("Clicked Submit button\n")
		fmt.Printf("Login successful!\n")
	}

	fmt.Printf("Detecting the wakeup button element to determine if we are on the developer portal homepage...\n")
	if err := chromedp.Run(ctx, chromedp.WaitVisible(`document.querySelector("body > dps-app").shadowRoot.querySelector("div > main > dps-home-auth").shadowRoot.querySelector("div > div > div.instance-widget > dps-instance-sidebar").shadowRoot.querySelector("div > div.dps-instance-sidebar-content.dps-instance-sidebar-content-instance-info > div.dps-instance-sidebar-content-btn-group > dps-button").shadowRoot.querySelector("button")`, chromedp.ByJSPath)); err != nil {
		return fmt.Errorf("could not find shadow element (header status bar): %v", err)
	} else {
		fmt.Printf("Element found\n")
	}

	fmt.Printf("Sleep for a %d seconds for the render of the nodes...\n", 5)
	time.Sleep(5 * time.Second)

	var res int
	if err := chromedp.Run(ctx, chromedp.EvaluateAsDevTools(`(function(){document.querySelector("body > dps-app").shadowRoot.querySelector("div > main > dps-home-auth").shadowRoot.querySelector("div > div > div.instance-widget > dps-instance-sidebar").shadowRoot.querySelector("div > div.dps-instance-sidebar-content.dps-instance-sidebar-content-instance-info > div.dps-instance-sidebar-content-btn-group > dps-button").shadowRoot.querySelector("button").click();return 1;})()`, &res)); err != nil {
		return fmt.Errorf("could not click on shadow element (button): %v", err)
	} else {
		fmt.Println("Clicked on the Wakeup instance button! Exiting...")
	}

	return nil
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
