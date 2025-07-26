// Global state
let issuerDID = '';
let holderDID = '';
let verifierDID = '';
let issuedCredential = null;
let currentPresentation = null;

// API base URL
const API_BASE = window.location.origin;

// Utility functions
function log(message, type = 'info') {
    const logsDiv = document.getElementById('demo-logs');
    if (!logsDiv.style.display || logsDiv.style.display === 'none') {
        logsDiv.style.display = 'block';
    }
    
    const timestamp = new Date().toLocaleTimeString();
    const logEntry = document.createElement('div');
    logEntry.className = `log-entry ${type}`;
    logEntry.innerHTML = `[${timestamp}] ${message}`;
    logsDiv.appendChild(logEntry);
    logsDiv.scrollTop = logsDiv.scrollHeight;
}

function updateFlowStep(step) {
    // Reset all steps
    for (let i = 1; i <= 4; i++) {
        document.getElementById(`step-${i}`).classList.remove('active');
    }
    // Activate current step
    document.getElementById(`step-${step}`).classList.add('active');
}

function showResponse(elementId, data, isError = false) {
    const element = document.getElementById(elementId);
    element.style.display = 'block';
    element.className = `response-area ${isError ? 'error' : 'success'}`;
    element.textContent = JSON.stringify(data, null, 2);
}

function updateStatus(elementId, status, text) {
    const element = document.getElementById(elementId);
    element.className = `status-badge ${status}`;
    element.textContent = text;
}

async function apiCall(endpoint, options = {}) {
    try {
        const response = await fetch(`${API_BASE}${endpoint}`, {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers
            },
            ...options
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || `HTTP ${response.status}`);
        }
        
        return data;
    } catch (error) {
        console.error('API call failed:', error);
        throw error;
    }
}

// Issuer functions
async function setupIssuer() {
    try {
        log('Setting up issuer (Government ID Authority)...', 'info');
        updateFlowStep(1);
        
        const response = await apiCall('/api/issuer/setup', {
            method: 'POST',
            body: JSON.stringify({ method: 'example' })
        });
        
        issuerDID = response.did;
        document.getElementById('issuer-did').value = issuerDID;
        document.getElementById('trusted-issuers').value = issuerDID;
        
        updateStatus('issuer-status', 'success', 'Setup Complete');
        showResponse('issuer-response', response);
        
        log(`✅ Issuer setup complete. DID: ${issuerDID}`, 'success');
        
    } catch (error) {
        log(`❌ Issuer setup failed: ${error.message}`, 'error');
        updateStatus('issuer-status', 'error', 'Setup Failed');
        showResponse('issuer-response', { error: error.message }, true);
    }
}

async function issueCredential() {
    try {
        if (!issuerDID) {
            throw new Error('Please setup issuer first');
        }
        
        if (!holderDID) {
            throw new Error('Please setup holder first');
        }
        
        log('Issuing digital ID credential...', 'info');
        updateFlowStep(2);
        
        const claims = [
            { key: 'firstName', value: document.getElementById('claim-firstname').value },
            { key: 'lastName', value: document.getElementById('claim-lastname').value },
            { key: 'dateOfBirth', value: document.getElementById('claim-dob').value },
            { key: 'nationality', value: document.getElementById('claim-nationality').value },
            { key: 'address', value: document.getElementById('claim-address').value },
            { key: 'idNumber', value: document.getElementById('claim-idnumber').value }
        ];
        
        const response = await apiCall('/api/issuer/credentials', {
            method: 'POST',
            body: JSON.stringify({
                issuerDid: issuerDID,
                subjectDid: holderDID,
                claims: claims
            })
        });
        
        issuedCredential = response.credential;
        document.getElementById('presentation-credential-id').value = response.credentialId;
        
        showResponse('issuer-response', response);
        log(`✅ Credential issued successfully. ID: ${response.credentialId}`, 'success');
        
        // Auto-store the credential for the holder
        await storeCredentialForHolder(response.credential);
        
    } catch (error) {
        log(`❌ Credential issuance failed: ${error.message}`, 'error');
        showResponse('issuer-response', { error: error.message }, true);
    }
}

// Holder functions
async function setupHolder() {
    try {
        log('Setting up holder (Citizen)...', 'info');
        
        const response = await apiCall('/api/holder/setup', {
            method: 'POST',
            body: JSON.stringify({ method: 'example' })
        });
        
        holderDID = response.did;
        document.getElementById('holder-did').value = holderDID;
        document.getElementById('credential-subject-did').value = holderDID;
        
        updateStatus('holder-status', 'success', 'Setup Complete');
        showResponse('holder-response', response);
        
        log(`✅ Holder setup complete. DID: ${holderDID}`, 'success');
        
    } catch (error) {
        log(`❌ Holder setup failed: ${error.message}`, 'error');
        updateStatus('holder-status', 'error', 'Setup Failed');
        showResponse('holder-response', { error: error.message }, true);
    }
}

