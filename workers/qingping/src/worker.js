export default {
  async fetch(request, env) {
    if (request.method !== 'POST') {
      return new Response('Method not allowed', { status: 405 });
    }

    try {
      const body = await request.json();
      console.log('Received request body:', body);

      // Signature validation
      const timestamp = body.signature.timestamp;
      const token = body.signature.token;
      const receivedSig = body.signature.signature;

      console.log('Received signature:', receivedSig);

      // Validate timestamp (5-minute window)
      const now = Math.floor(Date.now() / 1000);
      if (now - timestamp > 300) {
        console.log('Request expired: Timestamp too old');
        return new Response('Expired request', { status: 400 });
      }

      // Generate signature
      const encoder = new TextEncoder();
      const key = await crypto.subtle.importKey(
        'raw',
        encoder.encode(env.QP_APP_SECRET),
        { name: 'HMAC', hash: 'SHA-256' },
        false,
        ['verify']
      );

      const data = encoder.encode(timestamp + token);
      const signature = await crypto.subtle.sign('HMAC', key, data);
      const expectedSig = Array.from(new Uint8Array(signature))
        .map(b => b.toString(16).padStart(2, '0'))
        .join('');

      console.log('Generated expected signature:', expectedSig);

      // Compare signatures
      if (receivedSig !== expectedSig) {
        console.log('Signature mismatch: Invalid request');
        return new Response('Invalid signature', { status: 401 });
      }

      // Process payload
      console.log('Received valid payload:', body.payload);

      return new Response(JSON.stringify({ status: 'ok' }), {
        headers: { 'Content-Type': 'application/json' }
      });

    } catch (err) {
      console.error('Error processing request:', err);
      return new Response('Invalid request', { status: 400 });
    }
  }
}
