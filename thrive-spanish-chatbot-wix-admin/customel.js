class PlacementChatElement extends HTMLElement {
  constructor() {
    super();
    this.attachShadow({ mode: "open" });
    this.messages = [];
  }

  static get observedAttributes() {
    return ["data-assistant-message", "data-messages"];
  }

  connectedCallback() {
    this.shadowRoot.innerHTML = `
       <style>
        .chat-box {
            height: 400px;
            overflow: auto;
            border: 1px solid #ccc;
            padding: 20px;
            background-color: #f9f9f9;
            border-radius: 10px;
            box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
        }

        .message {
            margin-bottom: 10px;
            padding: 12px;
            border-radius: 20px;
            max-width: 70%;
            word-wrap: break-word;
            line-height: 1.4;
            font-size: 14px;
        }

        .user {
            align-self: flex-end;
            background-color: #007bff;
            color: white;
            margin-left: auto;
        }

        .assistant {
            align-self: flex-start;
            background-color: #e9e9eb;
            color: #333;
        }

        .input-area {
            display: flex;
            margin-top: 20px;
        }

        .input-area input {
            flex: 1;
            padding: 12px;
            border: 1px solid #ccc;
            border-radius: 20px;
            font-size: 14px;
            outline: none;
        }

        .input-area button {
            padding: 12px 20px;
            background-color: #007bff;
            color: white;
            border: none;
            border-radius: 20px;
            cursor: pointer;
            margin-left: 10px;
            font-size: 14px;
            transition: background-color 0.3s ease;
        }

        .input-area button:hover {
            background-color: #0056b3;
        }
        </style>
        <div class="chat-box" id="chatBox"></div>
        <div class="input-area">
          <input type="text" id="userMessage" placeholder="Enter your message" />
          <button id="sendButton">Send</button>
        </div>
      `;

    this.chatBox = this.shadowRoot.getElementById("chatBox");
    this.userMessage = this.shadowRoot.getElementById("userMessage");
    this.sendButton = this.shadowRoot.getElementById("sendButton");

    this.sendButton.addEventListener(
      "click",
      this.handleSendMessage.bind(this)
    );

    // Process any pending messages that were queued before the chatBox was available
    this.messages.forEach((message) => this.addMessageToChatBox(message));
    this.messages = [];
  }

  attributeChangedCallback(name, oldValue, newValue) {
    switch (name) {
      case "data-assistant-message":
        this.addAssistantMessage(newValue);
        break;
      case "data-messages":
        this.messages = JSON.parse(newValue);
        this.renderMessages();
        break;
    }
  }

  handleSendMessage() {
    const userMessageContent = this.userMessage.value.trim();
    if (!userMessageContent) return;
    this.userMessage.value = "";

    const userMessage = { role: "user", content: userMessageContent };
    this.messages.push(userMessage);
    this.renderMessages();

    const event = new CustomEvent("messageSent", {
      detail: userMessage,
      bubbles: true,
      composed: true,
    });
    this.dispatchEvent(event);
  }

  addAssistantMessage(content) {
    const assistantMessage = { role: "assistant", content };
    this.messages.push(assistantMessage);
    this.renderMessages();
  }
  renderMessages() {
    this.chatBox.innerHTML = "";
    this.messages.forEach((message) => {
      const messageElement = document.createElement("div");
      messageElement.className = `${message.role} message`;
      messageElement.innerHTML = `
          <strong>${message.role === "user" ? "You" : "Assistant"}:</strong> ${
        message.content
      }
        `;
      this.chatBox.appendChild(messageElement);
    });
    this.chatBox.scrollTop = this.chatBox.scrollHeight;
  }
}

customElements.define("placement-chat-element", PlacementChatElement);
