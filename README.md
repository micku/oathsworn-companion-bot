# Oathsowrn Companion Bot

This is a simple Telegram bot that helps my party to play the story phases async. We use Telegram to keep the progress and, as much as we love rolling dice, sometimes we do not have them with us.

## How to host

At the time of writing, the image is only built for ARM64 as I'm deploying it to a K3s cluster made of Raspberry Pi nodes.

Here is a sample Kubernetes deployment:

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: oathsworn-companion-bot
  labels:
    app: oathsworn-companion-bot
spec:
  replicas: 1
  selector:
    matchLabels:
      app: oathsworn-companion-bot
  template:
    metadata:
      labels:
        app: oathsworn-companion-bot
    spec:
      containers:
        - name: oathsworn-companion-bot
          image: ghcr.io/micku/oathsworn-companion-bot:main
          imagePullPolicy: Always
          envFrom:
            - secretRef: 
                name: oath-bot-secrets
```

Here is the secrets definition:

```
apiVersion: v1
kind: Secret
metadata:
  name: oath-bot-secrets
data:
  TG_TOKEN: ***
  TG_ALLOWED_CHATS: ***
```

Here you need to customize the 2 values:
- `TG_TOKEN`: the bot token provided by the BotFather upon bot creation;
- `TG_ALLOWED_CHATS`: empty for allowing the bot to answer to any chat, otherwise a comma separated list of chat IDs. To obtain the chat ID, you can initially deploy without any restriction and then read the chat IDs from the logs.

## How it works

For now, it supports only 2 commands.

`/roll <dice>` rolls the combination of dice provided as a string of `w`, `y`, `r` or `b`. Ie. `/roll bbry` would roll 2 black dice, 1 red and 1 yellow. This automatically rolls crits.

`/reroll <dice index>` must be used as a reply to a rolled dice message passing a space separated list of indexes.
