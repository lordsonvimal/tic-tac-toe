import { For } from "solid-js";
import "./tic-tac-toe.scss";

type Props = {
  isPlayerTurn: boolean,
  moves: Record<number, string>,
  onTurn: (num: number) => void
};

const cells = [0, 1, 2, 3, 4, 5, 6, 7, 8];

export function TicTacToe(props: Props) {

  const handleClick = (cell: number) => {
    if (props.moves[cell]) return;
    // Make move
    props.onTurn(cell);
  }

  return (
    <div class="tic-tac-toe">
      <For each={cells}>
        {(item, index) => {
          const canClick = () => !props.moves[index()] && props.isPlayerTurn;
          return <button classList={{ cell: true, clickable: canClick() }} onClick={() => canClick() ? handleClick(index()) : null}>{props.moves[item] || ""}</button>
        }}
      </For>
    </div>
  );
}