async function storeCredentialForHolder(credential) {
    try {
        log('Storing credential for holder...', 'info');
        
        const response = await apiCall('/api/holder/credentials', {
            method: 'POST',
            body: JSON.stringify({ credential: credential })
        });
        
        log(`✅ Credential stored successfully`, 'success');
        
    } catch (error) {
        log(`❌ Failed to store credential: ${error.message}`, 'error');
    }
}

async function listCredentials() {
    try {
        if (!holderDID) {
            throw new Error('Please setup holder first');
        }
        
        log('Listing stored credentials...', 'info');
        
        const response = await apiCall(`/api/holder/credentials/list?holderDid=${encodeURIComponent(holderDID)}`);
        
        showResponse('credentials-list', response);
        log(`✅ Found ${response.credentials.length} credential(s)`, 'success');
        
    } catch (error) {
        log(`❌ Failed to list credentials: ${error.message}`, 'error');
        showResponse('credentials-list', { error: error.message }, true);
    }
}

async function createPresentation() {
    try {
        if (!holderDID) {
            throw new Error('Please setup holder first');
        }
        
        const credentialId = document.getElementById('presentation-credential-id').value;
        if (!credentialId) {
            throw new Error('Please enter credential ID');
        }
        
        log('Creating selective disclosure presentation...', 'info');
        updateFlowStep(3);
        
        // Get revealed attributes based on checkboxes
        const revealedAttributes = [];
        if (document.getElementById('reveal-dob').checked) revealedAttributes.push('dateOfBirth');
        if (document.getElementById('reveal-nationality').checked) revealedAttributes.push('nationality');
        if (document.getElementById('reveal-firstname').checked) revealedAttributes.push('firstName');
        if (document.getElementById('reveal-lastname').checked) revealedAttributes.push('lastName');
        if (document.getElementById('reveal-address').checked) revealedAttributes.push('address');
        if (document.getElementById('reveal-idnumber').checked) revealedAttributes.push('idNumber');
        
        const response = await apiCall('/api/holder/presentations', {
            method: 'POST',
            body: JSON.stringify({
                holderDid: holderDID,
                credentialIds: [credentialId],
                selectiveDisclosure: [{
                    credentialId: credentialId,
                    revealedAttributes: revealedAttributes
                }]
            })
        });
        
        currentPresentation = response.presentation;
        document.getElementById('presentation-json').value = JSON.stringify(response.presentation, null, 2);
        
        showResponse('holder-response', response);
        log(`✅ Presentation created with revealed attributes: ${revealedAttributes.join(', ')}`, 'success');
        log(`📄 Hidden attributes: ${['firstName', 'lastName', 'address', 'idNumber'].filter(attr => !revealedAttributes.includes(attr)).join(', ')}`, 'warning');
        
    } catch (error) {
        log(`❌ Presentation creation failed: ${error.message}`, 'error');
        showResponse('holder-response', { error: error.message }, true);
    }
}

// Verifier functions
async function setupVerifier() {
    try {
        log('Setting up verifier (Cinema)...', 'info');
        
        const response = await apiCall('/api/verifier/setup', {
            method: 'POST',
            body: JSON.stringify({ method: 'example' })
        });
        
        verifierDID = response.did;
        document.getElementById('verifier-did').value = verifierDID;
        
        updateStatus('verifier-status', 'success', 'Setup Complete');
        showResponse('verifier-response', response);
        
        log(`✅ Verifier setup complete. DID: ${verifierDID}`, 'success');
        
    } catch (error) {
        log(`❌ Verifier setup failed: ${error.message}`, 'error');
        updateStatus('verifier-status', 'error', 'Setup Failed');
        showResponse('verifier-response', { error: error.message }, true);
    }
}

