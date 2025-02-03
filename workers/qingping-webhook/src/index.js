export default {
  async fetch(request, env) {
    try {
      const url = new URL(request.url);
      const path = url.pathname;
      if (path === '/api/data') {
        // Return the saved device data
        const timestamp = Math.floor(Date.now() / 1000);

        return new Response(JSON.stringify(timestamp || {}), {
          status: 200,
          headers: { 'Content-Type': 'application/json' },
        });
      } else if (path === '/api/data-push') {
        try {
          const body = await request.json();
          console.log('✅ Received body:', body);

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

          return new Response(JSON.stringify({ status: 'ok' }), {
            headers: { 'Content-Type': 'application/json' }
          });
        } catch (error) {
          console.error('Error processing request:', error);
          return new Response('Invalid JSON', { status: 400 });
        }
      }
    } catch (error) {
      return new Response(JSON.stringify({ error: error.message }), {
        status: 500,
        headers: { 'Content-Type': 'application/json' },
      });
    }
  },
};