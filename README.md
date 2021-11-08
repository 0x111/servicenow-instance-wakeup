# servicenow-instance-wakeup

At the first sight, this project may look like something you might have seen before.
A promise of never waking up your instance anymore. If you are in a search of something like that, you are in the wrong place!

This app symbolizes everything, that the process of waking up your instance should be!
A simple task of issuing one command to wake up your instance. No more time spent waiting on redirects.

Simply run the app and do something useful in the meanwhile instead of repetitive steps to wake your instance.

----
**Disclaimer**: This app is not made to keep your instance awake. This app was created due to the amount of time it takes for you, to wake up your instance. 

It takes you approximately three to four minutes (sometimes more), to wake up your instance right now. 90% of this time, is spent with waiting for redirects, filling out the username and password, waiting for some more redirects and then pushing one button.

This app can be used, to reduce the time to do the manual steps required to wake up your instance. I do not condone keeping any instance awake just for the sake of it. This app only clicks the wake up button if existing (e.g. if the instance is already hibernated). This is not generating any actions that would keep your instance awake. It is simply emulating the manual tasks that you would do anyways if your instance would go to sleep. The app does not even work or does anything if your instance is already awake.

This is not a software to keep any instance awake. Please respect that! There were attempts, to make this a tool to keep it awake, like you can see in [#10](https://github.com/0x111/servicenow-instance-wakeup/issues/10) but this was rejected. Simply said, please do not categorize this app as something, that is jeopardizing the free PDI program. It has nothing to do with it. 

----

All of you who work with dev instances, you know what is this about.
Dev instances expire after a specific time period.

With this software, you do not need to log in manually to the developer portal anymore to wake up your instance.

All this app is doing is taking your login credentials, logging you into the developer portal and then waking up your instance.

The app accepts cli parameters but if this is not something you would like to do, then you can have a `config.json` file in the same directory as the program based on the sample file here.

So an example for the cli app would be something like this:
```
program -username some@email -password somepassword -debug false
```

This would mean that the program will be started with the specified username and password values and with debugging disabled.
We left out the headless option, since that is something you will probably not use if on desktop but only on a server for example.

If we do not want to always specify parameters, we could use the config file in this repository, some basic config would look a lot like this:
```
{
  "username": "developer@email",
  "password": "password",
  "timeout": 60,
  "headless": false,
  "debug": false
}
```

## Docker

To simplify cross platform delivery, the following docker image is capable of waking your ServiceNow Developer Instance.

**Note**: This image does not contain anything, that would run the app periodically or automate it in any way.
You run the docker image, the app starts and after the specified timeout it does exit and thus, the docker container does the same too.
You should only run servicenow-instance-wakeup when you need to wake up your instance to try something. It is not recommended to run this
periodically.


### Docker hub

You can use the pre-built docker image from the [Docker hub](https://hub.docker.com/r/ruthless/servicenow-instance-wakeup) or you can pull the image.
```
docker pull ruthless/servicenow-instance-wakeup
```

### :new: Github Packages
Due to recent changes to Docker Hub and their operations model, I will now publish new images on github packages too :tada:
You can pull this image with the command below:
```
docker pull ghcr.io/0x111/servicenow-instance-wakeup:latest
```
This will pull the latest image, for other version see more on the right.
Eventually I will most likely wind down docker hub and use github packages only. (The timeline mostly depends on how and what kind of changes they will make in the near future...)

### Build from source

In order to build this image from source yourself, follow these steps. 

Clone this github repository and change into the cloned directory.

Run the build process:
```bash
docker build --rm -f "Dockerfile" -t servicenowinstancewakeup:latest "."
```

After the build is done, you are ready to run your own docker image:
```bash
docker run -e USERNAME='YOUR_USERNAME@YOUR_DOMAIN.com' -e PASSWORD='YOUR_SERVICENOW_DEVELOPER_PASSWORD' servicenow-instance-wakeup
```
The DEBUG and HEADLESS environment variables are available in this container should you need them. 

By default, the following environment variables are set:
```bash
DEBUG = false
HEADLESS = true
````

You can configure any of the existing cli flags in the docker container using the env variables. A full example could be something like this:
```bash
docker run -e USERNAME='YOUR_USERNAME@YOUR_DOMAIN.com' -e PASSWORD='YOUR_SERVICENOW_DEVELOPER_PASSWORD' -e DEBUG=`true` -e HEADLESS='false` servicenowinstancewakeup
```

If you have any bugs or suggestions please open an issue or pull request!

Every contribution is welcome!
