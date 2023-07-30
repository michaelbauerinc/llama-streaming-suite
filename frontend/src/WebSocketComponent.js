import React, { useEffect } from 'react';
import WebSocket from 'ws';

const WebSocketComponent = () => {
    useEffect(() => {
        const options = {
            headers: {
                'Connection': 'upgrade',
                'Upgrade': 'websocket',
                // 'Origin': 'http://0.0.0.0:8080'
            },
        };

        const ws = new WebSocket('ws://0.0.0.0:8080/', options);

        ws.on('open', () => {
            console.log('Connected to WebSocket server');
        });

        wws.on('message', (data) => {
            console.log('Received data:', data.toString('utf8'));
            // renderHTML(data.toString('utf8'));
        });

        ws.on('close', () => {
            console.log('Disconnected from WebSocket server');
        });

        ws.on('error', (error) => {
            console.error('WebSocket error:', error);
        });

        return () => {
            // Clean up WebSocket connection when the component is unmounted
            ws.close();
        };
    }, []);

    function renderHTML(htmlData) {
        // Assuming there's a div element with the id 'app' in your HTML file
        const appElement = document.getElementById('app');
        appElement.innerHTML = htmlData;
    }

    return (
        <div>
            {/* This component doesn't render anything directly */}
        </div>
    );
};

export default WebSocketComponent;
