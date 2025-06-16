import {
  MediaType,
  Operation,
  Response,
  PathItem,
  Responses,
  String,
} from "fluid-oas";

export const API_DOCS_ENDPOINT = PathItem.addMethod({
  get: Operation.addSummary("API documentation.").addResponses(
    Responses({
      "200": Response.addDescription("API Documentation Page.").addContents({
        "text/html": MediaType.addSchema(String),
      }),
    }),
  ),
});
