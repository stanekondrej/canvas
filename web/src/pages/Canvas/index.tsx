import { useEffect } from "preact/hooks";
import type { Stroke, Message } from "./worker.ts";
import ServerWorker from "./worker?worker";

/**
 * Call this from a context where the canvas element has been initialized!
 */
const draw = (s: Stroke, ctx: CanvasRenderingContext2D) => {
  if (s.length < 2) {
    console.warn(`Invalid stroke ${s}; less than two points`);
  }

  for (let i = 0; i < s.length - 1; ++i) {
    const c1 = s[i];
    const c2 = s[i + 1];
  }
};

const messageHandler = (msg: Message, ctx: CanvasRenderingContext2D) => {
  console.log(`Received message from worker (type: ${msg.type})`);

  switch (msg.type) {
    case "update":
      draw(msg.data, ctx);

      break;
    case "checkpoint":
      break;
  }
};

export const Canvas = () => {
  useEffect(() => {
    const w = new ServerWorker();
    const ctx = (
      document.querySelector("#canvas")! as HTMLCanvasElement
    ).getContext("2d");

    w.onmessage = (e: MessageEvent) => messageHandler(e.data, ctx);
  }, []);

  return (
    <canvas
      id="canvas"
      style="position: absolute; left: 0; top: 0; width: 100svw; height: 100svh;"
    />
  );
};
