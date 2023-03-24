Ponder
-------
###### OpenAI Powerd Chat Tool  
[![main](https://github.com/seemywingz/ponder/actions/workflows/dockerBuildX.yml/badge.svg?branch=main)](https://github.com/seemywingz/ponder/actions/workflows/dockerBuildX.yml)


# Install
```bash
go install github.com/seemywingz/ponder
```

# Usage
## Setup Your Environment
The [OpenAI API](https://platform.openai.com/docs/api-reference/authentication) uses API keys for authentication.  
Visit your [API Keys](https://platform.openai.com/account/api-keys) page to retrieve the API key you'll use in your requests.


#### ℹ️ This is required for both docker and cli usage
```bash
export OPENAI_API_KEY={YOUR OPENAI API KEY}
```


## A single thought
### CLI
```bash
ponder chat --prompt "AI is Amazing"
```
### Docker
```bash
docker run -e OPENAI_API_KEY=$OPENAI_API_KEY disciplesofai/ponder:latest chat --prompt "Ai is Amazing"
```
#### Example Output
```bash
AI is indeed amazing. It has the potential to revolutionize many industries, from healthcare to transportation. AI can help automate mundane tasks, allowing humans to focus on more complex and creative tasks. AI can also be used to analyze large amounts of data quickly and accurately, helping to make better decisions. AI can also be used to create personalized experiences for customers, making them feel more connected to a company or product.
```

## A small chat
### CLI
```bash
ponder chat --loop
```
### Docker
```bash
docker run -it -e OPENAI_API_KEY=$OPENAI_API_KEY disciplesofai/ponder:latest chat --loop
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
### CLI
```bash
ponder image -p "watercolor of a corgie"
```
### Docker
```bash
docker run -e OPENAI_API_KEY=$OPENAI_API_KEY disciplesofai/ponder:latest image -p "watercolor of a corgie"
```
#### Example Ouput
```bash
🖼  Creating Image...
🌐 Image URL: https://oaidalleapiprodscus.blob.core.windows.net/private/org-RCMQxIXre0Olhs0AvLVp672o/user-F1wdcIVNf2VrRqBRD0JWUczI/img-AWku5cm91XAv32jj27XWXZBE.png?st=2023-03-05T05%3A19%3A33Z&se=2023-03-05T07%3A19%3A33Z&sp=r&sv=2021-08-06&sr=b&rscd=inline&rsct=image/png&skoid=6aaadede-4fb3-4698-a8f6-684d7786b067&sktid=a48cca56-e6da-484e-a814-9c849652bcb3&skt=2023-03-05T01%3A25%3A44Z&ske=2023-03-06T01%3A25%3A44Z&sks=b&skv=2021-08-06&sig=xs9vSD0nA0mkxyulHEKABn5cbWH%2B6YOpab25yTAU/nc%3D
```