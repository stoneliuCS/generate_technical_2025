import { Array, Integer, Object, String } from "fluid-oas";
import { UUID } from ".";

// ALIEN INVASION API RESPONSE
export const ALIEN_INVASION = Array.addItems(
  Object.addProperties({
    challengeID: UUID.addDescription("Unique identifier for the challenge"),
    aliens: Array.addItems(
      Object.addProperties({
        hp: Integer.addMinimum(1).addMaximum(3),
        atk: Integer.addMinimum(1).addMaximum(3),
      }).addRequired(["hp", "atk"]),
    ),
    hp: Integer.addMinimum(50).addMaximum(100),
  }).addRequired(["aliens", "hp", "challengeID"]),
);

export const ALIEN_INVASION_ANSWER = Array.addItems(
  Object.addProperties({
    challengeID: UUID.addDescription("Unique Identifier for the challenge."),
    state: Object.addProperties({
      remainingHP: Integer,
      remainingAliens: Integer,
      commands: Array.addItems(
        String.addEnums([
          "volley",
          "alienAttack",
          "focusedShot",
          "focusedVolley",
        ]),
      ),
    }).addRequired(["remainingHP", "remainingAliens", "commands"]),
  }).addRequired(["state"]),
);
