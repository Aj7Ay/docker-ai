package learning

import (
	"strings"
)

type Topic struct {
	Keywords []string
	Content  string
}

var topics = []Topic{
	{
		Keywords: []string{"how to install docker", "install docker on ubuntu", "docker installation script"},
		Content: "Here is a full script to install Docker on Ubuntu, based on the official documentation.\n\nYou can run this script directly on a fresh Ubuntu system to get Docker up and running.\n\n" +
			"```bash\n" +
			"#!/bin/bash\n" +
			"# Uninstall old versions\n" +
			"for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do sudo apt-get remove $pkg; done\n\n" +
			"# Add Docker's official GPG key:\n" +
			"sudo apt-get update\n" +
			"sudo apt-get install -y ca-certificates curl\n" +
			"sudo install -m 0755 -d /etc/apt/keyrings\n" +
			"sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc\n" +
			"sudo chmod a+r /etc/apt/keyrings/docker.asc\n\n" +
			"# Add the repository to Apt sources:\n" +
			"echo \\\n" +
			"  \"deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \\\n" +
			"  $(. /etc/os-release && echo \"$VERSION_CODENAME\") stable\" | \\\n" +
			"  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null\n" +
			"sudo apt-get update\n\n" +
			"# Install the latest version\n" +
			"sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin\n\n" +
			"# Verify the installation\n" +
			"sudo docker run hello-world\n" +
			"```\n\n" +
			"For more details, see the official docs: https://docs.docker.com/engine/install/ubuntu/",
	},
	{
		Keywords: []string{"check running containers", "ps", "check container", "list containers", "check exited containers", "check all containers"},
		Content: `To see your containers, you use the 'docker ps' command. Think of it like asking, "What are my containers doing?"

- **To see only RUNNING containers:**
  'docker ps'
  This is the default. It shows you a list of all your active containers.

- **To see ALL containers (running and stopped/exited):**
  Use the '-a' flag, which stands for "all".
  'docker ps -a'
  This is how you check for containers that you ran in the past that have since stopped. It's very useful for debugging!

The output will show you:
- Their unique ID (like a name tag).
- What image they came from (which box of toys they are).
- Their status (running, exited, or created).
- Their name (you can give them custom names!).`,
	},
	{
		Keywords: []string{"dockerfile", "what is a dockerfile", "write dockerfile"},
		Content: `A Dockerfile is like a recipe for baking a cake. The "cake" is your application's container image (a perfect, ready-to-run package of your app).

The Dockerfile has a list of step-by-step instructions for Docker to follow:
1.  **Start with a base ingredient** (like a plain version of Linux).
2.  **Add your application code** (the special ingredients).
3.  **Install any other software** your app needs to run (like sugar and flour).
4.  **Tell Docker how to run your app** when someone starts the container (how to serve the cake).

When you give this recipe to Docker (using 'docker build'), it follows the steps precisely and creates a perfect image for you, every single time.`,
	},
	{
		Keywords: []string{"from", "run", "cmd", "dockerfile layer"},
		Content: `FROM, RUN, and CMD are the most common and important instructions in your Dockerfile recipe.

- **FROM**: This is ALWAYS the first step. It tells Docker what base image to start with.
  - Think of it as choosing the plate for your cake.
  - Example: 'FROM python:3.9' means "Start with a clean plate that already has Python 3.9 on it."

- **RUN**: This instruction tells Docker to execute a command to set things up. You use it to install software or prepare your application.
  - Think of this as adding and mixing ingredients onto your plate. Each 'RUN' command creates a new "layer" on top of the previous one.
  - Example: 'RUN pip install -r requirements.txt' tells Docker to install all the Python packages your app needs.

- **CMD**: This instruction tells the container what command to run when it starts up.
  - Think of this as the final instruction on the recipe card: "To serve, slice and enjoy!" It's the main purpose of your container.
  - Example: 'CMD ["python", "./app.py"]' tells the container to start your Python application. You usually only have one CMD instruction at the very end.`,
	},
	{
		Keywords: []string{"compose", "docker-compose", "multiple container", "multi-container"},
		Content: `Imagine you're building a big Lego castle. You need the main castle, a fire-breathing dragon, and a brave knight. Each of these is a separate part of your final creation.

In Docker, sometimes your application is made of multiple parts (multiple containers) that need to work together. For example, a website (one container) might need a database (another container).

A 'docker-compose.yml' file is the MASTER INSTRUCTION BOOK for building and running your whole Lego scene at once.

Instead of running commands for each container separately, you just write down all the parts in this one file:
- A 'web' service that uses your website image.
- A 'db' service that uses a database image.
- How they should connect to each other.

Then, you just run one command: 'docker-compose up'

Docker reads your master instruction book and starts everything up together, in the right order. It's perfect for managing applications that have several moving parts!`,
	},
	{
		Keywords: []string{"docker run", "run container", "start container"},
		Content: `The 'docker run' command is how you bring a container to life from an image.

Imagine you have a recipe for a cake (the image). 'docker run' is like saying, "Okay, follow this recipe and bake me a cake right now!"

When you use 'docker run', you tell Docker:
- Which image (recipe) to use.
- Optionally, you can give the container a name, connect it to networks, and much more.

Example: 'docker run hello-world'
This command takes the 'hello-world' image, creates a new container from it, and runs it. The container prints its "Hello from Docker!" message and then stops.`,
	},
	{
		Keywords: []string{"build image", "create image", "docker build"},
		Content: `The 'docker build' command is what you use to create a Docker image from a Dockerfile.

If a Dockerfile is your recipe, 'docker build' is the process of actually baking the cake.

You run this command in the same directory as your Dockerfile. It will read your recipe (the Dockerfile), follow all the steps (like RUN commands), and at the end, it will produce a finished cake (your image).

Example: 'docker build -t my-first-app .'
- 'docker build': The command to start the build process.
- '-t my-first-app': This gives your new image a name (a "tag"). Here, we're calling it 'my-first-app'.
- '.': This tells Docker to look for the Dockerfile in the current directory.`,
	},
	{
		Keywords: []string{"check logs", "container logs", "view logs", "since", "until", "tail"},
		Content: `Sometimes you don't want to see ALL the logs from a container, which could be thousands of lines! Docker gives you some tools to see just the parts you want.

- **See the most recent logs:**
  Use the '--tail' flag to see just the last few lines. It's like asking for the last page of a book.
  'docker logs --tail 50 my-container' (shows the last 50 lines of logs for 'my-container')

- **See logs from a certain time:**
  Use the '--since' flag to see logs that happened after a specific time.
  'docker logs --since 10m my-container' (shows all logs from the last 10 minutes)

- **See logs before a certain time:**
  Use the '--until' flag to see logs that happened before a specific time.
  'docker logs --until 10m my-container' (shows all logs from the beginning until 10 minutes ago)

You can even combine them! 'docker logs --since 1h --until 10m my-container' shows logs from an hour ago up until 10 minutes ago.`,
	},
	{
		Keywords: []string{"watch ps", "stats", "live stats"},
		Content: `The 'docker ps' command is great, but it only shows you what's running at that exact moment. It doesn't update automatically.

Docker doesn't have a built-in 'watch' command. However, you can use other tools for a live view:

1.  **For a live list of containers (on Linux/macOS):**
    You can use your terminal's 'watch' command to run 'docker ps' every few seconds.
    'watch docker ps' (This will refresh the list every 2 seconds by default).

2.  **For live resource usage (CPU, Memory):**
    This is what you usually want to "watch". Docker has a great command for this!
    'docker stats'
    This command gives you a live, updating table showing how much CPU, memory, and network your containers are using. It's like a Task Manager for your containers!`,
	},
	{
		Keywords: []string{"inspect", "details", "low-level info"},
		Content: `The 'docker inspect' command is like putting a container (or image, or network) under a microscope.

It gives you ALL the low-level details about a Docker object in a big JSON format. This is super useful for debugging or for automation when you need to find a specific piece of information, like a container's IP address.

How to use it:
- 'docker inspect my-container' (shows details for the container named 'my-container')
- 'docker inspect my-image' (shows details for the image named 'my-image')

Most Docker objects (containers, images, networks, volumes, etc.) can be inspected to get the full story about them.`,
	},
	{
		Keywords: []string{"list all", "ls", "show all"},
		Content: `Docker has a consistent way to show you all the things of a certain type you have on your system. Most of the time, the command is 'ls' (which stands for 'list').

Here are some common ones:
- 'docker container ls' (or 'docker ps'): Lists your running containers. Add '-a' to see all containers, even stopped ones!
- 'docker image ls' (or 'docker images'): Lists all your images.
- 'docker network ls': Lists all your Docker networks.
- 'docker volume ls': Lists all your Docker volumes (where you store persistent data).
- 'docker context ls': Lists your Docker contexts (different environments you can connect to).
- 'docker compose ls': Lists your Docker Compose projects.

If you ever wonder, "How many networks do I have?", just try 'docker network ls'!`,
	},
	{
		Keywords: []string{"system info", "disk usage", "docker version"},
		Content: `It's useful to know what's going on with your Docker system as a whole.

- **Check Docker's disk usage:**
  'docker system df'
  This command shows you how much space your images, containers, and volumes are taking up. It's great for figuring out what you need to clean up! ('df' stands for 'disk free').

- **Get general system info:**
  'docker system info' (or just 'docker info')
  This gives you a big report on your Docker installation, like how many containers and images you have, what operating system you're on, and more.

- **Check the version:**
  'docker version'
  This shows you the version of the Docker client (the command-line tool) and the Docker Engine (the server) you are running.`,
	},
	{
		Keywords: []string{"docker commands", "common docker commands", "explain docker commands", "basic docker commands", "list images", "check images", "see images"},
		Content: `Here are the most common Docker commands you'll use every day. Think of them like tools in your workshop.

- **'docker run'**: Starts a new container. (This is like turning on a new machine).
- **'docker ps'**: Lists your running containers. (See which machines are currently on).
- **'docker images'**: Lists all the images you have stored locally. (See all the blueprints you have for making machines).
- **'docker pull'**: Downloads an image from a remote registry like Docker Hub. (Get a new blueprint from the store).
- **'docker push'**: Uploads an image to a remote registry. (Share a blueprint you made with the store).
- **'docker stop'**: Stops a running container gracefully. (Turn off a machine).
- **'docker rm'**: Removes a stopped container. (Put a machine you've turned off back in the box).
- **'docker rmi'**: Removes an image. (Throw away a blueprint you don't need anymore).`,
	},
	{
		Keywords: []string{"buildx", "advanced build", "multi-platform", "target stage", "secret file"},
		Content: `'docker buildx build' is a modern, super-powered version of the 'docker build' command. It uses a new engine called BuildKit that unlocks a lot of advanced features.

Here are some of the most useful advanced options, explained simply:

- **--platform**: Build an image for different types of computers.
  - *Use Case*: You can build an image on your Mac that will work perfectly on a Linux server in the cloud.
  - *Example*: 'docker buildx build --platform linux/amd64,linux/arm64 .'

- **--target**: Build only one part of a multi-stage Dockerfile.
  - *Use Case*: If your Dockerfile has a 'testing' stage and a 'production' stage, you can choose to only build the 'testing' part to run your tests.
  - *Example*: 'docker buildx build --target testing .'

- **--secret**: Securely use secret files (like passwords or tokens) during the build.
  - *Use Case*: You can use a password to download a private file during the build, but the password itself won't be saved in the final image.
  - *Example*: 'docker buildx build --secret id=mysecret,src=mysecret.txt .'

- **--ssh**: Securely use SSH keys during the build.
  - *Use Case*: This is perfect for cloning private Git repositories as part of your build process.

- **--push**: Push the image directly to a registry when the build finishes.
  - *Use Case*: A great time-saver for your CI/CD pipelines. Build and publish in one step!`,
	},
	{
		Keywords: []string{"mcp", "model context protocol", "gordon", "mcp catalog"},
		Content: `Model Context Protocol (MCP) is a cool, open standard that lets AI models interact with the world through tools.

**What does it do?**
Imagine an AI is a brain in a jar. MCP is like giving that brain hands and eyes. It defines a way for an AI (like Docker's Gordon) to ask a special 'MCP server' to do things for it, like:
- Read a file from your computer.
- Clone a git repository.
- Check what's inside a container.

This gives the AI extra 'context' and 'functionality' beyond its built-in knowledge.

**Where can I find these tools?**
- **Docker MCP Catalog**: A collection of ready-to-use MCP servers that you can run as containers. (See: https://docs.docker.com/ai/mcp-catalog-and-toolkit/catalog/)
- **MCP Servers on GitHub**: You can find examples and build your own MCP servers from this repository. (See: https://github.com/docker/mcp-servers)

In short, MCP is the bridge that connects powerful language models to real-world tools and actions.`,
	},
	{
		Keywords: []string{"docker model runner", "model runner", "ai models", "run ai model", "pull model", "docker ai"},
		Content: `Docker Model Runner is a cool feature that makes it super easy to run powerful AI models on your own computer. Think of it as a special engine inside Docker just for AI!

**What can you do with it?**
- Download and run popular open-source AI models with simple commands.
- Build your own GenAI applications that run entirely on your machine.
- Avoid complex Python and dependency setups for many models.

**Key Commands:**
- **List available models:** 'docker model ls'
- **Download a model:** 'docker model pull ai/smollm2'
- **Run a model and chat with it:** 'docker model run ai/smollm2'
- **See the logs:** 'docker model logs'

**How to Get Started (on Ubuntu/Debian):**
First, you need to install the model runner plugin:
'sudo apt-get update && sudo apt-get install docker-model-plugin'

Then, test it to make sure it's working:
'docker model version'

For more details, see the official docs: https://docs.docker.com/ai/model-runner/`,
	},
}

func HandleQuery(query string) (string, error) {
	query = strings.ToLower(query)
	var bestTopic *Topic
	bestScore := 0

	for i := range topics {
		topic := &topics[i]
		score := 0
		for _, keyword := range topic.Keywords {
			if strings.Contains(query, keyword) {
				score++
			}
		}

		if score > 0 {
			if bestTopic == nil || score > bestScore {
				bestScore = score
				bestTopic = topic
			} else if score == bestScore {
				// If scores are equal, prefer the topic with fewer keywords (more specific)
				if len(topic.Keywords) < len(bestTopic.Keywords) {
					bestTopic = topic
				}
			}
		}
	}

	if bestTopic != nil {
		return bestTopic.Content, nil
	}

	return "I'm sorry, I don't have information on that topic in my local learning database. Please try a different query or switch to AI mode.", nil
} 