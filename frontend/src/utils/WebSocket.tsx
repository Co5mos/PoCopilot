import { useState, useEffect } from "react";

function useWebSocket(onMessage: Function) {
    const [ws, setWs] = useState<WebSocket | null>(null);
    const [shouldConnect, setShouldConnect] = useState(false);

    const initiateConnection = () => {
        setShouldConnect(true);
    };

    useEffect(() => {
        if (!shouldConnect) return;

        console.log("Connecting to WebSocket");
        const websocket = new WebSocket("ws://localhost:5555/ws");

        websocket.onopen = () => {
            console.log("WebSocket connection opened");
        };

        websocket.onmessage = (event) => {
            onMessage(event);
        };

        websocket.onerror = (error) => {
            console.error("WebSocket Error:", error);
        };

        websocket.onclose = (event) => {
            if (event.wasClean) {
                console.log(`Closed cleanly, code=${event.code}, reason=${event.reason}`);
            } else {
                console.error("Connection died");
            }
        };

        // Heartbeat
        const heartbeatInterval = setInterval(() => {
            if (websocket.readyState === WebSocket.OPEN) {
                console.log("Sending heartbeat");
                websocket.send("heartbeat");
            }
        }, 30000);

        setWs(websocket);

        return () => {
            clearInterval(heartbeatInterval);
            websocket.close();
        };
    }, [shouldConnect]);

    return { ws, initiateConnection };
}

export default useWebSocket;
