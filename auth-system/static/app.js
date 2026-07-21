// Go Auth System - Frontend Logic

document.addEventListener('DOMContentLoaded', () => {
  checkAuthStatus();
});

// Check active session on page load
async function checkAuthStatus() {
  try {
    const res = await fetch('/api/me');
    const data = await res.json();

    if (data.success && data.user) {
      renderDashboard(data.user);
    } else {
      renderAuthForms();
    }
  } catch (err) {
    console.error('Session check failed:', err);
    renderAuthForms();
  }
}

// Switch between Sign In and Sign Up tabs
function switchTab(tab) {
  const loginForm = document.getElementById('login-form');
  const signupForm = document.getElementById('signup-form');
  const tabLogin = document.getElementById('tab-login');
  const tabSignup = document.getElementById('tab-signup');

  if (tab === 'login') {
    loginForm.classList.remove('hidden');
    signupForm.classList.add('hidden');
    tabLogin.classList.add('active');
    tabSignup.classList.remove('active');
  } else {
    loginForm.classList.add('hidden');
    signupForm.classList.remove('hidden');
    tabLogin.classList.remove('active');
    tabSignup.classList.add('active');
  }
}

// Fill demo credentials
function fillDemoCredentials() {
  switchTab('login');
  document.getElementById('login-username').value = 'demo';
  document.getElementById('login-password').value = 'password123';
  showToast('Demo credentials loaded!', 'success');
}

// Toggle Password Field Visibility
function togglePasswordVisibility(inputId, btn) {
  const input = document.getElementById(inputId);
  if (input.type === 'password') {
    input.type = 'text';
    btn.innerHTML = `<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path>
      <line x1="1" y1="1" x2="23" y2="23"></line>
    </svg>`;
  } else {
    input.type = 'password';
    btn.innerHTML = `<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
      <circle cx="12" cy="12" r="3"></circle>
    </svg>`;
  }
}

// Password Strength Meter
function checkPasswordStrength(password) {
  const bar = document.getElementById('strength-bar');
  const text = document.getElementById('strength-text');

  if (!password) {
    bar.style.width = '0%';
    bar.style.backgroundColor = 'transparent';
    text.textContent = 'Password strength';
    return;
  }

  let score = 0;
  if (password.length >= 6) score++;
  if (password.length >= 10) score++;
  if (/[0-9]/.test(password)) score++;
  if (/[^a-zA-Z0-9]/.test(password)) score++;

  if (score <= 1) {
    bar.style.width = '25%';
    bar.style.backgroundColor = '#ef4444';
    text.textContent = 'Weak password';
    text.style.color = '#ef4444';
  } else if (score === 2 || score === 3) {
    bar.style.width = '65%';
    bar.style.backgroundColor = '#f59e0b';
    text.textContent = 'Moderate password';
    text.style.color = '#f59e0b';
  } else {
    bar.style.width = '100%';
    bar.style.backgroundColor = '#10b981';
    text.textContent = 'Strong password';
    text.style.color = '#10b981';
  }
}

// Login Handler
async function handleLogin(e) {
  e.preventDefault();
  const usernameOrEmail = document.getElementById('login-username').value.trim();
  const password = document.getElementById('login-password').value;

  try {
    const res = await fetch('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ usernameOrEmail, password }),
    });
    const data = await res.json();

    if (data.success) {
      showToast(data.message, 'success');
      renderDashboard(data.user);
    } else {
      showToast(data.message || 'Login failed', 'error');
    }
  } catch (err) {
    showToast('Network error, please try again.', 'error');
  }
}

// Signup Handler
async function handleSignup(e) {
  e.preventDefault();
  const fullName = document.getElementById('signup-fullname').value.trim();
  const username = document.getElementById('signup-username').value.trim();
  const email = document.getElementById('signup-email').value.trim();
  const password = document.getElementById('signup-password').value;

  try {
    const res = await fetch('/api/signup', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ fullName, username, email, password }),
    });
    const data = await res.json();

    if (data.success) {
      showToast(data.message, 'success');
      renderDashboard(data.user);
    } else {
      showToast(data.message || 'Signup failed', 'error');
    }
  } catch (err) {
    showToast('Network error, please try again.', 'error');
  }
}

// Logout Handler
async function handleLogout() {
  try {
    const res = await fetch('/api/logout', { method: 'POST' });
    const data = await res.json();
    showToast(data.message || 'Logged out', 'success');
    renderAuthForms();
  } catch (err) {
    showToast('Logout failed', 'error');
  }
}

// Render Dashboard View
function renderDashboard(user) {
  document.getElementById('auth-card').classList.add('hidden');
  const dash = document.getElementById('dashboard-card');
  dash.classList.remove('hidden');

  document.getElementById('user-avatar').textContent = (user.fullName || user.username).charAt(0).toUpperCase();
  document.getElementById('user-display-name').textContent = user.fullName || user.username;
  document.getElementById('user-id').textContent = user.id;
  document.getElementById('user-username').textContent = `@${user.username}`;
  document.getElementById('user-email').textContent = user.email;

  const date = new Date(user.createdAt);
  document.getElementById('user-created').textContent = date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  });
}

// Render Auth Forms View
function renderAuthForms() {
  document.getElementById('dashboard-card').classList.add('hidden');
  document.getElementById('auth-card').classList.remove('hidden');
}

// Toast Notifications
function showToast(message, type = 'info') {
  const container = document.getElementById('toast-container');
  const toast = document.createElement('div');
  toast.className = `toast ${type}`;

  const iconSvg = type === 'success'
    ? `<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg>`
    : `<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>`;

  toast.innerHTML = `<span class="toast-icon">${iconSvg}</span> <span>${message}</span>`;
  container.appendChild(toast);

  setTimeout(() => {
    toast.style.opacity = '0';
    setTimeout(() => toast.remove(), 300);
  }, 3500);
}
