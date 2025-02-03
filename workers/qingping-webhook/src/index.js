import { deviceData } from './deviceData.js';

export default {
  async fetch(request, env) {
    try {
      // Route requests
      const url = new URL(request.url);
      const path = url.pathname;

      if (path === '/api/data') {
        return new Response(deviceData, {
          status: 200,
          headers: { 'Content-Type': 'application/json' },
        });
      } else if (path === '/api/data-push') {
        // Signature validation
        try {
          const body = await request.json();

          console.log('Received request body:', body);
          // Signature validation
          const timestamp = body.signature.timestamp;
          const token = body.signature.token;
          const receivedSig = body.signature.signature;

          // Validate timestamp (5-minute window)
          const now = Math.floor(Date.now() / 1000);
          if (now - timestamp > 300) {
            console.log('Request expired: Timestamp too old');
            return new Response('Expired request', { status: 400 });
          }
          // Generate signature
          const encoder = new TextEncoder();
          const secret = await crypto.subtle.importKey(
            'raw',
            encoder.encode(env.QP_APP_SECRET),
            { name: 'HMAC', hash: 'SHA-256' },
            false,
            ['sign']
          );
          const data = encoder.encode(timestamp + token);
          const signature = await crypto.subtle.sign('HMAC', secret, data);

          const expectedSig = Array.from(new Uint8Array(signature))
            .map(b => b.toString(16).padStart(2, '0'))
            .join('');

          console.log('Received signature:', receivedSig);
          console.log('Expected signature:', expectedSig);
          // Compare signatures
          if (receivedSig !== expectedSig) {
            console.log('Signature mismatch: Invalid request');
            return new Response('Invalid signature', { status: 401 });
          }
          // Process payload
          deviceData = body.payload;
          console.log('Received device data:', deviceData);
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