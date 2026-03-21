import { exec } from "child_process";
import { promisify } from "util";
import { NextResponse } from "next/server";

const execAsync = promisify(exec);

// Get CLI path - assumes idops binary is in PATH or use full path
const CLI_PATH = process.env.IDOPS_CLI_PATH || "idops";

export async function GET() {
  try {
    const { stdout } = await execAsync(`${CLI_PATH} docker --json`);
    const containers = JSON.parse(stdout);
    return NextResponse.json({ containers });
  } catch (error) {
    console.error("Docker list error:", error);
    return NextResponse.json(
      { error: "Failed to list containers", containers: [] },
      { status: 500 }
    );
  }
}
