import {
  String,
  Object,
} from "fluid-oas";


export const NGROK_URL_SUBMISSION = Object.addProperties({
  url: String
    .addFormat("uri")
});