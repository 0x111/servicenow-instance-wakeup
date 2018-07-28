# servicenow-instance-wakeup

This app is here to help you with wakeing up your instance if needed.

All of you who work with dev instances, you know what is this about.
Dev instances expire after a specific time period.

With this program, you do not need to log in manually to the developer portal anymore to wake up your instance.

All this app is doing is taking your login credentials, logging you into the developer portal and then waking up your instance.

You can use this both ways but only choose one.

The app accepts cli parameters but if this is not something you would like to do, then you can have a `config.json` file in the same directory as the program based on the sample file here.

So an example for the cli app would be something like this:
```
program -username=some@email -password=somepassword -debug=false
```

This would mean that the program will be started with the specified username and password values and with debugging disabled.
We left out the headless option, since that is something you will probably not use if on desktop but only on a server for example.

If we do not want to always specify parameters, we could use the config file in this repository, some basic config would look a lot like this:
```
{
  "username": "developer@email",
  "password": "password",
  "headless": false,
  "debug": false
}
```

If you have any bugs or suggestions please open an issue or pull request.
