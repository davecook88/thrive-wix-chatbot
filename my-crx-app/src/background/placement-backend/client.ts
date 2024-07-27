class PlacementBackendClient {
  token: string;
  baseUrl: string;
  apiAuthHeader = 'X-Chrome-Ext-Token';
  constructor() {
    if (!process.env.API_TOKEN) {
      throw new Error('API token not found');
    }
    this.token = process.env.API_TOKEN;
    if (!process.env.API_BASE_URL) {
      throw new Error('API base URL not found');
    }
    this.baseUrl = process.env.API_BASE_URL;
  }

  private callApi<T>(path: string, method: 'GET' | 'POST', body?: any): Promise<T> {
    return fetch(`${this.baseUrl}/ext${path}`, {
      method,
      headers: {
        'Content-Type': 'application/json',
        [this.apiAuthHeader]: this.token,
      },
      body: body ? JSON.stringify(body) : undefined,
    }).then((response) => {
      if (!response.ok) {
        throw new Error(`API request failed: ${response.statusText}`);
      }
      return response.json();
    });
  }

  testApi(): Promise<{ message: string }> {
    return this.callApi<{ message: string }>('/test', 'GET');
  }
}
