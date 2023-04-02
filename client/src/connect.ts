type Options = {
  onOpen?: (event: Event) => void
  onClose?: (event: CloseEvent) => void
  onError?: (event: Event) => void
  onMessage?: (event: MessageEvent<any>) => void
};

export function connect(url: string, options: Options = {}) {
  const websocket = new WebSocket(url);
  websocket.onopen = evt => {
    if (options.onOpen) options.onOpen(evt);
  }

  websocket.onclose = evt => {
    if (options.onClose) options.onClose(evt);
    websocket.close();
  }

  websocket.onmessage = evt => {
    if (options.onMessage) options.onMessage(evt);
  }

  websocket.onerror = evt => {
    if (options.onError) options.onError(evt);
    websocket.close();
  }

  return websocket;
}
