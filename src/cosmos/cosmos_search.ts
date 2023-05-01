import { CosmosClient } from "@azure/cosmos";
import Course from "./../models/course";
import { sendWebhookErrorMessage } from "../webhook/send_message";
import {
  COSMOS_CONTAINER_ID,
  COSMOS_DATABASE_ID,
  COSMOS_KEY,
  COSMOS_ENDPOINT,
} from "../config";

const cosmosClient = new CosmosClient({
  endpoint: COSMOS_ENDPOINT as string,
  key: COSMOS_KEY as string,
});

export const findCoursesByCodeAndSection = async (
  courseCode: string,
  section: string
): Promise<Course[]> => {
  try {
    const { database } = await cosmosClient.databases.createIfNotExists({
      id: COSMOS_DATABASE_ID as string,
    });
    const { container } = await database.containers.createIfNotExists({
      id: COSMOS_CONTAINER_ID as string,
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
    sendWebhookErrorMessage("findCoursesByCodeAndSection", error);
    return [];
  }
};

export const findCoursesByDepartmentCourseCodeAndSection = async (
  department: string,
  courseCode: string,
  section: string
): Promise<Course[]> => {
  try {
    const { database } = await cosmosClient.databases.createIfNotExists({
      id: COSMOS_DATABASE_ID as string,
    });
    const { container } = await database.containers.createIfNotExists({
      id: COSMOS_CONTAINER_ID as string,
    });

    // Query the container for courses with matching roll
    const query = `SELECT * FROM c WHERE c['Dept.'] = '${department}' AND CONTAINS(c['Course Code'], '${courseCode}') AND CONTAINS(c['Section'], '${section}')`;

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
    sendWebhookErrorMessage(
      "findCoursesByDepartmentCourseCodeAndSection",
      error
    );
    return [];
  }
};
