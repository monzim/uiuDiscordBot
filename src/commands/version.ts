import { SlashCommandBuilder } from "discord.js";

module.exports = {
  data: new SlashCommandBuilder()
    .setName("version")
    .setDescription("Get upcoming makeup class info!"),
};
