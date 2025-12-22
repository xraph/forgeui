/**
 * ForgeUI Bridge - JavaScript client for calling Go functions
 * @version 1.0.0
 * @license MIT
 */

class ForgeBridge {
  constructor(config = {}) {
    this.config = {
      endpoint: config.endpoint || '/api/bridge',
      timeout: config.timeout || 30000,
      maxRetries: config.maxRetries || 3,
      retryDelay: config.retryDelay || 1000,
      csrf: config.csrf || null,
      ...config
    };

    this.requestId = 0;
  }

  /**
   * Call a single function
   * @param {string} method - Function name
   * @param {object} params - Function parameters
   * @returns {Promise<any>} - Function result
   */
  async call(method, params = {}) {
    const request = {
      jsonrpc: '2.0',
      id: String(++this.requestId),
      method,
      params
    };

    const response = await this._sendRequest(request);
    
    if (response.error) {
      throw new BridgeError(response.error);
    }

    return response.result;
  }

  /**
   * Call multiple functions in a batch
   * @param {Array<{method: string, params: object}>} calls - Array of calls
   * @returns {Promise<Array<any>>} - Array of results
   */
  async callBatch(calls) {
    const requests = calls.map((call, index) => ({
      jsonrpc: '2.0',
      id: String(++this.requestId),
      method: call.method,
      params: call.params || {}
    }));

    const responses = await this._sendRequest(requests, true);
    
    return responses.map(response => {
      if (response.error) {
        throw new BridgeError(response.error);
      }
      return response.result;
    });
  }

  /**
   * Stream results from a function using SSE
   * @param {string} method - Function name
   * @param {object} params - Function parameters
   * @param {function} onData - Callback for each data chunk
   * @param {function} onError - Callback for errors
   * @returns {function} - Cleanup function
   */
  stream(method, params = {}, onData, onError) {
    const url = new URL(`${this.config.endpoint}/stream`, window.location.origin);
    url.searchParams.set('method', method);
    url.searchParams.set('params', JSON.stringify(params));

    const eventSource = new EventSource(url.toString());

    eventSource.onmessage = (event) => {
      try {
        const chunk = JSON.parse(event.data);
        
        if (chunk.error) {
          if (onError) onError(new BridgeError(chunk.error));
          eventSource.close();
          return;
        }

        if (onData) onData(chunk.data);

        if (chunk.done) {
          eventSource.close();
        }
      } catch (err) {
        if (onError) onError(err);
        eventSource.close();
      }
    };

    eventSource.onerror = (err) => {
      if (onError) onError(err);
      eventSource.close();
    };

    // Return cleanup function
    return () => eventSource.close();
  }

  /**
   * Internal method to send HTTP requests
   * @private
   */
  async _sendRequest(data, isBatch = false) {
    const headers = {
      'Content-Type': 'application/json'
    };

    // Add CSRF token if configured
    if (this.config.csrf) {
      headers['X-CSRF-Token'] = this.config.csrf;
    }

    let lastError;
    for (let attempt = 0; attempt < this.config.maxRetries; attempt++) {
      try {
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), this.config.timeout);

        const response = await fetch(this.config.endpoint, {
          method: 'POST',
          headers,
          body: JSON.stringify(data),
          credentials: 'include',
          signal: controller.signal
        });

        clearTimeout(timeoutId);

        if (!response.ok) {
          throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }

        return await response.json();
      } catch (err) {
        lastError = err;

        // Don't retry on certain errors
        if (err.name === 'AbortError' || err.message.includes('401') || err.message.includes('403')) {
          throw err;
        }

        // Wait before retrying
        if (attempt < this.config.maxRetries - 1) {
          await new Promise(resolve => setTimeout(resolve, this.config.retryDelay * (attempt + 1)));
        }
      }
    }

    throw lastError;
  }

  /**
   * Set CSRF token
   * @param {string} token - CSRF token
   */
  setCSRF(token) {
    this.config.csrf = token;
  }

  /**
   * Get CSRF token from cookie
   * @returns {string|null} - CSRF token
   */
  static getCSRFFromCookie(cookieName = 'csrf_token') {
    const cookies = document.cookie.split(';');
    for (const cookie of cookies) {
      const [name, value] = cookie.trim().split('=');
      if (name === cookieName) {
        return decodeURIComponent(value);
      }
    }
    return null;
  }
}

/**
 * Bridge error class
 */
class BridgeError extends Error {
  constructor(error) {
    super(error.message);
    this.name = 'BridgeError';
    this.code = error.code;
    this.data = error.data;
  }
}

// Export for module systems
if (typeof module !== 'undefined' && module.exports) {
  module.exports = { ForgeBridge, BridgeError };
}

// Global export
if (typeof window !== 'undefined') {
  window.ForgeBridge = ForgeBridge;
  window.BridgeError = BridgeError;
}

