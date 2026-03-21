import { exec } from "child_process";
import { promisify } from "util";
import { NextResponse } from "next/server";

const execAsync = promisify(exec);
const CLI_PATH = process.env.IDOPS_CLI_PATH || "idops";

export async function POST(request: Request) {
  try {
    const { action, containerId } = await request.json();

    let command: string;
    switch (action) {
      case "start":
        command = `${CLI_PATH} docker start ${containerId}`;
        break;
      case "stop":
        command = `${CLI_PATH} docker stop ${containerId}`;
        break;
      case "restart":
        command = `${CLI_PATH} docker restart ${containerId}`;
        break;
      case "remove":
        command = `${CLI_PATH} docker rm ${containerId}`;
        break;
      default:
        return NextResponse.json({ error: "Invalid action" }, { status: 400 });
    }

    const { stdout, stderr } = await execAsync(command);
    return NextResponse.json({
      success: true,
      message: stdout || stderr || `${action} completed`,
    });
  } catch (error) {
    console.error("Docker action error:", error);
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    return NextResponse.json(
      { error: errorMessage },
      { status: 500 }
    );
  }
}
