# porygon
**_"Roughly 20 years ago, it was artificially created, utilizing the latest technology of the time."_**

Porygon is a reimagining of Discordopole for Golbat, written in go with massive amounts of input from GPT-4 to help make my spaghetti code even more spaghetti!

![image](https://i.imgur.com/Q7jKuVY.png)


**Note:** Comparitively to Discordopole the features are amazingly lackluster. This will create a simple board like so featuring daily stats, utilising both the database and API and update based on the interval you define within the config file, that's it (for now).

**A wise Jabes once said**


![image](https://i.imgur.com/ZOsk45B.png)

I tried to heed this warning as best I could, providing min/max lon/lat config options, a configurable refresh interval and the ability to disable active counts so you can tax your system as much or as little as you want.

# Requirements

[go 1.21](https://go.dev/doc/install)

# Installation

1. Git clone the repo `git clone https://github.com/roundaboutluke/porygon.git`
2. `cp default.toml config.toml` & adjust config.toml accordingly
3. `go build .`
5. `pm2 start ./porygon --name porygon`

# Further Customisation

Thanks to @lenisko you can now customise Porygon's localisation and layout. Simply `cp templates/current.json templates/current.override.json` and edit accordingly, using the examples in current.json. Note that some of the more generic emojis are contained within this.

# Updating

1. `git pull`
3. `go build .`
3. `pm2 restart porygon`

# Discord Permissions

Porygon requires your bot have the server permissions **Send Messages**, **Read Message History** and **Embed Links** to function.

# Important

I don't really know what I'm doing which is probably evident to anyone that does looking at my spaghetti, but it works for me so hopefully it does for you too. People that have more of a clue have made some great improvements to the codebase.

# Todo

Many things, I've started tracking them in Issues, please feel free to add more!
