import { SlashCommandBuilder } from "discord.js";

module.exports = {
  data: new SlashCommandBuilder()
    .setName("holiday")
    .setDescription("Replies with upcoming holiday!"),
};
