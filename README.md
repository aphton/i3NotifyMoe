# i3NotifyMoe
An i3blocks taskbar integration for notify.moe!

## What does it do?
This little tool goes to notify.moe/api and fetches curently airing & watching animes for the user (env variable "I3_NOTIFY_MOE_USERNAME").

It then categorizes the result into these 3 lists:
- released - contains episodes that have already aired
- today - contains episodes that will air today, but haven't yet
- unreleased - contains episodes which air tomorrow or later

It then outputs something like "1/3" into stdout, which indicates that there is 1 episode ready to be consumed today!

The tool is configured to run once when i3blocks starts and every time you do a left-click on it in the taskbar!
Right click will pop up a desktop notification showing you those 3 lists in a neatly formatted way!

## Screenshots
Taskbar

![taskbar](https://github.com/aphton/i3NotifyMoe/blob/master/screenshot-1-taskbar.png)

Notification

![notification](https://github.com/aphton/i3NotifyMoe/blob/master/screenshot-2-desktop-notification.png)


## How to build?
Two commands to execute:

    go get github.com/go-floki/jade
    go install i3NotifyMoe
    
## How to configure i3blocks?
Just paste this into your i3blocks config file (locate in i.e. /etc/i3blocks.conf)

    [anime]
    command=<path to executable>
    I3_NOTIFY_MOE_USERNAME=<your username on notify.moe (case sensitive!)>
    markup=pango
    interval=once
