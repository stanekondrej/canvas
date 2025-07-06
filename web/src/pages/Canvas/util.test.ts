import * as util from "./util.ts";
import { expect, test } from "vitest";

test("convert degrees to radians and back", () => {
  const deg = 94;
  const rad = util.degToRad(deg);
  const degAgain = util.radToDeg(rad);

  expect(degAgain).toBe(deg);
});
