# BBS+ Interface Implementation

This document describes the new BBS+ interface implementation that allows switching between different BBS+ providers.

## Overview

The BBS+ interface provides a unified API for different BBS+ implementations:

- **Simple Provider**: Basic implementation for testing and development
- **Production Provider**: Full BLS12-381 cryptographic implementation
- **Aries Provider**: Hyperledger Aries Framework Go integration (placeholder)

## Quick Start

### 1. Create a BBS Service

```go
package main

import (
    "fmt"
    "github.com/lugondev/bbs-selective-disclosure-example/pkg/bbs"
)

func main() {
    // Create simple service for testing
    service, err := bbs.NewSimpleBBSService()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Provider: %s\n", service.GetProvider())
    fmt.Printf("Production Ready: %t\n", service.IsProductionReady())
}
```

### 2. Switch Between Providers

```go
// Start with simple provider
simpleService, _ := bbs.NewSimpleBBSService()

// Switch to production provider
config := bbs.DefaultConfig()
productionService, err := bbs.SwitchProvider(
    simpleService, 
    bbs.ProviderProduction, 
    config,
)
```

### 3. Use with Custom Configuration

```go
config := &bbs.Config{
    Provider:         bbs.ProviderProduction,
    EnableLogging:    true,
    ConstantTimeOps:  true,
    SecureMemory:     true,
    OperationTimeout: 30 * time.Second,
}

service, err := bbs.NewBBSService(bbs.ProviderProduction, config)
```

## Providers

### Simple Provider

- **Use Case**: Testing, development, demos
- **Security**: Demo level (NOT cryptographically secure)
- **Performance**: Fast
- **Production Ready**: No

```go
service, err := bbs.NewSimpleBBSService()
```

### Production Provider

- **Use Case**: Production deployments
- **Security**: High (BLS12-381 based)
- **Performance**: Good
- **Production Ready**: Yes

```go
service, err := bbs.NewProductionBBSService()
```

### Aries Provider (Placeholder)

- **Use Case**: Enterprise, interoperability
- **Security**: High (industry standard)
- **Performance**: Good
- **Production Ready**: Yes (when implemented)

```go
ariesConfig := &bbs.AriesConfig{
    KMSType:         "local",
    StorageProvider: "mem",
    CryptoSuite:     "BLS12381G2",
}

service, err := bbs.NewAriesBBSService(ariesConfig)
```

## Core Operations

All providers implement the same interface:

### 1. Key Generation

```go
keyPair, err := service.GenerateKeyPair()
if err != nil {
    return err
}

// Validate the key pair
err = service.ValidateKeyPair(keyPair)
```

### 2. Signing

```go
messages := [][]byte{
    []byte("John Doe"),
    []byte("25"),
    []byte("Engineer"),
}

signature, err := service.Sign(keyPair.PrivateKey, messages)
```

### 3. Verification

```go
err = service.Verify(keyPair.PublicKey, signature, messages)
```

### 4. Selective Disclosure

```go
// Create proof revealing only first and third messages
revealedIndices := []int{0, 2}
nonce := []byte("unique-nonce")

proof, err := service.CreateProof(
    signature, 
    keyPair.PublicKey, 
    messages, 
    revealedIndices, 
    nonce,
)

// Verify proof
revealedMessages := [][]byte{messages[0], messages[2]}
err = service.VerifyProof(keyPair.PublicKey, proof, revealedMessages, nonce)
```

## Configuration

### Default Configuration

```go
config := bbs.DefaultConfig()
// Provider: production
// EnableLogging: true
// ConstantTimeOps: true
// SecureMemory: true
// OperationTimeout: 30s
```

### Custom Configuration

```go
config := &bbs.Config{
    Provider:         bbs.ProviderProduction,
    EnableLogging:    false,
    ConstantTimeOps:  true,
    SecureMemory:     true,
    OperationTimeout: 60 * time.Second,
}
```

### Aries Configuration

```go
config := &bbs.Config{
    Provider: bbs.ProviderAries,
    AriesConfig: &bbs.AriesConfig{
        KMSType:         "remote",
        StorageProvider: "leveldb",
        CryptoSuite:     "BLS12381G2",
        RemoteKMSURL:    "https://kms.example.com",
        AuthToken:       "your-auth-token",
    },
}
```

## Advanced Features

### Service Wrapper with Metrics

```go
config := bbs.DefaultConfig()
config.EnableLogging = true

baseService, _ := bbs.NewSimpleBBSService()
wrapper := bbs.NewServiceWrapper(baseService, config)

// Use wrapper - automatically tracks metrics
keyPair, _ := wrapper.GenerateKeyPair()

// Get performance metrics
metrics := wrapper.GetMetrics()
fmt.Printf("Total Operations: %d\n", metrics.TotalOperations)
fmt.Printf("Success Rate: %.2f%%\n", metrics.SuccessRate*100)
```

