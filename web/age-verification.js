// Age Verification Demo JavaScript

// Global variables to store demo data
let governmentDID = '';
let citizenDID = '';
let issuedCredentialID = '';

// API base URL
const API_BASE = '/api';

// Utility functions
function log(message, type = 'info') {
    const timestamp = new Date().toLocaleTimeString();
    const logEntry = document.createElement('div');
    logEntry.className = `log-entry ${type}`;
    logEntry.textContent = `[${timestamp}] ${message}`;
    
    const logsContainer = document.getElementById('demo-logs');
    if (logsContainer) {
        logsContainer.appendChild(logEntry);
        logsContainer.scrollTop = logsContainer.scrollHeight;
        logsContainer.style.display = 'block';
    }
    
    console.log(`[${type.toUpperCase()}] ${message}`);
}

function showResponse(elementId, data, isError = false) {
    const element = document.getElementById(elementId);
    if (element) {
        element.style.display = 'block';
        element.className = `response-area ${isError ? 'error' : 'success'}`;
        element.textContent = typeof data === 'string' ? data : JSON.stringify(data, null, 2);
    }
}

function updateStatus(elementId, status, text) {
    const element = document.getElementById(elementId);
    if (element) {
        element.className = `status-badge ${status}`;
        element.textContent = text;
    }
}

function setFlowStep(stepNumber) {
    // Reset all steps
    for (let i = 1; i <= 4; i++) {
        const step = document.getElementById(`step-${i}`);
        if (step) {
            step.classList.remove('active');
        }
    }
    
    // Activate current step
    const activeStep = document.getElementById(`step-${stepNumber}`);
    if (activeStep) {
        activeStep.classList.add('active');
    }
}

// 1. Setup Government Authority
async function setupGovernment() {
    try {
        setFlowStep(1);
        log('Setting up Government Authority...', 'info');
        
        const response = await fetch(`${API_BASE}/issuer/setup`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                provider: 'example'
            })
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        governmentDID = data.did;
        
        document.getElementById('government-did').value = governmentDID;
        document.getElementById('enhanced-citizen-did').value = citizenDID; // Auto-fill if citizen already setup
        
        updateStatus('government-status', 'success', 'Ready');
        showResponse('government-response', data);
        log('✓ Government Authority setup complete', 'success');

    } catch (error) {
        log(`✗ Government setup failed: ${error.message}`, 'error');
        updateStatus('government-status', 'error', 'Failed');
        showResponse('government-response', error.message, true);
    }
}

// 2. Setup Citizen
async function setupCitizen() {
    try {
        setFlowStep(1);
        log('Setting up Citizen...', 'info');
        
        const response = await fetch(`${API_BASE}/holder/setup`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                provider: 'example'
            })
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        citizenDID = data.did;
        
        document.getElementById('citizen-did').value = citizenDID;
        document.getElementById('enhanced-citizen-did').value = citizenDID; // Auto-fill
        document.getElementById('verification-citizen-did').value = citizenDID; // Auto-fill
        
        updateStatus('citizen-status', 'success', 'Ready');
        showResponse('citizen-response', data);
        log('✓ Citizen setup complete', 'success');

    } catch (error) {
        log(`✗ Citizen setup failed: ${error.message}`, 'error');
        updateStatus('citizen-status', 'error', 'Failed');
        showResponse('citizen-response', error.message, true);
    }
}

// 3. Issue Enhanced Digital ID Credential
async function issueEnhancedCredential() {
    try {
        setFlowStep(2);
        log('Issuing enhanced digital ID credential...', 'info');
        
        const requestData = {
            issuerDid: governmentDID || document.getElementById('government-did').value,
            subjectDid: document.getElementById('enhanced-citizen-did').value,
            firstName: document.getElementById('enhanced-firstname').value,
            lastName: document.getElementById('enhanced-lastname').value,
            dateOfBirth: document.getElementById('enhanced-dob').value,
            nationality: document.getElementById('enhanced-nationality').value,
            address: document.getElementById('enhanced-address').value,
            idNumber: document.getElementById('enhanced-idnumber').value
        };

        const response = await fetch(`${API_BASE}/age-verification/credential`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestData)
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        issuedCredentialID = data.credential.id;
        
        // Auto-fill credential ID for verification
        document.getElementById('verification-credential-id').value = issuedCredentialID;
        
        showResponse('government-response', data);
        log(`✓ Enhanced credential issued (Age: ${data.currentAge} years)`, 'success');
        log(`  Credential ID: ${issuedCredentialID}`, 'info');
        
        // Display age verification capabilities
        const capabilities = data.ageVerification;
        let capabilitiesText = 'Age Verification Capabilities:\n';
        Object.entries(capabilities).forEach(([key, value]) => {
            capabilitiesText += `  ${key}: ${value ? '✅' : '❌'}\n`;
        });
        
        const capabilitiesElement = document.getElementById('age-capabilities');
        if (capabilitiesElement) {
            capabilitiesElement.style.display = 'block';
            capabilitiesElement.className = 'response-area success';
            capabilitiesElement.textContent = capabilitiesText;
        }

    } catch (error) {
        log(`✗ Enhanced credential issuance failed: ${error.message}`, 'error');
        showResponse('government-response', error.message, true);
    }
}

