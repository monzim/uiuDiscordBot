import { SlashCommandBuilder } from "discord.js";

module.exports = {
  data: new SlashCommandBuilder()
    .setName("version")
    .setDescription("Shows the current version of the bot!"),
};
