import { SlashCommandBuilder, Interaction } from "discord.js";

module.exports = {
  data: new SlashCommandBuilder()
    .setName("installment")
    .setDescription("Replies with upcoming installment!"),
};
