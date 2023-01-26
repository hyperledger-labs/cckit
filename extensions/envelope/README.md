# Hyperledger Fabric chaincode kit (CCKit) - Envelope extension

The **envelope extension** allows to pass additional data with the request to check user signature and protect from replay attacks.
The extension is implemented as cckit middleware.

## Envelope structure

The envelope is a structure that consists of the following attributes:

- public_key - signer public key (bytes)
- signature - payload signature (bytes)
- nonce - number is given for replay protection (string)
- hash_to_sign - payload hash (bytes)
- hash_func - function used for hashing (string)
- deadline - signature is not valid after deadline, is not mandatory (timestamp)
- channel - chaincode channel
- chaincode - name of chaincode
- method - name of chaincode method

Channel + chaincode + method are used as domain separator to prevent replay attack from other domains (EIP-2612).

## How the envelope works

### Sign the message and send the envelope (usually on the frontend)

1. Create the nonce
2. Create the message to sign (include message payload and other attributes such as nonce, channel, chaincode, etc.)
3. Hash the message (sha-256)
4. Sign the hash (ed25519)
5. Create the envelope
6. Encode the envelope to base64
7. Send the envelope in request header (X-Envelop) if the cckit gateway is used

### Receive and transform the envelope (cckit gateway)

1. Receive the envelope from request header
2. Decode it from base64 and pack it as a third chaincode argument

### Verify the signature (chaincode)

1. Get the envelope as third argument
2. Unmarshall the envelope
3. Check envelope attributes (deadline, channel, chaincode, etc)
4. Check the nonce
5. Recreate the hash from the original message
6. Verify the signature using pubkey from the envelope and the hash from the previous step

The sha-256 algorithm is used to hash the message.

The signature is generated using Ed25519.

## How to use the envelope on chaincode

```
r := router.New(name).Use(envelope.Verify)
```

### Create an envelope using js

```
const { sign } = require('tweetnacl');
const { createHash } = require('crypto');

// CREATE ENVELOPE

function createEnvelope(
  payload,
  nonce,
  channel,
  chaincode,
  method,
  deadline,
  keys
) {
  const pk = Buffer.from(keys.publicKey).toString('hex');

  // make message to sign
  const b1 = Buffer.from(JSON.stringify(payload));
  const b2 = Buffer.from(nonce);
  const b3 = Buffer.from(channel);
  const b4 = Buffer.from(chaincode);
  const b5 = Buffer.from(method);
  const b6 = Buffer.from(deadline);
  const b7 = Buffer.from(pk);
  const arr = [b1, b2, b3, b4, b5, b6, b7];
  const bb = Buffer.concat(arr);

  // hash the message
  const hashed = createHash('sha256').update(bb).digest();

  // sign the hash
  const signature = sign.detached(hashed, keys.secretKey);

  // make the envelope
  const envelope = {
    hash_func: 'SHA256',
    hash_to_sign: hashed.toString('hex'),
    nonce: nonce,
    channel: channel,
    method: method,
    chaincode: chaincode,
    deadline: deadline,
    public_key: Buffer.from(keys.publicKey).toString('hex'),
    signature: Buffer.from(signature).toString('hex'),
  };

  return JSON.stringify(envelope);
}

// MAIN

const payload = {
  symbol: 'GLD',
  decimals: '8',
  name: 'Gold digital asset',
  type: 'DM',
  underlying_asset: 'gold',
  issuer_id: 'GLDINC',
};

const chaincode = 'envelope-chaincode';
const channel = 'envelope-channel';
const method = 'invokeWithEnvelope';
const nonce = String(new Date().getTime());
const deadline = new Date(new Date().getTime() + 86400000).toISOString(); // if without deadline then use new Date(0).toISOString()

const keys = sign.keyPair();

const envelope = createEnvelope(
  payload,
  nonce,
  channel,
  chaincode,
  method,
  deadline,
  keys
);
```
