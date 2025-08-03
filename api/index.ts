import { Component, Info, OpenApiV3, Path } from "fluid-oas";
import { REGISTER_ENDPOINT } from "./paths/register";
import { API_DOCS_ENDPOINT, SPEC_ENDPOINT } from "./paths/docs";
import { HEALTHCHECK_ENDPOINT } from "./paths/healthcheck";
import { MEMBER_ENDPOINT } from "./paths/member";
import {
  ALIEN_CHALLENGE_ENDPOINT,
  ALIEN_FRONTEND_CHALLENGE_ENDPOINT,
  SUBMIT_ENDPOINT,
  NGROK_ENDPOINT,
} from "./paths/challenge";
import { BASE_ALIEN_SCHEMA, DETAILED_ALIEN_SCHEMA } from "./schema/index.ts";
import { ALIEN_INVASION } from "./schema/alien.ts";

let oas = OpenApiV3.addOpenApiVersion("3.1.0")
  .addInfo(
    Info.addTitle(
      "Generate Fall 2025 Software Engineering Challenge",
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

export const COMPONENT = Component.addSchemas({
  BaseAlien: BASE_ALIEN_SCHEMA,
  DetailedAlien: DETAILED_ALIEN_SCHEMA,
  AlienInvasion: ALIEN_INVASION,
});

export const COMPONENT_MAPPINGS = COMPONENT.createMappings();

oas = oas.addComponents(COMPONENT);

// Write the openapi specification.
oas.writeOASSync("../openapi.json");
