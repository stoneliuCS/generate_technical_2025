import {
  String,
  Object,
} from "fluid-oas";


export const NGROK_URL_SUBMISSION = Object.addProperties({
  url: String
    .addFormat("uri")
    .addPattern(/^https:\/\/[a-z0-9-]+\.ngrok\.io$/),
});