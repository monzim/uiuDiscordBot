import { SlashCommandBuilder, Interaction } from "discord.js";

module.exports = {
  data: new SlashCommandBuilder()
    .setName("donation")
    .setDescription("Replies with Donation info!"),
};
