import dotenv from "dotenv";

dotenv.config();

const discordToken = process.env.DISCORD_TOKEN;
const clientID = process.env.CLIENT_ID;

export { discordToken, clientID };
