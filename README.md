<p style="color:red" align="center">
<pre style="color:#ff5f5f;font-weight:bold;text-align:center" align="center">
    ____  _          ____     
   / __ )(_)________/ / /_  __
  / __  / / ___/ __  / / / / /
 / /_/ / / /  / /_/ / / /_/ / 
/_____/_/_/   \__,_/_/\__, /  
                     /____/   
2.0!!
</pre>
</p>

# 🐦 A terminal based chat client.

Birdly is a **P2P** _(Peer-To-Peer)_ chat client that allows you to chat with your friends in a terminal. It uses **Websockets** to establish a connection between two clients and then allows them to chat with each other.

It's built using **Go** and some [charmbracelet's](https://github.com/charmbracelet/) libraries like **bubbletea** and **lipgloss**. 

## 📦 Installation

You can install Birdly by running the following command:

```bash
go get -u github.com/Matelz/Birdly

go build -o birdly <path-to-installation>
```

## 🚀 Usage

To start Birdly, you need to run the following command:

```bash
./birdly.exe <command> ...flags
```

| Command | Description |
|---------|-------|
| host    | Use to host a server, this **will** open a listener on the port you specified (default: 8080) |
| connect | Use to connect to a server, using this with no flags will connect to the localhost |

| Flag | Description |
|------|-------------|
| --host | Use to set the ip of the server you are connecting (default: 127.0.0.1) |
| --port | Use to set the port of the server you are connecting or hosting (default: 8080) |
| --name | Use to set the your username in the chat (default: anon) |

**Example:**
```bash
./birdly.exe host --name=john --port=25565
```

```bash
./birdly.exe connect --name=doe --host=192.168.0.5 --port=25565
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](https://choosealicense.com/licenses/mit/) file for details.

> disclaimer: this is not meant to be a finished product, it's just a project for me to learn more about Go and networking
