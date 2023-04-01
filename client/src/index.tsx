import { render } from "solid-js/web";
import { GameRoom } from "./GameRoom";

render(
  () => (
    <div>
      <h1>Tic Tac Toe</h1>
        <GameRoom />
    </div>
  ),
  document.getElementById("root") as HTMLElement
);
