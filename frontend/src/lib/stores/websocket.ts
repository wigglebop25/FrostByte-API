/**
 * WebSocket Store — Real-Time Communication Layer
 *
 * Manages a persistent WebSocket connection to the backend with automatic
 * reconnection, exponential backoff, token expiry detection, and automatic
 * token refresh. Designed for the FrostByte real-time order monitoring system.
 *
 * @module stores/websocket
 */

import { writable } from "svelte/store";
import axios from "axios";
import { env } from "$env/dynamic/public";

/** Maximum consecutive reconnect attempts before giving up. */
const MAX_RETRIES = 12;
/** Initial reconnect delay in milliseconds. */
const BASE_DELAY = 3000;
/** Maximum reconnect delay in milliseconds (60 seconds). */
const MAX_DELAY = 60000;

/**
 * Decodes a JWT payload and checks if the token is expired.
 * Returns true if the token expires within the given buffer (seconds).
 */
function isTokenExpired(token: string, bufferSecs = 30): boolean {
    try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        return (payload.exp * 1000) < (Date.now() + bufferSecs * 1000);
    } catch {
        return true;
    }
}

/**
 * Attempts to refresh the access token using the stored refresh token.
 * Returns the new access token or null if refresh fails.
 */
async function refreshAccessToken(): Promise<string | null> {
    const refreshToken = typeof localStorage !== 'undefined' ? localStorage.getItem('refresh_token') : null;
    if (!refreshToken) return null;

    try {
        let baseURL = env.PUBLIC_API_URL || "/api/v1";
        if (!baseURL.includes('/v1')) {
            baseURL = baseURL.endsWith('/') ? baseURL + 'api/v1' : baseURL + '/api/v1';
        }
        const res = await axios.post(`${baseURL}/auth/refresh`, { refresh_token: refreshToken });
        const newToken = res.data.token;
        localStorage.setItem('token', newToken);
        window.dispatchEvent(new CustomEvent('token-refreshed', { detail: newToken }));
        return newToken;
    } catch {
        return null;
    }
}

export function createWebSocketStore() {
    const { subscribe, set, update } = writable({
        connected: false,
        data: null as any,
        error: null as any
    });

    let socket: WebSocket | null = null;
    let pingInterval: any;
    let reconnectTimeout: any;
    let urlRef = "";
    let tokenRef = "";
    let retryCount = 0;

    /** Clears all active timers to prevent orphaned intervals. */
    function cleanup() {
        if (pingInterval) clearInterval(pingInterval);
        if (reconnectTimeout) clearTimeout(reconnectTimeout);
        pingInterval = null;
        reconnectTimeout = null;
    }

    /** Calculates exponential backoff delay with jitter. */
    function getBackoffDelay(): number {
        const delay = Math.min(BASE_DELAY * Math.pow(2, retryCount), MAX_DELAY);
        const jitter = delay * 0.2 * Math.random();
        return delay + jitter;
    }

    /**
     * Schedules a reconnection attempt with exponential backoff.
     * Checks token expiry and refreshes if needed before reconnecting.
     */
    function scheduleReconnect() {
        if (retryCount >= MAX_RETRIES) {
            console.warn("WS: Max retries reached, stopping reconnect.");
            update(s => ({ ...s, connected: false, error: { message: 'Max retries exceeded' } }));
            return;
        }

        const delay = getBackoffDelay();
        console.log(`WS: Reconnecting in ${Math.round(delay / 1000)}s (attempt ${retryCount + 1}/${MAX_RETRIES})`);

        reconnectTimeout = setTimeout(async () => {
            let freshToken = typeof localStorage !== 'undefined'
                ? localStorage.getItem('token') || tokenRef
                : tokenRef;

            // If the token is expired or about to expire, refresh it first.
            if (isTokenExpired(freshToken)) {
                console.log("WS: Token expired, attempting refresh...");
                const newToken = await refreshAccessToken();
                if (newToken) {
                    freshToken = newToken;
                    console.log("WS: Token refreshed successfully.");
                } else {
                    console.warn("WS: Token refresh failed, cannot reconnect.");
                    retryCount++;
                    scheduleReconnect();
                    return;
                }
            }

            tokenRef = freshToken;
            retryCount++;
            connect(urlRef, freshToken);
        }, delay);
    }

    /**
     * Establishes a WebSocket connection to the backend.
     * Safely handles duplicate calls, stale sockets, and automatic reconnection.
     */
    function connect(url: string, token: string) {
        if (!token) return;

        urlRef = url;
        tokenRef = token;

        // If the token is already expired on initial connect, try to refresh first.
        if (isTokenExpired(token)) {
            console.log("WS: Token expired before connect, scheduling refresh...");
            scheduleReconnect();
            return;
        }

        // Allow connection only if no socket exists or the current one is closed/closing.
        if (socket) {
            if (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING) return;
            socket.onclose = null;
            socket.onerror = null;
            socket = null;
        }

        cleanup();

        try {
            const wsUrl = `${url}?token=${token}`;
            console.log("WS: Connecting to", url);
            socket = new WebSocket(wsUrl);

            socket.onopen = () => {
                retryCount = 0; // Reset backoff on successful connection.
                update(s => ({ ...s, connected: true, error: null }));
                console.log("WS: Connected");

                // Send application-level keep-alive pings every 30 seconds.
                pingInterval = setInterval(() => {
                    if (socket?.readyState === WebSocket.OPEN) {
                        socket.send(JSON.stringify({ type: 'ping' }));
                    }
                }, 30000);
            };

            socket.onmessage = (event) => {
                try {
                    const message = JSON.parse(event.data);
                    update(s => ({ ...s, data: { ...message, _ts: Date.now() } }));
                } catch (e) {
                    console.error("WS: Parse Error", e);
                }
            };

            socket.onclose = (e) => {
                update(s => ({ ...s, connected: false }));
                console.log("WS: Disconnected", e.code, e.reason);
                cleanup();
                socket = null;
                scheduleReconnect();
            };

            socket.onerror = (e) => {
                console.error("WS: Error Event");
                update(s => ({ ...s, error: e }));
                socket?.close();
            };

        } catch (e) {
            console.error("WS: Init Error", e);
            socket = null;
            scheduleReconnect();
        }
    }

    return {
        subscribe,
        connect,
        /** Resets the retry counter, allowing reconnection after max retries. */
        resetRetries: () => { retryCount = 0; },
        sendMessage: (msg: any) => {
            if (socket?.readyState === WebSocket.OPEN) {
                socket.send(JSON.stringify(msg));
            }
        },
        disconnect: () => {
            cleanup();
            retryCount = MAX_RETRIES; // Prevent auto-reconnect after intentional disconnect.
            if (socket) {
                socket.onclose = null;
                socket.close();
                socket = null;
            }
            update(s => ({ ...s, connected: false }));
        }
    };
}

export const ws = createWebSocketStore();