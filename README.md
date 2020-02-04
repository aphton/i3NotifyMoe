# i3NotifyMoe
An i3blocks taskbar integration for notify.moe!

## What does it do?
This little tool goes to notify.moe/api and fetches curently airing & watching animes
for the user (set in environment variable "I3_NOTIFY_MOE_USERNAME").

It then splits the list into these 3 categories:
- released - contains episodes that have already aired
- today - contains episodes that will air today, but haven't yet
- unreleased - contains episodes which air tomorrow or later

And then it outputs something like "1/3" into stdout, which indicates that there is 1 episode ready to be consumed!

The tool is configured to run once when i3blocks starts and every time you do a left-click on it in the taskbar!
Right click will pop up a desktop notification showing you those 3 lists in a nicely formatted way!

## Screenshots
Taskbar

![taskbar](https://github.com/aphton/i3NotifyMoe/blob/master/screenshot-1-taskbar.png)

Notification

![notification](https://github.com/aphton/i3NotifyMoe/blob/master/screenshot-2-desktop-notification.png)


## How to build?
Two commands to execute:

    go get github.com/go-floki/jade
    go build
    
## How to configure i3blocks?
Just paste this into your i3blocks config file (locate in i.e. /etc/i3blocks.conf)

    [anime]
    command=<path to executable>
    I3_NOTIFY_MOE_USERNAME=<your username on notify.moe (case sensitive!)>
    label=<(unicode-) label of your choice>
    markup=pango
    interval=once
