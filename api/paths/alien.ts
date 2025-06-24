import {
  MediaType,
  Operation,
  Parameter,
  PathItem,
  Response,
  Responses,
} from "fluid-oas";
import { ALIEN_INVASION } from "../schema/alien";
import { UUID } from "../schema";

export const ALIEN_ENDPOINT = PathItem.addMethod({
  get: Operation.addParameters([
    Parameter.schema
      .addIn("path")
      .addRequired(true)
      .addName("id")
      .addSchema(UUID),
  ]).addResponses(
    Responses({
      "200": Response.addDescription(
        "Successfully gotten alien invasion data",
      ).addContents({
        "application/json": MediaType.addSchema(ALIEN_INVASION),
      }),
    }),
  ),
});
