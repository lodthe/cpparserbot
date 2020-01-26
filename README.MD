# Crypto parser
![Travis](https://travis-ci.com/LoDThe/cpparserbot.svg?branch=master)

This is a Telegram bot that provides with relevant information about crypto currency prices from different exchangers (such as Binance, Exmo).

![Message example](https://i.ibb.co/zVTFpd2/image.png)

## Stack

Heroku is used to start bot in docker. Chosen CI is Travis.

## How to run it

First of all, you should initialize environment variables:

| Name   | Description     |
|------------|-----------|
| TELEGRAM_BOT_TOKEN       | Telegram bot token ([@BotFather](https://teleg.run/BotFather) knows it) |
| TELEGRAM_CHANNEL_CHAT_ID       | Telegram channel ID where bot will send log-information      |
| BINANCE_API_KEY       | Binance API key      |
| BINANCE_SECRET_KEY       | Binance secret key (can be obtained in Binance settings)      |

The next thing to do is adding supported Binance currency pairs in `config.go` file:

```go
BinancePairs = []model.Pair{
    {"BTC", "USDT"},
    {"BNB", "BTC"},
    {"ETH", "BTC"},
    {"XRP", "BTC"},
    {"DASH", "BTC"},
    {"ETC", "BTC"},
}
```

## How to try it

Only russian language is supported.

You can try it here: [@cpparser_bot](https://teleg.run/cpparser_bot)