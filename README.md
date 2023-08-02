# discord-bot-slash-commands-go
This is a quick and dirty proof of concept to teach myself how to build a discord bot with slash commands
in golang using [discordgo](https://github.com/bwmarrin/discordgo).

I'm passing in the tokens via environment variables using [envconfig](https://github.com/kelseyhightower/envconfig).
I just have a simple little bash script that exports the variables and then runs the bot:
```
#!/bin/bash

export BOT_BOT_TOKEN="<Bot Token Here>"
export BOT_GUILD_ID="<Guild Token Here>"

./discord-bot-slash-commands-go
```
