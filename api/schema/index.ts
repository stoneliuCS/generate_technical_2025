import { Component } from "fluid-oas";
import { ALIEN } from "./alien";

export const COMPONENT = Component.addSchemas({ alien: ALIEN })
export const COMPONENT_MAPPINGS = COMPONENT.createMappings()
