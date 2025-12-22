/**
 * ForgeUI Bridge - Alpine.js integration
 * @requires Alpine.js v3+
 * @requires ForgeBridge
 */

function AlpineBridgePlugin(Alpine) {
  // Create global bridge instance
  const bridge = new ForgeBridge({
    endpoint: window.BRIDGE_ENDPOINT || '/api/bridge',
    csrf: ForgeBridge.getCSRFFromCookie()
  });

  // $bridge magic property
  Alpine.magic('bridge', () => ({
    /**
     * Call a bridge function
     * @param {string} method - Function name
     * @param {object} params - Function parameters
     * @returns {Promise<any>}
     */
    async call(method, params) {
      return bridge.call(method, params);
    },

    /**
     * Call multiple functions in batch
     * @param {Array} calls - Array of {method, params}
     * @returns {Promise<Array>}
     */
    async batch(calls) {
      return bridge.callBatch(calls);
    },

    /**
     * Stream results from a function
     * @param {string} method - Function name
     * @param {object} params - Function parameters
     * @param {function} onData - Data callback
     * @param {function} onError - Error callback
     * @returns {function} - Cleanup function
     */
    stream(method, params, onData, onError) {
      return bridge.stream(method, params, onData, onError);
    }
  }));

  // $go shorthand for bridge.call
  Alpine.magic('go', () => (method, params) => bridge.call(method, params));

  // $goBatch shorthand for bridge.callBatch
  Alpine.magic('goBatch', () => (calls) => bridge.callBatch(calls));

  // $goStream shorthand for bridge.stream
  Alpine.magic('goStream', () => (method, params, onData, onError) => 
    bridge.stream(method, params, onData, onError)
  );

  // x-bridge directive for declarative calls
  Alpine.directive('bridge', (el, { expression, modifiers }, { evaluateLater, effect }) => {
    const event = modifiers.length > 0 ? modifiers[0] : 'click';
    
    el.addEventListener(event, async () => {
      try {
        // Parse expression as method:params
        const [method, ...paramsParts] = expression.split(':');
        const paramsStr = paramsParts.join(':');
        
        let params = {};
        if (paramsStr) {
          // Evaluate params expression
          const getParams = evaluateLater(paramsStr);
          getParams(value => params = value);
        }

        // Add loading state
        el.disabled = true;
        el.classList.add('loading');

        // Call bridge
        const result = await bridge.call(method.trim(), params);

        // Dispatch custom event with result
        el.dispatchEvent(new CustomEvent('bridge:success', { 
          detail: { method, params, result },
          bubbles: true
        }));

        // Remove loading state
        el.disabled = false;
        el.classList.remove('loading');
      } catch (err) {
        // Dispatch error event
        el.dispatchEvent(new CustomEvent('bridge:error', { 
          detail: { error: err },
          bubbles: true
        }));

        // Remove loading state
        el.disabled = false;
        el.classList.remove('loading');

        console.error('Bridge call failed:', err);
      }
    });
  });

  // Bridge store for global state
  Alpine.store('bridge', {
    loading: false,
    error: null,
    lastResult: null,

    async call(method, params) {
      this.loading = true;
      this.error = null;

      try {
        const result = await bridge.call(method, params);
        this.lastResult = result;
        return result;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    clearError() {
      this.error = null;
    }
  });

  // Helper for form submission
  Alpine.magic('bridgeForm', () => ({
    async submit(formEl, method) {
      const formData = new FormData(formEl);
      const params = Object.fromEntries(formData.entries());

      return bridge.call(method, params);
    }
  }));
}

// Auto-register with Alpine if available
if (typeof window !== 'undefined' && window.Alpine) {
  window.Alpine.plugin(AlpineBridgePlugin);
}

// Export for module systems
if (typeof module !== 'undefined' && module.exports) {
  module.exports = AlpineBridgePlugin;
}

// Global export
if (typeof window !== 'undefined') {
  window.AlpineBridgePlugin = AlpineBridgePlugin;
}

