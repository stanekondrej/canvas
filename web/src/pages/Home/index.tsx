import "@fontsource-variable/open-sans";

import { useLocation } from "preact-iso";
import { useEffect } from "preact/hooks";
import { Point } from "./Point.tsx";
import { Link } from "./Link.tsx";

export function Home() {
  const { route } = useLocation();

  useEffect(() => {
    addEventListener("keydown", (e: KeyboardEvent) => {
      e.preventDefault();

      if (e.key !== " ") {
        return;
      }

      route("/canvas");
    });
  });

  return (
    <div class="flex flex-col items-center gap-12">
      <h1 class="text-5xl text-gray-700 my-36">
        A global <span class="font-bold text-black">Canvas</span>
      </h1>

      <div class="flex flex-col items-center my-4 py-2 border-t-1 border-b-1 border-t-gray-300 border-b-gray-300">
        <a class="text-3xl font-black text-center" href="/canvas">
          Start drawing
          <p>
            <span class="text-gray-700 text-lg font-normal italic">
              (Or, just press space)
            </span>
          </p>
        </a>
      </div>

      <div class="grid grid-cols-3 w-full">
        <Point image="/bolt.svg" title="Real-time updates">
          Want to draw with a friend and create something together? That's easy
          with Canvas - using a lightweight communication mechanism (
          <Link
            href="https://websocket.org/guides/websocket-protocol/"
            target="_blank"
          >
            WebSockets
          </Link>
          ), your edits will be reflected globally in real-time while
          performance stays virtually unaffected.
        </Point>

        <Point image="/shield.svg" title="Reliable">
          If your connection somehow manages to drop some packets, and this
          could happen, then the application will detect it and will heal the
          state after a few seconds.
        </Point>

        <Point image="/thought.svg" title="Made by a smart person">
          Oh yeah totally. I am to be trusted. I also know everything about
          programming, and everything as a whole ackshually. (pls check this
          project out on{" "}
          <Link href="https://github.com/stanekondrej/canvas" target="_blank">
            GitHub
          </Link>
          )
        </Point>
      </div>
    </div>
  );
}
