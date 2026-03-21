import { exec } from "child_process";
import { promisify } from "util";
import { NextResponse } from "next/server";

const execAsync = promisify(exec);
const CLI_PATH = process.env.IDOPS_CLI_PATH || "idops";

export async function GET() {
  try {
    const { stdout } = await execAsync(`${CLI_PATH} nginx list --json`);
    const configs = JSON.parse(stdout);
    return NextResponse.json({ configs });
  } catch (error) {
    console.error("Nginx list error:", error);
    return NextResponse.json(
      { error: "Failed to list nginx configs", configs: [] },
      { status: 500 }
    );
  }
}
