/**
 * API Client — Axios-Based HTTP Layer
 *
 * Provides a pre-configured Axios instance with automatic API versioning,
 * JWT authentication headers, and transparent access token refresh on 401 responses.
 * All outbound requests target the versioned backend endpoint (/api/v1).
 *
 * @module utils/api
 */

import axios from "axios";
import { env } from "$env/dynamic/public";

const api = axios.create({
    baseURL: env.PUBLIC_API_URL || "/api/v1",
    headers: {
        "Content-Type": "application/json"
    }
});

/**
 * API Versioning Middleware
 * Ensures all requests are directed to the correct versioned endpoint.
 * Dynamically appends the /v1 prefix if absent from the environment configuration.
 */
if (api.defaults.baseURL) {
    let url = api.defaults.baseURL;
    
    if (!url.includes('/v1')) {
        if (url.endsWith('/api')) {
            url += '/v1';
        } else if (url.endsWith('/api/')) {
            url += 'v1';
        } else {
            if (url.startsWith('http')) {
                 if (!url.endsWith('/')) url += '/';
                 url += 'api/v1';
            } else {
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

if (typeof window !== "undefined") {
    /** Request interceptor: attaches the current JWT access token to all outbound requests. */
    api.interceptors.request.use((config) => {
        const token = localStorage.getItem("token");
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    });

    let isRefreshing = false;
    let failedQueue: { resolve: (token: string) => void; reject: (err: any) => void }[] = [];

    /** Processes queued requests after a successful token refresh. */
    function processQueue(error: any, token: string | null = null) {
        failedQueue.forEach(prom => {
            if (error) {
                prom.reject(error);
            } else {
                prom.resolve(token!);
            }
        });
        failedQueue = [];
    }

    /**
     * Response interceptor: intercepts 401 responses and attempts a transparent
     * token refresh using the stored refresh token. Queues concurrent requests
     * during the refresh cycle to prevent duplicate refresh calls.
     */
    api.interceptors.response.use(
        (response) => response,
        async (error) => {
            const originalRequest = error.config;

            if (error.response?.status === 401 && !originalRequest._retry && !originalRequest.url?.includes('/auth/')) {
                if (isRefreshing) {
                    return new Promise((resolve, reject) => {
                        failedQueue.push({ resolve, reject });
                    }).then((token) => {
                        originalRequest.headers.Authorization = `Bearer ${token}`;
                        return api(originalRequest);
                    });
                }

                originalRequest._retry = true;
                isRefreshing = true;

                const refreshToken = localStorage.getItem('refresh_token');
                if (!refreshToken) {
                    isRefreshing = false;
                    return Promise.reject(error);
                }

                try {
                    const res = await axios.post(
                        `${api.defaults.baseURL}/auth/refresh`,
                        { refresh_token: refreshToken }
                    );
                    const newToken = res.data.token;
                    localStorage.setItem('token', newToken);

                    // Notify auth store reactively so WebSocket reconnects with the new token.
                    window.dispatchEvent(new CustomEvent('token-refreshed', { detail: newToken }));

                    processQueue(null, newToken);
                    originalRequest.headers.Authorization = `Bearer ${newToken}`;
                    return api(originalRequest);
                } catch (refreshError) {
                    processQueue(refreshError, null);
                    localStorage.removeItem('token');
                    localStorage.removeItem('refresh_token');
                    localStorage.removeItem('user');
                    window.location.href = '/login';
                    return Promise.reject(refreshError);
                } finally {
                    isRefreshing = false;
                }
            }

            return Promise.reject(error);
        }
    );
}

export default api;
