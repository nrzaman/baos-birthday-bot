# Baos Birthday Bot

This bot will automatically post birthday reminders of people in the How Bout Baos Discord server. It will post reminders at the beginning of every month, including all of the birthdays for that month, and on the day of each person's birthday.

## Usage

This bot was developed on macOS and is currently only supported on macOS and Linux operating systems.

### 1. Prerequisites
- Install [Homebrew](https://brew.sh/)
- Install `go` by typing the following in a Terminal window:
```
brew install go
```
- Install Docker by typing the following in a Terminal window:
```
brew install docker
```
- Install Colima (optional for macOS) by typing the following in a Terminal window:
```
brew install colima
```
- Install `kubectl` by typing the following in a Terminal window:
```
brew install kubectl
```
- Install `helm` by typing the following in a Terminal window:
```
brew install helm
```

#### Set environment variables
1. Grab a [Discord Bot Token](https://discordgsm.com/guide/how-to-get-a-discord-bot-token) from your Discord server.
2. Grab a [Discord Channel ID](https://support.discord.com/hc/en-us/articles/206346498-Where-can-I-find-my-User-Server-Message-ID#h_01HRSTXPS5FMK2A5SMVSX4JW4E) from the Discord channel that you'd like the bot to post reminders to.

### 2. Build and Run
```bash
# Build the code
go build -o bot .

# Run the bot
./bot -t [DISCORD_TOKEN]
```

### 3. Discord Slash Commands
#### `/month`
**Description:** List all birthdays in the current month

**Example:**
```
User: /month
Bot: Alice, March 15
     Bob, March 20
```

---

#### `/all`
**Description:** List all configured birthdays

**Example:**
```
User: /all
Bot: Alice, March 15
     Bob, March 20
     Charlotte, April 21
     ... (all birthdays)
```

---

#### `/next`
**Description:** Show the next upcoming birthday

**Example:**
```
User: /next
Bot: Next birthday: Alice on March 15 (in 3 days)
```

**Special cases:**
- If today is someone's birthday: `(Today! 🎉)`
- If tomorrow: `(Tomorrow!)`
- Multiple birthdays on same day: Shows all names

---

#### Backward Compatibility

The bot still supports legacy text commands:

| Slash Command | Legacy Command |
|---------------|----------------|
| `/month` | `!month` |
| `/all` | `!all` |
| `/next` | `!next` |

**Note:** Legacy commands will show a tip to use slash commands instead.

### 4. Deployment
Please note that this bot is currently deployed on an in-house server running a Kubernetes cluster.
The below steps assume a similar setup.

#### Prerequisites
1. Copy / create a Kubernetes config file to be able to use `kubectl`
2. Run the following command to confirm `kubectl` is running:
```bash
kubectl get po
```
3. Create `baos-birthday-bot-values.yml` in the root directory by typing:
```bash
vim baos-birthday-bot-values.yml
```
4. Add the following information to `baos-birthday-bot-values.yml` and save it
```bash
discord:
  # Discord bot token - REQUIRED
  token: "[DISCORD TOKEN]"
  # Discord channel ID where messages will be sent - REQUIRED
  channelId: "[DISCORD CHANNEL ID]"
```

#### Deployment Steps
1. Run the following commands in a Terminal window from the root `baos-birthday-bot` directory:

```bash
# Build the JAR file
make build

# Build the Docker image with a version tag
make docker-build VERSION=[VERSION NUMBER]

# Push the Docker image with a version tag
make docker-push VERSION=[VERSION NUMBER]
```

2. Update `./helm/values.yaml` under the `image:` section, next to `tag:` to match the version number specified above 
3. Run the following command to deploy:
```bash
helm upgrade --install aos-birthday-bot ./helm --values baos-birthday-bot.yml
```
4. Verify deployment by running the following command and observing that `baos-birthday-bot` exists in the list of services running and reads `1/1` under the `READY` column:
```bash
kubectl get po
```
5. Deployment may also be verified in Discord directly by typing one of the slash commands listed above in a channel

#### View Logs
1. Run the following command to grab the pod name:
```bash
kubectl get po
```
2. Copy the `baos-birthday-bot` pod name
3. Run the following command to view the logs:
```bash
kubectl logs baos-birthday-bot-[IDENTIFIER]
```

## Troubleshooting

### Building in VS Code, unrecognized dependencies
1. Open the Command Palette in VS Code by using the keyboard shortcut `Cmd + Shift + P`.
2. Run the following command: `Java: Reload Projects`.
3. If the above command does not resolve the issue, please try the following command: `Java: Clean Java Language Server Workspace`.

### Docker command not found

```bash
# Make sure Docker CLI is installed
brew install docker docker-compose

# Add to your shell profile (~/.zshrc or ~/.bashrc)
export PATH="/usr/local/bin:$PATH"
```

## Additional Resources

- [Docker Documentation](https://docs.docker.com/)
