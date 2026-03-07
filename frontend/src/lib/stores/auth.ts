/**
 * Authentication Store — JWT Session Management
 *
 * Manages user authentication state, token persistence via localStorage,
 * and automatic access token refresh using the backend's refresh token endpoint.
 * Provides reactive state for RBAC-based UI rendering.
 *
 * @module stores/auth
 */

import { writable } from 'svelte/store';

interface User {
    username: string;
    roles: { name: string; permissions: string }[];
}

function createAuthStore() {
    const { subscribe, set, update } = writable<{ user: User | null; token: string | null }>({
        user: null,
        token: null
    });

    return {
        subscribe,
        /** Restores session state from localStorage on application mount. */
        init: () => {
            if (typeof localStorage !== 'undefined') {
                const token = localStorage.getItem('token');
                const userStr = localStorage.getItem('user');
                if (token && userStr) {
                    set({ token, user: JSON.parse(userStr) });
                }
            }
        },
        /** Persists authentication credentials after successful login. */
        login: (token: string, user: User, refreshToken?: string) => {
            localStorage.setItem('token', token);
            localStorage.setItem('user', JSON.stringify(user));
            if (refreshToken) {
                localStorage.setItem('refresh_token', refreshToken);
            }
            set({ token, user });
        },
        /** Updates the access token after a successful refresh without changing user data. */
        updateToken: (token: string) => {
            localStorage.setItem('token', token);
            update(s => ({ ...s, token }));
        },
        /** Clears all session data and returns user to unauthenticated state. */
        logout: () => {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            localStorage.removeItem('refresh_token');
            set({ user: null, token: null });
        }
    };
}

export const auth = createAuthStore();
