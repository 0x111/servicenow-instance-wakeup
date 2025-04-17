package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Configuration struct {
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
	var seconds int64
	configuration := &Configuration{}

	flag.StringVar(&configuration.Username, "username", "", "write the username/email with which you are logging in to the developers account")
	flag.StringVar(&configuration.Password, "password", "", "write the password with which you are logging in to the developers account")
	flag.BoolVar(&configuration.ChromeHeadless, "headless", false, "bool, if we need headless mode with chrome or not, default:false")
	flag.BoolVar(&configuration.Debug, "debug", false, "bool, if you want debug output or not, default:false")
	flag.StringVar(&configFile, "config", "", "Provide the config file name, it can be a relative path or a full path, e.g. /home/user/servicenow-config.json or just simply 'config.json'")
	flag.Int64Var(&seconds, "timeout", 60, "Set the timeout after which the app should exit. This is a number in seconds, default:60")
	flag.Parse()

	// Read config into struct if exists
	if configFile != "" {
		log.Println("Your flags will be ignored and replaced by the values in the config file you specified...")
		log.Printf("Loading config file under the path [%s]", configFile)
		configuration = readConfig(configFile)
	}

	if configuration == nil || len(configuration.Username) == 0 || len(configuration.Password) == 0 {
		log.Println("No username or password provided. Use the -username and -password flags to set the username or password. e.g. program -username user@email.tld or setup a config.json with the details")
		os.Exit(1)
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.DisableGPU,
	)

	log.Printf("Starting the app with debug=%t/headless=%t/account=%s", configuration.Debug, configuration.ChromeHeadless, configuration.Username)

	// navigate to a page, wait for an element, click
	if !configuration.Debug {
		log.SetOutput(ioutil.Discard)
	}

	if configuration.ChromeHeadless {
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

	timeout = time.Duration(seconds) * time.Second

	err = wakeUpInstance(ctx, configuration.Username, configuration.Password, timeout)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func wakeUpInstance(ctx context.Context, username string, password string, timeout time.Duration) error {
	var cancel func()
	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	initialURL := "https://developer.servicenow.com/userlogin.do?relayState=https%3A%2F%2Fdeveloper.servicenow.com%2Fdev.do%23!%2Fhome%3Fwu%3Dtrue"

	// setting viewport to a wider screen so we see the wakeup button
	if err := chromedp.Run(ctx, chromedp.EmulateViewport(1920, 1280)); err != nil {
		return fmt.Errorf("could not set viewport: %v", err)
	} else {
		fmt.Printf("Successfully set the viewport...\n")
	}

	fmt.Printf("Navigating to the webpage: %s\n", initialURL)
	// first navigate to the sso login page
	if err := chromedp.Run(ctx, chromedp.Navigate(initialURL)); err != nil {
		return fmt.Errorf("could not navigate to the SSO login page: %v", err)
	} else {
		fmt.Printf("Successfully navigated to the webpage...\n")
	}

	fmt.Printf("Searching for the logo element...\n")
	if err := chromedp.Run(ctx, chromedp.WaitVisible(`logo`, chromedp.ByID)); err != nil {
		return fmt.Errorf("could not detect logo element: %v", err)
	} else {
		fmt.Printf("Found logo element\n")
	}

	fmt.Printf("Filling out the username field...\n")
	if err := chromedp.Run(ctx, chromedp.SendKeys(`username`, username, chromedp.ByID)); err != nil {
		return fmt.Errorf("could not fill out the username: %v", err)
	} else {
		fmt.Printf("Filled username field with %s\n", username)
	}

	fmt.Printf("Clicking the next button...\n")
	if err := chromedp.Run(ctx, chromedp.Click(`identify-submit`, chromedp.ByID)); err != nil {
		return fmt.Errorf("could not click the next button: %v", err)
	} else {
		fmt.Printf("Clicked Next button\n")
	}

	fmt.Printf("Searching for the password field...\n")
	if err := chromedp.Run(ctx, chromedp.WaitVisible(`password`, chromedp.ByID)); err != nil {
		return fmt.Errorf("could not detect password element: %v", err)
	} else {
		fmt.Printf("Found password field\n")
	}

	fmt.Printf("Filling out the password field...\n")
	if err := chromedp.Run(ctx, chromedp.SendKeys(`password`, password, chromedp.ByID)); err != nil {
		return fmt.Errorf("could not fill out the password: %v", err)
	} else {
		fmt.Printf("Filled password field with your password ******\n")
	}

	fmt.Printf("Clicking the Sign In button...\n")
	if err := chromedp.Run(ctx, chromedp.Click(`challenge-authenticator-submit`, chromedp.ByID)); err != nil {
		return fmt.Errorf("could not click Sign In button: %v", err)
	} else {
		fmt.Printf("Clicked Sign In button\n")
		fmt.Printf("Login successful!\n")
	}

	fmt.Printf("Setting cookies to pre-confirm cookie modal!\n")

	// set the cookies to remove the iframe modal
	if err := chromedp.Run(ctx, setcookies(
		"notice_preferences", "0",
		"notice_gdpr_prefs", "0",
		"cmapi_gtm_bl", "ga-ms-ua-ta-asp-bzi-sp-awct-cts-csm-img-flc-fls-mpm-mpr-m6d-tc-tdc",
		"cmapi_cookie_privacy", "permit 1 required")); err != nil {
		return fmt.Errorf("could not set cookies: %v", err)
	} else {
		fmt.Printf("Cookies set!\n")
	}

	if err := chromedp.Run(ctx, chromedp.WaitVisible(`document.querySelector("body > dps-app").shadowRoot.querySelector("div > main > dps-home-auth-quebec").shadowRoot.querySelector("div > section:nth-child(1) > div > dps-page-header > div:nth-child(1)")`, chromedp.ByJSPath)); err != nil {
		return fmt.Errorf("could not find start building button: %v", err)
	} else {
		fmt.Printf("Start building button found\n")
	}

	fmt.Printf("Instance wakeup initiated successfully, your instance should be awake pretty soon!\n")

	return nil
}

// Read the config file if required and load the json to the struct
func readConfig(config string) *Configuration {
	// Load the specified config file from the path provided
	jsonFile, err := os.Open(config)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	configurationParameters := Configuration{}

	err = json.Unmarshal(byteValue, &configurationParameters)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &configurationParameters
}

// setcookies returns a task to navigate to a host with the passed cookies set
// on the network request.
func setcookies(cookies ...string) chromedp.Tasks {
	if len(cookies)%2 != 0 {
		panic("length of cookies must be divisible by 2")
	}

	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// create cookie expiration
			expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
			// add cookies to chrome
			for i := 0; i < len(cookies); i += 2 {
				err := network.SetCookie(cookies[i], cookies[i+1]).
					WithExpires(&expr).
					WithDomain("developer.servicenow.com").
					WithHTTPOnly(false).
					WithSecure(true).
					Do(ctx)

				if err != nil {
					return err
				}
			}
			return nil
		}),
	}
}
