export let tokenCache = {
    token: null,
    expiresAt: 0,
};

export async function fetchToken(env) {
    console.log('🔑 Fetching new token...');
    const authURL = `${env.QP_AUTH_API_BASE}/token`;
    const appKey = env.QP_APP_KEY;
    const appSecret = env.QP_APP_SECRET;

    const authString = btoa(`${appKey}:${appSecret}`);
    const formData = new URLSearchParams();
    
    formData.append('grant_type', 'client_credentials');
    formData.append('scope', 'device_full_access');

    // Send request to get token
    const response = await fetch(authURL, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
            'Authorization': `Basic ${authString}`,
        },
        body: formData.toString(),
    });

    if (!response.ok) {
        throw new Error(`Failed to get token: ${response.status}`);
    }

    const result = await response.json();
    
    tokenCache = {
        token: result.access_token,
        expiresAt: Date.now() + result.expires_in * 1000,
    };
    console.log('Token fetched:', tokenCache);
    return tokenCache.token;
}
