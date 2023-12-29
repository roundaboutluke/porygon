# porygon
**_"Roughly 20 years ago, it was artificially created, utilizing the latest technology of the time."_**

Porygon is a reimagining of Discordopole for Golbat, written in go with massive amounts of input from GPT-4 to help make my spaghetti code even more spaghetti!

![image](https://i.imgur.com/OXv3jZ3.png)


**Note:** Comparitively to Discordopole the features are amazingly lackluster. This will create a simple board like so featuring daily stats, utilising both the database and API and update based on the interval you define within the config file, that's it (for now).

**A wise Jabes once said**


![image](https://i.imgur.com/ZOsk45B.png)

I tried to heed this warning as best I could, providing min/max lon/lat config options and a configurable refresh interval so you can tax your system as much or as little as you want.


# Installation

1. Git clone the repo
2. copy `config.example` to `config.json` and fill out accordingly _(refresh interval is in **seconds** - you'll also need to link to your own emojis, there are some in /emojis you can add to your server if you don't already have them)_
3. `go build Porygon.go`
4. `pm2 start ./Porygon --name porygon`

# Important

I don't really know what I'm doing which is probably evident to anyone that does looking at this, but it works for me so hopefully it does for you too. My map is teeny tiny and I trust Jabes wholely so don't be surprised if this is incredibly taxing for you big mappers.

# Todo

I want to build some (optional) stats from Blissey in, but I haven't even got Blissey running yet 😂
