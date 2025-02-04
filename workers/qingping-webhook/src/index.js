export default {
  async fetch(request, env) {
    console.log('📡 Request received');

    const url = new URL(request.url);
    const path = url.pathname;

    if (path === '/api/data' && request.method === 'GET') {
      // Existing GET handler
      return new Response(JSON.stringify({ message: 'Hello World!' }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' }
      });
    } else if (path === '/api/vi/data-push' && request.method === 'POST') {
      try {
        // Handle POST request
        // TODO: Verify the signature

        // const timestamp = body.signature.timestamp;
        // const token = body.signature.token;
        // const receivedSig = body.signature.signature;

        // const now = Math.floor(Date.now() / 1000);
        // if (now - timestamp > 3000) {
        //   console.log('Request expired: Timestamp too old');
        //   return new Response('Expired request', { status: 400 });
        // }

        // const encoder = new TextEncoder();
        // const secret = await crypto.subtle.importKey(
        //   'raw',
        //   encoder.encode(env.QP_APP_SECRET),
        //   { name: 'HMAC', hash: 'SHA-256' },
        //   false,
        //   ['sign']
        // );
        // const data = encoder.encode(timestamp + token);
        // const signature = await crypto.subtle.sign('HMAC', secret, data);

        // const expectedSig = Array.from(new Uint8Array(signature))
        //   .map(b => b.toString(16).padStart(2, '0'))
        //   .join('');

        // if (receivedSig !== expectedSig) {
        //   console.log('Signature mismatch: Invalid request');
        //   return new Response('Invalid signature', { status: 401 });
        // }

        // // Save the payload to env
        // env.DEVICE_DATA = body.payload;
        // console.log('Received device data:', env.DEVICE_DATA);

        const body = await request.json();
        console.log('✅ Received body:', body);

        return new Response(JSON.stringify({ message: 'Data received successfully!' }), {
          status: 200,
          headers: { 'Content-Type': 'application/json' }
        });
      } catch (error) {
        console.error('Error processing request:', error);
        return new Response('Invalid request', { status: 400 });
      }
    }

    // Return 404 for unmatched routes
    return new Response('Not Found', { status: 404 });
  },
};