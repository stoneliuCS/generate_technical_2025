import {
  Array,
  Integer,
  MediaType,
  Object,
  Operation,
  PathItem,
  Response,
  Responses,
} from "fluid-oas";
import { ALIEN } from "../schema/alien";
import { COMPONENT_MAPPINGS } from "../schema";

export const ALIEN_ENDPOINT = PathItem.addMethod({
  get: Operation.addResponses(
    Responses({
      "200": Response.addDescription(
        "Successfully gotten alien invasion data",
      ).addContents({
        "application/json": MediaType.addSchema(
          Object.addProperties({
            waves: Array.addItems(COMPONENT_MAPPINGS.get(ALIEN!)),
            budget: Integer,
            health: Integer,
          }).addRequired(["waves", "budget", "health"]),
        ),
      }),
    }),
  ),
});
