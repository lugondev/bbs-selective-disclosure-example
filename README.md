# BBS+ Selective Disclosure Example

A complete implementation example of **Selective Disclosure** using **BBS+** signatures in the context of **Decentralized Identifiers (DIDs)**, **Verifiable Credentials (VCs)**, and **Verifiable Presentations (VPs)** using Golang.

## 🎯 Objective

This project demonstrates how selective disclosure works in practice, allowing:
- **Issuer** to create and sign Verifiable Credentials
- **Holder** to store credentials and create presentations that only disclose necessary information
- **Verifier** to check authenticity without seeing the hidden information

## ✨ New: BBS+ Interface System

This project now includes a **flexible BBS+ interface system** that allows switching between different implementations:

- **🔧 Simple Provider**: Basic implementation for testing and development
- **🛡️ Production Provider**: Full BLS12-381 cryptographic implementation
- **🏢 Aries Provider**: Hyperledger Aries Framework Go integration (ready for implementation)

### Quick Interface Demo

```bash
# Run the interface demo
make interface

# Try all interface features
make interface-all
```

📖 **[Complete Interface Documentation](docs/BBS_INTERFACE.md)**

## 🏗️ Architecture

```
/bbs-selective-disclosure-example
├── cmd/
│   ├── demo/                    # CLI demo application
│   ├── server/                  # HTTP server with web UI
│   └── interface_demo/          # BBS+ interface demonstration
├── interfaces/
│   └── http/                    # HTTP handlers and DTOs
├── web/                         # Web UI files
├── pkg/
│   ├── bbs/                     # BBS+ cryptographic operations & interfaces
│   ├── did/                     # DID management
│   └── vc/                      # Verifiable Credentials & Presentations
├── internal/
│   ├── issuer/                  # Issuer use cases
│   ├── holder/                  # Holder use cases
│   └── verifier/                # Verifier use cases
├── test/
│   ├── integration/             # Integration tests
│   └── unit/                    # Unit tests
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

### 3. Run BBS+ Interface Demo
```bash
# Run basic interface demo
make interface

# Run configuration demo
make interface-config

# Run provider switching demo
make interface-switching

# Run all interface demos
make interface-all

# Or using go run directly
go run ./cmd/interface_demo
go run ./cmd/interface_demo config
go run ./cmd/interface_demo switching
go run ./cmd/interface_demo all
```

### 4. Run HTTP Server with Web UI
```bash
# Method 1: Using Makefile
make run-server

# Method 2: Direct `go run`
make server

# Method 3: Build and run
make build-server
./bin/server

# Run on custom port
./bin/server -port 3000
```

The server will start on `http://localhost:8089` by default and provide:
- 🌐 **Web UI**: Interactive demo interface at `http://localhost:8089`
- 📡 **REST API**: HTTP endpoints at `http://localhost:8089/api/*`
- 🏥 **Health Check**: Status endpoint at `http://localhost:8089/health`

### 5. Run CLI Demo
```bash
# Method 1: Using Makefile
make run-demo

# Method 2: Direct `go run`
make demo

# Method 3: Build and run
make build
./bin/demo
```

### 5. Run tests
```bash
# Run all tests
make test

# Only integration tests
make test-integration

# Test with coverage report
make test-coverage
```

## 🌐 Web UI Features

The web interface provides an interactive demonstration of the BBS+ selective disclosure flow:

### 🎬 **Demo Flow**
1. **Setup Entities**: Initialize Issuer (Government), Holder (Citizen), and Verifier (Cinema)
2. **Issue Credential**: Government issues a Digital ID with multiple claims
3. **Create Presentation**: Citizen creates selective disclosure proof revealing only necessary information
4. **Verify**: Cinema verifies age and nationality without seeing personal details

### 🚀 **Quick Demo**
- Click "Run Full Demo" to execute the complete flow automatically
- Watch the execution logs to understand each step
- See how privacy is preserved through selective disclosure

### 🔧 **Manual Testing**
- Use individual sections to test specific scenarios
- Modify revealed attributes to see different privacy outcomes
- Test verification with different requirements and trusted issuers

## 📡 API Endpoints

### Issuer API
- `POST /api/issuer/setup` - Setup issuer with DID
- `POST /api/issuer/credentials` - Issue verifiable credential
- `POST /api/issuer/verify` - Verify credential

### Holder API
- `POST /api/holder/setup` - Setup holder with DID
- `POST /api/holder/credentials` - Store received credential
- `GET /api/holder/credentials/list` - List stored credentials
- `POST /api/holder/presentations` - Create selective disclosure presentation

### Verifier API
- `POST /api/verifier/setup` - Setup verifier with DID
- `POST /api/verifier/verify` - Verify presentation
- `POST /api/verifier/verification-request` - Create verification request
- `GET /api/verifier/presentations` - List verified presentations

### Utility API
- `GET /health` - Health check

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

## � BBS+ Interface Features

### Multiple Provider Support
- **🔧 Simple Provider**: Fast, basic implementation for testing
- **🛡️ Production Provider**: Secure BLS12-381 implementation for production
- **🏢 Aries Provider**: Hyperledger Aries Framework Go integration

### Key Capabilities
1. **Provider Switching**: Runtime switching between implementations
2. **Performance Metrics**: Built-in operation tracking and benchmarking
3. **Configuration Management**: Flexible configuration for different environments
4. **Migration Support**: Easy migration between providers
5. **Security Features**: Constant-time operations and secure memory management

### Usage Examples
```go
// Create production service
service, err := bbs.NewProductionBBSService()

// Switch providers at runtime
newService, err := bbs.SwitchProvider(
    currentService, 
    bbs.ProviderAries, 
    config,
)

// Compare providers
comparisons := bbs.CompareProviders()
for provider, info := range comparisons {
    fmt.Printf("%s: %s security, %t production ready\n", 
        provider, info.SecurityLevel, info.ProductionReady)
}
```

## �🔒 Privacy Features

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
- **Web UI**: HTML, CSS, JavaScript (Vanilla)
- **API**: REST with JSON
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
5. **HTTPS**: Use HTTPS in production environments.
6. **CORS**: Configure CORS properly for production.

### Simplified Components
- BBS+ signing and proofs are simplified for demonstration purposes.
- DID resolution uses in-memory storage.
- Does not implement the full W3C VC/VP specifications.
- CORS is enabled for all origins (development only).

## 📈 Future Enhancements

- [ ] Integration with Hyperledger Aries BBS+
- [ ] Support for multiple DID methods
- [ ] Zero-knowledge proof optimizations
- [x] Web-based demo interface ✅
- [ ] Performance benchmarks
- [ ] Production deployment guides
- [ ] Authentication and authorization
- [ ] Persistent storage options
- [ ] Docker containerization

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
