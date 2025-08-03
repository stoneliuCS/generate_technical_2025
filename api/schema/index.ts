import { Integer, Object, String } from "fluid-oas";

// Reusable models used throughout the api specification
export const BASE_ALIEN_SCHEMA = Object.addProperties({
  atk: Integer.addMinimum(1).addMaximum(3),
  hp: Integer.addMinimum(1).addMaximum(3),
}).addRequired(["atk", "hp"]);

export const ALIEN_TYPE_SCHEMA = String.addEnums([
  "Regular",
  "Elite",
  "Boss",
]).addDescription("Type of alien species");

export const DETAILED_ALIEN_SCHEMA = Object.addProperties({
  id: String,
  base_alien: BASE_ALIEN_SCHEMA,
  first_name: String,
  last_name: String,
  type: ALIEN_TYPE_SCHEMA,
  spd: Integer,
  profile_url: String.addFormat("uri"),
}).addRequired([
  "id",
  "base_alien",
  "first_name",
  "last_name",
  "type",
  "spd",
  "profile_url",
]);

export const UUID = String.addFormat("uuid")
  .addDescription("Unique identifer for the registered participant.")
  .addExample("17aa5a93-73fc-4f8c-9977-2994481213be");

export const ERROR = Object.addProperties({
  message: String,
}).addRequired(["message"]);

export const EMAIL = String.addFormat("email")
  .addDescription("Must be a valid Northeastern email address.")
  .addExample("foobar@northeastern.edu");
export const NUID = String.addFormat("nuid")
  .addDescription("Valid nuid, must be 9 digits in length")
  .addExample("123456789");

export const MEMBER_DETAILS = Object.addProperties({
  email: EMAIL,
  nuid: NUID,
}).addRequired(["email", "nuid"]);

export const ID_RESPONSE = Object.addProperties({
  id: String.addFormat("uuid"),
})
  .addDescription(
    "Unique identifier associated with the registered northeastern email.",
  )
  .addRequired(["id"]);
