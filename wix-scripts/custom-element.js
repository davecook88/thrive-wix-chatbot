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
 * @property {HTMLDivElement | null} loginPreview - The login preview element.
 * @property {HTMLDivElement | null} chatContainer - The chat container element.
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
        :host {
          display: block;
          font-family: Arial, sans-serif;
          font-size: 16px;
        }
        .chat-container {
          margin: 0 auto;
          background-color: #FFF6F3;
          border-radius: 10px;
          box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
          padding: 10px;
          max-width: 100%;
          box-sizing: border-box;
        }
        .chat-box {
          height: 70vh;
          overflow-y: auto;
          margin-bottom: 10px;
        }
        .message {
          display: flex;
          align-items: flex-start;
          margin-bottom: 15px;
        }
        .avatar {
          width: 40px;
          height: 40px;
          border-radius: 50%;
          overflow: hidden;
          margin-right: 10px;
          flex-shrink: 0;
        }
        .avatar img {
          width: 100%;
          height: 100%;
          object-fit: cover;
        }
        .message-content {
          background-color: #FFFFFF;
          border-radius: 18px;
          padding: 10px;
          max-width: calc(100% - 50px);
          box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
          word-wrap: break-word;
        }
        .user {
          flex-direction: row-reverse;
        }
        .user .message-content {
          background-color: #4A90E2;
          color: white;
        }
        .input-area {
          display: flex;
          flex-wrap: wrap;
        }
        .input-area textarea {
          flex: 1 1 200px;
          min-height: 40px;
          padding: 10px;
          border: 1px solid #E0E0E0;
          border-radius: 20px;
          font-size: 14px;
          resize: none;
          margin-right: 10px;
          margin-bottom: 10px;
        }
        .input-area button {
          flex: 0 0 auto;
          padding: 10px 20px;
          background-color: #4A90E2;
          color: white;
          border: none;
          border-radius: 20px;
          font-size: 16px;
          cursor: pointer;
        }
        .login-preview {
          display: none;
          background-color: #FFF6F3;
          border-radius: 10px;
          box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
          padding: 20px;
          text-align: center;
          cursor: pointer;
        }
        .login-preview h2 {
          margin-top: 0;
          color: #4A90E2;
        }
        .login-preview p {
          margin-bottom: 20px;
        }
        .login-preview button {
          background-color: #4A90E2;
          color: white;
          border: none;
          border-radius: 20px;
          padding: 10px 20px;
          font-size: 16px;
          cursor: pointer;
        }
        .login-preview button:hover {
          background-color: #3A7BC8;
        }
        .friendly-diego {
          width: 80px;
          height: 120px;
          margin: 0 auto 20px;
          border-radius: 50%;
          overflow: hidden;
        }
        .friendly-diego img {
          width: 100%;
          height: 100%;
          object-fit: contain;
        }
        @media (max-width: 480px) {
          :host {
            font-size: 14px;
          }
          .chat-box {
            height: 60vh;
          }
          .message-content {
            max-width: calc(100% - 40px);
          }
          .avatar {
            width: 30px;
            height: 30px;
          }
          .input-area textarea {
            flex: 1 1 100%;
            margin-right: 0;
          }
          .input-area button {
            flex: 1 1 100%;
          }
          .friendly-diego {
            width: 60px;
            height: 60px;
          }
            
        }
      </style>
      <div class="chat-container" id="chat-container">
        <div class="chat-box" id="chatBox"></div>
        <div class="input-area">
          <textarea id="userMessage" placeholder="Type your message..."></textarea>
          <button id="sendButton">Send</button>
        </div>
      </div>
      <div class="login-preview" id="loginPreview">
        <div class="friendly-diego">
          <img src="https://static.wixstatic.com/media/2c423c_41a985cc2ceb436995c0f9e83a44872c~mv2.png/v1/fill/w_860,h_866,al_c,q_90,usm_0.66_1.00_0.01,enc_auto/con-gorra-naranja.png" alt="Friendly Diego">
        </div>
        <div class="input-area">
          <textarea id="userMessage" placeholder="Sign in to chat to Diego"></textarea>
          <button id="signInButton">Sign In</button>
        </div>
      </div>
    `;
    this.loginPreview = this.shadowRoot.getElementById("loginPreview");
    this.chatBox = this.shadowRoot.getElementById("chatBox");
    this.userMessage = this.shadowRoot.getElementById("userMessage");
    this.sendButton = this.shadowRoot.getElementById("sendButton");
    this.chatContainer = this.shadowRoot.getElementById("chat-container");

    this.loginPreview.addEventListener(
      "click",
      this.handleLoginClick.bind(this)
    );

    this.userMessage.addEventListener("keydown", (e) => {
      if (e.key === "Enter" && !e.shiftKey) {
        e.preventDefault();
        this.handleSendMessage();
      }
    });

    this.sendButton.addEventListener(
      "click",
      this.handleSendMessage.bind(this)
    );

    this.messages.forEach((message) => this.addMessageToChatBox(message));
    this.messages = [];
    this.renderMessages();
  }

  hide() {
    if (
      !this.chatBox ||
      !this.userMessage ||
      !this.sendButton ||
      !this.loginPreview
    )
      return;
    this.chatBox.style.display = "none";
    this.userMessage.style.display = "none";
    this.sendButton.style.display = "none";
    this.loginPreview.style.display = "block";
    this.chatContainer.style.display = "none";
  }

  show() {
    if (
      !this.chatBox ||
      !this.userMessage ||
      !this.sendButton ||
      !this.loginPreview
    )
      return;

    this.chatBox.style.display = "block";
    this.userMessage.style.display = "block";
    this.sendButton.style.display = "block";
    this.loginPreview.style.display = "none";
    this.chatContainer.style.display = "block";
  }

  handleLoginClick() {
    // Dispatch a custom event that the parent application can listen for
    const event = new CustomEvent("loginRequested", {
      bubbles: true,
      composed: true,
    });
    this.dispatchEvent(event);
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

    // Replace plain URLs with <a> tags
    content = content.replace(
      /(https?:\/\/[^\s]+)/g,
      `<a href="$1" target="_blank">$1</a>`
    );

    // Replace newline characters with <br> tags
    content = content.replace(/\n/g, "<br>");

    return content;
  }

  renderMessages() {
    if (!this.messages?.length) {
      this.hide();
      this.userMessage.disabled = true;
      this.userMessage.placeholder = "Sign in to chat";
      this.sendButton.disabled = true;
      return;
    } else {
      this.show();
      this.userMessage.disabled = false;
      this.userMessage.placeholder = "Type your message...";
      this.sendButton.disabled = false;
    }
    this.chatBox.innerHTML = "";
    this.messages?.forEach((message) => {
      if (!message?.content?.trim()) return;
      const messageElement = document.createElement("div");
      messageElement.className = `${message.role} message`;
      if (message.role === "assistant") {
        messageElement.innerHTML = `
          <div class="avatar">
            <img src="https://static.wixstatic.com/media/2c423c_a8d3dafa446d4c72b591d49e83071b8d~mv2.jpg/v1/fill/w_78,h_75,al_c,q_80,usm_0.66_1.00_0.01,enc_auto/2c423c_a8d3dafa446d4c72b591d49e83071b8d~mv2.jpg" alt="Diego avatar">
          </div>
          <div class="message-content">
            <strong>Diego:</strong> ${this.formatMessage(message.content)}
          </div>
        `;
      } else {
        messageElement.innerHTML = `
          <div class="message-content">
            <strong>You:</strong> ${this.formatMessage(message.content)}
          </div>
        `;
      }
      this.chatBox.appendChild(messageElement);
    });
    if (this.messages[this.messages.length - 1].role === "user") {
      this.addAssistantMessage("Thinking...");
    }
    this.chatBox.scrollTop = this.chatBox.scrollHeight;
  }
}

customElements.define("placement-chat-element", PlacementChatElement);
