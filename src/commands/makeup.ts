import { SlashCommandBuilder } from "discord.js";

module.exports = {
  data: new SlashCommandBuilder()
    .setName("makeup")
    .setDescription("Get upcoming makeup class info!"),
};
