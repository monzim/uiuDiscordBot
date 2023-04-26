import { CosmosClient } from "@azure/cosmos";
import * as fs from "fs";
import * as config from "../src/config.json";

const filesPath: string[] = [
  "data/spring_23_final/bscse_bseee.json",
  "data/spring_23_final/bba.json",
];
const { endpoint, key, databaseId, containerId } = config;
const client = new CosmosClient({ endpoint, key });

const insertFileData = async (filePath: string) => {
  // Read the JSON data from file
  const jsonData = JSON.parse(fs.readFileSync(filePath, "utf8"));

  let count = 0;

  try {
    const { database } = await client.databases.createIfNotExists({
      id: databaseId,
    });
    const { container } = await database.containers.createIfNotExists({
      id: containerId,
    });

    for (const data of jsonData) {
      count++;

      await container.items
        .create(data)
        .then((result) => {
          console.log(
            `Count: ${count} Inserted item with id: ${result.resource.id}`
          );
        })
        .catch((err) => {
          console.log(`Count: ${count} Error inserting data: ${err}`);
        });
    }

    console.log(`Total data inserted: ${count}`);
  } catch (error) {
    console.log("🚀 ~ insertReadDataData ~ error:", error);
  }
};

async function main() {
  for (const filePath of filesPath) {
    console.log(`Inserting data from ${filePath}`);
    await insertFileData(filePath);
  }
}

main();
