import { sendMessage, getMessages } from "backend/chat.web.js";
import { authentication } from "wix-members-frontend";
import { members } from "wix-members.v2";

/**
 * @typedef {Object} Message
 * @property {string} role - The role of the message sender. Can be either "user" or "assistant".
 * @property {string} content - The message content.
 */

$w.onReady(function () {
  // Load existing messages from local storage

  $w("#wixChat1")?.hide();
  try {
    $w("#quickActionBar1")?.hide();
  } catch {
    //
  }

  const chatElement = $w("#customElement1");

  getMessages().then((messages) => {
    chatElement.setAttribute("data-messages", JSON.stringify(messages));
  });
  chatElement.on("messageSent", async (e) => {
    const messages = await sendMessage(e.detail.content);

    chatElement.setAttribute("data-messages", JSON.stringify(messages));
  });

  authentication.onLogin(hideElements);
  authentication.onLogout(hideElements);
});

/**
*	Adds an event handler that runs when the element is clicked.
	[Read more](https://www.wix.com/corvid/reference/$w.ClickableMixin.html#onClick)
*	 @param {$w.MouseEvent} event
*/
export function onLoginButtonClick(event) {
  authentication.promptLogin();
}

export async function hideElements() {
  const chatbox = $w("#customElement1");
  const loginButton = $w("#login-button");
  const currentMember = await members.getCurrentMember();
  if (!currentMember) {
    chatbox.hide();
    loginButton.show();
  } else {
    chatbox.show();
    loginButton.hide();
  }
}
