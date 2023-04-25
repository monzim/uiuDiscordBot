import { CosmosClient } from "@azure/cosmos";
import * as config from "./config.json"; // Import the config file
import Course from "./models/course";

const { endpoint, key, databaseId, containerId } = config;

// Connect to Cosmos DB account
const cosmosClient = new CosmosClient({ endpoint, key });

export const findCoursesByCodeAndSection = async (
  courseCode: string,
  section: string
): Promise<Course[]> => {
  try {
    const { database } = await cosmosClient.databases.createIfNotExists({
      id: databaseId,
    });
    const { container } = await database.containers.createIfNotExists({
      id: containerId,
    });

    // Query the container for courses with matching roll
    const query = `SELECT * FROM c WHERE CONTAINS(c['Course Code'], '${courseCode}') AND CONTAINS(c['Section'], '${section}')`;

    const { resources: courses } = await container.items
      .query(query)
      .fetchAll();

    // Map the query result to Course objects
    const mappedCourses: Course[] = courses.map((course: any) => {
      return {
        Dept: course["Dept."],
        CourseCode: course["Course Code"],
        CourseTitle: course["Course Title"],
        Section: course.Section,
        Teacher: course.Teacher,
        ExamDate: course["Exam Date"],
        ExamTime: course["Exam Time"],
        Room: course.Room,
        id: course.id,
        _rid: course._rid,
        _self: course._self,
        _etag: course._etag,
        _attachments: course._attachments,
        _ts: course._ts,
      };
    });
    console.log(
      "🚀 ~ file: cosmos_search.ts:51 ~ constmappedCourses:Course[]=courses.map ~ mappedCourses:",
      mappedCourses.length
    );

    return mappedCourses;
  } catch (error) {
    console.log("🚀 ~ file: cosmos_search.ts:58 ~ error:", error);
    return [];
  }
};
