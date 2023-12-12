Ponder
-------
###### OpenAI Powered Chat Tool    
[![main](https://github.com/seemywingz/ponder/actions/workflows/dockerBuildX.yml/badge.svg?branch=main)](https://github.com/seemywingz/ponder/actions/workflows/dockerBuildX.yml)


# Install
```bash
go install github.com/seemywingz/ponder
```

# Usage
## Setup Your Environment
The [OpenAI API](https://platform.openai.com/docs/api-reference/authentication) uses API keys for authentication.  
Visit your [API Keys](https://platform.openai.com/account/api-keys) page to retrieve the API key you'll use in your requests.


### Required Environment Variables
###### ℹ️ These Environment Variables are required for both docker and cli usage
###### ℹ️ You can omit keys for unused API endpoints
```bash
OPENAI_API_KEY={YOUR OPENAI API KEY}
PRINTIFY_API_KEY={YOUR PRINTIFY API KEY}
DISCORD_API_KEY={YOUR DISCORD BOT API KEY}
DISCORD_PUB_KEY={YOUR DISCORD BOT PUB KEY}
```



## Ponder a single thought
### CLI
```bash
ponder --prompt "AI is Amazing"
```
### Docker
#### Running ponder in docker is exactly the same, but you have to provide the env vars when running
```bash
docker run -e OPENAI_API_KEY=$OPENAI_API_KEY  ghcr.io/seemywingz/ponder:latest --prompt "Ai is Amazing"
```
#### Example Output
```bash
AI is indeed amazing. It has the potential to revolutionize many industries, from healthcare to transportation. AI can help automate mundane tasks, allowing humans to focus on more complex and creative tasks. AI can also be used to analyze large amounts of data quickly and accurately, helping to make better decisions. AI can also be used to create personalized experiences for customers, making them feel more connected to a company or product.
```

## A small chat
```bash
ponder --convo
```
#### Example Ouput
```bash
You: 
Hello, Ponder

Ponder: 
Hello there! How can I help you?

You: 
you already have thank you

Ponder: 
Thank you for your kind words!
```

## Image Generation
```bash
ponder image -p "a ferocious cat with wings and fire"
```
#### Example Ouput
```bash
🖼  Creating Image...
🌐 Image URL: https://oaidalleapiprodscus.blob.core.windows.net/private/org-RCMQxIXre0Olhs0AvLVp672o/user-F1wdcIVNf2VrRqBRD0JWUczI/img-B4gaFhJQFl25authc5zMdw3T.png?st=2023-12-12T19%3A42%3A45Z&se=2023-12-12T21%3A42%3A45Z&sp=r&sv=2021-08-06&sr=b&rscd=inline&rsct=image/png&skoid=6aaadede-4fb3-4698-a8f6-684d7786b067&sktid=a48cca56-e6da-484e-a814-9c849652bcb3&skt=2023-12-12T05%3A22%3A04Z&ske=2023-12-13T05%3A22%3A04Z&sks=b&skv=2021-08-06&sig=RteaU2hpHlz5VElxgxdwUahGHoQmy6SEAVdpsjDbt%2Bg%3D
```

### You can always refer to the `--help` menu as well.
```yaml
        Ponder
        GitHub: https://github.com/seemywingz/ponder
        App Version: v0.3.0

  Ponder uses OpenAI's API to generate text responses to user input.
  You can use Ponder as a Discord chat bot or to generate images using the DALL-E API.
  Or whatever else you can think of...

Usage:
  ponder [flags]
  ponder [command]

Available Commands:
  adventure   lets you dive into a captivating text adventure
  chat        Open ended chat with OpenAI
  completion  Generate the autocompletion script for the specified shell
  discord-bot Discord Chat Bot Integration
  help        Help about any command
  image       Generate an image from a prompt
  printify    Interact with the Printify API

Flags:
      --config string   config file
  -c, --convo           Conversational Style chat
  -h, --help            help for ponder
  -x, --perform         Attempt to perform the response as cli command
  -p, --prompt string   Prompt AI generation
  -s, --say             Say text out loud (MacOS only)
  -v, --verbose         verbose output

Use "ponder [command] --help" for more information about a command.
```