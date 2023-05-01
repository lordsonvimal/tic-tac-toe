import { render } from "solid-js/web";
import { GameRoom } from "./GameRoom";
import "./index.scss";
import { Connection } from "./connection/connection";

render(
  () => (
    <Connection>
      <h1>Tic Tac Toe</h1>
      <GameRoom />
    </Connection>
  ),
  document.getElementById("root") as HTMLElement
);
