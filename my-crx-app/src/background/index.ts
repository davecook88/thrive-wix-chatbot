import { RUNTIME_MESSAGES } from '../constants';

console.log('background is running');

chrome.runtime.onMessage.addListener(async (request) => {
  const apiClient = new PlacementBackendClient();
  switch (request.type) {
    case RUNTIME_MESSAGES.TEST:
      console.log('Received test message');
      const res = await apiClient.testApi();
      console.log('Test API response:', res);
      return res;
  }
});
