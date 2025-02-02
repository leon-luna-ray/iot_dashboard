export default {
  async fetch(request, env) {
    if (request.method !== 'GET') {
      return new Response('Method not allowed', { status: 405 });
    }
    
    let body;
    try {
      body = await request.json();
    } catch (err) {
      console.error("Failed to parse JSON:", err);
      return new Response("Bad Request: Invalid JSON", { status: 400 });
    }
    
    console.log('Received request body:', body);
    
    // Signature validation
    try {
      const timestamp = body.signature.timestamp;
      // const token = body.signature.token;
      // const receivedSig = body.signature.signature;
    
      // console.log('Received signature:', receivedSig);
    
      // // Validate timestamp (5-minute window)
      // const now = Math.floor(Date.now() / 1000);
      // if (now - timestamp > 300) {
      //   console.log('Request expired: Timestamp too old');
      //   return new Response('Expired request', { status: 400 });
      // }
    
      // // Generate signature
      // const encoder = new TextEncoder();
      // const key = await crypto.subtle.importKey(
      //   'raw',
      //   encoder.encode(env.QP_APP_SECRET),
      //   { name: 'HMAC', hash: 'SHA-256' },
      //   false,
      //   ['verify']
      // );
    
      // const data = encoder.encode(timestamp + token);
      // const signature = await crypto.subtle.sign('HMAC', key, data);
      // const expectedSig = Array.from(new Uint8Array(signature))
      //   .map(b => b.toString(16).padStart(2, '0'))
      //   .join('');
    
      // console.log('Generated expected signature:', expectedSig);
    
      // // Compare signatures
      // if (receivedSig !== expectedSig) {
      //   console.log('Signature mismatch: Invalid request');
      //   return new Response('Invalid signature', { status: 401 });
      // }
    
      // // Process payload
      // console.log('Received valid payload:', body.payload);
    } catch (err) {
      console.error("Error processing signature:", err);
      return new Response("Internal Server Error", { status: 500 });
    }
    
    return new Response(JSON.stringify({ status: 'ok' }), {
      headers: { 'Content-Type': 'application/json' }
    });
  }
};