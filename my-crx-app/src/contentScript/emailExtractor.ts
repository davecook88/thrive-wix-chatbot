// emailExtractor.ts

export function extractEmails(text: string): string[] {
  // Regular expression to match email addresses
  const emailRegex = /[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}/g;

  // Find all matches in the text
  const matches = text.match(emailRegex);

  // Return unique email addresses
  return matches ? Array.from(new Set(matches)) : [];
}

export function parseWixForwardedEmail(body: string): {
  senderName: string;
  senderEmail: string;
  recipientEmails: string[];
  wixSubmitterEmail: string;
} {
  const lines = body.split('\n');
  let senderName = '';
  let senderEmail = '';
  const recipientEmails: string[] = [];
  let wixSubmitterEmail = '';

  // Extract sender information from the first line
  const senderMatch = lines[0].match(/(.+) <(.+?)>/);
  if (senderMatch) {
    senderName = senderMatch[1].trim();
    senderEmail = senderMatch[2];
  }

  // Extract recipient emails
  const toLine = lines.find((line) => line.startsWith('to '));
  if (toLine) {
    const emails = extractEmails(toLine);
    recipientEmails.push(...emails);
  }

  // Extract Wix submitter email
  const wixLine = lines.find((line) => line.includes('Email:'));
  if (wixLine) {
    const wixEmail = extractEmails(wixLine);
    if (wixEmail.length > 0) {
      wixSubmitterEmail = wixEmail[0];
    }
  }

  return {
    senderName,
    senderEmail,
    recipientEmails,
    wixSubmitterEmail,
  };
}
