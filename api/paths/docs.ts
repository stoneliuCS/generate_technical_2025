import {
  MediaType,
  Operation,
  Response,
  PathItem,
  Responses,
  String,
} from "fluid-oas";
import { ERROR } from "../schema";

export const API_DOCS_ENDPOINT = PathItem.addMethod({
  get: Operation.addSummary("API documentation.").addResponses(
    Responses({
      "200": Response.addDescription("API Documentation Page.").addContents({
        "text/html": MediaType.addSchema(String),
      }),
      "500": Response.addDescription(
        "API Docs could not be found.",
      ).addContents({
        "application/json": MediaType.addSchema(ERROR),
      }),
    }),
  ),
});

export const SPEC_ENDPOINT = PathItem.addMethod({
  get: Operation.addSummary("Challenge Specification").addResponses(
    Responses({
      "200": Response.addDescription(
        "Challenge Specification Page.",
      ).addContents({
        "text/html": MediaType.addSchema(String),
      }),
      "500": Response.addDescription(
        "Error cannot serve specification.",
      ).addContents({
        "application/json": MediaType.addSchema(ERROR),
      }),
    }),
  ),
});
