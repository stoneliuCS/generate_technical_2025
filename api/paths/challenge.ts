import {
  Array,
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
      "404": Response.addDescription("ID not found.").addContents({
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
        "404": Response.addDescription("ID not found.").addContents({
          "application/json": MediaType.addSchema(ERROR),
        }),
        "500": Response.addDescription("Internal Server Error").addContents({
          "application/json": MediaType.addSchema(ERROR),
        }),
      }),
    ),
});

export const ALIEN = Object.addProperties({
  id: UUID.addDescription("UUID of the alien."),
  name: String.addDescription("Name of the alien.").addExample("BillyBobJoe"),
  type: String.addDescription("The rank of the alien.").addEnums([
    "Regular",
    "Elite",
    "Boss",
  ]),
  stats: Object.addProperties({
    atk: Integer.addMinimum(1).addMaximum(3),
    hp: Integer.addMinimum(1).addMaximum(3),
    spd: Integer.addMinimum(1).addMaximum(3),
  }).addDescription("Combat description of the alien.").addRequired(["atk", "hp", "spd"]),
}).addRequired(["id", "name", "type", "stats"]);

// BEGIN ALIEN FRONTEND CHALLENGE ENDPOINT
export const ALIEN_FRONTEND_CHALLENGE_ENDPOINT = PathItem.addMethod({
  get: Operation.addParameters([
    ID_PARAMETER,
    Parameter.schema
      .addIn("query")
      .addName("limit")
      .addDescription("Limit of the pagination.")
      .addSchema(Integer),
    Parameter.schema
      .addIn("query")
      .addName("offset")
      .addDescription("Offset of the pagination.")
      .addSchema(Integer),
  ]).addResponses(
    Responses({
      "200": Response.addDescription(
        "Successfully retrieved alien data.",
      ).addContents({
        "application/json": MediaType.addSchema(Array.addItems(ALIEN)),
      }),
      "400": Response.addDescription("Bad Request.").addContents({
        "application/json": MediaType.addSchema(ERROR),
      }),
      "500": Response.addDescription("Internal Server Error.").addContents({
        "application/json": MediaType.addSchema(ERROR),
      }),
    }),
  ),
});
