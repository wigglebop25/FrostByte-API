import axios from "axios";
import { env } from "$env/dynamic/public";

const api = axios.create({
    baseURL: env.PUBLIC_API_URL || "/api/v1",
    headers: {
        "Content-Type": "application/json"
    }
});

/**
 * API Versioning Middleware: Ensures all requests are directed to the correct versioned endpoint.
 * This logic dynamically appends the /v1 prefix if it's missing from the environment configuration,
 * providing compatibility with the versioned backend router.
 */
if (api.defaults.baseURL) {
    let url = api.defaults.baseURL;
    
    // If it's just the domain or /api, ensure it has /v1
    if (!url.includes('/v1')) {
        if (url.endsWith('/api')) {
            url += '/v1';
        } else if (url.endsWith('/api/')) {
            url += 'v1';
        } else {
            // If it doesn't end with /api, we assume it needs the full prefix if it's relative
            // or just append /api/v1 if it's an absolute domain
            if (url.startsWith('http')) {
                 if (!url.endsWith('/')) url += '/';
                 url += 'api/v1';
            } else {
                 // Relative path
                 if (url === '/' || url === '') {
                     url = '/api/v1';
                 } else {
                     if (!url.endsWith('/')) url += '/';
                     url += 'v1';
                 }
            }
        }
    }
    
    api.defaults.baseURL = url;
}

// Interceptor to add JWT token to requests
if (typeof window !== "undefined") {
api.interceptors.request.use((config) => {
const token = localStorage.getItem("token");
if (token) {
config.headers.Authorization = `Bearer ${token}`;
}
return config;
});
}

export default api;
