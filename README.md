# 🚀 GoGut 💬 - Your AI Genie

> Unleash the power of artificial intelligence to streamline your command line experience.

![Intro](docs/_assets/intro.gif)

## What is GoGut ?

`GoGut` (your AI) is an assistant for your terminal, using [OpenAI ChatGPT](https://chat.openai.com/) to build and run commands for you. You just need to describe them in your everyday language, it will take care or the rest.

If you want to use OpenAI compatibility service/local models, you can for example use [Ollama](https://ollama.com/)

You have any questions on random topics in mind? You can also ask `GoGut`, and get the power of AI without leaving `/home`.

It is already aware of your:
- operating system & distribution
- username, shell & home directory
- preferred editor

And you can also give any supplementary preferences to fine tune your experience.

## Quick start

To install `GoGut`, simply run:

```shell
curl -sS https://raw.githubusercontent.com/bmichalkiewicz/gogut/main/install.sh | bash
```

At first run, it will ask you for an [OpenAI API key](https://platform.openai.com/account/api-keys), and use it to create the configuration file in `~/.gogut/config.yaml`.
