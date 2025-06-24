import { Info, OpenApiV3, Path } from "fluid-oas";
import { REGISTER_ENDPOINT } from "./paths/register";
import { API_DOCS_ENDPOINT } from "./paths/docs";
import { HEALTHCHECK_ENDPOINT } from "./paths/healthcheck";
import { ALIEN_ENDPOINT } from "./paths/alien";
import { COMPONENT } from "./schema";

let oas = OpenApiV3.addOpenApiVersion("3.1.0")
  .addInfo(
    Info.addTitle(
      "Generate Backend Software Engineering Challenge 2025",
    ).addVersion("1.0.0"),
  )
  .addPaths(
    Path.addEndpoints({
      "/api/v1/register": REGISTER_ENDPOINT,
      "/healthcheck": HEALTHCHECK_ENDPOINT,
      "/api/v1/challenge/{id}/aliens": ALIEN_ENDPOINT,
      "/": API_DOCS_ENDPOINT,
    }),
  );

oas = oas.addComponents(COMPONENT);

// Write the openapi specification.
oas.writeOASSync("../openapi.json");
