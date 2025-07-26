# BBS+ Selective Disclosure Example

A complete implementation example of **Selective Disclosure** using **BBS+** signatures in the context of **Decentralized Identifiers (DIDs)**, **Verifiable Credentials (VCs)**, and **Verifiable Presentations (VPs)** using Golang.

## 🎯 Objective

This project demonstrates how selective disclosure works in practice, allowing:
- **Issuer** to create and sign Verifiable Credentials
- **Holder** to store credentials and create presentations that only disclose necessary information
- **Verifier** to check authenticity without seeing the hidden information

## 🏗️ Architecture

```
/bbs-selective-disclosure-example
├── cmd/
│   └── demo/                    # CLI demo application
├── pkg/
│   ├── bbs/                     # BBS+ cryptographic operations
│   ├── did/                     # DID management
│   └── vc/                      # Verifiable Credentials & Presentations
├── internal/
│   ├── issuer/                  # Issuer use cases
│   ├── holder/                  # Holder use cases
│   └── verifier/                # Verifier use cases
├── test/
│   ├── integration/             # Integration tests
│   └── unit/                    # Unit tests (future)
├── docs/                        # Documentation
└── Makefile                     # Build automation
```

## 🚀 Installation and Usage

### System Requirements
- Go 1.21+
- Make (optional, for using the Makefile)

### 1. Clone repository
```bash
git clone <repository-url>
cd bbs-selective-disclosure-example
```

### 2. Install dependencies
```bash
go mod download
# or
make deps
```

### 3. Run demo
```bash
# Method 1: Using Makefile
make run-demo

# Method 2: Direct `go run`
make demo

# Method 3: Build and run
make build
./bin/demo
```

### 4. Run tests
```bash
# Run all tests
make test

# Only integration tests
make test-integration

# Test with coverage report
make test-coverage
```

## 📝 Demo Scenario

The demo illustrates a real-world scenario:

### 🏛️ **Setup Phase**
1. **Government** (Issuer) creates a DID and BBS+ keys
2. **Citizen** (Holder) creates a personal DID
3. **Movie Theater** (Verifier) creates a DID for verification

### 📋 **Credential Issuance**
The government issues a "Digital ID Card" with the following information:
```json
{
    "firstName": "An",
    "lastName": "Nguyen Van", 
    "dateOfBirth": "2000-01-20",
    "nationality": "Vietnamese",
    "address": "123 Nguyen Trai St, Ho Chi Minh City",
    "idNumber": "123456789"
}
```

### 🎬 **Selective Disclosure**
The movie theater needs to verify:
- ✅ **Age** (to check for 18+)
- ✅ **Nationality** (to apply ticket pricing)
- ❌ **NOT needed**: Name, address, ID number

The citizen creates a VP disclosing only `dateOfBirth` and `nationality`.

### 🔍 **Verification**
The movie theater verifies:
- Data integrity
- Signature from a trusted issuer
- Only receives the permitted information

## 🧪 Test Coverage

### Integration Tests
- **Full Lifecycle**: Complete DID → VC → VP workflow
- **Multiple Credentials**: Presenting from multiple credentials
- **Verification Failures**: Untrusted issuers, missing claims
- **DID Operations**: Creation, resolution, validation
- **BBS+ Operations**: Key generation, signing, proof creation

### Running Specific Tests
```bash
# Test full lifecycle
go test -v ./test/integration -run TestFullLifecycle

# Test multiple credentials
go test -v ./test/integration -run TestMultipleCredentialsPresentation

# Test verification failures
go test -v ./test/integration -run TestVerificationFailures
```

## 🔧 Development

### Code Quality
```bash
# Format code
make fmt

# Vet code
make vet

# Lint (if golangci-lint is installed)
make lint

# Verify project structure
make verify
```

### Development Setup
```bash
# Setup complete development environment
make dev-setup
```

## 📚 Core Concepts

### 🔐 **BBS+ Signatures**
- Allows signing multiple messages
- Creates selective disclosure proofs
- Does not reveal hidden information

### 🆔 **Decentralized Identifiers (DIDs)**
- Decentralized identifiers
- Self-sovereign identity
- Cryptographic verification

### 📄 **Verifiable Credentials (VCs)**
- Digital attestations
- Tamper-evident
- Cryptographically secure

### 🎪 **Verifiable Presentations (VPs)**
- Selective disclosure container
- Privacy-preserving
- Zero-knowledge proofs

## 🔒 Privacy Features

### Selective Disclosure Benefits
1. **Minimum Data Disclosure**: Only reveal necessary information
2. **Privacy Protection**: Hide sensitive information
3. **Cryptographic Integrity**: Ensure data integrity
4. **Non-repudiation**: Cannot be denied

### Example Privacy Scenario
```
Original Credential:
├── firstName: "An"              [HIDDEN]
├── lastName: "Nguyen Van"       [HIDDEN] 
├── dateOfBirth: "2000-01-20"    [REVEALED] → Age: 25
├── nationality: "Vietnamese"     [REVEALED]
├── address: "123 Nguyen Trai"   [HIDDEN]
└── idNumber: "123456789"        [HIDDEN]

Verifier only sees:
✓ Age: 25 (calculated from dateOfBirth)
✓ Nationality: Vietnamese
✗ Does not know name, address, or ID number
```

## 🛠️ Tech Stack

- **Language**: Go 1.21+
- **Cryptography**: Ed25519, BBS+ (simplified implementation)
- **Testing**: testify/assert, testify/require
- **Build**: Make, Go modules
- **Architecture**: Clean Architecture, Domain-Driven Design

## ⚠️ Important Notes

### Production Considerations
1. **BBS+ Implementation**: Use a production-ready library like Hyperledger Aries for production.
2. **Key Management**: Implement secure key storage.
3. **DID Methods**: Use production DID methods (e.g., did:web, did:ion).
4. **Cryptographic Security**: Audit cryptographic implementations.

### Simplified Components
- BBS+ signing and proofs are simplified for demonstration purposes.
- DID resolution uses in-memory storage.
- Does not implement the full W3C VC/VP specifications.

## 📈 Future Enhancements

- [ ] Integration with Hyperledger Aries BBS+
- [ ] Support for multiple DID methods
- [ ] Zero-knowledge proof optimizations
- [ ] Web-based demo interface
- [ ] Performance benchmarks
- [ ] Production deployment guides

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Implement changes with tests
4. Run `make test` to verify
5. Submit a pull request

## 📄 License

[Add your license here]

## 📞 Support

If you have questions or need support, please create an issue in the repository.
