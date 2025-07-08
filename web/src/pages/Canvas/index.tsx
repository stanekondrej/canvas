import { useEffect } from "preact/hooks";
import type { Stroke, Message, ControlMessage } from "./worker.ts";
import ServerWorker from "./worker?worker";
import { useLocation } from "preact-iso";
import StrokeInputHandler from "./input_handler.ts";

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
  const { route } = useLocation();

  useEffect(() => {
    const w = new ServerWorker();
    const canvas = document.querySelector("#canvas")! as HTMLCanvasElement;
    const ctx = canvas.getContext("2d");

    const inputHandler = new StrokeInputHandler();
    let insideOfStroke = false;
    canvas.addEventListener("mouseup", (e: MouseEvent) => {
      if (e.button !== 0) return;
      insideOfStroke = false;

      const s = inputHandler.finish();

      w.postMessage({ type: "stroke", data: s } satisfies ControlMessage)
    })
    canvas.addEventListener("mousedown", (e: MouseEvent) => {
      if (e.button !== 0 /* e.g. the main button */) return;
      insideOfStroke = true;
      inputHandler.add({ x: e.clientX, y: e.clientY })
    })
    canvas.addEventListener("mousemove", (e: MouseEvent) => {
      if (!insideOfStroke) return;

      inputHandler.add({ x: e.clientX, y: e.clientY })
    })

    w.onmessage = (e: MessageEvent) => messageHandler(e.data, ctx);

    addEventListener("keydown", (e: KeyboardEvent) => {
      if (e.key !== "Escape") {
        return;
      }

      w.postMessage({ type: "close", data: null } satisfies ControlMessage);
      route("/");
    });
  }, []);

  return <canvas id="canvas" class="absolute left-0 top-0 w-screen h-screen" />;
};
