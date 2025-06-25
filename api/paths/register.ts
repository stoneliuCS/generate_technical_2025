import {
  MediaType,
  Operation,
  PathItem,
  RequestBody,
  Response,
  Responses,
} from "fluid-oas";
import { ERROR, ID_RESPONSE, MEMBER_DETAILS } from "../schema";

export const REGISTER_ENDPOINT = PathItem.addSummary(
  "Register your Northeastern email address and grab your token",
).addMethod({
  post: Operation.addRequestBody(
    RequestBody.addContents({
      "application/json": MediaType.addSchema(MEMBER_DETAILS),
    }),
  ).addResponses(
    Responses({
      "201": Response.addDescription(
        "Successfully registered your Northeastern email",
      ).addContents({
        "application/json": MediaType.addSchema(ID_RESPONSE),
      }),
      "400": Response.addDescription(
        "Invalid northeastern email address or nuid provided.",
      ).addContents({
        "application/json": MediaType.addSchema(ERROR),
      }),
      "409": Response.addDescription(
        "User has already been registered.",
      ).addContents({
        "application/json": MediaType.addSchema(ERROR),
      }),
      "500": Response.addDescription("Internal server error.").addContents({
        "application/json": MediaType.addSchema(ERROR),
      }),
    }),
  ),
});
