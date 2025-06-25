import {
  PathItem,
  Operation,
  Response,
  String,
  Responses,
  MediaType,
  Object,
} from "fluid-oas";
import { ERROR } from "../schema";

export const HEALTHCHECK_ENDPOINT = PathItem.addMethod({
  get: Operation.addResponses(
    Responses({
      "200": Response.addDescription("Server is Healthy!").addContents({
        "application/json": MediaType.addSchema(
          Object.addProperties({
            message: String.addEnums(["OK"]),
          }),
        ),
      }),
      "500": Response.addDescription("Server is not Healthy!").addContents({
        "application/json": MediaType.addSchema(ERROR),
      }),
    }),
  ),
});
