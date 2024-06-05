# ssh-chat

ssh-chat is a simple SSH chat application coded in Golang.

## Installation

To install the required dependencies, run the following commands:

```bash
go get "github.com/gliderlabs/ssh" # simplifies the SSH server creation
go get "golang.org/x/term"          # allows manipulation of the user terminal 
```

## Usage
To start the server, run the following command:
```bash
go run main.go
```
Then, open a terminal and type:
```bash
ssh username@localhost -p 3000
```
Replace username with your desired username.
This will connect you to the SSH chat server running on localhost using port 3000.(or any other port you set)

## Security

### Client Connection via Password
- When connecting to the SSH server, clients are prompted for a password.(that u can set)
- The server authenticates clients using these passwords.

### Server Generating SSH Key Pair
- To enhance security, the server can generate a pair of SSH keys:
  ```bash
  ssh-keygen -t rsa -b 2048
If you want to bypass host key checking, you can use the following command :
```bash
 ssh -o "StrictHostKeyChecking=no" -p 3000 username@localhost
````
