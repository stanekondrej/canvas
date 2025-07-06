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

type Coordinate = {
  x: number;
  y: number;
};
export type Stroke = Coordinate[];

let state: Stroke[] = [];

const SERVER_URL = "ws://localhost:9999";
const TARGET_ORIGIN = "*";

const computeHash = async (x: string): Promise<string> => {
  let hash = new Sha256();
  hash.update(x);
  return (await hash.digest()).toString();
};

const main = () => {
  const s = new WebSocket(SERVER_URL);

  s.onclose = (_e) => {
    console.error("The server has closed the connection.");
  };

  s.onmessage = (e: MessageEvent<string>) => {
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

        // if this introduces lag, well idk
        let incomingHash: string;
        computeHash(msg.data.toString()).then((h) => (incomingHash = h));
        let presentHash: string;
        computeHash(state.toString()).then((h) => (presentHash = h));

        if (incomingHash === presentHash) {
          break;
        }
        console.log(
          "Hashes of local and remote data don't match, restoring from checkpoint",
        );
        state = msg.data;
        postMessage(
          {
            type: "checkpoint",
            data: state,
          } satisfies Message,
          TARGET_ORIGIN,
        );

        break;
      case "update":
        console.log("Received an update:", msg.data);
        postMessage(
          { type: "update", data: msg.data } satisfies Message,
          TARGET_ORIGIN,
        );

        break;
    }
  };

  s.onopen = (_e) => {
    console.log("Server connected");
  };
};

main();
