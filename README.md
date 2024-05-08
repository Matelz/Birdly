<p style="color:red" align="center">
<pre style="color:#ff5f5f;font-weight:bold;text-align:center" align="center">
    ____  _          ____     
   / __ )(_)________/ / /_  __
  / __  / / ___/ __  / / / / /
 / /_/ / / /  / /_/ / / /_/ / 
/_____/_/_/   \__,_/_/\__, /  
                     /____/   
</pre>
</p>

# 🐦 A terminal based chat client.

Birdly is a **P2P** _(Peer-To-Peer)_ chat client that allows you to chat with your friends in a terminal. It uses **TCP** sockets to establish a connection between two clients and then allows them to chat with each other.

It is built using **Go** and some [charmbracelet's](https://github.com/charmbracelet/) libraries like **bubbletea** and **lipgloss** witch makes it more **_gramourous_**.

## 📦 Installation

You can install Birdly by running the following command:

```bash
go get -u github.com/Matelz/Birdly

go build -o birdly <path-to-installation>
```

## 🚀 Usage

To start Birdly, you need to run the following command:

```bash
./birdly
```

or by running the **.exe**

To host a server you need to choose the **Host** option on the main menu, it will **print your IP and Port** to the chat so your friends can connect.

To connect to a server you need to choose the **Connect** option on the main menu, then you need to input the **IP and Port** of the server you want to connect to.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](https://choosealicense.com/licenses/mit/) file for details.

> disclaimer: this is not meant to be a finished product, it's just a project for me to learn more about Go and networking
