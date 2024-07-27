<script lang="ts">
  import '../app.css';
  import { extractEmails } from './emailExtractor';
  import CustomerInfoToolbar from './CustomInfoToolbar.svelte';
  import { CustomerInformation } from './types';
  import { RUNTIME_MESSAGES } from '../constants';

  export let subject: string | undefined;
  export let sender: string | undefined;
  export let body: string | undefined;

  $: emailAddresses = body ? extractEmails(body) : [sender || ''];

  // Mock data for the toolbar (replace with actual data from your Wix CRM)
  const customerInfo: CustomerInformation = {
    spanishLevel: 'Intermediate',
    notes: 'Prefers morning classes. Interested in business Spanish.',
    upcomingClasses: [
      { date: '2023-06-28', topic: 'Spanish for Business' },
      { date: '2023-07-05', topic: 'Advanced Conversation' },
    ],
  };

  chrome.runtime.sendMessage({ type: RUNTIME_MESSAGES.TEST }, (response) => {
    console.log('Content script received response:', response);
  });
</script>

<div class="w-full">
  <CustomerInfoToolbar
    {emailAddresses}
    spanishLevel={customerInfo.spanishLevel}
    notes={customerInfo.notes}
    upcomingClasses={customerInfo.upcomingClasses}
  />
</div>
