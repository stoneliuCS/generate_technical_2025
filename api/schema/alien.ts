import { Array, Integer, Object, String } from "fluid-oas";

export const WAVE = Array.addItems(
  Object.addProperties({
    aliens: Array.addItems(
      Object.addProperties({
        type: String.addEnums(["regular", "swift", "boss"]),
        count: Integer.addMinimum(0),
      }).addRequired(["type", "count"]),
    ),
  }),
).addExample([
  { type: "regular", count: 50 },
  { type: "swift", count: 10 },
  { type: "boss", count: 2 },
]);

export const ALIEN_TYPES = Object.addProperties({
  regular: Object.addProperties({
    hp: Integer.addEnums([2]),
    id: Integer.addEnums([1]),
    atk: Integer.addEnums([3]),
  }),
  swift: Object.addProperties({
    hp: Integer.addEnums([1]),
    id: Integer.addEnums([2]),
    atk: Integer.addEnums([5]),
  }),
  boss: Object.addProperties({
    hp: Integer.addEnums([10]),
    id: Integer.addEnums([3]),
    atk: Integer.addEnums([10]),
  }),
});

export const WEAPON_TYPES = Object.addProperties({
  turret: Object.addProperties({
    atk: Integer.addEnums([1]),
    cost: Integer.addEnums([10]),
  }),
  machineGun: Object.addProperties({
    atk: Integer.addEnums([3]),
    cost: Integer.addEnums([30]),
  }),
  rayGun: Object.addProperties({
    atk: Integer.addEnums([5]),
    cost: Integer.addEnums([50]),
  }),
});

// ALIEN INVASION API RESPONSE
export const ALIEN_INVASION = Object.addProperties({
  waves: WAVE,
  alienTypes: ALIEN_TYPES,
  budget: Integer.addEnums([100]),
  wallDurability: Integer.addEnums([100]),
});

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
