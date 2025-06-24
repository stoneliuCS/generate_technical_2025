import { Integer, Object } from "fluid-oas";

export const ALIEN = Object.addProperties({
  HP: Integer.addMinimum(1),
  ATK: Integer.addMinimum(1),
  SPD: Integer.addMinimum(1),
}).addRequired(["HP", "ATK", "SPD"]);

export const WEAPON = Object.addProperties({
  COST: Integer.addMinimum(1),
  ATK: Integer.addMinimum(1),
  SPD: Integer.addMinimum(1),
}).addRequired(["COST", "ATK", "SPD"]);
