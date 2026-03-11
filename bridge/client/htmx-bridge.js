/**
 * ForgeUI Bridge - HTMX Integration
 *
 * Provides:
 * - Automatic CSRF token injection for HTMX requests
 * - bridgeFnURL() helper for generating bridge function URLs
 * - Custom "forge-bridge" HTMX extension for error handling
 * - bridge:error and bridge:success custom events
 */
(function () {
  'use strict';

  var BRIDGE_FN_PREFIX = window.BRIDGE_FN_ENDPOINT || '/api/bridge/fn/';

  /**
   * Generate a bridge function URL for use with hx-get/hx-post attributes.
   * @param {string} method - The bridge function name (e.g., "dashboard.getOverview")
   * @param {Object} [params] - Optional query parameters for GET requests
   * @returns {string} The full URL
   */
  window.bridgeFnURL = function (method, params) {
    var url = BRIDGE_FN_PREFIX + method;
    if (params && typeof params === 'object') {
      var query = Object.keys(params)
        .filter(function (k) { return params[k] !== undefined && params[k] !== null; })
        .map(function (k) { return encodeURIComponent(k) + '=' + encodeURIComponent(params[k]); })
        .join('&');
      if (query) {
        url += '?' + query;
      }
    }
    return url;
  };

  // Auto-inject CSRF token into HTMX requests
  if (typeof htmx !== 'undefined') {
    // Inject CSRF token from cookie into all HTMX requests
    document.body.addEventListener('htmx:configRequest', function (evt) {
      var csrfCookieName = window.BRIDGE_CSRF_COOKIE || 'csrf_token';
      var csrfHeaderName = window.BRIDGE_CSRF_HEADER || 'X-CSRF-Token';

      var match = document.cookie.match(new RegExp('(?:^|;\\s*)' + csrfCookieName + '=([^;]+)'));
      if (match) {
        evt.detail.headers[csrfHeaderName] = match[1];
      }
    });

    // Define the forge-bridge HTMX extension
    htmx.defineExtension('forge-bridge', {
      onEvent: function (name, evt) {
        if (name === 'htmx:responseError') {
          var detail = {
            status: evt.detail.xhr ? evt.detail.xhr.status : 0,
            message: evt.detail.xhr ? evt.detail.xhr.responseText : 'Request failed',
            target: evt.detail.elt
          };

          // Dispatch bridge:error event on the target element
          evt.detail.elt.dispatchEvent(new CustomEvent('bridge:error', {
            bubbles: true,
            detail: detail
          }));

          // Also dispatch on document for global error handling
          document.dispatchEvent(new CustomEvent('bridge:error', {
            detail: detail
          }));
        }

        if (name === 'htmx:afterSwap') {
          evt.detail.elt.dispatchEvent(new CustomEvent('bridge:success', {
            bubbles: true,
            detail: { target: evt.detail.elt }
          }));
        }
      }
    });
  }
})();
