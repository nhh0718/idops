import { exec } from "child_process";
import { promisify } from "util";
import { NextResponse } from "next/server";

const execAsync = promisify(exec);
const CLI_PATH = process.env.IDOPS_CLI_PATH || "idops";

export async function GET(request: Request) {
  try {
    const { searchParams } = new URL(request.url);
    const containerId = searchParams.get("containerId");

    if (!containerId) {
      return NextResponse.json(
        { error: "Container ID required" },
        { status: 400 }
      );
    }

    const { stdout } = await execAsync(`${CLI_PATH} docker logs ${containerId}`);
    return NextResponse.json({ logs: stdout });
  } catch (error) {
    console.error("Docker logs error:", error);
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    return NextResponse.json(
      { error: errorMessage, logs: "" },
      { status: 500 }
    );
  }
}
