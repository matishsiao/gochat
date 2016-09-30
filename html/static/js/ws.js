Chat = function(container,user) {
  var message_store = [];
  var self = this;
  self.ws = null;
  self.hidden = false;
  self.user = user;
  self.channel = "Public";
  var msgbox = document.getElementById(container);
  if(msgbox == null){
    console.error("message box can't not be find.");
    return;
  }
  msgbox.className = "messagebox";

  Chat.prototype.CreateMessage = function(channel,message) {
      var msg = {
        user: self.user,
        timestamp: new Date().getTime(),
        channel:channel,
        message:message
      };
      return msg;
  }

  Chat.prototype.Add = function(message) {
    message_store.push(message);
    self.Render();
  }

  Chat.prototype.Init = function() {
    var welcomeMsg = self.CreateMessage("Global","Welcome to GoChat");
    welcomeMsg.user = "System";
    message_store.push(welcomeMsg);

    self.BoardRender();
  }
  Chat.prototype.BoardRender = function() {
    msgbox.innerHTML = '<div id="message_setting" class="message_setting">'+
    '<button id="switch_btn"> '+((self.hidden)?"+":"-")+' </button>'+
    '</div>';
    if(!self.hidden){
      msgbox.innerHTML += '<div id="message_ctx" class="message_ctx"></div>'+
      '<hr class="hr_line"></hr>'+
      '<div class="inputbox"><input type="text" size="20" name="send_box" id="send_box" placeholder="Your message">'+
      '<button type="button" id="send_btn" name="send_btn">Send</button></div>';
      document.getElementById("send_btn").addEventListener("click", self.Send);
      document.getElementById("send_box").addEventListener("keyup", self.Change);
    }

    document.getElementById("switch_btn").addEventListener("click", self.SwitchWindow);
    msgbox.style.top = window.innerHeight - msgbox.offsetHeight - 4 + 'px';
    msgbox.style.left = window.innerWidth - msgbox.offsetWidth - 4 + 'px';
    window.onresize = self.BoardRender;
    self.Render();
  }

  Chat.prototype.SwitchWindow = function() {
    self.hidden = (self.hidden)?false:true;
    if(self.hidden){
      msgbox.style.height = "26px";
    } else {
      msgbox.style.height = "300px";
    }
    self.BoardRender();
  }

  Chat.prototype.Render = function() {
    if(!self.hidden){
      var msg_ctx = "";
      for(idx in message_store) {
        msg_ctx += self.RenderMessage(message_store[idx]);
      }
      var msg_ctx_box = document.getElementById('message_ctx');
      msg_ctx_box.innerHTML = msg_ctx;
      msg_ctx_box.scrollTop = msg_ctx_box.scrollHeight
    }
  }

  Chat.prototype.RenderMessage = function(msg) {
    msg_ctx = '<div class="message">'+
               '<div class="message_title" style="'+((msg.user == self.user)?"background-color: #b0ff7b;":"background-color: #c6f104;")+'"><b>'+msg.user+'</b>:<label>'+
               timeConverter(msg.timestamp)+'</label></div>'+
               '<div class="message_content"><b>'+msg.message+'</div></div>';
    return msg_ctx;
  }

  Chat.prototype.Connect = function(url) {
    console.log(url);
    self.ws = new WebSocket(url);
    self.ws.onopen = self.onOpen;
    self.ws.onmessage = self.onMessage;
    self.ws.onerror = self.onError;
    self.ws.onclose = self.onClose;
  }
  Chat.prototype.onOpen = function(event) {
    console.log("onOpen:",event);
  }
  Chat.prototype.onMessage = function(event) {
    console.log("onMessage:",event);
    if(event.data != ""){
      var msg = JSON.parse(event.data);
      if(msg != null){
        console.log("receive Message:",msg);
        self.Add(msg);
      }
    }
  }
  Chat.prototype.onClose = function(event) {
    console.log("onClose:",event);
  }
  Chat.prototype.onError = function(event) {
    console.log("onError:",event);
  }

  Chat.prototype.Change = function(event) {
    var unicode = -1;
  	if(event != null){
  		unicode=event.keyCode? event.keyCode : event.charCode;
  	}
  	var message = document.getElementById("send_box").value;
  	if(message != ""){
      if(unicode == 13 || event == null){
  		    self.Send();
    	}
    }
  }

  Chat.prototype.Send = function() {
    var message = document.getElementById('send_box').value;
    console.log("Send",self.channel);
    var msg = self.CreateMessage(self.channel, message);
    if(self.ws != null && self.ws.readyState == self.ws.OPEN) {
      self.ws.send(JSON.stringify(msg));
    } else {
      msg = self.CreateMessage(self.channel,"message can't send. reason:no connection can used.");
    }
    document.getElementById('send_box').value = "";
    self.Add(msg);
  }

  Chat.prototype.MessageEncoder = function(channel,message) {
      var msg = {
        user: self.user,
        timestamp: new Date().getTime(),
        channel:channel,
        message:message
      };
      return JSON.stringify(msg);
  }



  Chat.prototype.Status = function() {
    switch(self.ws.readyState){
      case self.ws.OPEN:
        return true;
    }
    return false;
  }
  self.Init();
}