// 4. Verify Age (Privacy-Preserving)
async function verifyAge() {
    try {
        setFlowStep(3);
        log('Starting privacy-preserving age verification...', 'info');
        
        const serviceType = document.getElementById('service-type').value;
        const minAge = parseInt(document.getElementById('min-age-required').value);
        const credentialId = document.getElementById('verification-credential-id').value;
        const holderDid = document.getElementById('verification-citizen-did').value;

        log(`Service: ${serviceType}, Required age: ${minAge}+`, 'info');
        
        const requestData = {
            holderDid: holderDid,
            credentialId: credentialId,
            minAge: minAge,
            serviceType: serviceType,
            requiredClaims: ['nationality', 'documentType']
        };

        const response = await fetch(`${API_BASE}/age-verification/verify`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestData)
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        
        setFlowStep(4);
        
        // Display verification results
        let resultHTML = `
            <div style="margin-bottom: 15px;">
                <h3>${data.accessGranted ? '🎉 ACCESS GRANTED' : '❌ ACCESS DENIED'}</h3>
                <p><strong>Service:</strong> ${serviceType} (${minAge}+ required)</p>
                <p><strong>Message:</strong> ${data.message}</p>
            </div>
            
            <div style="margin-bottom: 15px;">
                <h4>🔍 What the Service Sees:</h4>
                <pre>${JSON.stringify(data.revealedClaims, null, 2)}</pre>
            </div>
            
            <div style="margin-bottom: 15px;">
                <h4>🔒 What Remains Private:</h4>
                <ul>
        `;
        
        data.hiddenAttributes.forEach(attr => {
            resultHTML += `<li>${attr}</li>`;
        });
        
        resultHTML += `
                </ul>
            </div>
            
            <div>
                <h4>🛡️ Privacy Protection: ${data.privacyProtected ? '✅ ACTIVE' : '❌ INACTIVE'}</h4>
                <p>Your exact age, birth date, and personal details remain completely hidden.</p>
            </div>
        `;
        
        const resultElement = document.getElementById('verification-result');
        if (resultElement) {
            resultElement.style.display = 'block';
            resultElement.className = `response-area ${data.accessGranted ? 'success' : 'error'}`;
            resultElement.innerHTML = resultHTML;
        }
        
        log(data.accessGranted ? '✓ Age verification successful - Access granted' : '✗ Age verification failed - Access denied', 
            data.accessGranted ? 'success' : 'error');
        log(`Privacy protected: ${data.hiddenAttributes.length} attributes remain hidden`, 'info');

    } catch (error) {
        log(`✗ Age verification failed: ${error.message}`, 'error');
        showResponse('verification-result', error.message, true);
    }
}

// 5. Show Age Capabilities
function showAgeCapabilities() {
    const capabilitiesElement = document.getElementById('age-capabilities');
    if (capabilitiesElement.style.display === 'none' || !capabilitiesElement.style.display) {
        capabilitiesElement.style.display = 'block';
        capabilitiesElement.className = 'response-area success';
        capabilitiesElement.innerHTML = `
            <h4>🎯 Your Age Verification Capabilities</h4>
            <p>Based on your enhanced digital ID, you can prove:</p>
            <ul>
                <li>✅ Age over 13 (Social Media, PG-13 content)</li>
                <li>✅ Age over 16 (Some regional content)</li>
                <li>✅ Age over 18 (Gaming, R-rated movies)</li>
                <li>✅ Age over 21 (Alcohol purchase)</li>
                <li>✅ Age over 25 (Some adult services)</li>
                <li>❌ Age over 65 (Senior discounts) - not applicable</li>
            </ul>
            <p><strong>Privacy Guarantee:</strong> Your exact age will never be revealed!</p>
        `;
    } else {
        capabilitiesElement.style.display = 'none';
    }
}

