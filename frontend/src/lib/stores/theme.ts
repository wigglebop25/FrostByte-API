/**
 * Theme Store — Application-Wide Dark/Light Mode Management
 *
 * Provides a reactive Svelte store that manages the user's theme preference
 * across sessions via localStorage. Supports three modes: light, dark, and
 * system (auto-detect from OS preference). Applies the resolved theme by
 * toggling the `.dark` class on `<html>` for Tailwind CSS compatibility.
 *
 * @module stores/theme
 */

import { writable } from 'svelte/store';
import { browser } from '$app/environment';

export type Theme = 'light' | 'dark' | 'system';

/** Detects the operating system's preferred color scheme. */
function getSystemTheme(): 'light' | 'dark' {
    if (browser && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        return 'dark';
    }
    return 'light';
}

/**
 * Creates the singleton theme store with localStorage persistence
 * and real-time OS preference detection.
 */
function createThemeStore() {
    const stored = browser ? (localStorage.getItem('theme') as Theme) : null;
    const initial: Theme = stored || 'system';

    const { subscribe, set, update } = writable<Theme>(initial);

    /** Applies the resolved theme to the document root element. */
    function applyTheme(theme: Theme) {
        if (!browser) return;
        const resolved = theme === 'system' ? getSystemTheme() : theme;
        document.documentElement.classList.toggle('dark', resolved === 'dark');
        document.documentElement.setAttribute('data-theme', resolved);
    }

    /** Persists the selected theme and applies it immediately. */
    function setTheme(theme: Theme) {
        if (browser) {
            localStorage.setItem('theme', theme);
        }
        set(theme);
        applyTheme(theme);
    }

    /**
     * Initializes the theme on application mount.
     * Restores the persisted preference and registers a listener for
     * OS-level color scheme changes when the user has selected 'system'.
     */
    function init() {
        if (!browser) return;
        const saved = (localStorage.getItem('theme') as Theme) || 'system';
        set(saved);
        applyTheme(saved);

        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
            const current = localStorage.getItem('theme') as Theme;
            if (current === 'system' || !current) {
                applyTheme('system');
            }
        });
    }

    return {
        subscribe,
        set: setTheme,
        init
    };
}

export const theme = createThemeStore();
