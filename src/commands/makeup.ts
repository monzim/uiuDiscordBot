import { SlashCommandBuilder } from "discord.js";

module.exports = {
  data: new SlashCommandBuilder()
    .setName("makeup")
    .setDescription("Replies with upcoming makeup!"),
};
