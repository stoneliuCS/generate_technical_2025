import {
  Boolean,
  Integer,
  MediaType,
  Object,
  OneOf,
  Operation,
  Parameter,
  PathItem,
  RequestBody,
  Response,
  Responses,
  String,
} from "fluid-oas";
import { ALIEN_INVASION, ALIEN_INVASION_ANSWER } from "../schema/alien";
import { ERROR, UUID } from "../schema";

const ID_PARAMETER = Parameter.schema
  .addIn("path")
  .addRequired(true)
  .addName("id")
  .addSchema(UUID);

export const ALIEN_CHALLENGE_ENDPOINT = PathItem.addMethod({
  get: Operation.addParameters([ID_PARAMETER]).addResponses(
    Responses({
      "200": Response.addDescription(
        "Successfully gotten alien invasion data",
      ).addContents({
        "application/json": MediaType.addSchema(ALIEN_INVASION),
      }),
      "401": Response.addDescription(
        "Invalid ID. Are you sure you are using the id that you got upon registration?",
      ).addContents({
        "application/json": MediaType.addSchema(ERROR),
      }),
      "500": Response.addDescription("Internal Server Error.").addContents({
        "application/json": MediaType.addSchema(ERROR),
      }),
    }),
  ),
});

export const SUBMIT_RESPONSE = OneOf(
  Object.addProperties({ valid: Boolean, score: Integer }),
  Object.addProperties({ valid: Boolean, reason: String }),
);

export const SUBMIT_ENDPOINT = PathItem.addMethod({
  post: Operation.addParameters([ID_PARAMETER])
    .addRequestBody(
      RequestBody.addContents({
        "application/json": MediaType.addSchema(ALIEN_INVASION_ANSWER),
      }),
    )
    .addResponses(
      Responses({
        "200": Response.addDescription(
          "Verify submission against testing server oracle.",
        ).addContents({
          "application/json": MediaType.addSchema(SUBMIT_RESPONSE),
        }),
        "400": Response.addDescription("Malformed Submission").addContents({
          "application/json": MediaType.addSchema(ERROR),
        }),
        "500": Response.addDescription("Internal Server Error").addContents({
          "application/json": MediaType.addSchema(ERROR),
        }),
      }),
    ),
});
