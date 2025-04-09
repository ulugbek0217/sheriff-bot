# Sheriff - telegram user bot

## What does Sheriff do?
#### - Sheriff Bot watches for group messages and collects users' data: user ID, first name, second name, username, phone number, about, personal channel id

## How to run?
#### 1. Install requirements:
- golang-migrate/migrate
- sqlc
- docker
- golang 1.24.1
- make

#### 2. Prepare needed tools:
- ```make prepare```
#### 3. Copy from .env.example to .env in config folder.
#### 4. Fill the .env with correct data.
#### 5. Install go requirements and run:
- ```go mod tidy```
- ```go run main.go```

<br><br>
## Notes
This bot uses postgresql as database inside docker container. Check postgresql container status, wether it is running or not, before running the bot. ADMINS field in .env.example should be filled with corresponding user ids who will have access to use the bot. It can be filled with more than 1 id separated by comma (,). The bot sends notification to the first ID when it is started.
