import { sendMessage, getMessages } from "backend/chat.web.js";

/**
 * @typedef {Object} Message
 * @property {string} role - The role of the message sender. Can be either "user" or "assistant".
 * @property {string} content - The message content.
 */

$w.onReady(function () {
  // Load existing messages from local storage

  const chatElement = $w("#customElement1");
  /**
   * Fetches messages and sets them as a data attribute on the chat element.
   * @returns {Promise<Message[]>}
   */
  getMessages().then((messages) => {
    chatElement.setAttribute("data-messages", JSON.stringify(messages));
  });
  chatElement.on("messageSent", async (e) => {
    const messages = await sendMessage(e.detail.content);

    chatElement.setAttribute("data-messages", JSON.stringify(messages));
  });

  // Handle form submission
});
