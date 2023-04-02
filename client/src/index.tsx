import { render } from "solid-js/web";
import { GameRoom } from "./GameRoom";
import "./index.scss";

render(
  () => (
    <>
      <h1>Tic Tac Toe</h1>
        <GameRoom />
    </>
  ),
  document.getElementById("root") as HTMLElement
);
