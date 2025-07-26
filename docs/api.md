# BBS+ Selective Disclosure API Documentation

## Overview

This API provides endpoints for demonstrating BBS+ selective disclosure with verifiable credentials. The API follows REST principles and uses JSON for data exchange.

## Base URL

```
http://localhost:8089
```

## Authentication

Currently, no authentication is required for the demo endpoints.

## Content Type

All endpoints expect and return `application/json` content type.

## CORS

CORS is enabled for all origins in development mode.

---

## Health Check

### GET /health

Returns the health status of the API service.

**Response:**
```json
{
  "status": "healthy",
  "service": "BBS+ Selective Disclosure API",
  "version": "1.0.0"
}
```

---

## Issuer API

### POST /api/issuer/setup

Setup a new issuer with DID and cryptographic keys.

**Request Body:**
```json
{
  "method": "example"
}
```

**Response:**
```json
{
  "did": "did:example:issuer123",
  "status": "success"
}
```

### POST /api/issuer/credentials

Issue a new verifiable credential.

**Request Body:**
```json
{
  "issuerDid": "did:example:issuer123",
  "subjectDid": "did:example:holder456",
  "claims": [
    {
      "key": "firstName",
      "value": "An"
    },
    {
      "key": "lastName", 
      "value": "Nguyen Van"
    },
    {
      "key": "dateOfBirth",
      "value": "2000-01-20"
    },
    {
      "key": "nationality",
      "value": "Vietnamese"
    },
    {
      "key": "address",
      "value": "123 Nguyen Trai St, Ho Chi Minh City"
    },
    {
      "key": "idNumber",
      "value": "123456789"
    }
  ]
}
```

**Response:**
```json
{
  "credentialId": "vc:example:credential789",
  "credential": {
    "@context": ["https://www.w3.org/2018/credentials/v1"],
    "id": "vc:example:credential789",
    "type": ["VerifiableCredential"],
    "issuer": "did:example:issuer123",
    "issuanceDate": "2025-07-27T00:42:17Z",
    "credentialSubject": {
      "id": "did:example:holder456",
      "firstName": "An",
      "lastName": "Nguyen Van",
      "dateOfBirth": "2000-01-20",
      "nationality": "Vietnamese",
      "address": "123 Nguyen Trai St, Ho Chi Minh City",
      "idNumber": "123456789"
    },
    "proof": {
      "type": "BbsBlsSignature2020",
      "created": "2025-07-27T00:42:17Z",
      "verificationMethod": "did:example:issuer123#key-1",
      "proofPurpose": "assertionMethod",
      "proofValue": "..."
    }
  }
}
```

---

## Holder API

### POST /api/holder/setup

Setup a new holder with DID.

**Request Body:**
```json
{
  "method": "example"
}
```

**Response:**
```json
{
  "did": "did:example:holder456",
  "status": "success"
}
```

### POST /api/holder/credentials

Store a received verifiable credential.

**Request Body:**
```json
{
  "credential": {
    "@context": ["https://www.w3.org/2018/credentials/v1"],
    "id": "vc:example:credential789",
    "type": ["VerifiableCredential"],
    "issuer": "did:example:issuer123",
    "issuanceDate": "2025-07-27T00:42:17Z",
    "credentialSubject": {
      "id": "did:example:holder456",
      "firstName": "An",
      "lastName": "Nguyen Van",
      "dateOfBirth": "2000-01-20",
      "nationality": "Vietnamese",
      "address": "123 Nguyen Trai St, Ho Chi Minh City",
      "idNumber": "123456789"
    },
    "proof": {
      "type": "BbsBlsSignature2020",
      "created": "2025-07-27T00:42:17Z",
      "verificationMethod": "did:example:issuer123#key-1",
      "proofPurpose": "assertionMethod",
      "proofValue": "..."
    }
  }
}
```

**Response:**
```json
{
  "status": "success"
}
```

### GET /api/holder/credentials/list?holderDid={did}

List all stored credentials for a holder.

**Query Parameters:**
- `holderDid` (required): The DID of the holder

