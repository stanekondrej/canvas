import { useEffect } from "preact/hooks";
import { Sha256 } from "@aws-crypto/sha256-browser";

type Message =
  | {
      type: "update";
      data: Stroke;
    }
  | {
      type: "checkpoint";
      data: Stroke[];
    }
  | {
      type: "close";
      data: null;
    }
  | {
      type: "error";
      data: string;
    };

type Coordinate = {
  x: number;
  y: number;
};
type Stroke = Coordinate[];

let state: Stroke[] = [];

const SERVER_URL = "ws://localhost:9999";

const computeHash = async (x: string): Promise<string> => {
  let hash = new Sha256();
  hash.update(x);
  return (await hash.digest()).toString();
};

/**
 * Call this from a context where the canvas element has been initialized!
 */
const draw = (s: Stroke, ctx: CanvasRenderingContext2D) => {
  throw new Error("TBI");
};

export const Canvas = () => {
  useEffect(() => {
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

      const canvas: HTMLCanvasElement = document.querySelector("#canvas")!;
      const ctx = canvas.getContext("2d");
      switch (msg.type) {
        case "error":
          console.error("Error from server:", msg.data);
          break;
        case "close":
          // In the holy text of RFC 6455 section 5.5.1 (Close), the spec says that we
          // should echo a Close frame back. However, there's not really a good reason
          // for the server to close the WS connection, so we're going to assume that if
          // we receive a Close frame, it's already an echo. Not the best way to do it,
          // but it reduces the complexity at no real observable cost.
          s.close();
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
          // TODO: maybe reset the canvas?
          for (const stroke of state) {
            draw(stroke, ctx);
          }

          break;
        case "update":
          console.log("Received an update:", msg.data);
          draw(msg.data, ctx);

          break;
      }
    };

    s.onopen = (_e) => {
      console.log("Server connected");
    };
  }, []);

  return (
    <canvas
      id="canvas"
      style="position: absolute; left: 0; top: 0; width: 100svw; height: 100svh;"
    ></canvas>
  );
};
