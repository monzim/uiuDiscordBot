import { SlashCommandBuilder, Interaction } from "discord.js";

module.exports = {
  data: new SlashCommandBuilder()
    .setName("exam")
    .setDescription("Replies with Pong!")
    .addStringOption((option) =>
      option
        .setName("course")
        .setDescription("Enter the course code")
        .setRequired(true)
        .setMinLength(4)
        .setMaxLength(8)
    )

    .addStringOption((option) =>
      option
        .setName("section")
        .setDescription("Enter the section")
        .setRequired(true)
        .setMinLength(1)
        .setMaxLength(2)
    ),
};