**Response:**
```json
{
  "credentials": [
    {
      "@context": ["https://www.w3.org/2018/credentials/v1"],
      "id": "vc:example:credential789",
      "type": ["VerifiableCredential"],
      "issuer": "did:example:issuer123",
      "issuanceDate": "2025-07-27T00:42:17Z",
      "credentialSubject": {
        "id": "did:example:holder456",
        "firstName": "An",
        "lastName": "Nguyen Van",
        "dateOfBirth": "2000-01-20",
        "nationality": "Vietnamese",
        "address": "123 Nguyen Trai St, Ho Chi Minh City",
        "idNumber": "123456789"
      },
      "proof": {
        "type": "BbsBlsSignature2020",
        "created": "2025-07-27T00:42:17Z",
        "verificationMethod": "did:example:issuer123#key-1",
        "proofPurpose": "assertionMethod",
        "proofValue": "..."
      }
    }
  ]
}
```

### POST /api/holder/presentations

Create a selective disclosure presentation.

**Request Body:**
```json
{
  "holderDid": "did:example:holder456",
  "credentialIds": ["vc:example:credential789"],
  "selectiveDisclosure": [
    {
      "credentialId": "vc:example:credential789",
      "revealedAttributes": ["dateOfBirth", "nationality"]
    }
  ]
}
```

**Response:**
```json
{
  "presentationId": "vp:example:presentation101",
  "presentation": {
    "@context": ["https://www.w3.org/2018/credentials/v1"],
    "id": "vp:example:presentation101",
    "type": ["VerifiablePresentation"],
    "holder": "did:example:holder456",
    "verifiableCredential": [
      {
        "@context": ["https://www.w3.org/2018/credentials/v1"],
        "id": "vc:example:credential789",
        "type": ["VerifiableCredential"],
        "issuer": "did:example:issuer123",
        "issuanceDate": "2025-07-27T00:42:17Z",
        "credentialSubject": {
          "id": "did:example:holder456",
          "dateOfBirth": "2000-01-20",
          "nationality": "Vietnamese"
        },
        "proof": {
          "type": "BbsBlsSignatureProof2020",
          "created": "2025-07-27T00:42:17Z",
          "verificationMethod": "did:example:issuer123#key-1",
          "proofPurpose": "assertionMethod",
          "proofValue": "...",
          "nonce": "...",
          "revealedAttributes": [2, 3]
        }
      }
    ],
    "proof": {
      "type": "Ed25519Signature2020",
      "created": "2025-07-27T00:42:17Z",
      "verificationMethod": "did:example:holder456#key-1",
      "proofPurpose": "authentication",
      "proofValue": "..."
    }
  }
}
```

---

## Verifier API

### POST /api/verifier/setup

Setup a new verifier with DID.

**Request Body:**
```json
{
  "method": "example"
}
```

**Response:**
```json
{
  "did": "did:example:verifier789",
  "status": "success"
}
```

### POST /api/verifier/verify

Verify a verifiable presentation.

**Request Body:**
```json
{
  "presentation": {
    "@context": ["https://www.w3.org/2018/credentials/v1"],
    "id": "vp:example:presentation101",
    "type": ["VerifiablePresentation"],
    "holder": "did:example:holder456",
    "verifiableCredential": [
      {
        "@context": ["https://www.w3.org/2018/credentials/v1"],
        "id": "vc:example:credential789",
        "type": ["VerifiableCredential"],
        "issuer": "did:example:issuer123",
        "issuanceDate": "2025-07-27T00:42:17Z",
        "credentialSubject": {
          "id": "did:example:holder456",
          "dateOfBirth": "2000-01-20",
          "nationality": "Vietnamese"
        },
        "proof": {
          "type": "BbsBlsSignatureProof2020",
          "created": "2025-07-27T00:42:17Z",
          "verificationMethod": "did:example:issuer123#key-1",
          "proofPurpose": "assertionMethod",
          "proofValue": "...",
          "nonce": "...",
          "revealedAttributes": [2, 3]
        }
      }
    ],
    "proof": {
      "type": "Ed25519Signature2020",
      "created": "2025-07-27T00:42:17Z",
      "verificationMethod": "did:example:holder456#key-1",
      "proofPurpose": "authentication",
      "proofValue": "..."
    }
  },
  "requiredClaims": ["dateOfBirth", "nationality"],
  "trustedIssuers": ["did:example:issuer123"],
  "verificationNonce": "cinema-verification-1722041537"
}
```

