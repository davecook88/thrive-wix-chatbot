/**
 * @typedef {Object} Message
 * @property {string} role - The role of the message sender. Can be either "user" or "assistant".
 * @property {string} content - The message content.
 */

/**
 * A custom element that displays a chat interface.
 * @extends HTMLElement
 * @property {Message[]} messages - The messages to display in the chat.
 * @property {HTMLDivElement | null} chatBox - The chat box element.
 * @property {HTMLInputElement | null} userMessage - The user message input element.
 * @property {HTMLButtonElement | null} sendButton - The send button element.
 */
class PlacementChatElement extends HTMLElement {
  constructor() {
    super();
    this.attachShadow({ mode: "open" });
    /** @type {Message[]} */
    this.messages = [];
    this.userMessage = null;
    this.chatBox = null;
    this.sendButton = null;
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
  
          .input-area textarea {
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
            <textarea id="userMessage" placeholder="Enter your message"></textarea>
            <button id="sendButton">Send</button>
          </div>
        `;

    this.chatBox = this.shadowRoot.getElementById("chatBox");
    this.userMessage = this.shadowRoot.getElementById("userMessage");
    this.userMessage.addEventListener("keydown", (e) => {
      if (e.key === "Enter") {
        this.handleSendMessage();
      }
    });

    this.sendButton = this.shadowRoot.getElementById("sendButton");

    this.sendButton.addEventListener(
      "click",
      this.handleSendMessage.bind(this)
    );

    // Process any pending messages that were queued before the chatBox was available
    this.messages.forEach((message) => this?.addMessageToChatBox(message));
    this.messages = [];
    this.renderMessages();
  }

  hide() {
    this.style.display = "none";
  }

  show() {
    this.style.display = "block";
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

  /*
   * Adds a message to the chat box.
   * @param {string} content - The message content.
   */
  addAssistantMessage(content) {
    const assistantMessage = { role: "assistant", content };
    this.messages.push(assistantMessage);
    this.renderMessages();
  }

  /**
   * Adds a message to the chat box.
   * @param {string} content - The message to add.
   */
  formatMessage(content) {
    content = content.replace(/\*\*(.*?)\*\*/g, "<strong>$1</strong>");

    // Replace italic syntax with <em> tags
    content = content.replace(/\*(.*?)\*/g, "<em>$1</em>");

    // Replace links with <a> tags
    content = content.replace(
      /\[(.*?)\]\((.*?)\)/g,
      `<a href='$2' target="_blank">$1</a>`
    );

    // Replace newline characters with <br> tags
    content = content.replace(/\n/g, "<br>");

    return content;
  }

  renderMessages() {
    console.log("Rendering messages", this?.messages);
    if (!this.messages?.length) {
      this.hide();
      this.userMessage.disabled = true;
      this.userMessage.placeholder = "Sign in to chat";
      this.sendButton.disabled = true;
      return;
    } else {
      this.show();
      this.userMessage.disabled = false;
      this.userMessage.placeholder = "Enter your message";
      this.sendButton.disabled = false;
    }
    this.chatBox.innerHTML = "";
    this.messages?.forEach((message) => {
      if (!message?.content?.trim()) return;
      console.log("Rendering message", message);
      const messageElement = document.createElement("div");
      messageElement.className = `${message.role} message`;
      messageElement.innerHTML = `
            <strong>${
              message.role === "user" ? "You" : "Diego"
            }:</strong> ${this.formatMessage(message.content)}
          `;
      this.chatBox.appendChild(messageElement);
    });
    // if last message is a user message, add a temporary assistant message
    // with a loading indicator
    if (this.messages[this.messages.length - 1].role === "user") {
      this.addAssistantMessage("Thinking...");
    }
    this.chatBox.scrollTop = this.chatBox.scrollHeight;
  }
}

customElements.define("placement-chat-element", PlacementChatElement);
