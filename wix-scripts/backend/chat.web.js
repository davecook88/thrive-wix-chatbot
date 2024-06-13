import { Permissions, webMethod } from "wix-web-module";
import { currentMember } from "wix-members-backend";
import { fetch } from "wix-fetch";

const API_URL =
  "https://9023-2806-2f0-63c0-9be8-f42f-125-e165-1343.ngrok-free.app";

/**
 * @typedef {Object} Message
 * @property {string} role - The role of the message sender. Can be either "user" or "assistant".
 * @property {string} content - The message content.
 */

/**
 * Sends a message to the chat API.
 * @param {string} message - The message content.
 * @returns {Promise<Message[]>} The updated list of messages.
 */
export const sendMessage = webMethod(
  Permissions.SiteMember,
  async (message) => {
    const member = await currentMember.getMember();

    const payload = { message };
    const res = await fetch(`${API_URL}/chat`, {
      body: JSON.stringify(payload),
      method: "POST",
      headers: {
        "X-Wix-Member-ID": member.contactId,
      },
    });
    /** @type {Message[]} */
    const resJson = await res.json();
    if (res.ok) {
      if (Array.isArray(resJson)) {
        return resJson;
      }
    } else {
      console.log("error", res.text());
      return "error";
    }
  }
);

/**
 * Retrieves the list of messages from the chat API.
 * @returns {Promise<Message[]>} The list of messages.
 */
export const getMessages = webMethod(Permissions.SiteMember, async () => {
  const member = await currentMember.getMember();

  const res = await fetch(`${API_URL}/chat`, {
    method: "GET",
    headers: {
      "X-Wix-Member-ID": member.contactId,
    },
  });
  /** @type {Message[]} */
  const resJson = await res.json();
  console.log(resJson);
  if (res.ok) {
    return resJson;
  } else {
    console.log("error", res.text());
    return [];
  }
});