async function verifyPresentation() {
    try {
        if (!verifierDID) {
            throw new Error('Please setup verifier first');
        }
        
        const presentationJson = document.getElementById('presentation-json').value;
        if (!presentationJson) {
            throw new Error('Please enter presentation JSON');
        }
        
        log('Verifying presentation (Cinema checking age & nationality)...', 'info');
        updateFlowStep(4);
        
        let presentation;
        try {
            presentation = JSON.parse(presentationJson);
        } catch (e) {
            throw new Error('Invalid presentation JSON');
        }
        
        const requiredClaims = document.getElementById('required-claims').value.split(',').map(s => s.trim());
        const trustedIssuers = document.getElementById('trusted-issuers').value.split(',').map(s => s.trim()).filter(s => s);
        
        const response = await apiCall('/api/verifier/verify', {
            method: 'POST',
            body: JSON.stringify({
                presentation: presentation,
                requiredClaims: requiredClaims,
                trustedIssuers: trustedIssuers,
                verificationNonce: `cinema-verification-${Date.now()}`
            })
        });
        
        showResponse('verifier-response', response);
        
        if (response.valid) {
            log(`✅ Presentation verification: PASSED`, 'success');
            log(`📊 Revealed claims: ${Object.keys(response.revealedClaims).join(', ')}`, 'info');
            
            // Age verification
            if (response.revealedClaims.dateOfBirth) {
                const age = calculateAge(response.revealedClaims.dateOfBirth);
                log(`🎂 Calculated age: ${age} years`, 'info');
                if (age >= 18) {
                    log(`✅ Age verification: PASSED (18+)`, 'success');
                } else {
                    log(`❌ Age verification: FAILED (under 18)`, 'error');
                }
            }
            
            // Nationality verification
            if (response.revealedClaims.nationality) {
                log(`🌍 Nationality: ${response.revealedClaims.nationality}`, 'info');
                log(`✅ Nationality verification: PASSED`, 'success');
            }
            
            log(`🔒 Privacy Protection: Cinema CANNOT see firstName, lastName, address, idNumber`, 'warning');
            
        } else {
            log(`❌ Presentation verification: FAILED`, 'error');
            if (response.errors && response.errors.length > 0) {
                response.errors.forEach(error => log(`   Error: ${error}`, 'error'));
            }
        }
        
    } catch (error) {
        log(`❌ Presentation verification failed: ${error.message}`, 'error');
        showResponse('verifier-response', { error: error.message }, true);
    }
}

// Utility functions
function calculateAge(dateOfBirth) {
    const birthDate = new Date(dateOfBirth);
    const today = new Date();
    let age = today.getFullYear() - birthDate.getFullYear();
    const monthDiff = today.getMonth() - birthDate.getMonth();
    
    if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birthDate.getDate())) {
        age--;
    }
    
    return age;
}

// Demo automation
async function runFullDemo() {
    try {
        log('🚀 Starting full demo automation...', 'info');
        
        // Clear previous state
        clearAll();
        
        // Step 1: Setup all entities
        log('📋 Step 1: Setting up all entities...', 'info');
        await setupIssuer();
        await new Promise(resolve => setTimeout(resolve, 500));
        
        await setupHolder();
        await new Promise(resolve => setTimeout(resolve, 500));
        
        await setupVerifier();
        await new Promise(resolve => setTimeout(resolve, 500));
        
        // Step 2: Issue credential
        log('📋 Step 2: Issuing digital ID credential...', 'info');
        await issueCredential();
        await new Promise(resolve => setTimeout(resolve, 500));
        
        // Step 3: Create presentation
        log('📋 Step 3: Creating selective disclosure presentation...', 'info');
        await createPresentation();
        await new Promise(resolve => setTimeout(resolve, 500));
        
        // Step 4: Verify presentation
        log('📋 Step 4: Verifying presentation...', 'info');
        await verifyPresentation();
        
        log('🎉 Full demo completed successfully!', 'success');
        log('👁️ Notice: Cinema can verify age (18+) and nationality without seeing personal details like name, address, or ID number', 'warning');
        
    } catch (error) {
        log(`❌ Demo automation failed: ${error.message}`, 'error');
    }
}

function clearAll() {
    // Reset global state
    issuerDID = '';
    holderDID = '';
    verifierDID = '';
    issuedCredential = null;
    currentPresentation = null;
    
    // Clear form fields
    document.getElementById('issuer-did').value = '';
    document.getElementById('holder-did').value = '';
    document.getElementById('verifier-did').value = '';
    document.getElementById('credential-subject-did').value = '';
    document.getElementById('presentation-credential-id').value = '';
    document.getElementById('presentation-json').value = '';
    document.getElementById('trusted-issuers').value = '';
    
    // Reset status badges
    updateStatus('issuer-status', 'pending', 'Not Setup');
    updateStatus('holder-status', 'pending', 'Not Setup');
    updateStatus('verifier-status', 'pending', 'Not Setup');
    
    // Hide response areas
    document.getElementById('issuer-response').style.display = 'none';
    document.getElementById('holder-response').style.display = 'none';
    document.getElementById('verifier-response').style.display = 'none';
    document.getElementById('credentials-list').style.display = 'none';
    
    // Clear logs
    const logsDiv = document.getElementById('demo-logs');
    logsDiv.innerHTML = '<div style="color: #63b3ed; margin-bottom: 10px;">📋 Demo Execution Log:</div>';
    logsDiv.style.display = 'none';
    
    // Reset flow indicator
    for (let i = 1; i <= 4; i++) {
        document.getElementById(`step-${i}`).classList.remove('active');
    }
    
    log('🧹 All data cleared', 'info');
}

// Initialize page
document.addEventListener('DOMContentLoaded', function() {
    log('🔐 BBS+ Selective Disclosure Demo UI loaded', 'info');
    log('💡 This demo shows how to verify age and nationality without revealing personal details', 'info');
    log('🎬 Cinema scenario: Verify 18+ age and nationality while keeping name, address, and ID number private', 'warning');
});
