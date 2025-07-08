import { Sha256 } from "@aws-crypto/sha256-browser";

/** The message that is received from the server (and update messages can also be sent)
 * to the server, of course)
 */
// TODO: make this use inheritance or something to make the code cleaner
export type Message =
  | {
      type: "update";
      data: Stroke;
    }
  | {
      type: "checkpoint";
      data: Stroke[];
    }
  | {
      type: "error";
      data: string;
    };

/** This message can be received from the main "thread", e.g. the UI.
 */
export type ControlMessage =
  | {
      type: "close";
      data: null;
    }
  | {
      type: "stroke";
      data: Stroke;
    };

export type Coordinate = {
  x: number;
  y: number;
};
export type Stroke = Coordinate[];

let state: Stroke[] = [];

const SERVER_URL = "ws://localhost:9999";

const computeHash = async (x: string): Promise<string> => {
  let hash = new Sha256();
  hash.update(x);
  return (await hash.digest()).toString();
};

const main = () => {
  const s = new WebSocket(SERVER_URL);

  onmessage = (e: MessageEvent<ControlMessage>) => {
    switch (e.data.type) {
      case "close":
        console.log("received close request from UI");

        s.close();
        break;
      case "stroke":
        console.log("received stroke message from UI");

        s.send(
          JSON.stringify({
            type: "update",
            data: e.data.data,
          } satisfies Message),
        );
    }
  };

  s.onclose = (e) => {
    console.error(`The server has closed the connection: ${e.reason}`);
  };

  s.onmessage = async (e: MessageEvent<string>) => {
    let msg: Message;
    try {
      msg = JSON.parse(e.data);
    } catch (e) {
      console.error("Invalid JSON:", e);
      return;
    }

    switch (msg.type) {
      case "error":
        console.error("Error from server:", msg.data);
        break;
      case "checkpoint":
        console.log("Received checkpoint");
        // TODO: add length check; this might help speed this up
        // on a similar note - this might not be necessary at all. i'm _pretty_ sure that
        // websockets integrity-check the frames anyway (because its TCP) and therefore
        // we don't need to exactly check if each number matches. instead a cheap length
        // check may be more suitable

        const incomingHash: string = await computeHash(msg.data.toString());
        const presentHash: string = await computeHash(state.toString());

        if (incomingHash === presentHash) {
          break;
        }
        console.warn(
          "Hashes of local and remote data don't match, restoring from checkpoint",
        );
        state = msg.data;
        postMessage({
          type: "checkpoint",
          data: state,
        } satisfies Message);

        break;
      case "update":
        console.log("Received an update:", msg.data);
        postMessage({ type: "update", data: msg.data } satisfies Message);

        break;
    }
  };

  s.onopen = (_e) => {
    console.log("Server connected");
  };
};

main();
