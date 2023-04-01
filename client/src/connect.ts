type Options = {
  onOpen?: (event: Event) => void
  onClose?: (event: Event) => void
  onError?: (event: Event) => void
  onMessage?: (event: Event) => void
};

export function connect(url: string, options: Options = {}) {
  const websocket = new WebSocket(url);
  websocket.onopen = evt => {
    console.log("Open");
    if (options.onOpen) options.onOpen(evt);
  }

  websocket.onclose = evt => {
    console.log("Close");
    if (options.onClose) options.onClose(evt);
    websocket.close();
  }

  websocket.onmessage = evt => {
    console.log("Response: ", evt.data);
    if (options.onMessage) options.onMessage(evt);
  }

  websocket.onerror = evt => {
    console.log("Error: ", evt);
    if (options.onError) options.onError(evt);
    websocket.close();
  }

  return websocket;
}
