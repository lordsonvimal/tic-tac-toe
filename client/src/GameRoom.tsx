import { connect } from "./connect";
import { createSignal, Match, Switch, onCleanup } from "solid-js";

const STATUS = {
  connecting: "CONNECTING",
  connected: "CONNECTED",
  paired: "PAIRED",
  inGame: "IN_GAME",
  ended: "ENDED"
};

export function GameRoom() {
  const [status, setStatus] = createSignal<typeof STATUS[keyof typeof STATUS]>(STATUS.connecting);

  const websocket = connect("ws://localhost:8080/ws");

  onCleanup(() => {
    websocket.close();
  });

  const closeWebsocket = () => {
    websocket.close();
  };

  const sendMessage = () => {
    websocket.send(JSON.stringify({ status: "Hello" }));
  };

  return (
    <>
      <button onClick={sendMessage}>SEND</button>
      <button onClick={closeWebsocket}>CLOSE</button>
      <Switch>
        <Match when={status() === STATUS.connecting}>
          <div>Connecting to game room</div>
        </Match>
        <Match when={status() === STATUS.connected}>
          <div>Connected</div>
        </Match>
        <Match when={status() === STATUS.paired}>
          <div>Paired</div>
        </Match>
        <Match when={status() === STATUS.ended}>
          <div>Ended</div>
        </Match>
      </Switch>
    </>
  )
}
