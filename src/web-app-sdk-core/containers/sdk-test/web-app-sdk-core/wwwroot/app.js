var baseUrl = location.origin + location.pathname.substring(0, location.pathname.lastIndexOf('/') + 1);

// Handle user type selection for both forms
function handleUserTypeSelection(userTypeId, credentialFieldsId, passwordLabelId, passwordInputId, usernameId) {
    const userType = document.getElementById(userTypeId).value;
    const credentialFields = document.getElementById(credentialFieldsId);
    const passwordLabel = document.getElementById(passwordLabelId);
    const passwordInput = document.getElementById(passwordInputId);
    const usernameElement = document.getElementById(usernameId);
    const usernameLabel = document.querySelector(`label[for="${usernameId}"]`);

    if (userType === 'DefaultWindows') {
        credentialFields.style.display = 'none';
    } else if (userType === 'External') {
        credentialFields.style.display = 'block';
        passwordLabel.textContent = 'Access Token:';
        passwordInput.type = 'text';
        passwordInput.placeholder = 'Enter access token';
        if (usernameElement) usernameElement.style.display = 'none';
        if (usernameLabel) usernameLabel.style.display = 'none';
    } else if (userType === 'Windows' || userType === 'Basic') {
        credentialFields.style.display = 'block';
        passwordLabel.textContent = 'Password:';
        passwordInput.type = 'password';
        passwordInput.placeholder = 'Enter password';
        if (usernameElement) usernameElement.style.display = 'block';
        if (usernameLabel) usernameLabel.style.display = 'block';
    }
}

// Handle user type selection for first form (if exists)
const userType = document.getElementById('userType');
if (userType) {
    userType.addEventListener('change', function() {
        handleUserTypeSelection('userType', 'credentialFields', 'passwordLabel', 'password', 'username');
    });
}

// Handle user type selection for second form (if exists)
const userTypeRC = document.getElementById('userTypeRC');
if (userTypeRC) {
    userTypeRC.addEventListener('change', function() {
        handleUserTypeSelection('userTypeRC', 'credentialFieldsRC', 'passwordLabelRC', 'passwordRC', 'usernameRC');
    });
}

// Helper function to show result
function showResult(resultId, success, data) {
    const resultDiv = document.getElementById(resultId);
    resultDiv.style.display = 'block';
    resultDiv.className = `result ${success ? 'success' : 'error'}`;
    
    if (success) {
        resultDiv.innerHTML = `
            <div class="result-field"><strong>Session ID:</strong> ${data.id}</div>
            <div class="result-field"><strong>Server URI:</strong> ${data.serverUri}</div>
            <div class="result-field"><strong>Access Token:</strong> ${data.token}</div>
        `;
    } else {
        resultDiv.innerHTML = `<strong>Error:</strong> ${data.error || 'Unknown error occurred'}`;
    }
}

// Helper function to handle loading state
function setLoading(loadingId, buttonElement, isLoading) {
    const loadingDiv = document.getElementById(loadingId);
    loadingDiv.style.display = isLoading ? 'block' : 'none';
    buttonElement.disabled = isLoading;
}

// Form 1: Server Configuration
document.getElementById('serverConfigForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    const button = this.querySelector('button');
    setLoading('loading1', button, true);

    const userType = document.getElementById('userType').value;
    let username = document.getElementById('username').value;
    let password = document.getElementById('password').value;

    // For External type, use password field as access token
    if (userType === 'External') {
        username = null;
    }

    try {
        const response = await fetch(`${baseUrl}api/session/create-with-server-config`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                serverUrl: document.getElementById('serverUrl').value,
                userType: userType,
                username: username,
                password: password
            })
        });

        const data = await response.json();
        showResult('result1', response.ok, data);
    } catch (error) {
        showResult('result1', false, { error: error.message });
    } finally {
        setLoading('loading1', button, false);
    }
});

// Form 2: Runtime Configuration
document.getElementById('runtimeConfigForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    const button = this.querySelector('button');
    setLoading('loading2', button, true);

    const userType = document.getElementById('userTypeRC').value;
    let username = document.getElementById('usernameRC').value;
    let password = document.getElementById('passwordRC').value;

    // For External type, use password field as access token
    if (userType === 'External') {
        username = null;
    }

    try {
        const response = await fetch(`${baseUrl}api/session/create-with-runtime-config`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                userType: userType,
                username: username,
                password: password
            })
        });

        const data = await response.json();
        showResult('result2', response.ok, data);
    } catch (error) {
        showResult('result2', false, { error: error.message });
    } finally {
        setLoading('loading2', button, false);
    }
});

// Form 3: Runtime Configuration with Default User
document.getElementById('defaultUserForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    const button = this.querySelector('button');
    setLoading('loading3', button, true);

    try {
        const response = await fetch(`${baseUrl}api/session/create-with-runtime-config-default-user`, {
            method: 'POST'
        });

        const data = await response.json();
        showResult('result3', response.ok, data);
    } catch (error) {
        showResult('result3', false, { error: error.message });
    } finally {
        setLoading('loading3', button, false);
    }
});