const ws = new WebSocket("ws://localhost:8080/ws");

ws.onopen = function() {
  console.log("WebSocket connection established.");
};

ws.onmessage = function(event) {
  const data = JSON.parse(event.data);
  const message = data.message;
  const chatBox = document.getElementById("chat-box");
  chatBox.innerHTML += `<p>${message}</p>`;
};

const chatForm = document.getElementById("chat-form");
chatForm.addEventListener("submit", function(event) {
  event.preventDefault();
  const messageInput = document.getElementById("message-input");
  const message = messageInput.value;
  ws.send(JSON.stringify({ message: message }));
  messageInput.value = "";
});
