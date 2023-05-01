import { SlashCommandBuilder, Interaction } from "discord.js";

module.exports = {
  data: new SlashCommandBuilder()
    .setName("installment")
    .setDescription("Get info about upcoming installment!"),
};
