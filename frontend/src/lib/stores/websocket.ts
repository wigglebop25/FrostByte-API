/**
 * WebSocket Store — Real-Time Communication Layer
 *
 * Manages a persistent WebSocket connection to the backend with automatic
 * reconnection, keep-alive ping frames, and reactive connection state.
 * Designed for the FrostByte real-time order monitoring system.
 *
 * @module stores/websocket
 */

import { writable } from "svelte/store";

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

    /** Clears all active timers to prevent orphaned intervals. */
    function cleanup() {
        if (pingInterval) clearInterval(pingInterval);
        if (reconnectTimeout) clearTimeout(reconnectTimeout);
        pingInterval = null;
        reconnectTimeout = null;
    }

    /**
     * Establishes a WebSocket connection to the backend.
     * Safely handles duplicate calls, stale sockets, and automatic reconnection.
     */
    function connect(url: string, token: string) {
        if (!token) return;

        urlRef = url;
        tokenRef = token;

        // Allow connection only if no socket exists or the current one is closed/closing.
        if (socket) {
            if (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING) return;
            // Clean up stale socket in CLOSING or CLOSED state before reconnecting.
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

                // Use the latest token from localStorage in case it was refreshed.
                reconnectTimeout = setTimeout(() => {
                    const freshToken = typeof localStorage !== 'undefined' 
                        ? localStorage.getItem('token') || tokenRef 
                        : tokenRef;
                    tokenRef = freshToken;
                    connect(urlRef, freshToken);
                }, 5000);
            };

            socket.onerror = (e) => {
                console.error("WS: Error", e);
                update(s => ({ ...s, error: e }));
                socket?.close();
            };

        } catch (e) {
            console.error("WS: Init Error", e);
            socket = null;
            reconnectTimeout = setTimeout(() => {
                const freshToken = typeof localStorage !== 'undefined'
                    ? localStorage.getItem('token') || tokenRef
                    : tokenRef;
                tokenRef = freshToken;
                connect(urlRef, freshToken);
            }, 5000);
        }
    }

    return {
        subscribe,
        connect,
        sendMessage: (msg: any) => {
            if (socket?.readyState === WebSocket.OPEN) {
                socket.send(JSON.stringify(msg));
            }
        },
        disconnect: () => {
            cleanup();
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