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
- deadline - signature is not valid after deadline (timestamp)
- domain_separator - to prevent replay attacks from other domains (bytes)

## How the envelope works

### Sign the message and send the envelope (usually on the frontend)

1. Create the message to sign
2. Create the nonce
3. Hash the message with the nonce (sha-256)
4. Sign the hash (ed25519)
5. Create the envelope
6. Encode the envelope to base64
7. Send the envelope in request header (X-Envelop) if the cckit gateway is used
8. Get envelope in the cckit gateway, decode it from base64 and pack it as a third chaincode argument

### Receive the envelope and verify the signature (chaincode)

1. Receive the envelope from request header
2. Decode the envelope from base64
3. Recreate the hash from the original message and the nonce (from envelope)
4. Verify the signature using pubkey from the envelope and the hash from the previous step

The sha-256 algorithm is used to hash the message.

The signature is generated using Ed25519.

## How to use the envelope on chaincode

```
r := router.New(name).Pre(envelope.Verify)
```
