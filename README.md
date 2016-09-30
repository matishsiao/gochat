# gochat

  GoChat is simplify web socket chat server with client.

# Version: 0.0.1

# Features
  
  Public channel chat

## Todo

  Private channel chat

  Administrator support
  
  
# Screenshot
  
  Chat box
  
  ![Chat Box](https://github.com/matishsiao/gochat/blob/master/images/chatbox.png)
  
  Send Message
  
  ![Send Message](https://github.com/matishsiao/gochat/blob/master/images/sendmessage.png)
  
  Minimize chat box
  
  ![Chat Box Minimize](https://github.com/matishsiao/gochat/blob/master/images/minimize.png)
  
  Send chinese words to chat box
  
  ![Chat Box for Chinese](https://github.com/matishsiao/gochat/blob/master/images/chat.png)

# Example
  ```
  <html>
    <link href="static/css/chat.css" rel="stylesheet">
    <body>
      <div id="message_box"></div>
    </body>
    <script src="static/js/ws.js"></script>
    <script src="static/js/util.js"></script>
    <script>
      var chat = new Chat("message_box","username");
      chat.Connect("ws://127.0.0.1:8080/socket");
    </script>
  </html>
  ```
# Build
  ```
  go get https://github.com/matishsiao/gochat/
  cd $GOPATH/github.com/matishsiao/gochat/
  go build
  ```
# Run
  ```
  ./gochat
  ```
# Configuration
  ```
  {
    "debug":false,
    "http":"127.0.0.1:8080", // http listen host with port
    "https":"127.0.0.1:4443",// https listen host with port
    "ssl":{ //https ssl key and crt file settting
      "key":"ssl/test.key",
      "crt":"ssl/test.crt"
    },
    "timeout":120,// no used. todo
    "connection_limit":100//no used. todo
  }

  ```
  
  
  