### Migration Between Providers

```go
migrationHelper := bbs.NewMigrationHelper(
    bbs.ProviderSimple, 
    bbs.ProviderProduction, 
    config,
)

// Validate migration is possible
err := migrationHelper.ValidateMigration()
if err == nil {
    // Perform migration
    err = migrationHelper.PerformMigration()
}
```

### Benchmarking Providers

```go
providers := []bbs.Provider{
    bbs.ProviderSimple,
    bbs.ProviderProduction,
}

results, err := bbs.BenchmarkProviders(providers, 5) // 5 messages
for provider, metrics := range results {
    fmt.Printf("%s: %d operations\n", provider, metrics.TotalOperations)
}
```

### Provider Comparison

```go
comparisons := bbs.CompareProviders()
for provider, comparison := range comparisons {
    fmt.Printf("%s:\n", provider)
    fmt.Printf("  Security: %s\n", comparison.SecurityLevel)
    fmt.Printf("  Production Ready: %t\n", comparison.ProductionReady)
    fmt.Printf("  Use Case: %s\n", comparison.RecommendedUse)
}
```

## Running Demos

### Build and Run Interface Demo

```bash
# Build demo
make build-interface

# Run basic demo
make interface

# Run configuration demo
make interface-config

# Run provider switching demo
make interface-switching

# Run all demos
make interface-all
```

### Go Run Commands

```bash
# Basic demo
go run ./cmd/interface_demo

# Configuration demo
go run ./cmd/interface_demo config

# Provider switching demo
go run ./cmd/interface_demo switching

# All demos
go run ./cmd/interface_demo all
```

## Hyperledger Aries Integration

To integrate Hyperledger Aries Framework Go:

### 1. Add Dependency

```bash
go get github.com/hyperledger/aries-framework-go
```

### 2. Import Packages

```go
import (
    "github.com/hyperledger/aries-framework-go/pkg/crypto/primitive/bbs12381g2pub"
    "github.com/hyperledger/aries-framework-go/pkg/kms"
    "github.com/hyperledger/aries-framework-go/pkg/kms/localkms"
)
```

### 3. Complete Implementation

Update `pkg/bbs/aries_adapter.go` to implement the actual Aries integration:

```go
func (a *AriesService) initializeAries() error {
    // Initialize KMS
    kmsStorage, err := mem.NewProvider()
    if err != nil {
        return err
    }
    
    a.kms, err = localkms.New("local-lock://primary", kmsProvider)
    if err != nil {
        return err
    }
    
    // Initialize BBS+ crypto suite
    a.crypto = bbs12381g2pub.New()
    
    return nil
}

func (a *AriesService) GenerateKeyPair() (*KeyPair, error) {
    keyID, pubKeyBytes, err := a.kms.CreateAndExportPubKeyBytes(kms.BLS12381G2Type)
    if err != nil {
        return err
    }
    
    // Export private key (in production, keep private key in KMS)
    privKeyBytes, err := a.kms.ExportPubKeyBytes(keyID)
    if err != nil {
        return err
    }
    
    return &KeyPair{
        PublicKey:  pubKeyBytes,
        PrivateKey: privKeyBytes,
    }, nil
}
```

## Testing

```bash
# Run all tests
make test

# Run only unit tests
make test-unit

# Run with coverage
make test-coverage
```

## Security Considerations

### Production Use

- Always use `ProviderProduction` or `ProviderAries` for production
- Enable `ConstantTimeOps` to prevent timing attacks
- Enable `SecureMemory` for sensitive data cleanup
- Use proper key management (Aries KMS recommended)
- Implement audit logging
- Regular security reviews

### Simple Provider Warning

⚠️ **WARNING**: The Simple Provider is NOT cryptographically secure and should NEVER be used in production. It's designed only for testing and development purposes.

## Performance

### Benchmarks

| Provider   | Key Gen | Signing | Verification | Proof Creation | Proof Verification |
|------------|---------|---------|--------------|----------------|-------------------|
| Simple     | ~1μs    | ~1μs    | ~1μs         | ~1μs           | ~1μs              |
| Production | ~300μs  | ~1-3ms  | ~2-5ms       | ~1ms           | ~10μs             |
| Aries      | TBD     | TBD     | TBD          | TBD            | TBD               |

*Note: Benchmarks are approximate and depend on hardware and message count.*

## License

This project is licensed under the MIT License - see the LICENSE file for details.
