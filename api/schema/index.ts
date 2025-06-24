import { Component, Object, String } from "fluid-oas";

// Reusable models used throughout the api specification
export const COMPONENT = Component.addSchemas({});
export const COMPONENT_MAPPINGS = COMPONENT.createMappings();
export const UUID = String.addFormat("uuid")
  .addDescription("Unique identifer for the registered participant.")
  .addExample("17aa5a93-73fc-4f8c-9977-2994481213be");
export const ERROR = Object.addProperties({
  message: String,
}).addRequired(["message"]);
