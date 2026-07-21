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
    btn.textContent = '🙈';
  } else {
    input.type = 'password';
    btn.textContent = '👁️';
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
  toast.innerHTML = `<span>${type === 'success' ? '✅' : '⚠️'}</span> ${message}`;
  container.appendChild(toast);

  setTimeout(() => {
    toast.style.opacity = '0';
    setTimeout(() => toast.remove(), 300);
  }, 3500);
}
