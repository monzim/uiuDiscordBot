import dotenv from "dotenv";

dotenv.config();

const discordToken = process.env.DISCORD_TOKEN;
const clientID = process.env.CLIENT_ID;
const WEBHOOK_TOKEN = process.env.WEBHOOK_TOKEN;
const WEBHOOK_ID = process.env.WEBHOOK_ID;

export { discordToken, clientID, WEBHOOK_TOKEN, WEBHOOK_ID };
