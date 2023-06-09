import { TicTacToe } from "./tic-tac-toe";
import { createSignal, Match, Switch, Show } from "solid-js";
import { CONNECTION_STATUS, ConnectionValues, useConnection } from "./connection";

const GAME_STATUS = {
  pending: "GAME_PENDING",
  started: "GAME_STARTED",
  ended: "GAME_ENDED",
  turn: "PLAYER_TURN",
  turnChange: "PLAYER_TURN_CHANGE"
};

type TicTacToe = {
  Connection: string,
  Id: string,
  Game: {
    Data: number,
    Player: Record<string, "X" | "O">,
    Status: typeof GAME_STATUS[keyof typeof GAME_STATUS],
    Turn: string
  },
  Status: typeof CONNECTION_STATUS[keyof typeof CONNECTION_STATUS]
};

export function GameRoom() {
  const [gameStatus, setGameStatus] = createSignal(GAME_STATUS.pending);
  const [playerId, setPlayerId] = createSignal("");
  const [opponentId, setOpponentId] = createSignal("");
  const [roomId, setRoomId] = createSignal("");
  const [isPlayerTurn, setIsPlayerTurn] = createSignal(false);
  const [playerShape, setPlayerShape] = createSignal("");
  const [moves, setMoves] = createSignal<Record<number, string >>({});
  const [winner, setWinner] = createSignal("");

  const onMessage = (data: any) => {
    switch(data.Game.Status) {
      case GAME_STATUS.turn: {
        // Receive other player move from server
        updateMoves(data.Game.Data, data.Game.Turn);
        return;
      }
      case GAME_STATUS.turnChange: {
        setIsPlayerTurn(data.Game.Turn === playerShape());
        return;
      }
      case GAME_STATUS.started: {
        setPlayerShape(data.Game.Player[playerId()]);
        setIsPlayerTurn(data.Game.Turn === playerShape());
        setGameStatus(GAME_STATUS.started);
        Object.keys(data.Game.Player).forEach(id => {
          if (id !== playerId()) {
            setOpponentId(id);
          }
        });
        return;
      }
    }
    switch(data.Status) {
      case CONNECTION_STATUS.connected: {
        if (playerId()) return; // Already has a player Id
        setPlayerId(data.Connection);
        setRoomId(data.Id);
        return;
      }
      case CONNECTION_STATUS.disconnected: {
        setPlayerId("");
        setRoomId("");
        return;
      }
    }
  };

  const { connectionStatus, sendMessage, subscribe } = useConnection() as ConnectionValues;

  subscribe("GAMEROOM", onMessage);

  const checkWinner = () => {
    const setEndGame = (match: string | null) => {
      setGameStatus(GAME_STATUS.ended);
      if (!match) {
        setWinner("DRAW");
        return;
      }
      if (match === playerShape()) setWinner("YOU WIN!!!");
      else setWinner(`YOU LOSE!!!`);
    }

    const isWin = (a: any, b: any, c: any) => {
      if (a) {
        if (a === b && a === c) {
          setEndGame(a);
          return true;
        }
      }
      return false;
    }

    if (isWin(moves()[0], moves()[1], moves()[2])) return;
    if (isWin(moves()[0], moves()[3], moves()[6])) return;
    if (isWin(moves()[0], moves()[4], moves()[8])) return;
    if (isWin(moves()[1], moves()[4], moves()[7])) return;
    if (isWin(moves()[2], moves()[5], moves()[8])) return;
    if (isWin(moves()[2], moves()[4], moves()[6])) return;
    if (isWin(moves()[3], moves()[4], moves()[5])) return;
    if (isWin(moves()[6], moves()[7], moves()[8])) return;

    if (Object.keys(moves()).length === 9) setEndGame(null);

    // console.log(`ROOM ID: ${roomId()}, PLAYER ID: ${playerId()}, SHAPE: ${playerShape()}, TURN: ${isPlayerTurn()}`);
  };

  const updateMoves = (num: number, shape: string) => {
    if (moves()[num]) return;
    const newMoves = {...moves(), [num]: shape};
    setMoves(newMoves);
    checkWinner();
  };

  const getTurnData = (num: number): TicTacToe => {
    return {
      Connection: playerId(),
      Game: {
        Data: num,
        Player: {},
        Status: GAME_STATUS.turn,
        Turn: playerShape()
      },
      Id: roomId(),
      Status: connectionStatus()
    };
  }

  const onTurn = (num: number) => {
    updateMoves(num, playerShape());
    sendMessage(getTurnData(num));
    setIsPlayerTurn(false);
  };

  const canPlayerPlay = () => isPlayerTurn() && gameStatus() === GAME_STATUS.started;

  return (
    <>
      <TicTacToe isPlayerTurn={canPlayerPlay()} moves={moves()} onTurn={onTurn} />
      <Show when={winner()}>{`${winner()}`}</Show>
      <Switch>
        <Match when={connectionStatus() === CONNECTION_STATUS.connecting}>
          <div>Connecting to game room</div>
        </Match>
        <Match when={connectionStatus() === CONNECTION_STATUS.connected}>
          <div>Connected to game server</div>
        </Match>
        <Match when={connectionStatus() === CONNECTION_STATUS.disconnected}>
          <div>Disconnected from game server. Refresh to try connecting again</div>
        </Match>
        <Match when={connectionStatus() === CONNECTION_STATUS.failed}>
          <div>Connection to game failed</div>
        </Match>
      </Switch>
    </>
  )
}
