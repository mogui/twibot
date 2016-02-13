
# What is it?

[twibot](http://github.com/mogui/twibot) is a light [Twitter](https://twitter.com) bot writtenin Go, that can execute configured commands by *Direct message* or *Mention* him with a tweet.   
You just have to create a dedicated twitter account (see tutorial [here](https://github.com/mogui/twibot/wiki/Account-Tutorial)) create a configuration file with authentication credential, the authorized twitter account, and the description of the **commands** your bot can execute.

Its main aim is have a secure channel to execute scripts on a remote server.

## Getting Started
Make sure you've properly set up [golang](https://golang.org/doc/install) environment, then to install twibot:

```bash
$ go get github.com/mogui/twibot
```

To run twibot you have to create a config file like this [one](https://github.com/mogui/twibot/blob/master/conf.json.example):
```
{
  "consumer_key": "YOUR TWITTER CONSUMER KEY",
  "consumer_secret": "YOUR TWITTER CONSUMER SECRET",
  "token": "YOUR TWITTER TOKEN",
  "token_secret": "YOUR TWITTER TOKEN SECRET",
  "authorized_account": ["mogui247", "asder"]
  "on_dm": [
    {
      "name": "Test Command",
      "match": "^test",
      "script": "ls",
      "reply": true,
      "case": true
    },
    {
      "name": "Test create file",
      "match": "new ([a-zA-Z]+)",
      "script": "touch {1}",
      "reply": true
    }
  ],
  "on_mentions": [
      ...
  ]
}
```
## Breakdown of config file

### Config object

| key                |  type  | required | description                                                  |
|:-------------------|:------:|:--------:|:-------------------------------------------------------------|
| consumer_key       | string |    x     | twitter api auth                                             |
| consumer_secret    | string |    x     | twitter api auth                                             |
| token              | string |    x     | twitter api auth                                             |
| token_secret       | string |    x     | twitter api auth                                             |
| screen_name        | string |    x     | `screen_name` of the twitter account used for the bot        |
| authorized_account | string |          | list of twitter `screen_name` authorized to use the bot      |
| on_dm              | array  |          | list of *Commands* scanned when a Direct Message is received |
| on_mentions        | array  |          | list of *Commands* scanned when a Mention is received        |

### Command object
| key    |  type  | required | description                                                           |
|:-------|:------:|:--------:|:----------------------------------------------------------------------|
| name   | string |    x     | Name of the command                                                   |
| match  | string |    x     | a regular expression that will be run against the reiceved tweet. *NOTE:*  you can use groups in the regex in order to pass theem as arguments to the script by using this noteation `{1}` just like the example in the json    |
| script | string |    x     | the executable script to be run when matched.  |
| reply  |  bool  |          | if the bot should reply with a DM after command exec [default: false] |
| case   |  bool  |          | if regex id case sensitive [default: false]                           |

The order of the *Commands* in the config file is important because they will be scanned sequentially and the first one to match will be fired

You can then run twibot as follows:

```bash
$ twibot --config conf.json --verbose
```

## Bundled commands

twibot comes with some bundled *special* commands over the On Direct Message channel, all of them are case insensitive:

| command  | description                                                |
|:---------|:-----------------------------------------------------------|
| ping     | will just reply PONG                                       |
| help / ? | will reply with a DM conatining all the available commands |


## TODO's
- [ ] tutorial for setting up the twitter account
- [ ] implement Mentions
- [ ] dockerization
- [ ] add other bundled commands ??
- [ ] common templates to run it at startup

## Contributing
Any form of contribution is welcome and appreciated, though I prefer the standard pull-request flow. Thanks

## Copyright
Copyright © 2016 Niko Usai. See LICENSE for details.   
