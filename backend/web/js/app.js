// StockApp Admin - Common JavaScript Utilities

const API_BASE = '/api/v1';

// Auth utilities
function getToken() {
    return localStorage.getItem('token');
}

function getUser() {
    const user = localStorage.getItem('user');
    return user ? JSON.parse(user) : null;
}

function isLoggedIn() {
    return !!getToken();
}

function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/web/index.html';
}

// API helper with auth
async function apiCall(endpoint, options = {}) {
    const token = getToken();
    const headers = {
        'Content-Type': 'application/json',
        ...(token && { 'Authorization': `Bearer ${token}` }),
        ...options.headers
    };

    const response = await fetch(`${API_BASE}${endpoint}`, {
        ...options,
        headers
    });

    if (response.status === 401) {
        logout();
        throw new Error('Sesión expirada');
    }

    // Handle 204 No Content (common for DELETE responses)
    if (response.status === 204) {
        return true;
    }

    const data = await response.json();

    if (!response.ok) {
        throw new Error(data.message || 'Error en la solicitud');
    }

    return data;
}

// Show toast notification
function showToast(message, type = 'success') {
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.textContent = message;
    document.body.appendChild(toast);

    setTimeout(() => {
        toast.remove();
    }, 3000);
}

// Format date
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('es-AR', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
}

// Format number
function formatNumber(num) {
    return new Intl.NumberFormat('es-AR').format(num);
}

// Check auth on protected pages
function requireAuth() {
    if (!isLoggedIn()) {
        window.location.href = '/web/index.html';
        return false;
    }
    return true;
}

// Initialize common components
document.addEventListener('DOMContentLoaded', () => {
    // Add logout handler
    const logoutBtn = document.getElementById('logoutBtn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', logout);
    }
});

// Export functions for use in pages
window.apiCall = apiCall;
window.getToken = getToken;
window.getUser = getUser;
window.isLoggedIn = isLoggedIn;
window.logout = logout;
window.showToast = showToast;
window.formatDate = formatDate;
window.formatNumber = formatNumber;
window.requireAuth = requireAuth;