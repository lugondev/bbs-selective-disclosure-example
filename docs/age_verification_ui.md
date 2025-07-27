# Age Verification Web UI

## Quick Start

```bash
# Start the age verification demo server
./start-age-verification-demo.sh

# Or manually
go run ./cmd/server --port=8089
```

Then visit: **http://localhost:8089/age-verification.html**

## Features

### 🎯 **Interactive Age Verification Flow**
1. **Setup Phase**: Government Authority, Citizen, Service Provider
2. **Enhanced ID Issuance**: Digital ID with age verification claims  
3. **Privacy-Preserving Verification**: Prove age without revealing details
4. **Access Decision**: Grant/deny based on age requirements

### 🛡️ **Privacy Protection Demonstration**
- **✅ Revealed**: Age verification (boolean), nationality, document type
- **🔒 Hidden**: Exact age, birth date, name, address, ID number

### 🎮 **Multiple Service Scenarios**
- **Gaming Platform** (18+): Access adult gaming content
- **Alcohol Store** (21+): Purchase alcoholic beverages  
- **Social Media** (13+): Create social media accounts
- **Movie Theater** (16+): Watch R-rated movies
- **Senior Services** (65+): Access senior discounts

### 🚀 **Demo Automation**
- **Quick Demo**: Run complete flow with one click
- **Step-by-Step**: Manual control over each verification step
- **Real-time Logs**: See detailed execution process
- **Privacy Metrics**: Track what information is protected

## API Endpoints

The UI connects to these REST endpoints:

```
POST /api/age-verification/credential  # Issue enhanced ID
POST /api/age-verification/verify      # Verify age (privacy-preserving)
GET  /api/age-verification/scenarios   # Get supported age scenarios
POST /api/age-verification/demo        # Run automated demo
```

## Privacy Achievements

### 🎯 **Zero-Knowledge Age Proof**
Prove age ≥ N without revealing exact age

### 🔒 **Personal Data Protection** 
Name, address, birth date remain completely private

### 🎪 **Unlinkable Verifications**
Each verification is unlinkable across services

### ✅ **Regulatory Compliance**
Meet legal age requirements while maximizing privacy

## Real-World Applications

- **E-commerce**: Age-restricted product purchases
- **Entertainment**: Gaming, streaming, movie theaters  
- **Financial Services**: Age verification for banking
- **Healthcare**: Age-appropriate medical services
- **Education**: Age-based course access

## Technical Implementation

- **BBS+ Signatures**: Enable selective disclosure
- **Boolean Claims**: Privacy-preserving age verification
- **W3C Standards**: Verifiable Credentials/Presentations
- **DID Method**: Decentralized identity management
- **Zero-Knowledge Proofs**: Age verification without data exposure
