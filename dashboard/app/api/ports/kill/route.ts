import { exec } from "child_process";
import { promisify } from "util";
import { NextResponse } from "next/server";

const execAsync = promisify(exec);
const CLI_PATH = process.env.IDOPS_CLI_PATH || "idops";

export async function POST(request: Request) {
  try {
    const { port } = await request.json();

    if (!port) {
      return NextResponse.json(
        { error: "Port required" },
        { status: 400 }
      );
    }

    const { stdout, stderr } = await execAsync(`${CLI_PATH} ports kill ${port}`);
    return NextResponse.json({
      success: true,
      message: stdout || stderr || `Killed process on port ${port}`,
    });
  } catch (error) {
    console.error("Port kill error:", error);
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    return NextResponse.json(
      { error: errorMessage },
      { status: 500 }
    );
  }
}
