import dotenv from "dotenv";

dotenv.config();

const DISCORD_TOKEN = process.env.DISCORD_TOKEN;
const CLIENT_ID = process.env.CLIENT_ID;
const WEBHOOK_TOKEN = process.env.WEBHOOK_TOKEN;
const WEBHOOK_ID = process.env.WEBHOOK_ID;

const COSMOS_ENDPOINT = process.env.COSMOS_ENDPOINT;
const COSMOS_KEY = process.env.COSMOS_KEY;
const COSMOS_DATABASE_ID = process.env.COSMOS_DATABASE_ID;
const COSMOS_CONTAINER_ID = process.env.COSMOS_CONTAINER_ID;

export {
  DISCORD_TOKEN,
  CLIENT_ID,
  WEBHOOK_TOKEN,
  WEBHOOK_ID,
  COSMOS_ENDPOINT,
  COSMOS_KEY,
  COSMOS_DATABASE_ID,
  COSMOS_CONTAINER_ID,
};
