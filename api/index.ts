import { Info, OpenApiV3, Path } from "fluid-oas";
import { REGISTER_ENDPOINT } from "./paths/register";
import { API_DOCS_ENDPOINT, SPEC_ENDPOINT } from "./paths/docs";
import { HEALTHCHECK_ENDPOINT } from "./paths/healthcheck";
import { COMPONENT } from "./schema";
import { MEMBER_ENDPOINT } from "./paths/member";
import {
  ALIEN_CHALLENGE_ENDPOINT,
  ALIEN_FRONTEND_CHALLENGE_ENDPOINT,
  SUBMIT_ENDPOINT,
  NGROK_ENDPOINT,
} from "./paths/challenge";

let oas = OpenApiV3.addOpenApiVersion("3.1.0")
  .addInfo(
    Info.addTitle(
      "Generate Backend Software Engineering Challenge 2025",
    ).addVersion("1.0.0"),
  )
  .addPaths(
    Path.addEndpoints({
      "/": API_DOCS_ENDPOINT,
      "/healthcheck": HEALTHCHECK_ENDPOINT,
      "/challenge": SPEC_ENDPOINT,
      "/api/v1/member/register": REGISTER_ENDPOINT,
      "/api/v1/member": MEMBER_ENDPOINT,
      "/api/v1/challenge/backend/{id}/aliens": ALIEN_CHALLENGE_ENDPOINT,
      "/api/v1/challenge/backend/{id}/aliens/submit": SUBMIT_ENDPOINT,
      "/api/v1/challenge/frontend/{id}/aliens":
        ALIEN_FRONTEND_CHALLENGE_ENDPOINT,
      "/api/v1/challenge/backend/{id}/ngrok/submit": NGROK_ENDPOINT,
    }),
  );

oas = oas.addComponents(COMPONENT);

// Write the openapi specification.
oas.writeOASSync("../openapi.json");