// 6. Load Age Scenarios
async function loadAgeScenarios() {
    try {
        log('Loading age verification scenarios...', 'info');
        
        const response = await fetch(`${API_BASE}/age-verification/scenarios`);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        
        let scenariosHTML = '<h3>🎯 Available Age Verification Scenarios</h3>';
        
        data.scenarios.forEach(scenario => {
            const isEligible = scenario.minAge <= 28; // Assuming 28-year-old user
            scenariosHTML += `
                <div class="age-scenario ${isEligible ? 'eligible' : 'not-eligible'}">
                    <div>
                        <h4>${scenario.service} (${scenario.minAge}+)</h4>
                        <p>${scenario.description}</p>
                        <small><strong>Privacy:</strong> ${scenario.privacyLevel}</small>
                    </div>
                    <div>
                        <span class="status-badge ${isEligible ? 'success' : 'error'}">
                            ${isEligible ? '✅ Eligible' : '❌ Not Eligible'}
                        </span>
                    </div>
                </div>
            `;
        });
        
        scenariosHTML += `
            <div style="margin-top: 20px; padding: 15px; background: #e8f5e8; border-radius: 8px;">
                <h4>🛡️ Privacy Benefits:</h4>
                <ul>
        `;
        
        data.privacy_benefits.forEach(benefit => {
            scenariosHTML += `<li>${benefit}</li>`;
        });
        
        scenariosHTML += '</ul></div>';
        
        const scenariosElement = document.getElementById('age-scenarios-list');
        if (scenariosElement) {
            scenariosElement.innerHTML = scenariosHTML;
        }
        
        log('✓ Age scenarios loaded successfully', 'success');

    } catch (error) {
        log(`✗ Failed to load age scenarios: ${error.message}`, 'error');
    }
}

// 7. Run Full Age Demo
async function runFullAgeDemo() {
    try {
        setFlowStep(1);
        log('🚀 Starting full age verification demo...', 'info');
        
        // Clear previous logs
        const logsContainer = document.getElementById('demo-logs');
        if (logsContainer) {
            logsContainer.innerHTML = '<div style="color: #63b3ed; margin-bottom: 10px;">📋 Age Verification Demo Log:</div>';
        }
        
        // Step 1: Setup Government
        log('Step 1: Setting up Government Authority...', 'info');
        await setupGovernment();
        await new Promise(resolve => setTimeout(resolve, 1000));
        
        // Step 2: Setup Citizen
        log('Step 2: Setting up Citizen...', 'info');
        await setupCitizen();
        await new Promise(resolve => setTimeout(resolve, 1000));
        
        // Step 3: Issue Enhanced Credential
        log('Step 3: Issuing enhanced digital ID...', 'info');
        await issueEnhancedCredential();
        await new Promise(resolve => setTimeout(resolve, 1500));
        
        // Step 4: Verify Age
        const serviceType = document.getElementById('demo-service-type').value;
        const ageMapping = {
            'gaming': 18,
            'alcohol': 21,
            'social': 13,
            'senior': 65
        };
        
        document.getElementById('service-type').value = serviceType;
        document.getElementById('min-age-required').value = ageMapping[serviceType];
        
        log(`Step 4: Verifying age for ${serviceType} service...`, 'info');
        await verifyAge();
        
        log('🎉 Full age verification demo completed!', 'success');
        log('🛡️ Privacy Achievement: Personal details remain protected while proving age eligibility', 'success');

    } catch (error) {
        log(`✗ Demo failed: ${error.message}`, 'error');
    }
}

// 8. Clear All Data
function clearAll() {
    // Reset global variables
    governmentDID = '';
    citizenDID = '';
    issuedCredentialID = '';
    
    // Clear all input fields
    document.getElementById('government-did').value = '';
    document.getElementById('citizen-did').value = '';
    document.getElementById('enhanced-citizen-did').value = '';
    document.getElementById('verification-credential-id').value = '';
    document.getElementById('verification-citizen-did').value = '';
    
    // Reset status badges
    updateStatus('government-status', 'pending', 'Not Setup');
    updateStatus('citizen-status', 'pending', 'Not Setup');
    
    // Hide response areas
    const responseAreas = ['government-response', 'citizen-response', 'verification-result', 'age-capabilities', 'demo-logs'];
    responseAreas.forEach(id => {
        const element = document.getElementById(id);
        if (element) {
            element.style.display = 'none';
        }
    });
    
    // Clear scenarios
    const scenariosElement = document.getElementById('age-scenarios-list');
    if (scenariosElement) {
        scenariosElement.innerHTML = '';
    }
    
    // Reset flow steps
    setFlowStep(1);
    
    log('All data cleared. Ready for new demo.', 'info');
}

// Initialize page
document.addEventListener('DOMContentLoaded', function() {
    log('Age Verification Demo initialized', 'info');
    setFlowStep(1);
    
    // Load age scenarios on page load
    loadAgeScenarios();
});
