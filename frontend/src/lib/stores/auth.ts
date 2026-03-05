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
        init: () => {
            if (typeof localStorage !== 'undefined') {
                const token = localStorage.getItem('token');
                const userStr = localStorage.getItem('user');
                if (token && userStr) {
                    set({ token, user: JSON.parse(userStr) });
                }
            }
        },
        login: (token: string, user: User) => {
            localStorage.setItem('token', token);
            localStorage.setItem('user', JSON.stringify(user));
            set({ token, user });
        },
        logout: () => {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            set({ user: null, token: null });
        }
    };
}

export const auth = createAuthStore();
