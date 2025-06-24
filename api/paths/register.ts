import {
  MediaType,
  Object,
  Operation,
  PathItem,
  RequestBody,
  Response,
  Responses,
  String,
} from "fluid-oas";
import { ERROR } from "../schema";

export const REGISTER_ENDPOINT = PathItem.addSummary(
  "Register your Northeastern email address and grab your token",
).addMethod({
  post: Operation.addRequestBody(
    RequestBody.addContents({
      "application/json": MediaType.addSchema(
        Object.addProperties({
          email: String.addFormat("email").addDescription(
            "Must be a valid Northeastern email address.",
          ),
          nuid: String.addFormat("nuid").addDescription("Valid nuid"),
        }).addRequired(["email", "nuid"]),
      ),
    }),
  ).addResponses(
    Responses({
      "201": Response.addDescription(
        "Successfully registered your Northeastern email",
      ).addContents({
        "application/json": MediaType.addSchema(
          Object.addProperties({
            id: String.addFormat("uuid"),
          })
            .addDescription(
              "Unique identifier associated with the registered northeastern email.",
            )
            .addRequired(["id"]),
        ),
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
