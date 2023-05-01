import { SlashCommandBuilder } from "discord.js";

module.exports = {
  data: new SlashCommandBuilder()
    .setName("exam")
    .setDescription("Replies with upcoming exam!"),
};
