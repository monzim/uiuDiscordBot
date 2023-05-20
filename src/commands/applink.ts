import { SlashCommandBuilder } from "discord.js";

module.exports = {
  data: new SlashCommandBuilder()
    .setName("applink")
    .setDescription("Get the Unizim App Public Beta Testing Link!"),
};
