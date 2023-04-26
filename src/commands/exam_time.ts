import { SlashCommandBuilder, Interaction } from "discord.js";

module.exports = {
  data: new SlashCommandBuilder()
    .setName("exam_time")
    .setDescription("Get course exam info like date, time, room, etc.")
    .addStringOption((option) =>
      option
        .setName("department")
        .setDescription("Choose the department")
        .setRequired(true)
        .addChoices(
          { name: "BBA", value: "BBA" },
          { name: "BSCSE", value: "BSCSE" },
          { name: "BSEEE", value: "BSEEE" }
        )
    )
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
