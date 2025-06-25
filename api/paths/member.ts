import {
  Integer,
  MediaType,
  Operation,
  Parameter,
  PathItem,
  Response,
  Responses,
  String,
} from "fluid-oas";
import { ERROR, ID_RESPONSE } from "../schema";

export const MEMBER_ENDPOINT = PathItem.addSummary(
  "Get the associated id from the northeastern email address and nuid incase you forgot.",
).addMethod({
  get: Operation.addParameters([
    Parameter.schema
      .addIn("query")
      .addName("email")
      .addDescription("Northeastern email address.")
      .addRequired(true)
      .addSchema(String.addFormat("email")),
    Parameter.schema
      .addIn("query")
      .addName("nuid")
      .addDescription("Northeastern NUID")
      .addRequired(true)
      .addSchema(Integer),
  ]).addResponses(
    Responses({
      "200": Response.addDescription(
        "Successfully retrieved id from email and nuid",
      ).addContents({
        "application/json": MediaType.addSchema(ID_RESPONSE),
      }),

      "400": Response.addDescription(
        "Invalid northeastern email address or nuid provided.",
      ).addContents({
        "application/json": MediaType.addSchema(ERROR),
      }),

      "500": Response.addDescription("Internal server error.").addContents({
        "application/json": MediaType.addSchema(ERROR),
      }),
    }),
  ),
});
