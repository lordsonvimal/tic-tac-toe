import { Accessor, createContext, createSignal, onCleanup, useContext } from "solid-js";
import { connect } from "../services/websocket";

export type ConnectionValues = {
  connectionStatus: Accessor<string>,
  sendMessage: <T>(data: T) => void,
  subscribe: (name: string, event: <T>(data: T) => void) => void,
  unsubscribe: (name: string) => boolean
};

const ConnectionContext = createContext<ConnectionValues>();

export const CONNECTION_STATUS = {
  connecting: "CONNECTION_CONNECTING",
  connected: "CONNECTION_CONNECTED", // You are connected
  connectedToRoom: "CONNECTION_CONNECTED_TO_ROOM", // Broadcast message when someone connected
  disconnected: "CONNECTION_DISCONNECTED",
  failed: "CONNECTION_FAILED"
};

type Props = {
  children: any
};

export function Connection(props: Props) {
  const [connectionStatus, setConnectionStatus] = createSignal<typeof CONNECTION_STATUS[keyof typeof CONNECTION_STATUS]>(CONNECTION_STATUS.connecting);
  const [connectionId, setConnectionId] = createSignal("");

  const subscribers = new Map<string, <T>(data: T) => void>();

  const onOpen = () => {
    setConnectionStatus(CONNECTION_STATUS.connected);
  };

  const onClose = () => {
    setConnectionStatus(CONNECTION_STATUS.disconnected);
  };

  const onError = () => {
    setConnectionStatus(CONNECTION_STATUS.failed);
  };

  const onMessage = (event: MessageEvent<string>) => {
    try {
      const data = JSON.parse(event.data);

      switch (data.Status) {
        case CONNECTION_STATUS.connected: {
          if (!connectionId()) {
            setConnectionStatus(CONNECTION_STATUS.connected);
            setConnectionId(data.Connection);
          }
          break;
        }
        case CONNECTION_STATUS.disconnected: {
          setConnectionStatus(CONNECTION_STATUS.disconnected);
          break;
        }
        default:
          break;
      }

      for (let notify of subscribers.values()) {
        notify(data);
      }

    } catch(e) {
      console.log("Unable to parse: ", event.data);
    }
  };

  // TODO: Move this hardcoded url to .env
  const websocket = connect("ws://localhost:3000/ws", { onOpen, onClose, onError, onMessage });

  onCleanup(() => {
    websocket.close();
    subscribers.clear();
  });

  const sendMessage = <T, >(data: T) => {
    try {
      const stringData = JSON.stringify(data);
      websocket.send(stringData);
    } catch(e) {
      console.error("Unable to stringify data", data);
    }
  };

  const subscribe = (name: string, event: <T>(data: T) => void) => {
    subscribers.set(name, event);
  };

  const unsubscribe = (name: string) => {
    const isSuccess = subscribers.delete(name);
    return isSuccess;
  };

  // Whenever a subscribe is used, make sure unsubscribe is used on cleanup / unmount.
  const connection = { connectionStatus, sendMessage, subscribe, unsubscribe };

  return <ConnectionContext.Provider value={connection}>{props.children}</ConnectionContext.Provider>
}

export function useConnection() { return useContext(ConnectionContext); }
