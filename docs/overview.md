# Overview: Selective Disclosure with BBS+, DID, VC, and VP in Golang

This project is a demonstration of how to use BBS+ Signatures to implement **Selective Disclosure** in the context of **Decentralized Identifiers (DIDs)**, **Verifiable Credentials (VCs)**, and **Verifiable Presentations (VPs)**.

## Objective

The main goal is to build a complete workflow in Golang, where:
1.  An **Issuer** creates a Verifiable Credential containing multiple claims.
2.  The **Issuer** signs this credential using the BBS+ algorithm.
3.  A **Holder** receives the credential and can create a Verifiable Presentation.
4.  In the Presentation, the **Holder** selectively discloses certain information to a **Verifier** without revealing the rest, while still ensuring the integrity and authenticity of the data.

This is a core feature for protecting privacy in decentralized identity systems.

## Key Concepts

*   **BBS+ Signatures**: A type of digital signature that allows signing a set of messages. Its special feature is that the signature holder can generate a proof for a subset of these messages without revealing the remaining ones.
*   **Decentralized Identifier (DID)**: A globally unique identifier that does not depend on any centralized organization. DIDs are used to identify the Issuer, Holder, and Verifier.
*   **Verifiable Credential (VC)**: A tamper-evident digital certificate containing claims about a subject. Examples: Driver's license, student ID card, certificate...
*   **Verifiable Presentation (VP)**: A data structure created by the Holder to present one or more VCs to a Verifier. The VP is where Selective Disclosure takes place.

## Workflow

The workflow of this example involves three main roles: **Issuer**, **Holder**, and **Verifier**.

1.  **Setup**
    *   The Issuer generates a BBS+ key pair (public key, private key).
    *   The Issuer publishes their public key via a DID Document.

2.  **Issuance**
    *   The Holder requests a VC from the Issuer (e.g., a "Digital Resident Card").
    *   The Issuer creates a VC containing multiple claims, for example:
        *   `firstName: "An"`
        *   `lastName: "Nguyen Van"`
        *   `dateOfBirth: "2000-01-20"`
        *   `nationality: "Vietnamese"`
    *   The Issuer uses their **BBS+ private key** to sign this set of claims and sends the signed VC to the Holder.

3.  **Presentation (Selective Disclosure)**
    *   A Verifier (e.g., a movie theater) requires the Holder to prove they are **over 18 years old** and have **Vietnamese nationality**, without needing to know their name.
    *   The Holder creates a Verifiable Presentation (VP).
    *   Using the original BBS+ signature, the Holder generates a **derived proof** that includes only the `dateOfBirth` and `nationality` claims. The `firstName` and `lastName` claims are hidden.
    *   The Holder sends the VP (containing the disclosed claims and the derived proof) to the Verifier.

4.  **Verification**
    *   The Verifier receives the VP from the Holder.
    *   The Verifier retrieves the Issuer's **BBS+ public key** (via the Issuer's DID).
    *   The Verifier uses this public key to verify the proof in the VP.
    *   If the verification is successful, the Verifier can trust that the `dateOfBirth` and `nationality` claims are true and were issued by the Issuer, without knowing about the existence of other claims.

## Project Structure (Proposed)

The project will be organized into functional modules:

```
/bbs-selective-disclosure-example
├── cmd/                # Entrypoints to run the example (CLI)
│   └── simple-flow/
├── pkg/
│   ├── bbs/            # Core logic for creating/verifying BBS+ signatures and proofs
│   ├── did/            # Utilities for simulating DID creation and resolution
│   └── vc/             # Structs and functions for working with VCs/VPs
├── internal/
│   ├── issuer/         # Logic for the Issuer role
│   ├── holder/         # Logic for the Holder role
│   └── verifier/       # Logic for the Verifier role
└── go.mod
```

This project will use appropriate Golang libraries to handle BBS+ cryptography and data structures related to DIDs/VCs.
