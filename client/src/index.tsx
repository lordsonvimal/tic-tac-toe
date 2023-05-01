import { render } from "solid-js/web";
import { Connection } from "./connection";
import { GameRoom } from "./GameRoom";
import "./index.scss";

render(
  () => (
    <Connection>
      <h1>Tic Tac Toe</h1>
      <GameRoom />
    </Connection>
  ),
  document.getElementById("root") as HTMLElement
);
