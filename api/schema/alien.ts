import { Array, Integer, Object, String } from "fluid-oas";

// ALIEN INVASION API RESPONSE
export const ALIEN_INVASION = Array.addItems(
  Object.addProperties({
    aliens: Array.addItems(
      Object.addProperties({
        hp: Integer.addMinimum(1).addMaximum(3),
        atk: Integer.addMinimum(1).addMaximum(3),
      }).addRequired(["hp", "atk"]),
    ),
    hp: Integer.addMinimum(50).addMaximum(100),
  }).addRequired(["aliens", "hp"]),
);

export const ALIEN_INVASION_ANSWER = Object.addProperties({
  state: Object.addProperties({
    remainingHP: Integer,
    remainingAliens: Integer,
    commands: Array.addItems(
      Array.addItems(
        String.addEnums([
          "volley",
          "alienAttack",
          "focusedShot",
          "focusedVolley",
        ]),
      ),
    ),
  }).addRequired(["remainingHP", "remainingAliens", "commands"]),
}).addRequired(["state"]);
