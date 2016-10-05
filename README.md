# gochat

  GoChat is simplify web socket chat server with client.

# Version: 0.0.1

# Features

  Public channel chat

  Administrator support

  Private channel chat

## Support Command

  1. /? // You can find help

  2. /list // Get all online user list

  3. /kick username // kick user(Administrator only)

  4. /join "channel name" // join channel

  5. /leave "channel name" // leave channel

  6. /change "channel name" // change used channel
  
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

  Normal User

  ```
  <html>
    <link href="static/css/chat.css" rel="stylesheet">
    <body>
      <div id="message_box"></div>
    </body>
    <script src="static/js/chat.js"></script>
    <script src="static/js/util.js"></script>
    <script>
      var chat = new Chat("message_box","username");
      chat.Connect("ws://127.0.0.1:8080/socket");
    </script>
  </html>
  ```

  Administrator

  ```
  <html>
    <link href="static/css/chat.css" rel="stylesheet">
    <body>
      <div id="message_box"></div>
    </body>
    <script src="static/js/chat.js"></script>
    <script src="static/js/util.js"></script>
    <script>
      var chat = new Chat("message_box","username","authtoken");
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

  Browser

    Normal User

    ```
      http://127.0.0.1:8080/
    ```

    Administrator

    ```
      http://127.0.0.1:8080/admin.html
    ```

# Configuration
  ```
  {
    "debug":false,
    "mode":"chat", // gochat mode: chat(public chat), TODO:service(for customer service)
    "authtoken":"authtoken", // Administrator authorization token
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
