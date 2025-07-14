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
  gunsPurchased: Array.addItems(
    Object.addProperties({
      type: String.addEnums(["turret", "machineGun", "rayGun"]),
    }),
  ),
  totalCost: Integer,
  assignments: Array.addItems(
    Object.addProperties({
      wave: Integer.addMinimum(1),
      gunQueues: Array.addItems(Array.addItems(Integer)),
      wallDurabilityRemaining: Integer,
    }),
  ),
});
