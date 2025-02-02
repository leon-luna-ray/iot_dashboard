let tokenCache = {
  token: null,
  expiresAt: 0,
};



async function fetchToken(env) {
  console.log('🔑 Fetching new token...');
  const authURL = env.QP_AUTH_API_BASE;
  const appKey = env.QP_APP_KEY;
  const appSecret = env.QP_APP_SECRET;

  const formData = new URLSearchParams();
  formData.append('grant_type', 'client_credentials');
  formData.append('scope', 'device_full_access');

  const authString = btoa(`${appKey}:${appSecret}`);
  console.log('Auth String:', authString);
  const response = await fetch(authURL, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
      'Authorization': `Basic ${authString}`,
    },
    body: formData.toString(),
  });

  console.log('Response:', response);

  if (!response.ok) {
    throw new Error(`Failed to get token: ${response.status}`);
  }

  const result = await response.json();
  tokenCache = {
    token: result.access_token,
    expiresAt: Date.now() + result.expires_in * 1000,
  };
  return tokenCache.token;
}

export default {
  async fetch(request, env) {
    const testResponse = new Response(tokenCache, {status: 200});
    try {
      // Check token validity
      if (!tokenCache.token || tokenCache.expiresAt < Date.now()) {
        await fetchToken(env); // Token refresh is now INSIDE the handler
      }

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