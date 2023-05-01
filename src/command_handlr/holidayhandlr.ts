import { CacheType, ChatInputCommandInteraction } from "discord.js";
import { sendWebhookErrorMessage } from "../webhook/send_message";

type Holiday = {
  index: number;
  start: Date;
  end: Date;
  isOneDay: boolean;
  message: string;
};

export const holidayHandlr = async (
  interaction: ChatInputCommandInteraction<CacheType>
) => {
  // Jun 26 – Jul 3, 2023 Mon - Mon Holiday: *Eid-ul-Azha
  // Jul 29, 2023 Sat Holiday: *Ashura
  // Aug 15, 2023 Tue Holiday: National Mourning Day
  // Sep 6, 2023 Wed Holiday: Janmashtami

  const installmentList: Holiday[] = [
    {
      index: 1,
      start: new Date("Jun 26, 2023"),
      end: new Date("Jul 3, 2023"),
      message: "Eid-ul-Azha",
      isOneDay: false,
    },
    {
      index: 2,
      start: new Date("Jul 29, 2023"),
      end: new Date("Jul 29, 2023"),
      message: "Ashura",
      isOneDay: true,
    },
    {
      index: 3,
      start: new Date("Aug 15, 2023"),
      end: new Date("Aug 15, 2023"),
      message: "National Mourning Day",
      isOneDay: true,
    },
    {
      index: 4,
      start: new Date("Sep 6, 2023"),
      end: new Date("Sep 6, 2023"),
      message: "Janmashtami",
      isOneDay: true,
    },
  ];

  await interaction
    .reply({
      content:
        `**😀 List of upcoming holidays** In total only ${installmentList.length} events ☹️\n\n` +
        installmentList
          .map((item) => {
            return `${item.index}. **${item.message}** ${
              item.isOneDay
                ? item.start.toDateString()
                : `from ${item.start.toDateString()} to ${item.end.toDateString()} in total ${Math.round(
                    (item.end.getTime() - item.start.getTime()) /
                      (1000 * 60 * 60 * 24)
                  )} days`
            }`;
          })
          .join("\n"),
    })
    .catch((err) => {
      console.log("🚀 ~ file: holidayhandlr.ts:74 ~ err:", err);

      sendWebhookErrorMessage("holidayHandlr", err);
      interaction.followUp("Error occurred while processing the command");
    });
};
