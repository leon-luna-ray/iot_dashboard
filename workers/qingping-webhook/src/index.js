import { fetchToken, tokenCache } from './tokenManager.js';

export default {
  async fetch(request, env) {
    const testResponse = new Response('hello world', { status: 200 });
    try {
      // Check token validity
      if (!tokenCache.token || tokenCache.expiresAt < Date.now()) {
        console.log('🔴 Token expired or not present, fetching new token...');
        await fetchToken(env);
      }
      const token = tokenCache.token;
      console.log('Token is valid:', token);
      // Route requests
      const url = new URL(request.url);
      const path = url.pathname;

      if (path === '/api/data') {
        // const apiResponse = await fetch(env.QINGPING_API_URL, {
        //   headers: { Authorization: `Bearer ${tokenCache.token}` },
        // });
        // return new Response(apiResponse.body, {
        //   status: apiResponse.status,
        //   headers: { 'Content-Type': 'application/json' },
        // });
        return testResponse;
      } else if (path === '/api/data-push') {
        const payload = await request.json();
        console.log('Received data push:', payload);
        return new Response(JSON.stringify({ success: true }), { status: 200 });
      } else {
        return new Response('Not Found', { status: 404 });
      }
    } catch (error) {
      return new Response(JSON.stringify({ error: error.message }), {
        status: 500,
        headers: { 'Content-Type': 'application/json' },
      });
    }
  },
};