import {
  Example,
  MediaType,
  Object,
  Operation,
  PathItem,
  RequestBody,
  Response,
  Responses,
  String,
} from "fluid-oas";

export const REGISTER_ENDPOINT = PathItem.addSummary(
  "Register your Northeastern email address and grab your token",
).addMethod({
  get: Operation.addRequestBody(
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
        "application/json": MediaType.addExamples({
          "Example-One": Example.addDescription(
            "Response with token.",
          ).addValue(
            JSON.stringify({
              message: "Successfully registered user!",
              token: "f81d4fae-7dec-11d0-a765-00a0c91e6bf6",
            }),
          ),
        }).addSchema(
          Object.addProperties({
            message: String,
            token: String.addFormat("uuid"),
          }).addDescription("Successfully returned response."),
        ),
      }),
    }),
  ),
});
