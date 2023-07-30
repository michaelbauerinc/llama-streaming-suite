import React, { useEffect, useState } from 'react';

function App() {
  const [receivedMessages, setReceivedMessages] = useState('');

  useEffect(() => {
    // ComponentDidMount equivalent
    const options = {
      headers: {
        'Connection': 'upgrade',
        'Upgrade': 'websocket',
        'Origin': 'http://0.0.0.0:8080'
      },
    };

    const socket = new WebSocket('ws://0.0.0.0:8080/');

    socket.onopen = () => {
      console.log('WebSocket connection opened.');
    };

    socket.onmessage = (event) => {
      const message = JSON.parse(event.data).message;
      setReceivedMessages((prevMessages) => prevMessages + message);
    };

    socket.onclose = () => {
      console.log('WebSocket connection closed.');
    };

    return () => {
      // ComponentWillUnmount equivalent
      socket.close();
    };
  }, []);

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
      <div style={{ textAlign: 'left', maxWidth: '600px', margin: '0 auto', padding: '20px', border: '1px solid #ccc' }}>
        <h1>WebSocket Example</h1>
        <div style={{ whiteSpace: 'pre-wrap' }}>{receivedMessages}</div>
      </div>
    </div>
  );
}

export default App;
