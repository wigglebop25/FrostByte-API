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
            const wsUrl = `${url}?token=${token}`;
            console.log("WS: Connecting to", url);
            socket = new WebSocket(wsUrl);

            socket.onopen = () => {
                update(s => ({ ...s, connected: true, error: null }));
                console.log("WS: Connected");

                if (reconnectTimeout) clearTimeout(reconnectTimeout);

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
                console.log("WS: Disconnected", e.reason);
                cleanup();
                
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