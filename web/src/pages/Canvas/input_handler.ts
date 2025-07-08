import type { Coordinate, Stroke } from "./worker";

export default class StrokeInputHandler {
  #currentStroke: Stroke | null;

  constructor() {
    this.#currentStroke = null;
  }

  add(coord: Coordinate) {
    if (!this.#currentStroke) this.#currentStroke = [];

    this.#currentStroke.push(coord);
  }

  finish(): Stroke {
    const s = this.#currentStroke;
    this.#currentStroke = null;

    return s;
  }
}
