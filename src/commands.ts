import fs from "fs";
import path from "path";
import { REST, Routes } from "discord.js";
import { clientID, discordToken } from "./config";

const initializeCommands = async () => {
  let commands = [];

  let commandsFile = fs
    .readdirSync(path.join(__dirname, "commands"))
    .filter((file) => file.endsWith(".js"));

  for (let file of commandsFile) {
    let command = require(`./commands/${file}`);
    commands.push(command.data.toJSON());
  }

  const rest = new REST({ version: "9" }).setToken(discordToken as string);

  await rest
    .put(Routes.applicationCommands(clientID as string), {
      body: commands,
    })
    .then(() => console.log("Successfully registered application commands."))
    .catch((e) =>
      console.error(`Failed to register application commands: ${e}`)
    );
};

export { initializeCommands };