**Response:**
```json
{
  "valid": true,
  "errors": [],
  "revealedClaims": {
    "dateOfBirth": "2000-01-20",
    "nationality": "Vietnamese"
  },
  "holderDid": "did:example:holder456",
  "issuerDids": ["did:example:issuer123"],
  "credentialTypes": ["VerifiableCredential"]
}
```

### POST /api/verifier/verification-request

Create a verification request template.

**Request Body:**
```json
{
  "requiredClaims": ["dateOfBirth", "nationality"],
  "trustedIssuers": ["did:example:issuer123"],
  "verificationNonce": "custom-nonce-123"
}
```

**Response:**
```json
{
  "requiredClaims": ["dateOfBirth", "nationality"],
  "trustedIssuers": ["did:example:issuer123"],
  "verificationNonce": "custom-nonce-123"
}
```

### GET /api/verifier/presentations?verifierDid={did}

List all verified presentations for a verifier.

**Query Parameters:**
- `verifierDid` (required): The DID of the verifier

**Response:**
```json
{
  "presentations": [
    {
      "@context": ["https://www.w3.org/2018/credentials/v1"],
      "id": "vp:example:presentation101",
      "type": ["VerifiablePresentation"],
      "holder": "did:example:holder456",
      "verifiableCredential": [...],
      "proof": {...}
    }
  ]
}
```

---

## Error Responses

All endpoints may return error responses in the following format:

```json
{
  "error": "Error message",
  "code": 400,
  "details": "Additional error details"
}
```

### Common HTTP Status Codes

- `200 OK`: Successful operation
- `400 Bad Request`: Invalid request body or parameters
- `405 Method Not Allowed`: HTTP method not supported for this endpoint
- `500 Internal Server Error`: Server-side error

---

## Demo Flow

### Complete Demo Flow Example

1. **Setup Issuer:**
   ```bash
   curl -X POST http://localhost:8089/api/issuer/setup \
     -H "Content-Type: application/json" \
     -d '{"method": "example"}'
   ```

2. **Setup Holder:**
   ```bash
   curl -X POST http://localhost:8089/api/holder/setup \
     -H "Content-Type: application/json" \
     -d '{"method": "example"}'
   ```

3. **Setup Verifier:**
   ```bash
   curl -X POST http://localhost:8089/api/verifier/setup \
     -H "Content-Type: application/json" \
     -d '{"method": "example"}'
   ```

4. **Issue Credential:**
   ```bash
   curl -X POST http://localhost:8089/api/issuer/credentials \
     -H "Content-Type: application/json" \
     -d '{
       "issuerDid": "did:example:issuer123",
       "subjectDid": "did:example:holder456",
       "claims": [
         {"key": "firstName", "value": "An"},
         {"key": "lastName", "value": "Nguyen Van"},
         {"key": "dateOfBirth", "value": "2000-01-20"},
         {"key": "nationality", "value": "Vietnamese"},
         {"key": "address", "value": "123 Nguyen Trai St, Ho Chi Minh City"},
         {"key": "idNumber", "value": "123456789"}
       ]
     }'
   ```

5. **Store Credential:** (Use credential from step 4)

6. **Create Presentation:**
   ```bash
   curl -X POST http://localhost:8089/api/holder/presentations \
     -H "Content-Type: application/json" \
     -d '{
       "holderDid": "did:example:holder456",
       "credentialIds": ["vc:example:credential789"],
       "selectiveDisclosure": [{
         "credentialId": "vc:example:credential789",
         "revealedAttributes": ["dateOfBirth", "nationality"]
       }]
     }'
   ```

7. **Verify Presentation:** (Use presentation from step 6)

---

## Privacy Features

### Selective Disclosure

The API demonstrates selective disclosure where:
- **Revealed**: Only necessary attributes (e.g., dateOfBirth, nationality)
- **Hidden**: Sensitive attributes (e.g., firstName, lastName, address, idNumber)
- **Integrity**: Cryptographic proof that hidden attributes exist and are valid
- **Trust**: Verification that credentials come from trusted issuers

### Example Privacy Scenario

```
Cinema Verification:
✅ Can verify: Age 18+ (from dateOfBirth)
✅ Can verify: Nationality (Vietnamese)
❌ Cannot see: First name, last name, address, ID number
✅ Can trust: Government issued the credential
✅ Can trust: All hidden attributes are valid and signed
```

This enables **zero-knowledge verification** where the verifier can confirm necessary properties without accessing sensitive personal information.
