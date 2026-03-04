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

    function connect(url: string, token: string) {
        if (socket && (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING)) return;

        urlRef = url;
        tokenRef = token;

        try {
            // Append token to query param for authentication
            const wsUrl = `${url}?token=${token}`;
            console.log("WS: Connecting to", url);
            socket = new WebSocket(wsUrl);

            socket.onopen = () => {
                update(s => ({ ...s, connected: true, error: null }));
                console.log("WS: Connected");

                // Clear any pending reconnect
                if (reconnectTimeout) clearTimeout(reconnectTimeout);

                // Keep-alive ping every 30s
                pingInterval = setInterval(() => {
                    if (socket?.readyState === WebSocket.OPEN) {
                        socket.send(JSON.stringify({ type: 'ping' }));
                    }
                }, 30000);
            };

            socket.onmessage = (event) => {
                try {
                    const message = JSON.parse(event.data);
                    // Update the store with the latest message
                    // We use a timestamp to ensure even identical messages trigger a change
                    update(s => ({ ...s, data: { ...message, _ts: Date.now() } }));
                } catch (e) {
                    console.error("WS: Parse Error", e);
                }
            };

            socket.onclose = (e) => {
                update(s => ({ ...s, connected: false }));
                console.log("WS: Disconnected", e.reason);
                cleanup();
                
                // Auto-reconnect after 5s
                reconnectTimeout = setTimeout(() => connect(urlRef, tokenRef), 5000);
            };

            socket.onerror = (e) => {
                console.error("WS: Error", e);
                update(s => ({ ...s, error: e }));
                socket?.close();
            };

        } catch (e) {
            console.error("WS: Init Error", e);
            reconnectTimeout = setTimeout(() => connect(urlRef, tokenRef), 5000);
        }
    }

    function cleanup() {
        if (pingInterval) clearInterval(pingInterval);
        pingInterval = null;
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
            if (reconnectTimeout) clearTimeout(reconnectTimeout);
            if (socket) {
                socket.onclose = null; // Prevent reconnect loop on manual disconnect
                socket.close();
                socket = null;
            }
        }
    };
}

export const ws = createWebSocketStore();