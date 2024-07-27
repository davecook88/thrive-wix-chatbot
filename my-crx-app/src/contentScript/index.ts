import App from './ContentScript.svelte';
let appContainer: HTMLDivElement | null = null;

let intervalId: number;
const appendElement = (
  targetElement: Element,
  {
    subject = 'No subject',
    sender = 'No sender',
    body = 'No body',
  }: {
    subject?: string;
    sender?: string;
    body?: string;
  },
) => {
  if (appContainer || !targetElement) {
    return;
  }
  // Step 1: Find the element

  if (targetElement) {
    // Step 2: Create a Svelte app container
    appContainer = document.createElement('div');

    // Step 3: Append the container to the found element
    targetElement.parentNode?.insertBefore(appContainer, targetElement);

    // Step 4: Instantiate the Svelte app
    new App({
      target: appContainer,
      props: {
        subject,
        sender,
        body,
      },
    });
    clearInterval(intervalId);
  } else {
    appContainer = null;
    console.error('Element with the specified data-message-id not found.');
  }
};

// Listen for changes in the DOM
const observer = new MutationObserver((mutations) => {
  mutations.forEach((mutation) => {
    if (mutation.type === 'childList') {
      // Check if an email view has been added
      const emailView = document.querySelector('div[role="main"]');
      if (!emailView) return;
      extractEmailInfo(emailView);
    }
  });
});

// Start observing the document body for changes
observer.observe(document.body, { childList: true, subtree: true });

function extractEmailInfo(emailView: Element) {
  if (!emailView) return;
  // Extract subject
  const subject = emailView.querySelector('h2[data-thread-perm-id]')?.textContent;

  // Extract sender
  const sender = emailView.querySelector('span[email]')?.getAttribute('email');

  // Extract body (this might need adjustment based on Gmail's current structure)
  const bodyElement = emailView.querySelector('div[data-message-id]');
  const body = bodyElement?.textContent;

  console.log('Subject:', subject);
  console.log('Sender:', sender);
  console.log('Body:', body);

  if (bodyElement)
    appendElement(bodyElement, {
      subject: subject || 'No subject',
      sender: sender || 'No sender',
      body: body || 'No body',
    });

  // Here you can do whatever you want with the extracted information
}

intervalId = window.setInterval(appendElement, 1000);

console.info('contentScript is running1');
