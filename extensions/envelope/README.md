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
